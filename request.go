package gokord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
)

// All error constants
var (
	ErrJSONUnmarshal = errors.New("json unmarshal")
	ErrStatusOffline = errors.New("you can't set your Status to offline")
	ErrUnauthorized  = errors.New("HTTP request was unauthorized. This could be because the provided token was not a bot token. Please add \"Bot \" to the start of your token. https://discord.com/developers/docs/reference#authentication-example-bot-token-authorization-header")
)

// RESTError stores error information about a request with a bad response code.
// Message is not always present, there are cases where api calls can fail
// without returning a json message.
type RESTError struct {
	Request      *http.Request
	Response     *http.Response
	ResponseBody []byte

	// Message may be nil.
	Message *discord.APIErrorMessage
}

// newRestError returns a new REST API error.
func newRestError(req *http.Request, resp *http.Response, body []byte) *RESTError {
	restErr := &RESTError{
		Request:      req,
		Response:     resp,
		ResponseBody: body,
	}

	// Attempt to decode the error and assume no message was provided if it fails
	var msg *discord.APIErrorMessage
	err := json.Unmarshal(body, &msg)
	if err == nil {
		restErr.Message = msg
	}

	return restErr
}

// Error returns a Rest API Error with its status code and body.
func (r RESTError) Error() string {
	base := fmt.Sprintf("[HTTP %d]", r.Response.StatusCode)
	if r.Message != nil {
		return fmt.Sprintf("%s %s\n%s", base, r.Message.Error(), r.ResponseBody)
	}
	return fmt.Sprintf("%s %s", base, r.ResponseBody)
}

// RateLimitError is returned when a request exceeds a rate limit and Session.ShouldRetryOnRateLimit is false.
// The request may be manually retried after waiting the duration specified by RetryAfter.
type RateLimitError struct {
	*event.RateLimit
}

// Error returns a rate limit error with rate limited endpoint and retry time.
func (e RateLimitError) Error() string {
	return "rate limit exceeded on " + e.URL + ", retrying after " + e.RetryAfter.String()
}

func unmarshal(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return errors.Join(ErrJSONUnmarshal, err)
	}

	return nil
}

// RESTSession is the part of the Session responsible for the REST API.
type RESTSession struct {
	identify *Identify
	logger   *slog.Logger
	// Should the session retry requests when rate limited.
	ShouldRetryOnRateLimit bool
	eventManager           *event.Manager
	// Max number of REST API retries.
	MaxRestRetries int
	// The http.Client used for REST requests.
	Client *http.Client
	// The UserAgent used for REST APIs.
	UserAgent string
	// Used to deal with rate limits.
	RateLimiter        *discord.RateLimiter
	emitRateLimitEvent func(ctx context.Context, evt *event.RateLimit)
}

func (s *RESTSession) Logger() *slog.Logger {
	return s.logger
}

func (s *RESTSession) Unmarshal(bytes []byte, i any) error {
	return unmarshal(bytes, i)
}

func (s *RESTSession) Request(method, urlStr string, data any, options ...discord.RequestOption) ([]byte, error) {
	return s.RequestWithBucketID(method, urlStr, data, strings.SplitN(urlStr, "?", 2)[0], options...)
}

func (s *RESTSession) RequestWithBucketID(method, urlStr string, data any, bucketID string, options ...discord.RequestOption) ([]byte, error) {
	var body []byte
	if data != nil {
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	return s.RequestRaw(method, urlStr, "application/json", body, bucketID, 0, options...)
}

func (s *RESTSession) RequestRaw(method, urlStr, contentType string, b []byte, bucketID string, sequence int, options ...discord.RequestOption) ([]byte, error) {
	if bucketID == "" {
		bucketID = strings.SplitN(urlStr, "?", 2)[0]
	}
	return s.RequestWithLockedBucket(method, urlStr, contentType, b, s.RateLimiter.LockBucket(bucketID), sequence, options...)
}

func (s *RESTSession) RequestWithLockedBucket(method, urlStr, contentType string, b []byte, bucket *discord.Bucket, sequence int, options ...discord.RequestOption) ([]byte, error) {
	s.logger.Debug(fmt.Sprintf("%s :: %s", method, urlStr))
	s.logger.Debug("PAYLOAD", "content", string(b))

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(b))
	if err != nil {
		err2 := bucket.Release(nil)
		if err2 != nil {
			err = errors.Join(err, err2)
		}
		return nil, err
	}

	// Not used on initial login..
	// TODO: Verify if a login, otherwise complain about no-token
	if s.identify.Token != "" {
		req.Header.Set("authorization", s.identify.Token)
	}

	// Discord's API returns a 400 Bad Request is Content-Type is set, but the
	// request body is empty.
	if b != nil {
		req.Header.Set("Content-Type", contentType)
	}

	// TODO: Make a configurable static variable.
	req.Header.Set("User-Agent", s.UserAgent)

	cfg := &discord.RequestConfig{
		ShouldRetryOnRateLimit: s.ShouldRetryOnRateLimit,
		MaxRestRetries:         s.MaxRestRetries,
		Client:                 s.Client,
		Request:                req,
	}
	for _, opt := range options {
		opt(cfg)
	}
	req = cfg.Request

	resp, err := cfg.Client.Do(req)
	if err != nil {
		err2 := bucket.Release(nil)
		if err2 != nil {
			err = errors.Join(err, err2)
		}
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			s.logger.Error("closing resp body", "error", err)
		}
	}()

	err = bucket.Release(resp.Header)
	if err != nil {
		return nil, err
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusNoContent:
	case http.StatusBadGateway:
		// Retry sending request if possible
		if sequence < cfg.MaxRestRetries {
			s.logger.Warn("failed, retrying...", "url", urlStr, "status", resp.Status)
			response, err = s.RequestWithLockedBucket(method, urlStr, contentType, b, s.RateLimiter.LockBucketObject(bucket), sequence+1, options...)
		} else {
			err = fmt.Errorf("exceeded max HTTP retries %s, %s", resp.Status, response)
		}
	case http.StatusTooManyRequests: // rate limiting
		rl := discord.TooManyRequests{}
		err = json.Unmarshal(response, &rl)
		if err != nil {
			s.logger.Error("rate limit unmarshal", "error", err)
			return nil, err
		}

		if cfg.ShouldRetryOnRateLimit {
			s.logger.Info("rate limited", "url", urlStr, "retry in", rl.RetryAfter)
			// background because it will never use the websocket -> this is an internal event
			s.emitRateLimitEvent(context.Background(), &event.RateLimit{TooManyRequests: &rl, URL: urlStr})

			time.Sleep(rl.RetryAfter)
			// we can make the above smarter
			// this method can cause longer delays than required

			response, err = s.RequestWithLockedBucket(method, urlStr, contentType, b, s.RateLimiter.LockBucketObject(bucket), sequence, options...)
		} else {
			err = &RateLimitError{&event.RateLimit{TooManyRequests: &rl, URL: urlStr}}
		}
	case http.StatusUnauthorized:
		if strings.Index(s.identify.Token, "Bot ") != 0 {
			s.logger.Error(ErrUnauthorized.Error())
			err = ErrUnauthorized
		}
		fallthrough
	default: // Error condition
		err = newRestError(req, resp, response)
	}

	return response, err
}

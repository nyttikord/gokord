package gokord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nyttikord/gokord/discord"
)

// RequestConfig is an HTTP request configuration.
type RequestConfig struct {
	Request                *http.Request
	ShouldRetryOnRateLimit bool
	MaxRestRetries         int
	Client                 *http.Client
}

// newRequestConfig returns a new HTTP request configuration based on parameters in Session.
func newRequestConfig(s *Session, req *http.Request) *RequestConfig {
	return &RequestConfig{
		ShouldRetryOnRateLimit: s.ShouldRetryOnRateLimit,
		MaxRestRetries:         s.MaxRestRetries,
		Client:                 s.Client,
		Request:                req,
	}
}

// RequestOption is a function which mutates request configuration.
// It can be supplied as an argument to any REST method.
type RequestOption func(cfg *RequestConfig)

// WithClient changes the HTTP client used for the request.
func WithClient(client *http.Client) RequestOption {
	return func(cfg *RequestConfig) {
		if client != nil {
			cfg.Client = client
		}
	}
}

// WithRetryOnRatelimit controls whether session will retry the request on rate limit.
func WithRetryOnRatelimit(retry bool) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.ShouldRetryOnRateLimit = retry
	}
}

// WithRestRetries changes maximum amount of retries if request fails.
func WithRestRetries(max int) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.MaxRestRetries = max
	}
}

// WithHeader sets a header in the request.
func WithHeader(key, value string) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.Request.Header.Set(key, value)
	}
}

// WithAuditLogReason changes audit Log reason associated with the request.
func WithAuditLogReason(reason string) RequestOption {
	return WithHeader("X-Audit-Log-Reason", reason)
}

// WithLocale changes accepted locale of the request.
func WithLocale(locale discord.Locale) RequestOption {
	return WithHeader("X-Discord-Locale", string(locale))
}

// WithContext changes context of the request.
func WithContext(ctx context.Context) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.Request = cfg.Request.WithContext(ctx)
	}
}

// Request is the same as RequestWithBucketID but the bucket id is the same as the urlStr
func (s *Session) Request(method, urlStr string, data interface{}, options ...RequestOption) ([]byte, error) {
	return s.RequestWithBucketID(method, urlStr, data, strings.SplitN(urlStr, "?", 2)[0], options...)
}

// RequestWithBucketID makes a (GET/POST/...) http.Request to Discord REST API with JSON data.
func (s *Session) RequestWithBucketID(method, urlStr string, data interface{}, bucketID string, options ...RequestOption) ([]byte, error) {
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

// RequestRaw makes a (GET/POST/...) Requests to Discord REST API.
// Preferably use the other request methods but this lets you send JSON directly if that's what you have.
//
// sequence is the sequence number, if it fails with a 502 it will retry with sequence+1 until it either succeeds or
// sequence >= Session.MaxRestRetries
func (s *Session) RequestRaw(method, urlStr, contentType string, b []byte, bucketID string, sequence int, options ...RequestOption) ([]byte, error) {
	if bucketID == "" {
		bucketID = strings.SplitN(urlStr, "?", 2)[0]
	}
	return s.RequestWithLockedBucket(method, urlStr, contentType, b, s.RateLimiter.LockBucket(bucketID), sequence, options...)
}

// RequestWithLockedBucket makes a request using a bucket that's already been locked
func (s *Session) RequestWithLockedBucket(method, urlStr, contentType string, b []byte, bucket *Bucket, sequence int, options ...RequestOption) ([]byte, error) {
	s.LogDebug("API REQUEST %8s :: %s\n", method, urlStr)
	s.LogDebug("API REQUEST  PAYLOAD :: [%s]\n", string(b))

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
	if s.Identify.Token != "" {
		req.Header.Set("authorization", s.Identify.Token)
	}

	// Discord's API returns a 400 Bad Request is Content-Type is set, but the
	// request body is empty.
	if b != nil {
		req.Header.Set("Content-Type", contentType)
	}

	// TODO: Make a configurable static variable.
	req.Header.Set("User-Agent", s.UserAgent)

	cfg := newRequestConfig(s, req)
	for _, opt := range options {
		opt(cfg)
	}
	req = cfg.Request

	for k, v := range req.Header {
		s.LogDebug("API REQUEST   HEADER :: [%s] = %+v\n", k, v)
	}

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
			s.LogDebug("error closing resp body: %v", err)
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

	s.LogDebug("API RESPONSE  STATUS :: %s\n", resp.Status)
	for k, v := range resp.Header {
		s.LogDebug("API RESPONSE  HEADER :: [%s] = %+v\n", k, v)
	}
	s.LogDebug("API RESPONSE    BODY :: [%s]\n\n\n", response)

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusNoContent:
	case http.StatusBadGateway:
		// Retry sending request if possible
		if sequence < cfg.MaxRestRetries {
			s.LogInfo("%s Failed (%s), Retrying...", urlStr, resp.Status)
			response, err = s.RequestWithLockedBucket(method, urlStr, contentType, b, s.RateLimiter.LockBucketObject(bucket), sequence+1, options...)
		} else {
			err = fmt.Errorf("exceeded max HTTP retries %s, %s", resp.Status, response)
		}
	case http.StatusTooManyRequests: // rate limiting
		rl := TooManyRequests{}
		err = json.Unmarshal(response, &rl)
		if err != nil {
			s.LogError("rate limit unmarshal error, %s", err)
			return nil, err
		}

		if cfg.ShouldRetryOnRateLimit {
			s.LogInfo("Rate Limiting %s, retry in %v", urlStr, rl.RetryAfter)
			s.handleEvent(rateLimitEventType, &RateLimit{TooManyRequests: &rl, URL: urlStr})

			time.Sleep(rl.RetryAfter)
			// we can make the above smarter
			// this method can cause longer delays than required

			response, err = s.RequestWithLockedBucket(method, urlStr, contentType, b, s.RateLimiter.LockBucketObject(bucket), sequence, options...)
		} else {
			err = &RateLimitError{&RateLimit{TooManyRequests: &rl, URL: urlStr}}
		}
	case http.StatusUnauthorized:
		if strings.Index(s.Identify.Token, "Bot ") != 0 {
			s.LogInfo("%s", ErrUnauthorized.Error())
			err = ErrUnauthorized
		}
		fallthrough
	default: // Error condition
		err = newRestError(req, resp, response)
	}

	return response, nil
}

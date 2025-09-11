package gokord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nyttikord/gokord/discord"
)

func unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrJSONUnmarshal, err)
	}

	return nil
}

func (s *Session) unmarshal(bytes []byte, i interface{}) error {
	return unmarshal(bytes, i)
}

func (s *Session) Request(method, urlStr string, data interface{}, options ...discord.RequestOption) ([]byte, error) {
	return s.RequestWithBucketID(method, urlStr, data, strings.SplitN(urlStr, "?", 2)[0], options...)
}

func (s *Session) RequestWithBucketID(method, urlStr string, data interface{}, bucketID string, options ...discord.RequestOption) ([]byte, error) {
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

func (s *Session) RequestRaw(method, urlStr, contentType string, b []byte, bucketID string, sequence int, options ...discord.RequestOption) ([]byte, error) {
	if bucketID == "" {
		bucketID = strings.SplitN(urlStr, "?", 2)[0]
	}
	return s.RequestWithLockedBucket(method, urlStr, contentType, b, s.RateLimiter.LockBucket(bucketID), sequence, options...)
}

func (s *Session) RequestWithLockedBucket(method, urlStr, contentType string, b []byte, bucket *discord.Bucket, sequence int, options ...discord.RequestOption) ([]byte, error) {
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

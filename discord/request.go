package discord

import (
	"context"
	"net/http"
)

// Requester is used to interact with the Discord API.
type Requester interface {
	// Request is the same as RequestWithBucketID but the bucket id is the same as the urlStr
	Request(method string, urlStr string, data interface{}, options ...RequestOption) ([]byte, error)
	// RequestWithBucketID makes a (GET/POST/...) http.Request to Discord REST API with JSON data.
	RequestWithBucketID(method string, urlStr string, data interface{}, bucketID string, options ...RequestOption) ([]byte, error)
	// RequestRaw makes a (GET/POST/...) Requests to Discord REST API.
	// Preferably use the other request methods but this lets you send JSON directly if that's what you have.
	//
	// sequence is the sequence number, if it fails with a 502 it will retry with sequence+1 until it either succeeds or
	// sequence >= Session.MaxRestRetries
	RequestRaw(method string, urlStr string, contentType string, data []byte, bucketID string, sequence int, options ...RequestOption) ([]byte, error)
	// RequestWithLockedBucket makes a request using a bucket that's already been locked
	RequestWithLockedBucket(method string, urlStr string, contentType string, data []byte, bucket *Bucket, sequence int, options ...RequestOption) ([]byte, error)
	// VoiceRegions returns the VoiceRegion
	VoiceRegions(options ...RequestOption) ([]*VoiceRegion, error)
	// Unmarshal is for unmarshalling body returned by the Discord API.
	Unmarshal(bytes []byte, i interface{}) error
}

// RequestConfig is an HTTP request configuration.
type RequestConfig struct {
	Request                *http.Request
	ShouldRetryOnRateLimit bool
	MaxRestRetries         int
	Client                 *http.Client
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
func WithLocale(locale Locale) RequestOption {
	return WithHeader("X-Discord-Locale", string(locale))
}

// WithContext changes context of the request.
func WithContext(ctx context.Context) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.Request = cfg.Request.WithContext(ctx)
	}
}

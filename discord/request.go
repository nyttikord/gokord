package discord

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/coder/websocket"
)

// RESTRequester is used to interact with the Discord API.
type RESTRequester interface {
	Logger() *slog.Logger
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
	// Unmarshal is for unmarshalling body returned by the Discord API.
	Unmarshal(bytes []byte, i interface{}) error
}

type WSRequester interface {
	// GatewayWriteStruct writes a struck as a json to Discord gateway.
	GatewayWriteStruct(context.Context, any) error
	// GatewayDial dials a new websocket connection.
	GatewayDial(context.Context, string, http.Header) (*websocket.Conn, *http.Response, error)
}

// RequestConfig is an HTTP request configuration.
type RequestConfig struct {
	Request                *http.Request
	ShouldRetryOnRateLimit bool
	MaxRestRetries         int
	Client                 *http.Client
}

// RequestOption is a function which modifies how the request is handled.
// It can be supplied as an argument to any REST method.
//
// You can call WithContext to use a context.Context during the request.
type RequestOption func(cfg *RequestConfig)

// WithClient changes the HTTP client used for the request.
func WithClient(client *http.Client) RequestOption {
	return func(cfg *RequestConfig) {
		if client != nil {
			cfg.Client = client
		}
	}
}

// WithRetryOnRatelimit controls whether the session should retry the request on rate limit.
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

// WithAuditLogReason changes audit log reason associated with the request.
func WithAuditLogReason(reason string) RequestOption {
	return WithHeader("X-Audit-Log-Reason", reason)
}

// WithLocale changes accepted locale of the request.
func WithLocale(locale Locale) RequestOption {
	return WithHeader("X-Discord-Locale", string(locale))
}

// WithContext changes context.Context of the request.
func WithContext(ctx context.Context) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.Request = cfg.Request.WithContext(ctx)
	}
}

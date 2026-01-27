package request

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/coder/websocket"
)

// REST is used to interact with the Discord REST API.
type REST interface {
	Logger() *slog.Logger
	// Request is the same as RequestWithBucketID but the bucket id is the same as the urlStr
	Request(ctx context.Context, method string, urlStr string, data any, option Config) ([]byte, error)
	// RequestWithBucketID makes a (GET/POST/...) http.Request to Discord REST API with JSON data.
	RequestWithBucketID(ctx context.Context, method string, urlStr string, data any, bucketID string, option Config) ([]byte, error)
	// RequestRaw makes a (GET/POST/...) Requests to Discord REST API.
	// Preferably use the other request methods but this lets you send JSON directly if that's what you have.
	//
	// sequence is the sequence number, if it fails with a 502 it will retry with sequence+1 until it either succeeds or
	// sequence >= Session.MaxRestRetries
	RequestRaw(ctx context.Context, method string, urlStr string, contentType string, data []byte, bucketID string, sequence uint, option Config) ([]byte, error)
	// RequestWithLockedBucket makes a request using a bucket that's already been locked
	RequestWithLockedBucket(ctx context.Context, method string, urlStr string, contentType string, data []byte, bucket *Bucket, sequence uint, option Config) ([]byte, error)
	// Unmarshal is for unmarshalling body returned by the Discord API.
	Unmarshal(bytes []byte, i any) error
}

// REST is used to interact with the Discord websocket API.
type Websocket interface {
	// GatewayWriteStruct writes a struck as a json to Discord gateway.
	GatewayWriteStruct(context.Context, any) error
	// GatewayDial dials a new websocket connection.
	GatewayDial(context.Context, string, http.Header) (*websocket.Conn, *http.Response, error)
}

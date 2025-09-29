package gokord

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	// ErrWSAlreadyOpen is thrown when you attempt to open a websocket that already is open.
	ErrWSAlreadyOpen = errors.New("web socket already opened")
	// ErrWSNotFound is thrown when you attempt to use a websocket that doesn't exist
	ErrWSNotFound = errors.New("no websocket connection exists")
	// ErrWSShardBounds is thrown when you try to use a shard ID that is more than the total shard count
	ErrWSShardBounds = errors.New("ShardID must be less than ShardCount")
)

func (s *Session) GatewayWriteStruct(v any) error {
	s.wsMutex.Lock()
	defer s.wsMutex.Unlock()
	if s.ws == nil {
		return ErrWSNotFound
	}
	return s.ws.WriteJSON(v)
}

func (s *Session) GatewayReady() bool {
	if s.ws == nil {
		return false
	}
	return s.DataReady
}

func (s *Session) GatewayDial(ctx context.Context, urlString string, headers http.Header) (*websocket.Conn, *http.Response, error) {
	return s.Dialer.DialContext(ctx, urlString, headers)
}

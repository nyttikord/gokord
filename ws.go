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
	s.RLock()
	defer s.RUnlock()
	if s.ws == nil {
		return ErrWSNotFound
	}

	s.wsMutex.Lock()
	err := s.ws.WriteJSON(v)
	s.wsMutex.Unlock()
	return err
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

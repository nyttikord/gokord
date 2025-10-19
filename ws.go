package gokord

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/coder/websocket"
)

var (
	// ErrWSAlreadyOpen is thrown when you attempt to open a websocket that already is open.
	ErrWSAlreadyOpen = errors.New("web socket already opened")
	// ErrWSNotFound is thrown when you attempt to use a websocket that doesn't exist
	ErrWSNotFound = errors.New("no websocket connection exists")
	// ErrWSShardBounds is thrown when you try to use a shard ID that is more than the total shard count
	ErrWSShardBounds = errors.New("ShardID must be less than ShardCount")
)

func (s *Session) GatewayWriteStruct(ctx context.Context, v any) error {
	s.wsMutex.Lock()
	defer s.wsMutex.Unlock()
	if s.ws == nil {
		return ErrWSNotFound
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.ws.Write(ctx, websocket.MessageText, b)
}

func (s *Session) GatewayDial(ctx context.Context, urlString string, headers http.Header) (*websocket.Conn, *http.Response, error) {
	return websocket.Dial(ctx, urlString, &websocket.DialOptions{HTTPHeader: headers})
}

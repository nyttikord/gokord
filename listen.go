package gokord

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/logger"
)

func (s *Session) setupListen(ctx context.Context) {
	if s.wsRead != nil {
		s.logger.Info("listen already running")
		return
	}
	ctx2, cancel := context.WithCancel(ctx)
	s.cancelWSRead = cancel

	wsRead := make(chan readResult)
	s.wsRead = wsRead
	go func() {
		s.logger.Info("listening started")
		err := s.listen(ctx2, wsRead)
		s.wsRead = nil
		select {
		case <-ctx2.Done():
			return
		default:
			s.logger.Warn("listening websocket", "error", err, "gateway", s.gateway)
			s.forceReconnect(ctx)
		}
	}()
}

// listen polls the websocket connection for data, it will stop when an error occurs.
func (s *Session) listen(ctx context.Context, c chan<- readResult) error {
	var messageType websocket.MessageType
	var message []byte
	var err error
	for err == nil {
		messageType, message, err = s.ws.Read(ctx)
		if err == nil {
			c <- readResult{messageType, message}
		}
	}
	return err
}

// syncListener must not be copied!
type syncListener struct {
	wg      sync.WaitGroup
	logger  *slog.Logger
	cancel  func()
	counter atomic.Uint32
}

func (s *syncListener) Add(fn func(free func())) {
	s.wg.Add(1)
	s.counter.Add(1)
	go fn(func() {
		s.wg.Done()
		s.counter.Store(s.counter.Load() - 1)
	})
}

func (s *syncListener) Wait(ctx context.Context) error {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	done := make(chan struct{}, 1)
	go func() {
		s.wg.Wait()
		s.logger.Debug("goroutines closed")
		done <- struct{}{}
	}()
	select {
	case <-done:
	case <-ctx2.Done():
		s.logger.Error("cannot close goroutines", "remaining", s.counter.Load())
		return ctx2.Err()
	}
	s.cancel = nil
	return nil
}

func (s *syncListener) Close(ctx context.Context) {
	if s.cancel == nil {
		s.logger.WarnContext(logger.NewContext(context.Background(), 1), "cancel func was already called (or was never set)")
		return
	}
	s.logger.Debug("closing goroutines")
	s.cancel()
}

type readResult struct {
	MessageType websocket.MessageType
	Message     []byte
}

func (r *readResult) getEvent() (*discord.Event, error) {
	return getGatewayEvent(r.MessageType, r.Message)
}

// dispatch the event received
func (r *readResult) dispatch(s *Session, ctx context.Context) error {
	e, err := r.getEvent()
	if err == nil {
		err = s.onGatewayEvent(ctx, e)
	}
	return err
}

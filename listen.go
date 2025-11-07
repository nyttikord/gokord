package gokord

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nyttikord/gokord/logger"
)

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

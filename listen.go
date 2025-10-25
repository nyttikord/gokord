package gokord

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// syncListener must not be copied!
type syncListener struct {
	wg     sync.WaitGroup
	logger *slog.Logger
	cancel func()
}

func (m *syncListener) Add(fn func(free func())) {
	m.wg.Add(1)
	go fn(func() {
		m.wg.Done()
	})
}

func (m *syncListener) Wait() {
	m.wg.Wait()
}

func (m *syncListener) Close(ctx context.Context) error {
	if m.cancel == nil {
		return nil
	}
	m.logger.Debug("closing goroutines")
	m.cancel()
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	done := make(chan struct{}, 1)
	go func() {
		m.Wait()
		m.logger.Debug("goroutines closed")
		done <- struct{}{}
	}()
	select {
	case <-done:
	case <-ctx2.Done():
		m.logger.Error("cannot close goroutines")
		return ctx2.Err()
	}
	m.cancel = nil
	return nil
}

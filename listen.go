package gokord

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/nyttikord/gokord/logger"
)

// syncListener must not be copied!
type syncListener struct {
	state  atomic.Uint32
	wg     sync.WaitGroup
	logger *slog.Logger
}

func (m *syncListener) Add(fn func(free func())) {
	m.state.Add(1)
	m.wg.Add(1)
	go fn(func() {
		if m.state.Load() == 0 {
			m.logger.WarnContext(logger.NewContext(context.Background(), 1), "goroutine already free from sync")
			return
		}
		m.state.Store(m.state.Load() - 1)
		m.wg.Done()
	})
}

func (m *syncListener) Wait() {
	m.wg.Wait()
}

func (m *syncListener) ForceFree() {
	m.wg.Add(-int(m.state.Load()))
	m.state.Store(0)
}

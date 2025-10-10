package event

import (
	"context"
	"log/slog"
	"sync"

	"github.com/nyttikord/gokord/bot"
)

type Manager struct {
	sync.RWMutex
	Logger func() *slog.Logger

	SyncEvents bool

	handlers     map[string][]*eventHandlerInstance
	onceHandlers map[string][]*eventHandlerInstance

	onInterface func(context.Context, any)
}

func NewManager(s bot.Session, onInterface func(context.Context, any)) *Manager {
	return &Manager{
		RWMutex:      sync.RWMutex{},
		Logger:       func() *slog.Logger { return s.Logger().With("module", "event") },
		SyncEvents:   false,
		handlers:     make(map[string][]*eventHandlerInstance),
		onceHandlers: make(map[string][]*eventHandlerInstance),
		onInterface:  onInterface,
	}
}

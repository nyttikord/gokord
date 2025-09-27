package event

import (
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

	onInterface func(any)
	onReady     func(*Ready)
}

func NewManager(s bot.Session, onInterface func(any), onReady func(*Ready)) *Manager {
	return &Manager{
		RWMutex:      sync.RWMutex{},
		Logger:       s.Logger,
		SyncEvents:   false,
		handlers:     make(map[string][]*eventHandlerInstance),
		onceHandlers: make(map[string][]*eventHandlerInstance),
		onInterface:  onInterface,
		onReady:      onReady,
	}
}

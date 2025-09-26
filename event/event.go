package event

import (
	"sync"

	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/logger"
)

type Manager struct {
	sync.RWMutex
	logger.Logger

	SyncEvents bool

	handlers     map[string][]*eventHandlerInstance
	onceHandlers map[string][]*eventHandlerInstance

	onInterface func(any)
	onReady     func(*Ready)
}

func NewManager(s bot.Session, onInterface func(any), onReady func(*Ready)) *Manager {
	return &Manager{
		RWMutex:      sync.RWMutex{},
		Logger:       s,
		SyncEvents:   false,
		handlers:     make(map[string][]*eventHandlerInstance),
		onceHandlers: make(map[string][]*eventHandlerInstance),
		onInterface:  onInterface,
		onReady:      onReady,
	}
}

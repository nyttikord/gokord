package event

import (
	"sync"

	"github.com/nyttikord/gokord/application/applicationapi"
	"github.com/nyttikord/gokord/channel/channelapi"
	"github.com/nyttikord/gokord/guild/guildapi"
	"github.com/nyttikord/gokord/interaction/interactionapi"
	"github.com/nyttikord/gokord/logger"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user/invite/inviteapi"
	"github.com/nyttikord/gokord/user/userapi"
)

// Session represents a gokord.Session used in events.
type Session interface {
	logger.Logger

	// EventManager returns the Manager used by the Session.
	EventManager() *Manager

	// ChannelAPI returns a channelapi.Requester to interact with the channel package.
	ChannelAPI() *channelapi.Requester
	// UserAPI returns an userapi.Requester to interact with the user package.
	UserAPI() *userapi.Requester
	// GuildAPI returns a guildapi.Requester to interact with the guild package.
	GuildAPI() *guildapi.Requester
	// InviteAPI returns an inviteapi.Requester to interact with the invite package.
	InviteAPI() *inviteapi.Requester
	// InteractionAPI returns an interactionapi.Requester to interact with the interaction package.
	InteractionAPI() *interactionapi.Requester
	// ApplicationAPI returns an applicationapi.Requester to interact with the application package.
	ApplicationAPI() *applicationapi.Requester

	// SessionState returns the state.Bot of the Session.
	SessionState() state.Bot
}

type Manager struct {
	sync.RWMutex
	logger.Logger

	SyncEvents bool

	handlers     map[string][]*eventHandlerInstance
	onceHandlers map[string][]*eventHandlerInstance

	onInterface func(any)
	onReady     func(*Ready)
}

func NewManager(s Session, onInterface func(any), onReady func(*Ready)) *Manager {
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

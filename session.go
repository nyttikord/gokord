package gokord

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/nyttikord/gokord/application/applicationapi"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/bot/botapi"
	"github.com/nyttikord/gokord/channel/channelapi"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/guild/guildapi"
	"github.com/nyttikord/gokord/interaction/interactionapi"
	"github.com/nyttikord/gokord/logger"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user/invite/inviteapi"
	"github.com/nyttikord/gokord/user/status"
	"github.com/nyttikord/gokord/user/userapi"
)

type mutex struct {
	sync.RWMutex
	logger *slog.Logger
}

func (m *mutex) Lock() {
	m.logger.DebugContext(logger.NewContext(context.Background(), 1), "locking")
	m.RWMutex.Lock()
}

func (m *mutex) Unlock() {
	m.logger.DebugContext(logger.NewContext(context.Background(), 1), "unlocking")
	m.RWMutex.Unlock()
}

// Session represents a connection to the Discord API.
type Session struct {
	mu     *mutex
	logger *slog.Logger

	// General configurable settings.

	MFA bool

	// Options of the Session.
	Options bot.Options

	// Identify is sent during initial handshake with the discord gateway.
	// https://discord.com/developers/docs/topics/gateway#identify
	Identify Identify

	// REST contains the Session interacting with the REST API.
	REST *RESTSession

	// Managed state object, updated internally with events when StateEnabled is true.
	sessionState *sessionState
	// Stores when the lastHeartbeatAck was received (UTC).
	lastHeartbeatAck atomic.Int64
	// Stores the lastHeartbeatSent (UTC).
	lastHeartbeatSent atomic.Int64
	// heartbeatInterval is the interval between two heartbeats
	heartbeatInterval time.Duration

	// Event handlers
	eventManager *event.Manager
	// The websocket connection.
	ws *websocket.Conn
	// Wait for listen goroutines to stop
	waitListen *syncListener
	// sequence tracks the current gateway api websocket sequence number.
	sequence *atomic.Int64
	// Stores sessions current Discord Resume Gateway.
	resumeGatewayURL string
	// Stores sessions current Discord Gateway.
	gateway string
	// Stores session ID of current Gateway connection.
	sessionID string
	// Used to make sure gateway websocket writes do not happen concurrently.
	wsMutex sync.Mutex
	// Used to receive result of ws.Read
	wsRead <-chan readResult
	// cancel wsRead
	cancelWSRead func()
	// true if the session is restarting
	restarting atomic.Bool

	// API with state.State

	userAPI    *userapi.Requester
	channelAPI *channelapi.Requester
	guildAPI   *guildapi.Requester

	UserStorage    state.Storage // UserStorage is the state.Storage used for the UserAPI
	ChannelStorage state.Storage // ChannelStorage is the state.Storage used for the ChannelAPI
	GuildStorage   state.Storage // GuildStorage is the state.Storage used for the GuildAPI
}

// GatewayBotResponse stores the data for the gateway/bot response.
type GatewayBotResponse struct {
	URL               string             `json:"url"`
	Shards            int                `json:"shards"`
	SessionStartLimit SessionInformation `json:"session_start_limit"`
}

// SessionInformation provides the information for max concurrency sharding.
type SessionInformation struct {
	Total          int `json:"total,omitempty"`
	Remaining      int `json:"remaining,omitempty"`
	ResetAfter     int `json:"reset_after,omitempty"`
	MaxConcurrency int `json:"max_concurrency,omitempty"`
}

// GatewayStatusUpdate is sent by the client to indicate a presence or status update.
// https://discord.com/developers/docs/topics/gateway#update-status-gateway-status-update-structure
type GatewayStatusUpdate struct {
	Since  int             `json:"since"`
	Game   status.Activity `json:"game"`
	Status string          `json:"status"`
	AFK    bool            `json:"afk"`
}

// OpenAndBlock calls Session.Open and block the program until an OS signal is received.
// It returns an error if Session.Open or Session.Close return an error.
// When this function returns, the session is already disconnected.
func (s *Session) OpenAndBlock(ctx context.Context) error {
	err := s.Open(ctx)
	if err != nil {
		return err
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sc:
		return s.Close(ctx)
	case <-ctx.Done():
		s.logger.Error("waiting to close", "error", ctx.Err())
		return s.Close(context.Background())
	}
}

// UserAPI returns an userapi.Requester to interact with the user package.
func (s *Session) UserAPI() *userapi.Requester {
	if s.userAPI == nil {
		s.logger.Debug("creating new user state")
		s.userAPI = &userapi.Requester{RESTRequester: s.REST, State: userapi.NewState(s.sessionState, s.UserStorage)}
	}
	return s.userAPI
}

// GuildAPI returns a guildapi.Requester to interact with the guild package.
func (s *Session) GuildAPI() *guildapi.Requester {
	if s.guildAPI == nil {
		s.logger.Debug("creating new guild state")
		s.guildAPI = &guildapi.Requester{
			RESTRequester: s.REST,
			WSRequester:   s,
			State:         guildapi.NewState(s.sessionState, s.GuildStorage),
		}
	}
	return s.guildAPI
}

// ChannelAPI returns a channelapi.Requester to interact with the channel package.
func (s *Session) ChannelAPI() *channelapi.Requester {
	if s.channelAPI == nil {
		s.logger.Debug("creating new channel state")
		s.channelAPI = &channelapi.Requester{RESTRequester: s.REST, State: channelapi.NewState(s.sessionState, s.ChannelStorage)}
	}
	return s.channelAPI
}

// InviteAPI returns an inviteapi.Requester to interact with the invite package.
func (s *Session) InviteAPI() *inviteapi.Requester {
	return &inviteapi.Requester{RESTRequester: s.REST}
}

// InteractionAPI returns an interactionapi.Requester to interact with the interaction package.
func (s *Session) InteractionAPI() *interactionapi.Requester {
	return &interactionapi.Requester{RESTRequester: s.REST, ChannelAPI: s.ChannelAPI}
}

// ApplicationAPI returns an applicationapi.Requester to interact with the application package.
func (s *Session) ApplicationAPI() *applicationapi.Requester {
	return &applicationapi.Requester{RESTRequester: s.REST}
}

// BotAPI returns a botapi.Requester to interact with the bot package.
func (s *Session) BotAPI() *botapi.Requester {
	return &botapi.Requester{RESTRequester: s.REST, WSRequester: s}
}

// EventManager returns the event.Manager used by the Session.
func (s *Session) EventManager() bot.EventManager {
	return s.eventManager
}

// SessionState returns the state.Bot of the Session.
func (s *Session) SessionState() state.Bot {
	return s.sessionState
}

// SetStateParams sets the state.Params for the state.State
func (s *Session) SetStateParams(params state.Params) {
	s.sessionState.params = params
}

// LastHeartbeatAck returns the time.Time of the last heartbeat ack received.
func (s *Session) LastHeartbeatAck() time.Time {
	last := s.lastHeartbeatAck.Load()
	return time.Unix(last/1000, (last%1000)*int64(time.Millisecond))
}

// LastHeartbeatAck returns the time.Time of the last heartbeat sent.
func (s *Session) LastHeartbeatSent() time.Time {
	last := s.lastHeartbeatSent.Load()
	return time.Unix(last/1000, (last%1000)*int64(time.Millisecond))
}

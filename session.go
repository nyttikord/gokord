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
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/interaction/interactionhandler"
	"github.com/nyttikord/gokord/logger"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user/status"
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

	// rest contains the Session interacting with the rest API.
	rest *RESTSession

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
	// Interaction handlers
	interactionManager *interactionhandler.Manager
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

	// Cached things

	memberState  *state.Member
	channelState *state.Channel
	guildState   *state.Guild

	MemberStorage  state.MemberStorage  // MemberStorage is the [state.Storage] used for [state.Member].
	ChannelStorage state.ChannelStorage // ChannelStorage is the [state.Storage] used for [state.Channel].
	GuildStorage   state.GuildStorage   // GuildStorage is the [state.Storage] used for [state.Guild].
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

// EventManager returns the event.Manager used by the Session.
func (s *Session) EventManager() bot.EventManager {
	return s.eventManager
}

// InteractionManager returns the *interactionhandler.Manager used by the Session.
func (s *Session) InteractionManager() *interactionhandler.Manager {
	return s.interactionManager
}

// GatewayAPI returns the API used to interact with the gateway.
func (s *Session) GatewayAPI() bot.GatewayAPI {
	return &wsAPI{
		logger:  s.logger.With("module", "gateway"),
		Session: s,
	}
}

func (s *Session) MemberState() *state.Member {
	if s.memberState == nil {
		s.memberState = state.NewMember(s, s.MemberStorage, &s.sessionState.params)
	}
	return s.memberState
}

func (s *Session) ChannelState() *state.Channel {
	if s.memberState == nil {
		s.channelState = state.NewChannel(s, s.ChannelStorage, &s.sessionState.params)
	}
	return s.channelState
}

func (s *Session) GuildState() *state.Guild {
	if s.memberState == nil {
		s.guildState = state.NewGuild(s, s.GuildStorage, &s.sessionState.params)
	}
	return s.guildState
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

// Logger returns the logger used by the Session.
func (s *Session) Logger() *slog.Logger {
	return s.logger
}

// NewContext returns a new context usable everywhere.
func (s *Session) NewContext(ctx context.Context) context.Context {
	return bot.NewContext(ctx, s.logger, s, s.rest)
}

package gokord

import (
	"context"
	"log/slog"
	"net/http"
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
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/guild/guildapi"
	"github.com/nyttikord/gokord/interaction/interactionapi"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user/invite/inviteapi"
	"github.com/nyttikord/gokord/user/status"
	"github.com/nyttikord/gokord/user/userapi"
	"github.com/nyttikord/gokord/voice"
)

// Session represents a connection to the Discord API.
type Session struct {
	*sync.RWMutex
	logger *slog.Logger

	// General configurable settings.

	MFA bool
	// Should the session reconnect the websocket on errors.
	ShouldReconnectOnError bool
	// Should voice connections reconnect on a session reconnect.
	ShouldReconnectVoiceOnSessionError bool
	// Should the session retry requests when rate limited.
	ShouldRetryOnRateLimit bool
	// Identify is sent during initial handshake with the discord gateway.
	// https://discord.com/developers/docs/topics/gateway#identify
	Identify Identify
	// Should state tracking be enabled.
	// State tracking is the best way for getting the users active guilds and the members of the guilds.
	StateEnabled bool
	// Whether to call event handlers synchronously.
	// e.g. false = launch event handlers in their own goroutines.
	SyncEvents bool

	// Exposed but should not be modified by Application.

	// Max number of REST API retries.
	MaxRestRetries int
	// Status stores the current status of the websocket connection this is being tested, may stay, may go away.
	status int32
	// Managed state object, updated internally with events when StateEnabled is true.
	sessionState *sessionState
	// The http.Client used for REST requests.
	Client *http.Client
	// The websocket.Dialer used for WebSocket connection.
	//Dialer *websocket.
	// The UserAgent used for REST APIs.
	UserAgent string
	// Stores the LastHeartbeatAck that was received (in UTC).
	LastHeartbeatAck time.Time
	// Stores the LastHeartbeatSent (in UTC).
	LastHeartbeatSent time.Time
	// Used to deal with rate limits.
	RateLimiter *discord.RateLimiter
	// heartbeatInterval is the interval between two heartbeats
	heartbeatInterval time.Duration

	// Event handlers
	eventManager *event.Manager
	// The websocket connection.
	ws *websocket.Conn
	// Cancel listen goroutines
	cancelListen func()
	// Wait for listen goroutines to stop
	waitListen sync.WaitGroup
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

	// API with state.State

	userAPI    *userapi.Requester
	channelAPI *channelapi.Requester
	guildAPI   *guildapi.Requester
	voiceAPI   *voice.Requester
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

func (s *Session) Logger() *slog.Logger {
	return s.logger
}

// UserAPI returns an userapi.Requester to interact with the user package.
func (s *Session) UserAPI() *userapi.Requester {
	if s.userAPI == nil {
		s.userAPI = &userapi.Requester{Requester: s, State: userapi.NewState(s.sessionState)}
	}
	return s.userAPI
}

// GuildAPI returns a guildapi.Requester to interact with the guild package.
func (s *Session) GuildAPI() *guildapi.Requester {
	if s.guildAPI == nil {
		s.guildAPI = &guildapi.Requester{API: s, State: guildapi.NewState(s.sessionState)}
	}
	return s.guildAPI
}

// ChannelAPI returns a channelapi.Requester to interact with the channel package.
func (s *Session) ChannelAPI() *channelapi.Requester {
	if s.channelAPI == nil {
		s.channelAPI = &channelapi.Requester{Requester: s, State: channelapi.NewState(s.sessionState)}
	}
	return s.channelAPI
}

// InviteAPI returns an inviteapi.Requester to interact with the invite package.
func (s *Session) InviteAPI() *inviteapi.Requester {
	return &inviteapi.Requester{Requester: s}
}

// InteractionAPI returns an interactionapi.Requester to interact with the interaction package.
func (s *Session) InteractionAPI() *interactionapi.Requester {
	return &interactionapi.Requester{API: s}
}

// ApplicationAPI returns an applicationapi.Requester to interact with the application package.
func (s *Session) ApplicationAPI() *applicationapi.Requester {
	return &applicationapi.Requester{Requester: s}
}

// BotAPI returns a botapi.Requester to interact with the bot package.
func (s *Session) BotAPI() *botapi.Requester {
	return &botapi.Requester{Requester: s}
}

// VoiceAPI returns a voice.Requester to interact with the voice package.
func (s *Session) VoiceAPI() *voice.Requester {
	return s.voiceAPI
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

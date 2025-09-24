package gokord

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nyttikord/gokord/application/applicationapi"
	"github.com/nyttikord/gokord/channel/channelapi"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/guild/guildapi"
	"github.com/nyttikord/gokord/interaction/interactionapi"
	"github.com/nyttikord/gokord/user/invite/inviteapi"
	"github.com/nyttikord/gokord/user/status"
	"github.com/nyttikord/gokord/user/userapi"
)

// Session represents a connection to the Discord API.
type Session struct {
	sync.RWMutex
	stdLogger

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

	// Sharding
	ShardID    int
	ShardCount int

	// Should state tracking be enabled.
	// State tracking is the best way for getting the users active guilds and the members of the guilds.
	StateEnabled bool

	// Whether to call event handlers synchronously.
	// e.g. false = launch event handlers in their own goroutines.
	SyncEvents bool

	// Exposed but should not be modified by Application.

	// Whether the Data Websocket is ready.
	//
	// Note: May be deprecated soon.
	DataReady bool

	// Max number of REST API retries.
	MaxRestRetries int

	// Status stores the current status of the websocket connection this is being tested, may stay, may go away.
	status int32

	// Stores a mapping of guild id's to VoiceConnection.
	VoiceConnections map[string]*VoiceConnection

	// Managed state object, updated internally with events when StateEnabled is true.
	State *State

	// The http.Client used for REST requests.
	Client *http.Client

	// The websocket.Dialer used for WebSocket connection.
	Dialer *websocket.Dialer

	// The UserAgent used for REST APIs.
	UserAgent string

	// Stores the LastHeartbeatAck that was received (in UTC).
	LastHeartbeatAck time.Time

	// Stores the LastHeartbeatSent (in UTC).
	LastHeartbeatSent time.Time

	// Used to deal with rate limits.
	RateLimiter *discord.RateLimiter

	// Event handlers
	eventManager *event.Manager

	// The websocket connection.
	wsConn *websocket.Conn

	// When nil, the session is not listening.
	listening chan any

	// sequence tracks the current gateway api websocket sequence number.
	sequence *int64

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

// Identify is sent during initial handshake with the discord gateway.
// https://discord.com/developers/docs/topics/gateway#identify
type Identify struct {
	Token          string              `json:"token"`
	Properties     IdentifyProperties  `json:"properties"`
	Compress       bool                `json:"compress"`
	LargeThreshold int                 `json:"large_threshold"`
	Shard          *[2]int             `json:"shard,omitempty"`
	Presence       GatewayStatusUpdate `json:"presence,omitempty"`
	Intents        discord.Intent      `json:"intents"`
}

// IdentifyProperties contains the "properties" portion of an Identify packet.
// https://discord.com/developers/docs/topics/gateway#identify-identify-connection-properties
type IdentifyProperties struct {
	OS              string `json:"$os"`
	Browser         string `json:"$browser"`
	Device          string `json:"$device"`
	Referer         string `json:"$referer"`
	ReferringDomain string `json:"$referring_domain"`
}

// UserAPI returns an userapi.Requester to interact with the user package.
func (s *Session) UserAPI() *userapi.Requester {
	if s.userAPI == nil {
		s.userAPI = &userapi.Requester{Requester: s, State: userapi.NewState(s.State)}
	}
	return s.userAPI
}

// GuildAPI returns a guildapi.Requester to interact with the guild package.
func (s *Session) GuildAPI() *guildapi.Requester {
	if s.guildAPI == nil {
		s.guildAPI = &guildapi.Requester{Requester: s, State: guildapi.NewState(s.State)}
	}
	return s.guildAPI
}

// ChannelAPI returns a channelapi.Requester to interact with the channel package.
func (s *Session) ChannelAPI() *channelapi.Requester {
	if s.channelAPI == nil {
		s.channelAPI = &channelapi.Requester{Requester: s, State: channelapi.NewState(s.State)}
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

func (s *Session) EventManager() *event.Manager {
	return s.eventManager
}

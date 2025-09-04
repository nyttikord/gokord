package gokord

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

// VERSION of DiscordGo, follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.29.0"

// New creates a new Discord session with provided token.
// If the token is for a bot, it must be prefixed with "Bot "
//
//	e.g. "Bot ..."
//
// Or if it is an OAuth2 token, it must be prefixed with "Bearer "
//
//	e.g. "Bearer ..."
func New(token string) (s *Session, err error) {

	// Create an empty Session interface.
	s = &Session{
		State:                              NewState(),
		Ratelimiter:                        NewRatelimiter(),
		StateEnabled:                       true,
		Compress:                           true,
		ShouldReconnectOnError:             true,
		ShouldReconnectVoiceOnSessionError: true,
		ShouldRetryOnRateLimit:             true,
		ShardID:                            0,
		ShardCount:                         1,
		MaxRestRetries:                     3,
		Client:                             &http.Client{Timeout: (20 * time.Second)},
		Dialer:                             websocket.DefaultDialer,
		UserAgent:                          "DiscordBot (https://github.com/bwmarrin/discordgo, v" + VERSION + ")",
		sequence:                           new(int64),
		LastHeartbeatAck:                   time.Now().UTC(),
	}

	// Initialize the Identify Package with defaults
	// These can be modified prior to calling Open()
	s.Identify.Compress = true
	s.Identify.LargeThreshold = 250
	s.Identify.Properties.OS = runtime.GOOS
	s.Identify.Properties.Browser = "DiscordGo v" + VERSION
	s.Identify.Intents = IntentsAllWithoutPrivileged
	s.Identify.Token = token
	s.Token = token

	return
}

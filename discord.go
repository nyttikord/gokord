package gokord

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/logger"
)

// VERSION of Gokord, follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.31.0"

// New creates a new Discord session with provided token.
// If the token is for a bot, it must be prefixed with "Bot "
//
//	e.g. "Bot ..."
//
// Or if it is an OAuth2 token, it must be prefixed with "Bearer "
//
//	e.g. "Bearer ..."
func New(token string) *Session {
	s := &Session{
		RateLimiter:                        discord.NewRateLimiter(),
		StateEnabled:                       true,
		ShouldReconnectOnError:             true,
		ShouldReconnectVoiceOnSessionError: true,
		ShouldRetryOnRateLimit:             true,
		ShardID:                            0,
		ShardCount:                         1,
		MaxRestRetries:                     3,
		Client:                             &http.Client{Timeout: 20 * time.Second},
		Dialer:                             websocket.DefaultDialer,
		UserAgent:                          "DiscordBot (https://github.com/nyttikord/gokord, v" + VERSION + ")",
		sequence:                           new(int64),
		LastHeartbeatAck:                   time.Now().UTC(),
		stdLogger:                          stdLogger{Level: logger.LevelInfo},
	}
	s.State = NewState(s)

	// Initialize the Identify Package with defaults
	// These can be modified prior to calling Open()
	s.Identify.Compress = true
	s.Identify.LargeThreshold = 250
	s.Identify.Properties.OS = runtime.GOOS
	s.Identify.Properties.Browser = "DiscordGo v" + VERSION
	s.Identify.Intents = discord.IntentsAllWithoutPrivileged
	s.Identify.Token = token

	return s
}

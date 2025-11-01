package gokord

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/logger"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/voice"
)

// VERSION of Gokord, follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.33.0"

// New creates a new Discord session with provided token.
// If the token is for a bot, it must be prefixed with "Bot "
//
//	e.g. "Bot ..."
//
// Or if it is an OAuth2 token, it must be prefixed with "Bearer "
//
//	e.g. "Bearer ..."
//
// See NewWithLogLevel to modify the default slog.Level.
// See NewWithLogger to set the default slog.Logger.
func New(token string) *Session {
	return NewWithLogLevel(token, slog.LevelInfo)
}

// NewWithLogLevel creates a new Discord session with provided token and set the slog.Level of the logger.
//
// See New for the full documentation.
// See NewWithLogger to set the default slog.Logger.
func NewWithLogLevel(token string, logLevel slog.Level) *Session {
	return NewWithLogger(token, slog.New(logger.New(os.Stdout, &logger.Options{Level: logLevel})))
}

// NewWithLogger creates a new Discord session with provided token and set the logger.
//
// See New for the full documentation.
// See NewWithLogLevel to modify the default slog.Level.
func NewWithLogger(token string, logger *slog.Logger) *Session {
	s := &Session{
		StateEnabled:                       true,
		ShouldReconnectOnError:             true,
		ShouldReconnectVoiceOnSessionError: true,
		LastHeartbeatAck:                   time.Now().UTC(),
		logger:                             logger,
		RWMutex:                            &sync.RWMutex{},
		waitListen:                         &syncListener{logger: logger.With("module", "ws")},
		UserStorage:                        &state.MapStorage[user.Member]{},
		ChannelStorage:                     &state.MapStorage[channel.Channel]{},
		GuildStorage:                       &state.MapStorage[guild.Guild]{},
	}
	s.sessionState = NewState(s).(*sessionState)
	s.eventManager = event.NewManager(s, s.onInterface)

	s.REST = &RESTSession{
		identify:               &s.Identify,
		logger:                 logger.With("module", "rest"),
		ShouldRetryOnRateLimit: true,
		eventManager:           s.eventManager,
		MaxRestRetries:         3,
		Client:                 &http.Client{Timeout: 20 * time.Second},
		UserAgent:              "DiscordBot (https://github.com/nyttikord/gokord, v" + VERSION + ")",
		RateLimiter:            discord.NewRateLimiter(),
		emitRateLimitEvent: func(ctx context.Context, rl *event.RateLimit) {
			s.eventManager.EmitEvent(ctx, s, event.RateLimitType, rl)
		},
	}

	s.voiceAPI = &voice.Requester{
		RESTRequester: s.REST,
		WSRequester:   s,
		Connections:   make(map[string]*voice.Connection),
	}

	// Initialize Identify with defaults values.
	// These can be modified prior to calling Open().
	s.Identify.Compress = true
	s.Identify.LargeThreshold = 250
	s.Identify.Properties.OS = runtime.GOOS
	s.Identify.Properties.Browser = "DiscordGo v" + VERSION
	s.Identify.Intents = discord.IntentsAllWithoutPrivileged
	s.Identify.Token = token
	s.Identify.Shard = &[2]int{0, 1}

	return s
}

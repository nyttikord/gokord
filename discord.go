package gokord

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/logger"
)

// VERSION of Gokord, follows Semantic Versioning. (http://semver.org/)
//
// APIVersion is appended after the version.
const VERSION = "0.31.0+v" + discord.APIVersion

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
		State:                              NewState(),
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

// MemberPermissions calculates the permissions for a user.Member.
// https://support.discord.com/hc/en-us/articles/206141927-How-is-the-permission-hierarchy-structured-
func MemberPermissions(guild *guild.Guild, channel *channel.Channel, userID string, roles []string) int64 {
	if userID == guild.OwnerID {
		return discord.PermissionAll
	}

	var perms int64
	for _, role := range guild.Roles {
		if role.ID == guild.ID {
			perms |= role.Permissions
			break
		}
	}

	for _, role := range guild.Roles {
		for _, roleID := range roles {
			if role.ID == roleID {
				perms |= role.Permissions
				break
			}
		}
	}

	if perms&discord.PermissionAdministrator == discord.PermissionAdministrator {
		perms |= discord.PermissionAll
	}

	// Apply @everyone overrides from the channel.
	for _, overwrite := range channel.PermissionOverwrites {
		if guild.ID == overwrite.ID {
			perms &= ^overwrite.Deny
			perms |= overwrite.Allow
			break
		}
	}

	var denies, allows int64
	// Member overwrites can override role overrides, so do two passes
	for _, overwrite := range channel.PermissionOverwrites {
		for _, roleID := range roles {
			if overwrite.Type == types.PermissionOverwriteRole && roleID == overwrite.ID {
				denies |= overwrite.Deny
				allows |= overwrite.Allow
				break
			}
		}
	}

	perms &= ^denies
	perms |= allows

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.Type == types.PermissionOverwriteMember && overwrite.ID == userID {
			perms &= ^overwrite.Deny
			perms |= overwrite.Allow
			break
		}
	}

	if perms&discord.PermissionAdministrator == discord.PermissionAdministrator {
		perms |= discord.PermissionAllChannel
	}

	return perms
}

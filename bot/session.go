// Package bot contains everything related to the bot.
//
// You can use [gokord.Session] which is the default implementation of [Session] working with the Discord Gateway.
package bot

import (
	"context"

	"github.com/nyttikord/gokord/channel/channelapi"
	"github.com/nyttikord/gokord/guild/guildapi"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user/status"
	"github.com/nyttikord/gokord/user/userapi"
)

// Session represents a bot session.
// Default implementation is gokord.Session which is using the gateway.
// You can create your own implementation to use webhooks.
type Session interface {
	// EventManager returns the EventManager used by the Session.
	EventManager() EventManager

	// ChannelAPI returns a channelapi.Requester to interact with the channel package.
	ChannelAPI() *channelapi.Requester
	// UserAPI returns an userapi.Requester to interact with the user package.
	UserAPI() *userapi.Requester
	// GuildAPI returns a guildapi.Requester to interact with the guild package.
	GuildAPI() *guildapi.Requester
	// GatewayAPI returns the [GatewayAPI] used by the [Session].
	// It can be nil if the [Session] does not support the gateway.
	// [gokord.Session] supports this API.
	GatewayAPI() GatewayAPI

	// SessionState returns the state.Bot of the Session.
	SessionState() state.Bot
}

// Options for the Session.
type Options struct {
	// Should the session reconnect the websocket on errors.
	ShouldReconnectOnError bool
	// Should the session retry requests when rate limited.
	ShouldRetryOnRateLimit bool
	// Max number of REST API retries.
	MaxRestRetries uint
	// Should state tracking be enabled.
	// State tracking is the best way for getting the users active guilds and the members of the guilds.
	StateEnabled bool
	// Whether to call event handlers synchronously.
	// e.g. false = launch event handlers in their own goroutines.
	SyncEvents bool
}

// EventManager handles events for the Session.
type EventManager interface {
	// AddHandler allows you to add an event handler that will be fired anytime the Discord WSAPI event that matches the
	// function fires.
	// The first parameter is a Session, and the second parameter is a pointer to a struct corresponding to the event for
	// which you want to listen.
	//
	// eg:
	//
	//	Session.AddHandler(func(s event.Session, m *discordgo.MessageCreate) {
	//	})
	//
	// or:
	//
	//	Session.AddHandler(func(s event.Session, m *discordgo.PresenceUpdate) {
	//	})
	//
	// List of events can be found at this page, with corresponding names in the library for each event:
	// https://discord.com/developers/docs/topics/gateway#event-names
	// There are also synthetic events fired by the library internally which are available for handling, like Connect,
	// Disconnect, and RateLimit.
	// events.go contains all the Discord WSAPI and synthetic events that can be handled.
	//
	// The return value of this method is a function, that when called will remove the event handler.
	AddHandler(any) func()
	// AddHandlerOnce allows you to add an event handler that will be fired the next time the Discord WSAPI event that
	// matches the function fires.
	//
	// See AddHandler for more details.
	AddHandlerOnce(any) func()
}

// UpdateStatusData is used in [GatewayAPI] to update bot status via the gateway.
type UpdateStatusData struct {
	IdleSince  *int               `json:"since"`
	Activities []*status.Activity `json:"activities"`
	AFK        bool               `json:"afk"`
	Status     string             `json:"status"`
}

// GatewayAPI is the interface used to communicate with the Gateway API.
type GatewayAPI interface {
	// UpdateGameStatus is used to update the [user.User]'s status.
	//
	// If idle > 0 then set status to idle.
	// If name != "" then set game.
	// if otherwise, set status to active, and no activity.
	UpdateGameStatus(ctx context.Context, idle bool, name string) error
	// UpdateWatchStatus is used to update the [user.User]'s watch status.
	//
	// If idle > 0 then set status to idle.
	// If name != "" then set movie/stream.
	// if otherwise, set status to active, and no activity.
	UpdateWatchStatus(ctx context.Context, idle bool, name string) error
	// UpdateStreamingStatus is used to update the [user.User]'s streaming status.
	//
	// If idle > 0 then set status to idle.
	// If name != "" then set game.
	// If name != "" and url != "" then set the status type to streaming with the URL set.
	// if otherwise, set status to active, and no game.
	UpdateStreamingStatus(ctx context.Context, idle bool, name string, url string) error
	// UpdateListeningStatus is used to set the [user.User] to "Listening to..."
	//
	// If name != "" then set to what user is listening to.
	// Else, set user to active and no activity.
	UpdateListeningStatus(ctx context.Context, name string) error
	// UpdateCustomStatus is used to update the [user.User]'s custom status.
	//
	// If state != "" then set the custom status.
	// Else, set user to active and remove the custom status.
	UpdateCustomStatus(ctx context.Context, state string) error
	// UpdateStatusComplex allows for sending the raw status update data.
	UpdateStatusComplex(ctx context.Context, usd UpdateStatusData) error

	// GatewayMembers requests [user.Member] from the gateway.
	// It responds with [event.GuildMembersChunk].
	//
	// query is a string that username starts with, leave empty to return every [user.Member].
	// limit is the maximum number of items to return, or 0 to request every [user.Member] matched.
	// nonce to identify the [event.GuildMembersChunk] response.
	// presences indicates whether to request presences of [user.Member].
	GatewayMembers(ctx context.Context, guildID, query string, limit int, nonce string, presences bool) error
	// GatewayMembersList requests [user.Member] from the gateway.
	// It responds with [event.GuildMembersChunk].
	//
	// userIDs are the [user.Member]'s IDs to fetch.
	// limit is the maximum number of items to return, or 0 to request every [user.Member] matched.
	// nonce to identify the [event.GuildMembersChunk] response.
	// presences indicates whether to request presences of [user.Member].
	GatewayMembersList(ctx context.Context, guildID string, userIDs []string, limit int, nonce string, presences bool) error
}

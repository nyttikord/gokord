// Package bot contains everything related to the bot.
//
// You can use gokord.Session which is the default implementation of Session working with the Discord Gateway.
package bot

import (
	"log/slog"

	"github.com/nyttikord/gokord/application/applicationapi"
	"github.com/nyttikord/gokord/bot/botapi"
	"github.com/nyttikord/gokord/channel/channelapi"
	"github.com/nyttikord/gokord/guild/guildapi"
	"github.com/nyttikord/gokord/interaction/interactionapi"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user/invite/inviteapi"
	"github.com/nyttikord/gokord/user/userapi"
	"github.com/nyttikord/gokord/voice"
)

// Session represents a bot session.
// Default implementation is gokord.Session which is using the gateway.
// You can create your own implementation to use webhooks.
type Session interface {
	// Logger returns the slog.Logger used by the Session.
	Logger() *slog.Logger

	// EventManager returns the EventManager used by the Session.
	EventManager() EventManager

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
	// BotAPI returns a botapi.Requester to interact with the bot package.
	BotAPI() *botapi.Requester
	// VoiceAPI returns a voice.Requester to interact with the voice package.
	VoiceAPI() *voice.Requester

	// SessionState returns the state.Bot of the Session.
	SessionState() state.Bot
}

// Options for the Session.
type Options struct {
	// Should the session reconnect the websocket on errors.
	ShouldReconnectOnError bool
	// Should voice connections reconnect on a session reconnect.
	ShouldReconnectVoiceOnSessionError bool
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

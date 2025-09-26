package bot

import (
	"github.com/nyttikord/gokord/application/applicationapi"
	"github.com/nyttikord/gokord/channel/channelapi"
	"github.com/nyttikord/gokord/guild/guildapi"
	"github.com/nyttikord/gokord/interaction/interactionapi"
	"github.com/nyttikord/gokord/logger"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user/invite/inviteapi"
	"github.com/nyttikord/gokord/user/userapi"
)

// Session represents a bot session.
// Default implementation is gokord.Session which is using the gateway.
// You can create your own implementation to use webhooks.
type Session interface {
	logger.Logger

	// EventManager returns the Manager used by the Session.
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
	// BotAPI returns a bot.Requester to interact with the bot package.
	BotAPI() *Requester

	// SessionState returns the state.Bot of the Session.
	SessionState() state.Bot
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
	// EmitEvent calls internal methods, fires handlers and fires the "any" event.
	//
	// NOTE: I don't know if this should be private, or not
	//EmitEvent(Session, string, any)
}

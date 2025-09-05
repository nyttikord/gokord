package discord

// Intent is the type of a Gateway Intent
// https://discord.com/developers/docs/topics/gateway#gateway-intents
type Intent int

// Constants for the different bit offsets of intents
const (
	IntentGuilds                      Intent = 1 << 0
	IntentGuildMembers                Intent = 1 << 1
	IntentGuildModeration             Intent = 1 << 2
	IntentGuildEmojis                 Intent = 1 << 3
	IntentGuildIntegrations           Intent = 1 << 4
	IntentGuildWebhooks               Intent = 1 << 5
	IntentGuildInvites                Intent = 1 << 6
	IntentGuildVoiceStates            Intent = 1 << 7
	IntentGuildPresences              Intent = 1 << 8
	IntentGuildMessages               Intent = 1 << 9
	IntentGuildMessageReactions       Intent = 1 << 10
	IntentGuildMessageTyping          Intent = 1 << 11
	IntentDirectMessages              Intent = 1 << 12
	IntentDirectMessageReactions      Intent = 1 << 13
	IntentDirectMessageTyping         Intent = 1 << 14
	IntentMessageContent              Intent = 1 << 15
	IntentGuildScheduledEvents        Intent = 1 << 16
	IntentAutoModerationConfiguration Intent = 1 << 20
	IntentAutoModerationExecution     Intent = 1 << 21
	IntentGuildMessagePolls           Intent = 1 << 24
	IntentDirectMessagePolls          Intent = 1 << 25

	// TODO: remove when compatibility is not needed

	IntentGuildBans Intent = IntentGuildModeration

	IntentsGuilds                 Intent = 1 << 0
	IntentsGuildMembers           Intent = 1 << 1
	IntentsGuildBans              Intent = 1 << 2
	IntentsGuildEmojis            Intent = 1 << 3
	IntentsGuildIntegrations      Intent = 1 << 4
	IntentsGuildWebhooks          Intent = 1 << 5
	IntentsGuildInvites           Intent = 1 << 6
	IntentsGuildVoiceStates       Intent = 1 << 7
	IntentsGuildPresences         Intent = 1 << 8
	IntentsGuildMessages          Intent = 1 << 9
	IntentsGuildMessageReactions  Intent = 1 << 10
	IntentsGuildMessageTyping     Intent = 1 << 11
	IntentsDirectMessages         Intent = 1 << 12
	IntentsDirectMessageReactions Intent = 1 << 13
	IntentsDirectMessageTyping    Intent = 1 << 14
	IntentsMessageContent         Intent = 1 << 15
	IntentsGuildScheduledEvents   Intent = 1 << 16

	IntentsAllWithoutPrivileged = IntentGuilds |
		IntentGuildBans |
		IntentGuildEmojis |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping |
		IntentGuildScheduledEvents |
		IntentAutoModerationConfiguration |
		IntentAutoModerationExecution

	IntentsAll = IntentsAllWithoutPrivileged |
		IntentGuildMembers |
		IntentGuildPresences |
		IntentMessageContent

	IntentsNone Intent = 0
)

// MakeIntent used to help convert a gateway intent value for use in the Identify structure;
// this was useful to help support the use of a pointer type when intents were optional.
// This is now a no-op, and is not necessary to use.
func MakeIntent(intents Intent) Intent {
	return intents
}

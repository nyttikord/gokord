package types

import "fmt"

// ApplicationCommand represents the type of gokord.ApplicationCommand.
type ApplicationCommand uint8

const (
	// ApplicationCommandChat is default command type. They are slash commands (i.e. called directly from the chat).
	ApplicationCommandChat ApplicationCommand = 1
	// ApplicationCommandUser adds command to user context menu.
	ApplicationCommandUser ApplicationCommand = 2
	// ApplicationCommandMessage adds command to message context menu.
	ApplicationCommandMessage ApplicationCommand = 3
)

// ApplicationCommandOption indicates the type of gokord.ApplicationCommandOption.
type ApplicationCommandOption uint8

const (
	ApplicationCommandOptionSubCommand      ApplicationCommandOption = 1
	ApplicationCommandOptionSubCommandGroup ApplicationCommandOption = 2
	ApplicationCommandOptionString          ApplicationCommandOption = 3
	ApplicationCommandOptionInteger         ApplicationCommandOption = 4
	ApplicationCommandOptionBoolean         ApplicationCommandOption = 5
	ApplicationCommandOptionUser            ApplicationCommandOption = 6
	ApplicationCommandOptionChannel         ApplicationCommandOption = 7
	ApplicationCommandOptionRole            ApplicationCommandOption = 8
	ApplicationCommandOptionMentionable     ApplicationCommandOption = 9
	ApplicationCommandOptionNumber          ApplicationCommandOption = 10
	ApplicationCommandOptionAttachment      ApplicationCommandOption = 11
)

func (t ApplicationCommandOption) String() string {
	switch t {
	case ApplicationCommandOptionSubCommand:
		return "SubCommand"
	case ApplicationCommandOptionSubCommandGroup:
		return "SubCommandGroup"
	case ApplicationCommandOptionString:
		return "String"
	case ApplicationCommandOptionInteger:
		return "Integer"
	case ApplicationCommandOptionBoolean:
		return "Boolean"
	case ApplicationCommandOptionUser:
		return "User"
	case ApplicationCommandOptionChannel:
		return "Channel"
	case ApplicationCommandOptionRole:
		return "Role"
	case ApplicationCommandOptionMentionable:
		return "Mentionable"
	case ApplicationCommandOptionNumber:
		return "Number"
	case ApplicationCommandOptionAttachment:
		return "Attachment"
	}
	return fmt.Sprintf("ApplicationCommandOption(%d)", t)
}

// ApplicationCommandPermission indicates whether a gokord.ApplicationCommandPermissions permission is user or role based.
type ApplicationCommandPermission uint8

// Application command permission types.
const (
	ApplicationCommandPermissionRole    ApplicationCommandPermission = 1
	ApplicationCommandPermissionUser    ApplicationCommandPermission = 2
	ApplicationCommandPermissionChannel ApplicationCommandPermission = 3
)

// Interaction indicates the type of gokord.Interaction event.
type Interaction uint8

// Interaction types
const (
	InteractionPing                           Interaction = 1
	InteractionApplicationCommand             Interaction = 2
	InteractionMessageComponent               Interaction = 3
	InteractionApplicationCommandAutocomplete Interaction = 4
	InteractionModalSubmit                    Interaction = 5
)

func (t Interaction) String() string {
	switch t {
	case InteractionPing:
		return "Ping"
	case InteractionApplicationCommand:
		return "ApplicationCommand"
	case InteractionMessageComponent:
		return "Message"
	case InteractionModalSubmit:
		return "ModalSubmit"
	}
	return fmt.Sprintf("Interaction(%d)", t)
}

// InteractionContext represents the context in which gokord.Interaction can be used or was triggered from.
type InteractionContext uint

const (
	// InteractionContextGuild indicates that interaction can be used within guilds.
	InteractionContextGuild InteractionContext = 0
	// InteractionContextBotDM indicates that interaction can be used within DMs with the bot.
	InteractionContextBotDM InteractionContext = 1
	// InteractionContextPrivateChannel indicates that interaction can be used within group DMs and DMs with other users.
	InteractionContextPrivateChannel InteractionContext = 2
)

// InteractionResponse is type of gokord.InteractionResponse.
type InteractionResponse uint8

// Interaction response types.
const (
	// InteractionResponsePong is for ACK ping event.
	InteractionResponsePong InteractionResponse = 1
	// InteractionResponseChannelMessageWithSource is for responding with a message, showing the user's input.
	InteractionResponseChannelMessageWithSource InteractionResponse = 4
	// InteractionResponseDeferredChannelMessageWithSource acknowledges that the event was received, and that a follow-up will come later.
	InteractionResponseDeferredChannelMessageWithSource InteractionResponse = 5
	// InteractionResponseDeferredMessageUpdate acknowledges that the message component interaction event was received, and message will be updated later.
	InteractionResponseDeferredMessageUpdate InteractionResponse = 6
	// InteractionResponseUpdateMessage is for updating the message to which message component was attached.
	InteractionResponseUpdateMessage InteractionResponse = 7
	// InteractionApplicationCommandAutocompleteResult shows autocompletion results. Autocomplete interaction only.
	InteractionApplicationCommandAutocompleteResult InteractionResponse = 8
	// InteractionResponseModal is for responding to an interaction with a modal window.
	InteractionResponseModal InteractionResponse = 9
)

package types

import "fmt"

// Command represents the type of gokord.Command.
type Command uint8

const (
	// CommandChat is default command type. They are slash commands (i.e. called directly from the chat).
	CommandChat Command = 1
	// CommandUser adds command to user context menu.
	CommandUser Command = 2
	// CommandMessage adds command to message context menu.
	CommandMessage Command = 3
)

// CommandOption indicates the type of gokord.CommandOption.
type CommandOption uint8

const (
	CommandOptionSubCommand      CommandOption = 1
	CommandOptionSubCommandGroup CommandOption = 2
	CommandOptionString          CommandOption = 3
	CommandOptionInteger         CommandOption = 4
	CommandOptionBoolean         CommandOption = 5
	CommandOptionUser            CommandOption = 6
	CommandOptionChannel         CommandOption = 7
	CommandOptionRole            CommandOption = 8
	CommandOptionMentionable     CommandOption = 9
	CommandOptionNumber          CommandOption = 10
	CommandOptionAttachment      CommandOption = 11
)

func (t CommandOption) String() string {
	switch t {
	case CommandOptionSubCommand:
		return "SubCommand"
	case CommandOptionSubCommandGroup:
		return "SubCommandGroup"
	case CommandOptionString:
		return "String"
	case CommandOptionInteger:
		return "Integer"
	case CommandOptionBoolean:
		return "Boolean"
	case CommandOptionUser:
		return "User"
	case CommandOptionChannel:
		return "Channel"
	case CommandOptionRole:
		return "Role"
	case CommandOptionMentionable:
		return "Mentionable"
	case CommandOptionNumber:
		return "Number"
	case CommandOptionAttachment:
		return "Attachment"
	}
	return fmt.Sprintf("CommandOption(%d)", t)
}

// CommandPermission indicates whether a gokord.ApplicationCommandPermissions permission is user or role based.
type CommandPermission uint8

// User command permission types.
const (
	CommandPermissionRole    CommandPermission = 1
	CommandPermissionUser    CommandPermission = 2
	CommandPermissionChannel CommandPermission = 3
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
		return "Command"
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

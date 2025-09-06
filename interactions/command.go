package interactions

import (
	"strconv"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
)

// Command represents an application's slash command.
type Command struct {
	ID                string                     `json:"id,omitempty"`
	ApplicationID     string                     `json:"application_id,omitempty"`
	GuildID           string                     `json:"guild_id,omitempty"`
	Version           string                     `json:"version,omitempty"`
	Type              types.ApplicationCommand   `json:"type,omitempty"`
	Name              string                     `json:"name"`
	NameLocalizations *map[discord.Locale]string `json:"name_localizations,omitempty"`

	// Note: DefaultPermission will be soon deprecated. Use DefaultMemberPermissions and Contexts instead.
	DefaultPermission        *bool  `json:"default_permission,omitempty"`
	DefaultMemberPermissions *int64 `json:"default_member_permissions,string,omitempty"`
	NSFW                     *bool  `json:"nsfw,omitempty"`

	// Deprecated: use Contexts instead.
	DMPermission     *bool                       `json:"dm_permission,omitempty"`
	Contexts         *[]types.InteractionContext `json:"contexts,omitempty"`
	IntegrationTypes *[]types.Integration        `json:"integration_types,omitempty"`

	// Note: Chat commands only.
	// Otherwise, it mustn't be set.
	Description string `json:"description,omitempty"`
	// Note: Chat commands only.
	// Otherwise, it mustn't be set.
	DescriptionLocalizations *map[discord.Locale]string `json:"description_localizations,omitempty"`
	// Note: Chat commands only.
	// Otherwise, it mustn't be set.
	Options []*CommandOption `json:"options"`
}

// CommandOption represents an option/subcommand/subcommands group.
type CommandOption struct {
	Type                     types.ApplicationCommandOption `json:"type"`
	Name                     string                         `json:"name"`
	NameLocalizations        map[discord.Locale]string      `json:"name_localizations,omitempty"`
	Description              string                         `json:"description,omitempty"`
	DescriptionLocalizations map[discord.Locale]string      `json:"description_localizations,omitempty"`
	// Note: This feature was on the API, but at some point developers decided to remove it.
	// So I commented it, until it will be officially on the docs.
	// Default     bool                              `json:"default"`

	ChannelTypes []types.Channel  `json:"channel_types"`
	Required     bool             `json:"required"`
	Options      []*CommandOption `json:"options"`

	// Note: mutually exclusive with Choices.
	Autocomplete bool `json:"autocomplete"`
	// Note: mutually exclusive with Autocomplete.
	Choices []*CommandOptionChoice `json:"choices"`
	// Minimal value of types.ApplicationCommandOptionInteger/types.ApplicationCommandOptionNumber.
	MinValue *float64 `json:"min_value,omitempty"`
	// Maximum value of types.ApplicationCommandOptionInteger/types.ApplicationCommandOptionNumber.
	MaxValue float64 `json:"max_value,omitempty"`
	// Minimum length of types.ApplicationCommandOptionString.
	MinLength *int `json:"min_length,omitempty"`
	// Maximum length of types.ApplicationCommandOptionString.
	MaxLength int `json:"max_length,omitempty"`
}

// CommandOptionChoice represents a slash CommandOption choice.
type CommandOptionChoice struct {
	Name              string                    `json:"name"`
	NameLocalizations map[discord.Locale]string `json:"name_localizations,omitempty"`
	Value             interface{}               `json:"value"`
}

// CommandPermissions represents a single user.User or guild.Role permission for a Command.
type CommandPermissions struct {
	ID         string                             `json:"id"`
	Type       types.ApplicationCommandPermission `json:"type"`
	Permission bool                               `json:"permission"`
}

// GuildAllChannelsID is a helper function which returns guild_id-1.
// It is used in CommandPermissions to target all the channels within a guild.Guild.
func GuildAllChannelsID(guild string) (id string, err error) {
	var v uint64
	v, err = strconv.ParseUint(guild, 10, 64)
	if err != nil {
		return
	}

	return strconv.FormatUint(v-1, 10), nil
}

// CommandPermissionsList represents a list of CommandPermissions, needed for serializing to JSON.
type CommandPermissionsList struct {
	Permissions []*CommandPermissions `json:"permissions"`
}

// GuildCommandPermissions represents all permissions for a single guild.Guild Command.
type GuildCommandPermissions struct {
	ID            string                `json:"id"`
	ApplicationID string                `json:"application_id"`
	GuildID       string                `json:"guild_id"`
	Permissions   []*CommandPermissions `json:"permissions"`
}

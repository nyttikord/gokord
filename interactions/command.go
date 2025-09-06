package interactions

import (
	"strconv"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
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

// CommandInteractionData contains the data of Command Interaction.
type CommandInteractionData struct {
	ID          string                          `json:"id"`
	Name        string                          `json:"name"`
	CommandType types.ApplicationCommand        `json:"type"`
	Resolved    *CommandInteractionDataResolved `json:"resolved"`

	// Slash command options
	Options []*CommandInteractionDataOption `json:"options"`
	// InviteTarget (user/message) id on which context menu command was called.
	// The details are stored in Resolved according to command type.
	TargetID string `json:"target_id"`
}

// GetOption finds and returns an CommandOption by its name.
func (d CommandInteractionData) GetOption(name string) (option *CommandInteractionDataOption) {
	for _, opt := range d.Options {
		if opt.Name == name {
			option = opt
			break
		}
	}

	return
}

// CommandInteractionDataResolved contains resolved data of Command execution.
type CommandInteractionDataResolved struct {
	Users map[string]*user.User `json:"users"`
	// Partial user.Member are missing user.User, Deaf and Mute fields.
	Members map[string]*user.Member `json:"members"`
	Roles   map[string]*guild.Role  `json:"roles"`
	// Partial channel.Channel only have ID, Name, Type and Permissions fields.
	Channels    map[string]*channel.Channel           `json:"channels"`
	Messages    map[string]*channel.Message           `json:"messages"`
	Attachments map[string]*channel.MessageAttachment `json:"attachments"`
}

// Type returns the type of InteractionData.
func (CommandInteractionData) Type() types.Interaction {
	return types.InteractionApplicationCommand
}

// CommandInteractionDataOption represents an option of a slash Command.
type CommandInteractionDataOption struct {
	Name string                         `json:"name"`
	Type types.ApplicationCommandOption `json:"type"`
	// Note: Contains the value specified by Type.
	Value   interface{}                     `json:"value,omitempty"`
	Options []*CommandInteractionDataOption `json:"options,omitempty"`

	// Note: autocomplete Interaction only.
	Focused bool `json:"focused,omitempty"`
}

// GetOption finds and returns an CommandOption by its name.
func (o CommandInteractionDataOption) GetOption(name string) (option *CommandInteractionDataOption) {
	for _, opt := range o.Options {
		if opt.Name == name {
			option = opt
			break
		}
	}

	return
}

// IntValue is a utility function for casting CommandOption value to integer
func (o CommandInteractionDataOption) IntValue() int64 {
	if o.Type != types.ApplicationCommandOptionInteger {
		panic("IntValue called on data option of type " + o.Type.String())
	}
	return int64(o.Value.(float64))
}

// UintValue is a utility function for casting CommandOption value to unsigned integer
func (o CommandInteractionDataOption) UintValue() uint64 {
	if o.Type != types.ApplicationCommandOptionInteger {
		panic("UintValue called on data option of type " + o.Type.String())
	}
	return uint64(o.Value.(float64))
}

// FloatValue is a utility function for casting CommandOption value to float
func (o CommandInteractionDataOption) FloatValue() float64 {
	if o.Type != types.ApplicationCommandOptionNumber {
		panic("FloatValue called on data option of type " + o.Type.String())
	}
	return o.Value.(float64)
}

// StringValue is a utility function for casting CommandOption value to string
func (o CommandInteractionDataOption) StringValue() string {
	if o.Type != types.ApplicationCommandOptionString {
		panic("StringValue called on data option of type " + o.Type.String())
	}
	return o.Value.(string)
}

// BoolValue is a utility function for casting CommandOption value to bool
func (o CommandInteractionDataOption) BoolValue() bool {
	if o.Type != types.ApplicationCommandOptionBoolean {
		panic("BoolValue called on data option of type " + o.Type.String())
	}
	return o.Value.(bool)
}

// ChannelValue is a utility function for casting CommandOption value to channel.Channel.
// s is the gokord.Session, if not nil, function additionally fetches all channel.Channel's data
func (o CommandInteractionDataOption) ChannelValue(s *gokord.Session) *channel.Channel {
	if o.Type != types.ApplicationCommandOptionChannel {
		panic("ChannelValue called on data option of type " + o.Type.String())
	}
	chanID := o.Value.(string)

	if s == nil {
		return &channel.Channel{ID: chanID}
	}

	ch, err := s.State.Channel(chanID)
	if err != nil {
		ch, err = s.Channel(chanID)
		if err != nil {
			return &channel.Channel{ID: chanID}
		}
	}

	return ch
}

// RoleValue is a utility function for casting CommandOption value to guild.Role.
// s is the gokord.Session, if not nil, function additionally fetches all role's data
func (o CommandInteractionDataOption) RoleValue(s *gokord.Session, gID string) *guild.Role {
	if o.Type != types.ApplicationCommandOptionRole && o.Type != types.ApplicationCommandOptionMentionable {
		panic("RoleValue called on data option of type " + o.Type.String())
	}
	roleID := o.Value.(string)

	if s == nil || gID == "" {
		return &guild.Role{ID: roleID}
	}

	r, err := s.State.Role(gID, roleID)
	if err != nil {
		roles, err := s.GuildRoles(gID)
		if err == nil {
			for _, r = range roles {
				if r.ID == roleID {
					return r
				}
			}
		}
		return &guild.Role{ID: roleID}
	}

	return r
}

// UserValue is a utility function for casting CommandOption value to user.User.
// s is the gokord.Session, if not nil, function additionally fetches all user.User's data
func (o CommandInteractionDataOption) UserValue(s *gokord.Session) *user.User {
	if o.Type != types.ApplicationCommandOptionUser && o.Type != types.ApplicationCommandOptionMentionable {
		panic("UserValue called on data option of type " + o.Type.String())
	}
	userID := o.Value.(string)

	if s == nil {
		return &user.User{ID: userID}
	}

	u, err := s.User(userID)
	if err != nil {
		return &user.User{ID: userID}
	}

	return u
}

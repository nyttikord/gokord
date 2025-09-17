package interaction

import (
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
)

// ChannelGetter represents a type fetching channel.Channel.
type ChannelGetter interface {
	// Channel returns the channel.Channel with the given ID.
	Channel(string, ...discord.RequestOption) (*channel.Channel, error)
}

type StateChannelGetter interface {
	Channel(string) (*channel.Channel, error)
}

// RoleGetter represents a type fetching guild.Role.
type RoleGetter interface {
	// Role returns the guild.Role in the given guild.Guild gID with the given rID.
	Role(gID string, rID string) (*guild.Role, error)
}

// RolesGetter represents a type fetching []*guild.Role.
type RolesGetter interface {
	// Roles returns the []*guild.Role in the given guild.Guild.
	Roles(string, ...discord.RequestOption) ([]*guild.Role, error)
}

// UserGetter represents a type fetching user.User.
type UserGetter interface {
	// User returns the user.User with the given ID.
	User(string, ...discord.RequestOption) (*user.User, error)
}

// CommandInteractionData contains the data of Command Interaction.
type CommandInteractionData struct {
	ID          string                          `json:"id"`
	Name        string                          `json:"name"`
	CommandType types.Command                   `json:"type"`
	Resolved    *CommandInteractionDataResolved `json:"resolved"`

	// Slash command options
	Options []*CommandInteractionDataOption `json:"options"`
	// InviteTarget (user/message) id on which context menu command was called.
	// The details are stored in Resolved according to command type.
	TargetID string `json:"target_id"`
}

// GetOption finds and returns an CommandOption by its name.
func (d CommandInteractionData) GetOption(name string) *CommandInteractionDataOption {
	for _, opt := range d.Options {
		if opt.Name == name {
			return opt
		}
	}
	return nil
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

// Type returns the type of Data.
func (CommandInteractionData) Type() types.Interaction {
	return types.InteractionApplicationCommand
}

// CommandInteractionDataOption represents an option of a slash Command.
type CommandInteractionDataOption struct {
	Name string              `json:"name"`
	Type types.CommandOption `json:"type"`
	// NOTE: Contains the value specified by Type.
	Value   interface{}                     `json:"value,omitempty"`
	Options []*CommandInteractionDataOption `json:"options,omitempty"`

	// NOTE: autocomplete Interaction only.
	Focused bool `json:"focused,omitempty"`
}

// GetOption finds and returns an CommandOption by its name.
func (o CommandInteractionDataOption) GetOption(name string) *CommandInteractionDataOption {
	for _, opt := range o.Options {
		if opt.Name == name {
			return opt
		}
	}
	return nil
}

// IntValue is a utility function for casting CommandOption value to integer.
func (o CommandInteractionDataOption) IntValue() int64 {
	if o.Type != types.CommandOptionInteger {
		panic("IntValue called on data option of type " + o.Type.String())
	}
	return int64(o.Value.(float64))
}

// UintValue is a utility function for casting CommandOption value to unsigned integer.
func (o CommandInteractionDataOption) UintValue() uint64 {
	if o.Type != types.CommandOptionInteger {
		panic("UintValue called on data option of type " + o.Type.String())
	}
	return uint64(o.Value.(float64))
}

// FloatValue is a utility function for casting CommandOption value to float.
func (o CommandInteractionDataOption) FloatValue() float64 {
	if o.Type != types.CommandOptionNumber {
		panic("FloatValue called on data option of type " + o.Type.String())
	}
	return o.Value.(float64)
}

// StringValue is a utility function for casting CommandOption value to string.
func (o CommandInteractionDataOption) StringValue() string {
	if o.Type != types.CommandOptionString {
		panic("StringValue called on data option of type " + o.Type.String())
	}
	return o.Value.(string)
}

// BoolValue is a utility function for casting CommandOption value to bool.
func (o CommandInteractionDataOption) BoolValue() bool {
	if o.Type != types.CommandOptionBoolean {
		panic("BoolValue called on data option of type " + o.Type.String())
	}
	return o.Value.(bool)
}

// ChannelValue is a utility function for casting CommandOption value to channel.Channel.
//
// s is a ChannelGetter (implemented by gokord.Session), if not nil, function additionally fetches all
// channel.Channel's data.
// state is another ChannelGetter representing the internal state of the application (implemented by gokord.State), if
// not nil, it is called before s.
func (o CommandInteractionDataOption) ChannelValue(s ChannelGetter, state StateChannelGetter) *channel.Channel {
	if o.Type != types.CommandOptionChannel {
		panic("ChannelValue called on data option of type " + o.Type.String())
	}
	chanID := o.Value.(string)

	if s == nil {
		return &channel.Channel{ID: chanID}
	}

	if state != nil {
		ch, err := state.Channel(chanID)
		if err == nil {
			return ch
		}
	}
	ch, err := s.Channel(chanID)
	if err != nil {
		return &channel.Channel{ID: chanID}
	}
	return ch
}

// RoleValue is a utility function for casting CommandOption value to guild.Role.
//
// gID is the guild.Guild ID containing the role.
// s is a RolesGetter (implemented by gokord.Session), if not nil, function additionally fetches all
// guild.Role's data.
// state is a RoleGetter representing the internal state of the application (implemented by gokord.State), if
// not nil, it is called before s.
func (o CommandInteractionDataOption) RoleValue(gID string, s RolesGetter, state RoleGetter) *guild.Role {
	if o.Type != types.CommandOptionRole && o.Type != types.CommandOptionMentionable {
		panic("RoleValue called on data option of type " + o.Type.String())
	}
	roleID := o.Value.(string)

	if s == nil || gID == "" {
		return &guild.Role{ID: roleID}
	}

	r, err := state.Role(gID, roleID)
	if err == nil {
		return r
	}
	roles, err := s.Roles(gID)
	if err == nil {
		for _, r = range roles {
			if r.ID == roleID {
				return r
			}
		}
	}
	return &guild.Role{ID: roleID}
}

// UserValue is a utility function for casting CommandOption value to user.User.
//
// s is a UserGetter (implemented by gokord.Session), if not nil, function additionally fetches all user.User's data.
func (o CommandInteractionDataOption) UserValue(s UserGetter) *user.User {
	if o.Type != types.CommandOptionUser && o.Type != types.CommandOptionMentionable {
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

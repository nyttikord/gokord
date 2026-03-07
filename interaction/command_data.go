package interaction

import (
	"context"
	"strconv"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user"
)

// CommandInteractionData contains the data of [ApplicationCommand] [Interaction].
type CommandInteractionData struct {
	ID          uint64                          `json:"id,string"`
	Name        string                          `json:"name"`
	CommandType types.Command                   `json:"type"`
	Resolved    *CommandInteractionDataResolved `json:"resolved"`

	// Slash [Command] options
	Options []*CommandInteractionDataOption `json:"options"`
	// InviteTarget (user/message) id on which context menu command was called.
	// The details are stored in Resolved according to command type.
	TargetID uint64 `json:"target_id,string"`
}

// GetOption finds and returns an [CommandOption] by its name.
//
// Deprecated: use [ApplicationCommand.OptionMap] to get the [CommandOptionMap].
func (d *CommandInteractionData) GetOption(name string) *CommandInteractionDataOption {
	for _, opt := range d.Options {
		if opt.Name == name {
			return opt
		}
	}
	return nil
}

// CommandInteractionDataResolved contains resolved data of [ApplicationCommand] execution.
type CommandInteractionDataResolved struct {
	Users map[uint64]*user.User `json:"users"`
	// Partial [user.Member] are missing [user.User], Deaf and Mute fields.
	Members map[uint64]*user.Member `json:"members"`
	Roles   map[uint64]*guild.Role  `json:"roles"`
	// Partial [channel.Channel] only have ID, Name, Type and Permissions fields.
	Channels    map[uint64]*channel.Channel           `json:"channels"`
	Messages    map[uint64]*channel.Message           `json:"messages"`
	Attachments map[uint64]*channel.MessageAttachment `json:"attachments"`
}

func (*CommandInteractionData) Type() types.Interaction {
	return types.InteractionApplicationCommand
}

// CommandInteractionDataOption represents an option of an [ApplicationCommand].
type CommandInteractionDataOption struct {
	Name string              `json:"name"`
	Type types.CommandOption `json:"type"`
	// NOTE: Contains the value specified by [Type].
	Value   any                             `json:"value,omitempty"`
	Options []*CommandInteractionDataOption `json:"options,omitempty"`

	// NOTE: autocomplete Interaction only.
	Focused bool `json:"focused,omitempty"`
}

// GetOption finds and returns an option by its name.
func (o CommandInteractionDataOption) GetOption(name string) *CommandInteractionDataOption {
	for _, opt := range o.Options {
		if opt.Name == name {
			return opt
		}
	}
	return nil
}

// IntValue is a utility function for casting option value to integer.
func (o CommandInteractionDataOption) IntValue() int64 {
	if o.Type != types.CommandOptionInteger {
		panic("IntValue called on data option of type " + o.Type.String())
	}
	return int64(o.Value.(float64))
}

// UintValue is a utility function for casting option value to unsigned integer.
func (o CommandInteractionDataOption) UintValue() uint64 {
	if o.Type != types.CommandOptionInteger {
		panic("UintValue called on data option of type " + o.Type.String())
	}
	return uint64(o.Value.(float64))
}

// FloatValue is a utility function for casting option value to float.
func (o CommandInteractionDataOption) FloatValue() float64 {
	if o.Type != types.CommandOptionNumber {
		panic("FloatValue called on data option of type " + o.Type.String())
	}
	return o.Value.(float64)
}

// StringValue is a utility function for casting option value to string.
func (o CommandInteractionDataOption) StringValue() string {
	if o.Type != types.CommandOptionString {
		panic("StringValue called on data option of type " + o.Type.String())
	}
	return o.Value.(string)
}

// BoolValue is a utility function for casting option value to bool.
func (o CommandInteractionDataOption) BoolValue() bool {
	if o.Type != types.CommandOptionBoolean {
		panic("BoolValue called on data option of type " + o.Type.String())
	}
	return o.Value.(bool)
}

// ChannelValue is a utility function for casting option value to [channel.Channel].
func (o CommandInteractionDataOption) ChannelValue(ctx context.Context, state *state.Channel) *channel.Channel {
	if o.Type != types.CommandOptionChannel {
		panic("ChannelValue called on data option of type " + o.Type.String())
	}
	chanID, err := strconv.ParseUint(o.Value.(string), 10, 64)
	if err != nil {
		panic(err)
	}

	if state != nil {
		ch, err := state.GetChannel(chanID)
		if err == nil {
			return ch
		}
	}
	ch, err := channel.Get(chanID).Do(ctx)
	if err != nil {
		return &channel.Channel{ID: chanID}
	}
	return ch
}

// RoleValue is a utility function for casting option value to [guild.Role].
func (o CommandInteractionDataOption) RoleValue(ctx context.Context, gID uint64, state *state.Guild) *guild.Role {
	if o.Type != types.CommandOptionRole && o.Type != types.CommandOptionMentionable {
		panic("RoleValue called on data option of type " + o.Type.String())
	}
	roleID, err := strconv.ParseUint(o.Value.(string), 10, 64)
	if err != nil {
		panic(err)
	}

	if gID == 0 {
		return &guild.Role{ID: roleID}
	}

	r, err := state.GetRole(gID, roleID)
	if err == nil {
		return r
	}
	roles, err := guild.ListRoles(gID).Do(ctx)
	if err == nil {
		for _, r = range roles {
			if r.ID == roleID {
				return r
			}
		}
	}
	return &guild.Role{ID: roleID}
}

// UserValue is a utility function for casting option value to [user.User].
func (o CommandInteractionDataOption) UserValue(ctx context.Context) *user.User {
	if o.Type != types.CommandOptionUser && o.Type != types.CommandOptionMentionable {
		panic("UserValue called on data option of type " + o.Type.String())
	}
	userID, err := strconv.ParseUint(o.Value.(string), 10, 64)
	if err != nil {
		panic(err)
	}

	u, err := user.Get(userID).Do(ctx)
	if err != nil {
		return &user.User{ID: userID}
	}
	return u
}

type CommandOptionMap map[string]*CommandInteractionDataOption

func (d *ApplicationCommand) OptionMap() CommandOptionMap {
	options := d.Data.Options
	optionMap := make(CommandOptionMap, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

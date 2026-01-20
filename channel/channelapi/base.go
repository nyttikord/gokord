// Package channelapi contains everything to interact with everything located in the channel package.
package channelapi

import (
	"context"
	"errors"
	"net/http"

	. "github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/logger"
	"github.com/nyttikord/gokord/user/invite"
)

var (
	ErrTooMuchStickers         = errors.New("too much stickers: cannot send more than 3 stickers")
	ErrTooMuchMessagesToDelete = errors.New("too much messages to delete: cannot delete more than 100 messages")
	ErrReplyNilMessageRef      = errors.New("reply attempted with nil message reference")
)

// Requester handles everything inside the channel package.
type Requester struct {
	REST
	State *State
}

// Channel returns the channel.Channel with the given ID.
func (r Requester) Channel(channelID string) Request[*Channel] {
	return NewData[*Channel](r, http.MethodGet, discord.EndpointChannel(channelID))
}

// ChannelEdit edits the given channel.Channel and returns the updated channel.Channel data.
func (r Requester) ChannelEdit(channelID string, data *Edit) Request[*Channel] {
	return NewData[*Channel](
		r, http.MethodPatch, discord.EndpointChannel(channelID),
	).WithData(data)
}

// ChannelDelete deletes the given channel.Channel.
func (r Requester) ChannelDelete(channelID string) Request[*Channel] {
	return NewData[*Channel](r, http.MethodDelete, discord.EndpointChannel(channelID))
}

// Typing broadcasts to all members that authenticated user.User is typing in the given channel.Channel.
func (r Requester) Typing(channelID string) Empty {
	req := NewSimple(r, http.MethodPost, discord.EndpointChannelTyping(channelID))
	return WrapAsEmpty(req)
}

// Invites returns all invite.Invite for the given channel.Channel.
func (r Requester) Invites(channelID string) Request[[]*invite.Invite] {
	return NewData[[]*invite.Invite](r, http.MethodDelete, discord.EndpointChannelInvites(channelID))
}

// InviteCreate creates a new invite.Invite for the given channel.Channel.
//
// NOTE: invite.Invite must have MaxAge, MaxUses and Temporary.
func (r Requester) InviteCreate(channelID string, i invite.Invite) Request[*invite.Invite] {
	uID := ""
	if i.TargetUser != nil {
		uID = i.TargetUser.ID
	}
	appID := ""
	if i.TargetApplication != nil {
		appID = i.TargetApplication.ID
	}
	data := struct {
		MaxAge            int                `json:"max_age"`
		MaxUses           int                `json:"max_uses"`
		Temporary         bool               `json:"temporary"`
		Unique            bool               `json:"unique"`
		TargetType        types.InviteTarget `json:"target_type"`
		TargetUser        string             `json:"target_user_id"`
		TargetApplication string             `json:"target_application_id"`
		//TargerUsers       []byte             `json:"target_users_file,omitempty"`
		Roles []string `json:"role_ids,omitempty"`
	}{i.MaxAge, i.MaxUses, i.Temporary, i.Unique, i.TargetType, uID, appID, i.Roles}

	if len(i.TargetUsersFile) > 0 {
		r.Logger().WarnContext(logger.NewContext(context.Background(), 1), "InviteCreate does not support yet TargetUsersFile")
	}

	return NewData[*invite.Invite](
		r, http.MethodPost, discord.EndpointChannelInvites(channelID),
	).WithData(data)
}

// PermissionSet creates a channel.PermissionOverwrite for the given channel.Channel.
//
// NOTE: This func name may be changed.
// Using Set instead of Create because you can both create a new override or update an override with this function.
func (r Requester) PermissionSet(channelID, targetID string, targetType types.PermissionOverwrite, allow, deny int64) Empty {
	data := struct {
		ID    string                    `json:"id"`
		Type  types.PermissionOverwrite `json:"type"`
		Allow int64                     `json:"allow,string"`
		Deny  int64                     `json:"deny,string"`
	}{targetID, targetType, allow, deny}

	req := NewSimple(r, http.MethodPut, discord.EndpointChannelPermission(channelID, targetID)).
		WithData(data).
		WithBucketID(discord.EndpointChannelPermission(channelID, ""))
	return WrapAsEmpty(req)
}

// PermissionDelete deletes a specific channel.PermissionOverwrite for the given channel.Channel.
//
// NOTE: Name of this func may change.
func (r Requester) PermissionDelete(channelID, targetID string) Empty {
	req := NewSimple(r, http.MethodPut, discord.EndpointChannelPermission(channelID, targetID)).
		WithBucketID(discord.EndpointChannelPermission(channelID, ""))
	return WrapAsEmpty(req)
}

// NewsFollow follows a news channel.Channel in the given channel.Channel.
//
// channelID is the Channel to follow.
// targetID is where the news Channel should post to.
func (r Requester) NewsFollow(channelID, targetID string) Request[*Follow] {
	data := struct {
		WebhookChannelID string `json:"webhook_channel_id"`
	}{targetID}

	return NewData[*Follow](
		r, http.MethodPost, discord.EndpointChannelFollow(channelID),
	).WithData(data)
}

// StageInstanceCreate creates and returns a new channel.Stage instance associated to a types.ChannelGuildStageVoice.
func (r Requester) StageInstanceCreate(data *StageInstanceParams) Request[*StageInstance] {
	return NewData[*StageInstance](
		r, http.MethodPost, discord.EndpointStageInstances,
	).WithData(data)
}

// StageInstance will retrieve a channel.Stage instance by the ID of the types.ChannelGuildStageVoice.
func (r Requester) StageInstance(channelID string) Request[*StageInstance] {
	return NewData[*StageInstance](
		r, http.MethodGet, discord.EndpointStageInstance(channelID),
	)
}

// StageInstanceEdit edits a channel.Stage instance by ID the types.ChannelGuildStageVoice.
func (r Requester) StageInstanceEdit(channelID string, data *StageInstanceParams) Request[*StageInstance] {
	return NewData[*StageInstance](
		r, http.MethodPatch, discord.EndpointStageInstance(channelID),
	).WithData(data)
}

// StageInstanceDelete deletes a channel.Stage instance by ID of the types.ChannelGuildStageVoice.
func (r Requester) StageInstanceDelete(ctx context.Context, channelID string) Empty {
	req := NewSimple(r, http.MethodGet, discord.EndpointStageInstance(channelID))
	return WrapAsEmpty(req)
}

// Package channelapi contains everything to interact with everything located in the channel package.
package channelapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/request"
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
	request.RESTRequester
	State *State
}

// Channel returns the channel.Channel with the given ID.
func (s Requester) Channel(channelID string) request.Request[*channel.Channel] {
	return request.NewSimpleData[*channel.Channel](s, http.MethodGet, discord.EndpointChannel(channelID))
}

// ChannelEdit edits the given channel.Channel and returns the updated channel.Channel data.
func (s Requester) ChannelEdit(channelID string, data *channel.Edit) request.Request[*channel.Channel] {
	return request.NewSimpleData[*channel.Channel](
		s, http.MethodPatch, discord.EndpointChannel(channelID),
	).WithData(data)
}

// ChannelDelete deletes the given channel.Channel.
func (s Requester) ChannelDelete(channelID string) request.Request[*channel.Channel] {
	return request.NewSimpleData[*channel.Channel](s, http.MethodDelete, discord.EndpointChannel(channelID))
}

// Typing broadcasts to all members that authenticated user.User is typing in the given channel.Channel.
func (s Requester) Typing(channelID string) request.EmptyRequest {
	req := request.NewSimple(s, http.MethodPost, discord.EndpointChannelTyping(channelID))
	return request.WrapAsEmpty(req)
}

// Invites returns all invite.Invite for the given channel.Channel.
func (s Requester) Invites(channelID string) request.Request[[]*invite.Invite] {
	return request.NewSimpleData[[]*invite.Invite](s, http.MethodDelete, discord.EndpointChannelInvites(channelID))
}

// InviteCreate creates a new invite.Invite for the given channel.Channel.
//
// NOTE: invite.Invite must have MaxAge, MaxUses and Temporary.
func (s Requester) InviteCreate(channelID string, i invite.Invite) request.Request[*invite.Invite] {
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
		s.Logger().WarnContext(logger.NewContext(context.Background(), 1), "InviteCreate does not support yet TargetUsersFile")
	}

	return request.NewSimpleData[*invite.Invite](
		s, http.MethodPost, discord.EndpointChannelInvites(channelID),
	).WithData(data)
}

// PermissionSet creates a channel.PermissionOverwrite for the given channel.Channel.
//
// NOTE: This func name may be changed.
// Using Set instead of Create because you can both create a new override or update an override with this function.
func (s Requester) PermissionSet(channelID, targetID string, targetType types.PermissionOverwrite, allow, deny int64) request.EmptyRequest {
	data := struct {
		ID    string                    `json:"id"`
		Type  types.PermissionOverwrite `json:"type"`
		Allow int64                     `json:"allow,string"`
		Deny  int64                     `json:"deny,string"`
	}{targetID, targetType, allow, deny}

	req := request.NewSimple(s, http.MethodPut, discord.EndpointChannelPermission(channelID, targetID)).
		WithData(data).
		WithBucketID(discord.EndpointChannelPermission(channelID, ""))
	return request.WrapAsEmpty(req)
}

// PermissionDelete deletes a specific channel.PermissionOverwrite for the given channel.Channel.
//
// NOTE: Name of this func may change.
func (s Requester) PermissionDelete(channelID, targetID string) request.EmptyRequest {
	req := request.NewSimple(s, http.MethodPut, discord.EndpointChannelPermission(channelID, targetID)).
		WithBucketID(discord.EndpointChannelPermission(channelID, ""))
	return request.WrapAsEmpty(req)
}

// NewsFollow follows a news channel.Channel in the given channel.Channel.
//
// channelID is the channel.Channel to follow.
// targetID is where the news channel.Channel should post to.
func (s Requester) NewsFollow(channelID, targetID string) request.Request[*channel.Follow] {
	data := struct {
		WebhookChannelID string `json:"webhook_channel_id"`
	}{targetID}

	return request.NewSimpleData[*channel.Follow](
		s, http.MethodPost, discord.EndpointChannelFollow(channelID),
	).WithData(data)
}

// StageInstanceCreate creates and returns a new Stage instance associated to a types.ChannelGuildStageVoice.
func (s Requester) StageInstanceCreate(data *channel.StageInstanceParams) request.Request[*channel.StageInstance] {
	return request.NewSimpleData[*channel.StageInstance](
		s, http.MethodPost, discord.EndpointStageInstances,
	).WithData(data)
}

// StageInstance will retrieve a Stage instance by the ID of the types.ChannelGuildStageVoice.
func (s Requester) StageInstance(channelID string) request.Request[*channel.StageInstance] {
	return request.NewSimpleData[*channel.StageInstance](
		s, http.MethodGet, discord.EndpointStageInstance(channelID),
	)
}

// StageInstanceEdit edits a Stage instance by ID the types.ChannelGuildStageVoice.
func (s Requester) StageInstanceEdit(channelID string, data *channel.StageInstanceParams) request.Request[*channel.StageInstance] {
	return request.NewSimpleData[*channel.StageInstance](
		s, http.MethodPatch, discord.EndpointStageInstance(channelID),
	).WithData(data)
}

// StageInstanceDelete deletes a Stage instance by ID of the types.ChannelGuildStageVoice.
func (s Requester) StageInstanceDelete(ctx context.Context, channelID string) request.EmptyRequest {
	req := request.NewSimple(s, http.MethodGet, discord.EndpointStageInstance(channelID))
	return request.WrapAsEmpty(req)
}

// Package channelapi contains everything to interact with everything located in the channel package.
package channelapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
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
	discord.RESTRequester
	State *State
}

// Channel returns the channel.Channel with the given ID.
func (s Requester) Channel(ctx context.Context, channelID string, options ...discord.RequestOption) (*channel.Channel, error) {
	body, err := s.Request(ctx, http.MethodGet, discord.EndpointChannel(channelID), nil, options...)
	if err != nil {
		return nil, err
	}

	var c channel.Channel
	return &c, s.Unmarshal(body, &c)
}

// ChannelEdit edits the given channel.Channel and returns the updated channel.Channel data.
func (s Requester) ChannelEdit(ctx context.Context, channelID string, data *channel.Edit, options ...discord.RequestOption) (*channel.Channel, error) {
	body, err := s.Request(ctx, http.MethodPatch, discord.EndpointChannel(channelID), data, options...)
	if err != nil {
		return nil, err
	}

	var c channel.Channel
	return &c, s.Unmarshal(body, &c)

}

// ChannelDelete deletes the given channel.Channel.
func (s Requester) ChannelDelete(ctx context.Context, channelID string, options ...discord.RequestOption) (*channel.Channel, error) {
	body, err := s.Request(ctx, http.MethodDelete, discord.EndpointChannel(channelID), nil, options...)
	if err != nil {
		return nil, err
	}

	var c channel.Channel
	return &c, s.Unmarshal(body, &c)
}

// Typing broadcasts to all members that authenticated user.User is typing in the given channel.Channel.
func (s Requester) Typing(ctx context.Context, channelID string, options ...discord.RequestOption) error {
	_, err := s.Request(ctx, http.MethodPost, discord.EndpointChannelTyping(channelID), nil, options...)
	return err
}

// Invites returns all invite.Invite for the given channel.Channel.
func (s Requester) Invites(ctx context.Context, channelID string, options ...discord.RequestOption) ([]*invite.Invite, error) {
	body, err := s.Request(ctx, http.MethodGet, discord.EndpointChannelInvites(channelID), nil, options...)
	if err != nil {
		return nil, err
	}

	var i []*invite.Invite
	return i, json.Unmarshal(body, &i)
}

// InviteCreate creates a new invite.Invite for the given channel.Channel.
//
// NOTE: invite.Invite must have MaxAge, MaxUses and Temporary.
func (s Requester) InviteCreate(ctx context.Context, channelID string, i invite.Invite, options ...discord.RequestOption) (*invite.Invite, error) {
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

	body, err := s.Request(ctx, http.MethodPost, discord.EndpointChannelInvites(channelID), data, options...)
	if err != nil {
		return nil, err
	}

	var m invite.Invite
	return &m, json.Unmarshal(body, &m)
}

// PermissionSet creates a channel.PermissionOverwrite for the given channel.Channel.
//
// NOTE: This func name may be changed.
// Using Set instead of Create because you can both create a new override or update an override with this function.
func (s Requester) PermissionSet(ctx context.Context, channelID, targetID string, targetType types.PermissionOverwrite, allow, deny int64, options ...discord.RequestOption) error {
	data := struct {
		ID    string                    `json:"id"`
		Type  types.PermissionOverwrite `json:"type"`
		Allow int64                     `json:"allow,string"`
		Deny  int64                     `json:"deny,string"`
	}{targetID, targetType, allow, deny}

	_, err := s.RequestWithBucketID(
		ctx,
		http.MethodPut,
		discord.EndpointChannelPermission(channelID, targetID),
		data,
		discord.EndpointChannelPermission(channelID, ""),
		options...,
	)
	return err
}

// PermissionDelete deletes a specific channel.PermissionOverwrite for the given channel.Channel.
//
// NOTE: Name of this func may change.
func (s Requester) PermissionDelete(ctx context.Context, channelID, targetID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		ctx,
		http.MethodDelete,
		discord.EndpointChannelPermission(channelID, targetID),
		nil,
		discord.EndpointChannelPermission(channelID, ""),
		options...,
	)
	return err
}

// NewsFollow follows a news channel.Channel in the given channel.Channel.
//
// channelID is the channel.Channel to follow.
// targetID is where the news channel.Channel should post to.
func (s Requester) NewsFollow(ctx context.Context, channelID, targetID string, options ...discord.RequestOption) (*channel.Follow, error) {
	endpoint := discord.EndpointChannelFollow(channelID)

	data := struct {
		WebhookChannelID string `json:"webhook_channel_id"`
	}{targetID}

	body, err := s.Request(ctx, http.MethodPost, endpoint, data, options...)
	if err != nil {
		return nil, err
	}

	var f channel.Follow
	return &f, json.Unmarshal(body, &f)
}

// StageInstanceCreate creates and returns a new Stage instance associated to a types.ChannelGuildStageVoice.
func (s Requester) StageInstanceCreate(ctx context.Context, data *channel.StageInstanceParams, options ...discord.RequestOption) (*channel.StageInstance, error) {
	body, err := s.Request(ctx, http.MethodPost, discord.EndpointStageInstances, data, options...)
	if err != nil {
		return nil, err
	}

	var si channel.StageInstance
	return &si, s.Unmarshal(body, &si)
}

// StageInstance will retrieve a Stage instance by the ID of the types.ChannelGuildStageVoice.
func (s Requester) StageInstance(ctx context.Context, channelID string, options ...discord.RequestOption) (*channel.StageInstance, error) {
	body, err := s.Request(ctx, http.MethodGet, discord.EndpointStageInstance(channelID), nil, options...)
	if err != nil {
		return nil, err
	}

	var si channel.StageInstance
	return &si, s.Unmarshal(body, &si)
}

// StageInstanceEdit edits a Stage instance by ID the types.ChannelGuildStageVoice.
func (s Requester) StageInstanceEdit(ctx context.Context, channelID string, data *channel.StageInstanceParams, options ...discord.RequestOption) (*channel.StageInstance, error) {
	body, err := s.Request(ctx, http.MethodPatch, discord.EndpointStageInstance(channelID), data, options...)
	if err != nil {
		return nil, err
	}

	var si channel.StageInstance
	return &si, s.Unmarshal(body, &si)
}

// StageInstanceDelete deletes a Stage instance by ID of the types.ChannelGuildStageVoice.
func (s Requester) StageInstanceDelete(ctx context.Context, channelID string, options ...discord.RequestOption) error {
	_, err := s.Request(ctx, http.MethodDelete, discord.EndpointStageInstance(channelID), nil, options...)
	return err
}

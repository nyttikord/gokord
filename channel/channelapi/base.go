// Package channelapi contains everything to interact with everything located in the channel package.
package channelapi

import (
	"context"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/logger"
	"github.com/nyttikord/gokord/user/invite"
)

var (
// ErrReplyNilMessageRef = errors.New("reply attempted with nil message reference")
)

// Requester handles everything inside the channel package.
type Requester struct {
	REST
	State *State
}

// Invites returns all invite.Invite for the given channel.Channel.
func (r Requester) Invites(channelID string) Request[[]*invite.Invite] {
	return NewData[[]*invite.Invite](http.MethodDelete, discord.EndpointChannelInvites(channelID))
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

	return NewData[*invite.Invite](http.MethodPost, discord.EndpointChannelInvites(channelID)).
		WithData(data)
}

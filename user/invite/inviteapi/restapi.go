// Package inviteapi contains everything to interact with everything located in the invite package.
package inviteapi

import (
	"net/http"
	"net/url"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/user/invite"
)

// Requester handles everything inside the invite package.
type Requester struct {
	REST
}

// Invite returns the invite.Invite.
func (r Requester) Invite(inviteID string) Request[*Invite] {
	return r.InviteComplex(inviteID, "", false, false)
}

// InviteWithCounts returns the invite.Invite including approximate user.Member counts.
func (r Requester) InviteWithCounts(inviteID string) Request[*Invite] {
	return r.InviteComplex(inviteID, "", true, false)
}

// InviteComplex returns the invite.Invite the given ID including specified fields.
//
// If specified, it includes specified guild scheduled event with guildScheduledEventID.
// withCounts indicates whether to include approximate user.Member counts.
// withExpiration indicates whether to include expiration time.
func (r Requester) InviteComplex(inviteID, guildScheduledEventID string, withCounts, withExpiration bool) Request[*Invite] {
	endpoint := discord.EndpointInvite(inviteID)
	v := url.Values{}
	if guildScheduledEventID != "" {
		v.Set("guild_scheduled_event_id", guildScheduledEventID)
	}
	if withCounts {
		v.Set("with_counts", "true")
	}
	if withExpiration {
		v.Set("with_expiration", "true")
	}

	if len(v) != 0 {
		endpoint += "?" + v.Encode()
	}

	return NewSimpleData[*Invite](
		r, http.MethodGet, endpoint,
	).WithBucketID(discord.EndpointInvite(""))
}

// InviteDelete deletes an existing invite.Invite.
func (r Requester) InviteDelete(inviteID string) Request[*Invite] {
	return NewSimpleData[*Invite](
		r, http.MethodDelete, discord.EndpointInvite(inviteID),
	).WithBucketID(discord.EndpointInvite(""))
}

// InviteAccept accepts an invite.Invite.
func (r Requester) InviteAccept(inviteID string) Request[*Invite] {
	return NewSimpleData[*Invite](
		r, http.MethodPut, discord.EndpointInvite(inviteID),
	).WithBucketID(discord.EndpointInvite(""))
}

// TargetUsers returns a CSV with a single column Users containing the user.User IDs targetted by the invite.Invite.
func (r Requester) TargetUsers(inviteID string) Request[[]byte] {
	return NewSimple(
		r, http.MethodPut, discord.EndpointInviteTargetUsers(inviteID),
	).WithBucketID(discord.EndpointInvite(""))
}

// TargetUsersUpdate updates the user.User allowed to see and accept this
// See Invite TargetUsers.
/*func (r Requester) TargetUsersUpdate(ctx context.Context, inviteID string, csvFile string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		ctx,
		http.MethodGet,
		discord.EndpointInviteTargetUsers(inviteID)+"?target_users_file="+url.PathEscape(csvFile),
		nil,
		discord.EndpointInvite(""),
		options...,
	)
	return err
}*/

// The Discord's documentation does not yet provide complete information.
// Check https://discord.com/developers/docs/resources/invite#get-target-users-job-status for more information.
func (r Requester) TargetUsersJobStatus(inviteID string) Request[*TargetUsersJobStatus] {
	return NewSimpleData[*TargetUsersJobStatus](
		r, http.MethodPut, discord.EndpointInviteTargetUsersJobStatus(inviteID),
	).WithBucketID(discord.EndpointInvite(""))
}

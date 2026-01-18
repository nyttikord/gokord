// Package inviteapi contains everything to interact with everything located in the invite package.
package inviteapi

import (
	"net/http"
	"net/url"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/user/invite"
)

// Requester handles everything inside the invite package.
type Requester struct {
	discord.RESTRequester
}

// Invite returns the invite.Invite.
func (r Requester) Invite(inviteID string, options ...discord.RequestOption) (*invite.Invite, error) {
	body, err := r.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointInvite(inviteID),
		nil,
		discord.EndpointInvite(""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var i invite.Invite
	return &i, r.Unmarshal(body, &i)
}

// InviteWithCounts returns the invite.Invite including approximate user.Member counts.
func (r Requester) InviteWithCounts(inviteID string, options ...discord.RequestOption) (*invite.Invite, error) {
	body, err := r.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointInvite(inviteID)+"?with_counts=true",
		nil,
		discord.EndpointInvite(""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var i invite.Invite
	return &i, r.Unmarshal(body, &i)
}

// InviteComplex returns the invite.Invite the given ID including specified fields.
//
// If specified, it includes specified guild scheduled event with guildScheduledEventID.
// withCounts indicates whether to include approximate user.Member counts.
// withExpiration indicates whether to include expiration time.
func (r Requester) InviteComplex(inviteID, guildScheduledEventID string, withCounts, withExpiration bool, options ...discord.RequestOption) (*invite.Invite, error) {
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

	body, err := r.RequestWithBucketID(http.MethodGet, endpoint, nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return nil, err
	}

	var i invite.Invite
	return &i, r.Unmarshal(body, &i)
}

// InviteDelete deletes an existing invite.Invite.
func (r Requester) InviteDelete(inviteID string, options ...discord.RequestOption) (*invite.Invite, error) {
	body, err := r.RequestWithBucketID(http.MethodDelete, discord.EndpointInvite(inviteID), nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return nil, err
	}

	var i invite.Invite
	return &i, r.Unmarshal(body, &i)
}

// InviteAccept accepts an invite.Invite.
func (r Requester) InviteAccept(inviteID string, options ...discord.RequestOption) (*invite.Invite, error) {
	body, err := r.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointInvite(inviteID),
		nil,
		discord.EndpointInvite(""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var i invite.Invite
	return &i, r.Unmarshal(body, &i)
}

// TargetUsers returns a CSV with a single column Users containing the user.User IDs targetted by the invite.Invite.
func (r Requester) TargetUsers(inviteID string, options ...discord.RequestOption) ([]byte, error) {
	return r.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointInviteTargetUsers(inviteID),
		nil,
		discord.EndpointInvite(""),
		options...,
	)
}

// TargetUsersUpdate updates the user.User allowed to see and accept this invite.Invite.
// See invite.Invite TargetUsers.
/*func (r Requester) TargetUsersUpdate(inviteID string, csvFile string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
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
func (r Requester) TargetUsersJobStatus(inviteID string, options ...discord.RequestOption) (*invite.TargetUsersJobStatus, error) {
	b, err := r.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointInviteTargetUsersJobStatus(inviteID),
		nil,
		discord.EndpointInvite(""),
		options...,
	)
	if err != nil {
		return nil, err
	}
	var js invite.TargetUsersJobStatus
	r.Unmarshal(b, &js)
	return &js, nil
}

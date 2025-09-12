package inviteapi

import (
	"net/http"
	"net/url"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/user/invite"
)

type Requester struct {
	discord.Requester
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

// InviteWithCounts returns the invite.Invite including approximate member counts.
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

// InviteComplex returns the invite.Invite the given invite including specified fields.
//
// If specified, it includes specified guild scheduled event with guildScheduledEventID.
// withCounts indicates whether to include approximate member counts or not.
// withExpiration indicates whether to include expiration time or not.
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

// InviteAccept accepts an Invite to a Guild or Channel
// inviteID : The invite code
func (r Requester) InviteAccept(inviteID string, options ...discord.RequestOption) (*invite.Invite, error) {
	body, err := r.RequestWithBucketID("POST", discord.EndpointInvite(inviteID), nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return nil, err
	}

	var i invite.Invite
	return &i, r.Unmarshal(body, &i)
}

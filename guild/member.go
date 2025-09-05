package guild

import (
	"encoding/json"
	"time"
)

// MemberParams stores data needed to update a user.Member
// https://discord.com/developers/docs/resources/guild#modify-guild-member
type MemberParams struct {
	// Value to set user.Member's nickname to.
	Nick string `json:"nick,omitempty"`
	// Array of Role.ID the user.Member is assigned.
	Roles *[]string `json:"roles,omitempty"`
	// ID of channel to move user.Member to (if they are connected to voice).
	// Set to "" to remove user from a voice channel.Channel.
	ChannelID *string `json:"channel_id,omitempty"`
	// Whether the user.Member is muted in voice channels.
	Mute *bool `json:"mute,omitempty"`
	// Whether the user.Member is deafened in voice channels.
	Deaf *bool `json:"deaf,omitempty"`
	// When the user.Member's timeout will expire and the user will be able to communicate in the guild again (up to 28
	// days in the future).
	// Set to time.Time{} to remove timeout.
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until,omitempty"`
}

// MarshalJSON is a helper function to marshal MemberParams.
func (p MemberParams) MarshalJSON() (res []byte, err error) {
	type guildMemberParams MemberParams
	v := struct {
		guildMemberParams
		ChannelID                  json.RawMessage `json:"channel_id,omitempty"`
		CommunicationDisabledUntil json.RawMessage `json:"communication_disabled_until,omitempty"`
	}{guildMemberParams: guildMemberParams(p)}

	if p.ChannelID != nil {
		if *p.ChannelID == "" {
			v.ChannelID = json.RawMessage(`null`)
		} else {
			res, err = json.Marshal(p.ChannelID)
			if err != nil {
				return
			}
			v.ChannelID = res
		}
	}

	if p.CommunicationDisabledUntil != nil {
		if p.CommunicationDisabledUntil.IsZero() {
			v.CommunicationDisabledUntil = json.RawMessage(`null`)
		} else {
			res, err = json.Marshal(p.CommunicationDisabledUntil)
			if err != nil {
				return
			}
			v.CommunicationDisabledUntil = res
		}
	}

	return json.Marshal(v)
}

// MemberAddParams stores data needed to add a user.Member to a Guild.
//
// Note: All fields are optional, except AccessToken.
type MemberAddParams struct {
	// Valid access_token for the user.
	AccessToken string `json:"access_token"`
	// Value to set users nickname to.
	Nick string `json:"nick,omitempty"`
	// A list of role ID's to set on the member.
	Roles []string `json:"roles,omitempty"`
	// Whether the user is muted.
	Mute bool `json:"mute,omitempty"`
	// Whether the user is deafened.
	Deaf bool `json:"deaf,omitempty"`
}

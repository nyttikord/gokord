package guild

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/user"
)

var ErrPruneDaysBounds = errors.New("the number of days should be more than or equal to 1")

// MemberParams stores data needed to update a [user.Member].
// https://discord.com/developers/docs/resources/guild#modify-guild-member
type MemberParams struct {
	// Nick of the [user.Member].
	Nick string `json:"nick,omitempty"`
	// Array of [Role.ID] the [user.Member] is assigned.
	Roles *[]string `json:"roles,omitempty"`
	// ChannelID to move [user.Member] to (if they are connected to voice).
	// Set to "" to remove user from a voice [channel.Channel].
	ChannelID *string `json:"channel_id,omitempty"`
	// Whether the [user.Member] is muted in voice [channel.Channel]s.
	Mute *bool `json:"mute,omitempty"`
	// Whether the [user.Member] is deafened in voice [channel.Channel]s.
	Deaf *bool `json:"deaf,omitempty"`
	// When the [user.Member]'s timeout will expire (up to 28 days in the future).
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

// MemberAddParams stores data needed to add a [user.User] to a [Guild].
//
// See [MemberParams] for full documentations.
//
// NOTE: All fields are optional, except AccessToken.
type MemberAddParams struct {
	// Valid access_token for the user.
	AccessToken string   `json:"access_token"`
	Nick        string   `json:"nick,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Mute        bool     `json:"mute,omitempty"`
	Deaf        bool     `json:"deaf,omitempty"`
}

// ListBans returns [Ban] in the given [Guild].
//
// limit is the limit of bans to return (max 1000).
// If not empty, all returned Ban will be before the ID specified by beforeID.
// If not empty, all returned Ban will be after the ID specified by afterID.
func ListBans(guildID string, limit int, beforeID, afterID string) Request[[]*Ban] {
	uri := discord.EndpointGuildBans(guildID)

	v := url.Values{}
	if limit != 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if afterID != "" {
		v.Set("after", afterID)
	}

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	return NewData[[]*Ban](http.MethodGet, uri)
}

// CreateBan bans the given [user.User] from the given [Guild].
//
// days is the number of days of previous comments to delete.
//
// NOTE: See BanCreateWithReason.
func CreateBan(guildID, userID string, days int) Request[*Ban] {
	return CreateBanWithReason(guildID, userID, "", days)
}

// GetBan finds a ban with given [Guild] and [user.User].
func GetBan(guildID, userID string) Request[*Ban] {
	return NewData[*Ban](http.MethodGet, discord.EndpointGuildBan(guildID, userID)).
		WithBucketID(discord.EndpointGuildBans(guildID))
}

// CreateBanWithReason bans the given [user.User] from the given [Guild] with the given reason.
func CreateBanWithReason(guildID, userID, reason string, days int) Request[*Ban] {
	uri := discord.EndpointGuildBan(guildID, userID)

	queryParams := url.Values{}
	if days > 0 {
		queryParams.Set("delete_message_days", strconv.Itoa(days))
	}
	if reason != "" {
		queryParams.Set("reason", reason)
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	return NewData[*Ban](http.MethodPut, uri).
		WithBucketID(discord.EndpointGuildBans(guildID))
}

// ListMembers returns a list of members for a [Guild].
//
// If afterID is set, every member ID will be after this.
// limit is the maximum number of members to return (max 1000).
func ListMembers(guildID string, afterID string, limit int) Request[[]*user.Member] {
	uri := discord.EndpointGuildMembers(guildID)

	v := url.Values{}

	if afterID != "" {
		v.Set("after", afterID)
	}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	return NewCustom[[]*user.Member](http.MethodGet, uri).
		WithPost(func(ctx context.Context, b []byte) ([]*user.Member, error) {
			var m []*user.Member
			err := Unmarshal(ctx, b, &m)
			if err != nil {
				return nil, err
			}
			// The returned object doesn't have the GuildID attribute so we will set it here.
			for _, mem := range m {
				mem.GuildID = guildID
			}
			return m, nil
		})
}

// SearchMembers returns a list of [user.Member] whose username or nickname starts with the provided string.
//
// limit is the maximum number of members to return (min 1, max 1000).
func SearchMembers(guildID, query string, limit int) Request[[]*user.Member] {
	uri := discord.EndpointGuildMembersSearch(guildID)

	queryParams := url.Values{}
	queryParams.Set("query", query)
	if limit > 1 {
		queryParams.Set("limit", strconv.Itoa(limit))
		uri += "?" + queryParams.Encode()
	}

	return NewData[[]*user.Member](http.MethodGet, uri)
}

// GetMember returns a [user.Member] of a [Guild].
func GetMember(guildID, userID string) Request[*user.Member] {
	return NewCustom[*user.Member](http.MethodGet, discord.EndpointGuildMember(guildID, userID)).
		WithBucketID(discord.EndpointGuildMembers(guildID)).
		WithPost(func(ctx context.Context, b []byte) (*user.Member, error) {
			var m user.Member
			err := Unmarshal(ctx, b, &m)
			if err != nil {
				return nil, err
			}
			// The returned object doesn't have the GuildID attribute so we will set it here.
			m.GuildID = guildID
			return &m, err
		})
}

// AddMember force joins a [user.User] to the [Guild].
func AddMember(guildID, userID string, data *MemberAddParams) Empty {
	req := NewSimple(http.MethodPut, discord.EndpointGuildMember(guildID, userID)).
		WithBucketID(discord.EndpointGuildMembers(guildID))
	return WrapAsEmpty(req)
}

// Kick the given [user.User] from the given [Guild].
func Kick(guildID, userID string) Empty {
	return KickWithReason(guildID, userID, "")
}

// KickWithReason removes the given [user.User] from the given [Guild] with the given reason.
func KickWithReason(guildID, userID, reason string) Empty {
	uri := discord.EndpointGuildMember(guildID, userID)
	if reason != "" {
		uri += "?reason=" + url.QueryEscape(reason)
	}

	req := NewSimple(http.MethodDelete, uri).
		WithBucketID(discord.EndpointGuildMembers(guildID))
	return WrapAsEmpty(req)
}

// UpdateMember with the given data and returns them.
func UpdateMember(guildID, userID string, data *MemberParams) Request[*user.Member] {
	return NewData[*user.Member](http.MethodPatch, discord.EndpointGuildMember(guildID, userID)).
		WithBucketID(discord.EndpointGuildMembers(guildID))
}

// MoveMember from one voice [channel.Channel] to another/none.
func MoveMember(guildID string, userID string, channelID *string) Empty {
	data := struct {
		ChannelID *string `json:"channel_id"`
	}{channelID}

	req := NewSimple(http.MethodPatch, discord.EndpointGuildMember(guildID, userID)).
		WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// UpdateMemberNickname updates the nickname of a [user.Member] in a [Guild].
//
// NOTE: To reset the nickname, set it to an empty string.
//
// NOTE: Use [UpdateCurrentMember] to modify the current member.
func UpdateMemberNickname(guildID, userID, nickname string) Empty {
	if userID == "@me" {
		/*r.Logger().WarnContext(
			logger.NewContext(context.Background(), 1),
			"this endpoint is deprecated for the current member, use MemberModifyCurrent instead",
		)*/
		return UpdateCurrentMember(guildID, nickname, "", "", "")
	}

	data := struct {
		Nick string `json:"nick"`
	}{nickname}

	req := NewSimple(http.MethodPatch, discord.EndpointGuildMember(guildID, userID)).
		WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// UpdateCurrentMember updates the nickname, the avatar, the banner and the bio of the current user.Member.
//
// NOTE: Set any parameter to "" to avoid modifying it.
func UpdateCurrentMember(guildID string, nick, avatar, banner, bio string) Empty {
	data := struct {
		Nick   string `json:"nick,omitempty"`
		Avatar string `json:"avatar,omitempty"`
		Banner string `json:"banner,omitempty"`
		Bio    string `json:"bio,omitempty"`
	}{nick, avatar, banner, bio}

	req := NewSimple(http.MethodPatch, discord.EndpointGuildMember(guildID, "@me")).
		WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// MuteMember (or unmute) a [user.Member] in a [Guild].
func MuteMember(guildID string, userID string, mute bool) Empty {
	data := struct {
		Mute bool `json:"mute"`
	}{mute}

	req := NewSimple(http.MethodPatch, discord.EndpointGuildMember(guildID, userID)).
		WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// Timeout a [user.Member] in a [Guild].
//
// NOTE: Set until to nil to remove timeout.
func Timeout(guildID string, userID string, until *time.Time) Empty {
	data := struct {
		CommunicationDisabledUntil *time.Time `json:"communication_disabled_until"`
	}{until}

	req := NewSimple(http.MethodPatch, discord.EndpointGuildMember(guildID, userID)).
		WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// DeafenMember in a [Guild].
func DeafenMember(guildID string, userID string, deaf bool) Empty {
	data := struct {
		Deaf bool `json:"deaf"`
	}{deaf}

	req := NewSimple(http.MethodPatch, discord.EndpointGuildMember(guildID, userID)).
		WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// PruneCount returns the number of user.Member that would be removed in a prune operation.
//
// Requires discord.PermissionKickMembers.
func PruneCount(guildID string, days uint32) Request[uint32] {
	if days <= 0 {
		return NewError[uint32](ErrPruneDaysBounds)
	}

	uri := discord.EndpointGuildPrune(guildID) + "?days=" + strconv.FormatUint(uint64(days), 10)
	return NewCustom[uint32](http.MethodGet, uri).
		WithPost(func(ctx context.Context, b []byte) (uint32, error) {
			var p struct {
				Pruned uint32 `json:"pruned"`
			}
			return p.Pruned, Unmarshal(ctx, b, &p)
		})
}

// Prune begins as prune operation.
// Returns the number of pruned members.
//
// Requires discord.PermissionKickMembers.
func Prune(guildID string, days uint32) Request[uint32] {
	if days <= 0 {
		return NewError[uint32](ErrPruneDaysBounds)
	}

	data := struct {
		days uint32
	}{days}

	return NewCustom[uint32](http.MethodGet, discord.EndpointGuildPrune(guildID)).
		WithData(data).
		WithPost(func(ctx context.Context, b []byte) (uint32, error) {
			var p struct {
				Pruned uint32 `json:"pruned"`
			}
			return p.Pruned, Unmarshal(ctx, b, &p)
		})
}

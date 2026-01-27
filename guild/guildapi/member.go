package guildapi

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/logger"
	"github.com/nyttikord/gokord/user"
)

// Bans returns guild.Ban in the given guild.Guild.
//
// limit is the limit of bans to return (max 1000).
// If not empty, all returned Ban will be before the ID specified by beforeID.
// If not empty, all returned Ban will be after the ID specified by afterID.
func (r Requester) Bans(guildID string, limit int, beforeID, afterID string) Request[[]*Ban] {
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

	return NewData[[]*Ban](r, http.MethodGet, uri)
}

// BanCreate bans the given user.User from the given guild.Guild.
//
// days is the number of days of previous comments to delete.
//
// NOTE: See BanCreateWithReason.
func (r Requester) BanCreate(guildID, userID string, days int) Request[*Ban] {
	return r.BanCreateWithReason(guildID, userID, "", days)
}

// Ban finds ban by given guild.Guild and user.User id and returns guild.Ban.
func (r Requester) Ban(guildID, userID string) Request[*Ban] {
	return NewData[*Ban](
		r, http.MethodGet, discord.EndpointGuildBan(guildID, userID),
	).WithBucketID(discord.EndpointGuildBans(guildID))
}

// BanCreateWithReason bans the given user.User from the given guild.Guild with the given reason.
func (r Requester) BanCreateWithReason(guildID, userID, reason string, days int) Request[*Ban] {
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

	return NewData[*Ban](
		r, http.MethodPut, uri,
	).WithBucketID(discord.EndpointGuildBans(guildID))
}

// Members returns a list of members for a guild.Guild.
// If afterID is set, every member ID will be after this.
// limit is the maximum number of members to return (max 1000).
func (r Requester) Members(guildID string, afterID string, limit int) Request[[]*user.Member] {
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

	return NewCustom[[]*user.Member](r, http.MethodGet, uri).
		WithPost(func(ctx context.Context, b []byte) ([]*user.Member, error) {
			var m []*user.Member
			err := r.Unmarshal(b, &m)
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

// MembersSearch returns a list of user.Member whose username or nickname starts with a provided string.
// limit is the maximum number of members to return (min 1, max 1000).
func (r Requester) MembersSearch(guildID, query string, limit int) Request[[]*user.Member] {
	uri := discord.EndpointGuildMembersSearch(guildID)

	queryParams := url.Values{}
	queryParams.Set("query", query)
	if limit > 1 {
		queryParams.Set("limit", strconv.Itoa(limit))
		uri += "?" + queryParams.Encode()
	}

	return NewData[[]*user.Member](r, http.MethodGet, uri)
}

// Member returns a user.Member of a guild.Guild.
func (r Requester) Member(guildID, userID string) Request[*user.Member] {
	return NewCustom[*user.Member](r, http.MethodGet, discord.EndpointGuildMember(guildID, userID)).
		WithBucketID(discord.EndpointGuildMembers(guildID)).
		WithPost(func(ctx context.Context, b []byte) (*user.Member, error) {
			var m user.Member
			err := r.Unmarshal(b, &m)
			if err != nil {
				return nil, err
			}
			// The returned object doesn't have the GuildID attribute so we will set it here.
			m.GuildID = guildID
			return &m, err
		})
}

// MemberAdd force joins a user.User to the guild.Guild with the given data.
func (r Requester) MemberAdd(guildID, userID string, data *MemberAddParams) Empty {
	req := NewSimple(
		r, http.MethodPut, discord.EndpointGuildMember(guildID, userID),
	).WithBucketID(discord.EndpointGuildMembers(guildID))
	return WrapAsEmpty(req)
}

// MemberKick kicks the given user.User from the given guild.Guild.
func (r Requester) MemberKick(guildID, userID string) Empty {
	return r.MemberKickWithReason(guildID, userID, "")
}

// MemberKickWithReason removes the given user.User from the given guild.Guild with the given reason.
func (r Requester) MemberKickWithReason(guildID, userID, reason string) Empty {
	uri := discord.EndpointGuildMember(guildID, userID)
	if reason != "" {
		uri += "?reason=" + url.QueryEscape(reason)
	}

	req := NewSimple(
		r, http.MethodDelete, uri,
	).WithBucketID(discord.EndpointGuildMembers(guildID))
	return WrapAsEmpty(req)
}

// MemberEdit edits a user.Member with the given data and returns them.
func (r Requester) MemberEdit(guildID, userID string, data *MemberParams) Request[*user.Member] {
	return NewData[*user.Member](
		r, http.MethodPatch, discord.EndpointGuildMember(guildID, userID),
	).WithBucketID(discord.EndpointGuildMembers(guildID))
}

// MemberMove moves a user.Member from one voice channel.Channel to another/none.
//
// NOTE: I am not entirely set on the name of this function, and it may change.
func (r Requester) MemberMove(guildID string, userID string, channelID *string) Empty {
	data := struct {
		ChannelID *string `json:"channel_id"`
	}{channelID}

	req := NewSimple(
		r, http.MethodPatch, discord.EndpointGuildMember(guildID, userID),
	).WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// MemberNickname updates the nickname of a user.Member in a guild.Guild.
//
// NOTE: To reset the nickname, set it to an empty string.
//
// NOTE: Use MemberModifyCurrent to modify the current member.
func (r Requester) MemberNickname(guildID, userID, nickname string) Empty {
	if userID == "@me" {
		r.Logger().WarnContext(
			logger.NewContext(context.Background(), 1),
			"this endpoint is deprecated for the current member, use MemberModifyCurrent instead",
		)
		return r.MemberModifyCurrent(guildID, nickname, "", "", "")
	}

	data := struct {
		Nick string `json:"nick"`
	}{nickname}

	req := NewSimple(
		r, http.MethodPatch, discord.EndpointGuildMember(guildID, userID),
	).WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// MemberModifyCurrent updates the nickname, the avatar, the banner and the bio of the current user.Member.
//
// NOTE: Set any parameter to "" to avoid modifying it.
func (r Requester) MemberModifyCurrent(guildID string, nick, avatar, banner, bio string) Empty {
	data := struct {
		Nick   string `json:"nick,omitempty"`
		Avatar string `json:"avatar,omitempty"`
		Banner string `json:"banner,omitempty"`
		Bio    string `json:"bio,omitempty"`
	}{nick, avatar, banner, bio}

	req := NewSimple(
		r, http.MethodPatch, discord.EndpointGuildMember(guildID, "@me"),
	).WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// MemberMute (un)mutes a user.Member in a guild.Guild.
func (r Requester) MemberMute(guildID string, userID string, mute bool) Empty {
	data := struct {
		Mute bool `json:"mute"`
	}{mute}

	req := NewSimple(
		r, http.MethodPatch, discord.EndpointGuildMember(guildID, userID),
	).WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// MemberTimeout times out a user.Member in a guild.Guild.
//
// NOTE: Set until to nil to remove timeout.
func (r Requester) MemberTimeout(guildID string, userID string, until *time.Time) Empty {
	data := struct {
		CommunicationDisabledUntil *time.Time `json:"communication_disabled_until"`
	}{until}

	req := NewSimple(
		r, http.MethodPatch, discord.EndpointGuildMember(guildID, userID),
	).WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// MemberDeafen server deafens a user.Member in a guild.Guild.
func (r Requester) MemberDeafen(guildID string, userID string, deaf bool) Empty {
	data := struct {
		Deaf bool `json:"deaf"`
	}{deaf}

	req := NewSimple(
		r, http.MethodPatch, discord.EndpointGuildMember(guildID, userID),
	).WithBucketID(discord.EndpointGuildMembers(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// MemberRoleAdd adds the specified guild.Role to a given user.Member.
func (r Requester) MemberRoleAdd(guildID, userID, roleID string) Empty {
	req := NewSimple(
		r, http.MethodPut, discord.EndpointGuildMemberRole(guildID, userID, roleID),
	).WithBucketID(discord.EndpointGuildMembers(guildID))
	return WrapAsEmpty(req)
}

// MemberRoleRemove removes the specified guild.Role to a given user.Member.
func (r Requester) MemberRoleRemove(guildID, userID, roleID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointGuildMemberRole(guildID, userID, roleID),
	).WithBucketID(discord.EndpointGuildMembers(guildID))
	return WrapAsEmpty(req)
}

// PruneCount returns the number of user.Member that would be removed in a prune operation.
//
// Requires discord.PermissionKickMembers.
func (r Requester) PruneCount(guildID string, days uint32) Request[uint32] {
	if days <= 0 {
		return NewError[uint32](ErrPruneDaysBounds)
	}

	uri := discord.EndpointGuildPrune(guildID) + "?days=" + strconv.FormatUint(uint64(days), 10)
	return NewCustom[uint32](r, http.MethodGet, uri).
		WithPost(func(ctx context.Context, b []byte) (uint32, error) {
			var p struct {
				Pruned uint32 `json:"pruned"`
			}
			return p.Pruned, r.Unmarshal(b, &p)
		})
}

// Prune begins as prune operation.
// Returns the number of pruned members.
//
// Requires discord.PermissionKickMembers.
func (r Requester) Prune(guildID string, days uint32) Request[uint32] {
	if days <= 0 {
		return NewError[uint32](ErrPruneDaysBounds)
	}

	data := struct {
		days uint32
	}{days}

	return NewCustom[uint32](r, http.MethodGet, discord.EndpointGuildPrune(guildID)).
		WithData(data).
		WithPost(func(ctx context.Context, b []byte) (uint32, error) {
			var p struct {
				Pruned uint32 `json:"pruned"`
			}
			return p.Pruned, r.Unmarshal(b, &p)
		})
}

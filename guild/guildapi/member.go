package guildapi

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
)

// Bans returns guild.Ban in the given guild.
//
// limit is the limit of bans to return (max 1000).
// If not empty, all returned guild.Ban will be before the ID specified by beforeID.
// If not empty, all returned guild.Ban will be after the ID specified by afterID.
func (r Requester) Bans(guildID string, limit int, beforeID, afterID string, options ...discord.RequestOption) ([]*guild.Ban, error) {
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

	body, err := r.Request(http.MethodGet, uri, nil, options...)
	if err != nil {
		return nil, err
	}

	var b []*guild.Ban
	return b, r.Unmarshal(body, &b)
}

// BanCreate bans the given user.User from the given guild.Guild.
//
// days is the number of days of previous comments to delete.
//
// NOTE: See BanCreate.
func (r Requester) BanCreate(guildID, userID string, days int, options ...discord.RequestOption) error {
	return r.BanCreateWithReason(guildID, userID, "", days, options...)
}

// Ban finds ban by given guild.Guild and user.User id and returns guild.Ban.
func (r Requester) Ban(guildID, userID string, options ...discord.RequestOption) (*guild.Ban, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuildBan(guildID, userID), nil, options...)
	if err != nil {
		return nil, err
	}

	var b guild.Ban
	return &b, r.Unmarshal(body, &b)
}

// BanCreateWithReason bans the given user.User from the given guild.Guild also providing a reason.
//
// NOTE: See BanCreate.
func (r Requester) BanCreateWithReason(guildID, userID, reason string, days int, options ...discord.RequestOption) error {
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

	_, err := r.RequestWithBucketID(http.MethodPut, uri, nil, discord.EndpointGuildBan(guildID, ""), options...)
	return err
}

// BanDelete unbans the given user.User from the given guild.Guild.
func (r Requester) BanDelete(guildID, userID string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildBan(guildID, userID),
		nil,
		discord.EndpointGuildBan(guildID, ""),
		options...,
	)
	return err
}

// Members returns a list of members for a guild.Guild.
// If afterID is set, every member ID will be after this.
// limit is the maximum number of members to return (max 1000).
func (r Requester) Members(guildID string, afterID string, limit int, options ...discord.RequestOption) ([]*user.Member, error) {

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

	body, err := r.Request(http.MethodGet, uri, nil, options...)
	if err != nil {
		return nil, err
	}

	var m []*user.Member
	return m, r.Unmarshal(body, &m)
}

// MembersSearch returns a list of user.Member whose username or nickname starts with a provided string.
// limit is the maximum number of members to return (min 1, max 1000).
func (r Requester) MembersSearch(guildID, query string, limit int, options ...discord.RequestOption) ([]*user.Member, error) {

	uri := discord.EndpointGuildMembersSearch(guildID)

	queryParams := url.Values{}
	queryParams.Set("query", query)
	if limit > 1 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}

	body, err := r.Request(http.MethodGet, uri+"?"+queryParams.Encode(), nil, options...)
	if err != nil {
		return nil, err
	}

	var m []*user.Member
	return m, r.Unmarshal(body, &m)
}

// Member returns a user.Member of a guild.Guild.
func (r Requester) Member(guildID, userID string, options ...discord.RequestOption) (*user.Member, error) {
	body, err := r.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildMember(guildID, userID),
		nil,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var m user.Member
	err = r.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	// The returned object doesn't have the GuildID attribute so we will set it here.
	m.GuildID = guildID
	return &m, err
}

// MemberAdd force joins a user.User to the guild.Guild with the given data.
func (r Requester) MemberAdd(guildID, userID string, data *guild.MemberAddParams, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodPut,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// MemberKick kicks the given user.User from the given guild.Guild.
func (r Requester) MemberKick(guildID, userID string, options ...discord.RequestOption) error {
	return r.MemberKickWithReason(guildID, userID, "", options...)
}

// MemberKickWithReason removes the given user.User from the given guild.Guild with the given reason.
func (r Requester) MemberKickWithReason(guildID, userID, reason string, options ...discord.RequestOption) error {
	uri := discord.EndpointGuildMember(guildID, userID)
	if reason != "" {
		uri += "?reason=" + url.QueryEscape(reason)
	}

	_, err := r.RequestWithBucketID(http.MethodDelete, uri, nil, discord.EndpointGuildMember(guildID, ""), options...)
	return err
}

// MemberEdit edits a user.Member with the given data and returns them.
func (r Requester) MemberEdit(guildID, userID string, data *guild.MemberParams, options ...discord.RequestOption) (*user.Member, error) {
	body, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var m user.Member
	return &m, r.Unmarshal(body, &m)
}

// MemberMove moves a user.Member from one voice channel.Channel to another/none.
//
// NOTE: I am not entirely set on the name of this function, and it may change.
func (r Requester) MemberMove(guildID string, userID string, channelID *string, options ...discord.RequestOption) error {
	data := struct {
		ChannelID *string `json:"channel_id"`
	}{channelID}

	_, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// MemberNickname updates the nickname of a user.Member in a guild.Guild.
//
// NOTE: To reset the nickname, set it to an empty string.
func (r Requester) MemberNickname(guildID, userID, nickname string, options ...discord.RequestOption) error {
	data := struct {
		Nick string `json:"nick"`
	}{nickname}

	if userID == "@me" {
		userID += "/nick"
	}

	_, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// MemberMute (un)mutes a user.Member in a guild.Guild.
func (r Requester) MemberMute(guildID string, userID string, mute bool, options ...discord.RequestOption) error {
	data := struct {
		Mute bool `json:"mute"`
	}{mute}

	_, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// MemberTimeout times out a user.Member in a guild.Guild.
//
// NOTE: Set until to nil to remove timeout.
func (r Requester) MemberTimeout(guildID string, userID string, until *time.Time, options ...discord.RequestOption) error {
	data := struct {
		CommunicationDisabledUntil *time.Time `json:"communication_disabled_until"`
	}{until}

	_, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// MemberDeafen server deafens a user.Member in a guild.Guild.
func (r Requester) MemberDeafen(guildID string, userID string, deaf bool, options ...discord.RequestOption) error {
	data := struct {
		Deaf bool `json:"deaf"`
	}{deaf}

	_, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// MemberRoleAdd adds the specified guild.Role to a given user.Member.
func (r Requester) MemberRoleAdd(guildID, userID, roleID string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodPut,
		discord.EndpointGuildMemberRole(guildID, userID, roleID),
		nil,
		discord.EndpointGuildMemberRole(guildID, "", ""),
		options...,
	)
	return err
}

// MemberRoleRemove removes the specified guild.Role to a given user.Member.
func (r Requester) MemberRoleRemove(guildID, userID, roleID string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildMemberRole(guildID, userID, roleID),
		nil,
		discord.EndpointGuildMemberRole(guildID, "", ""),
		options...,
	)
	return err
}

// PruneCount returns the number of user.Member that would be removed in a prune operation.
//
// Requires discord.PermissionKickMembers.
func (r Requester) PruneCount(guildID string, days uint32, options ...discord.RequestOption) (uint32, error) {
	if days <= 0 {
		return 0, ErrPruneDaysBounds
	}

	p := struct {
		Pruned uint32 `json:"pruned"`
	}{}

	uri := discord.EndpointGuildPrune(guildID) + "?days=" + strconv.FormatUint(uint64(days), 10)
	body, err := r.Request(http.MethodGet, uri, nil, options...)
	if err != nil {
		return 0, err
	}

	err = r.Unmarshal(body, &p)
	if err != nil {
		return 0, err
	}

	return p.Pruned, err
}

// Prune begins as prune operation.
// Returns the number of pruned members.
//
// Requires discord.PermissionKickMembers.
func (r Requester) Prune(guildID string, days uint32, options ...discord.RequestOption) (uint32, error) {
	if days <= 0 {
		return 0, ErrPruneDaysBounds
	}

	data := struct {
		days uint32
	}{days}

	p := struct {
		Pruned uint32 `json:"pruned"`
	}{}

	body, err := r.Request(http.MethodPost, discord.EndpointGuildPrune(guildID), data, options...)
	if err != nil {
		return 0, err
	}

	err = r.Unmarshal(body, &p)
	if err != nil {
		return 0, err
	}

	return p.Pruned, nil
}

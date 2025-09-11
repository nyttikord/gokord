package gokord

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/invite"
)

// UserGuilds returns an array of guild.UserGuild structures for all guilds.
//
// limit is the number of guilds that can be returned (max 200).
// If beforeID is set, it will return all guilds before this ID.
// If afterID is set, it will return all guilds after this ID.
// Set withCounts to true if you want to include approximate member and presence counts.
func (s *Session) UserGuilds(limit int, beforeID, afterID string, withCounts bool, options ...discord.RequestOption) ([]*guild.UserGuild, error) {
	v := url.Values{}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if withCounts {
		v.Set("with_counts", "true")
	}

	uri := discord.EndpointUserGuilds("@me")

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	body, err := s.RequestWithBucketID(http.MethodGet, uri, nil, discord.EndpointUserGuilds(""), options...)
	if err != nil {
		return nil, err
	}

	var ug []*guild.UserGuild
	return ug, s.Unmarshal(body, &ug)
}

// Guild returns the guild.Guild with the given guildID.
func (s *Session) Guild(guildID string, options ...discord.RequestOption) (*guild.Guild, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuild(guildID),
		nil,
		discord.EndpointGuild(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var g guild.Guild
	return &g, unmarshal(body, &g)
}

// GuildWithCounts returns the guild.Guild with the given guildID with approximate user.Member and status.Presence counts.
func (s *Session) GuildWithCounts(guildID string, options ...discord.RequestOption) (*guild.Guild, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuild(guildID)+"?with_counts=true",
		nil,
		discord.EndpointGuild(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var g guild.Guild
	return &g, unmarshal(body, &g)
}

// GuildPreview returns the guild.Preview for the given public guild.Guild guildID.
func (s *Session) GuildPreview(guildID string, options ...discord.RequestOption) (*guild.Preview, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildPreview(guildID),
		nil,
		discord.EndpointGuildPreview(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var gp guild.Preview
	return &gp, unmarshal(body, &gp)
}

// GuildCreate creates a new guild.Guild with the given name.
func (s *Session) GuildCreate(name string, options ...discord.RequestOption) (*guild.Guild, error) {
	data := struct {
		Name string `json:"name"`
	}{name}

	body, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointGuildCreate,
		data,
		discord.EndpointGuildCreate,
		options...,
	)
	if err != nil {
		return nil, err
	}

	var g guild.Guild
	return &g, unmarshal(body, &g)
}

// GuildEdit edits a guild.Guild with the given params.
func (s *Session) GuildEdit(guildID string, params *guild.Params, options ...discord.RequestOption) (*guild.Guild, error) {
	// Bounds checking for VerificationLevel, interval: [0, 4]
	if params.VerificationLevel != nil {
		val := *params.VerificationLevel
		if val < 0 || val > 4 {
			return nil, ErrVerificationLevelBounds
		}
	}

	// Bounds checking for regions
	if params.Region != "" {
		isValid := false
		regions, _ := s.VoiceRegions(options...)
		for _, r := range regions {
			if params.Region == r.ID {
				isValid = true
			}
		}
		if !isValid {
			var valid []string
			for _, r := range regions {
				valid = append(valid, r.ID)
			}
			return nil, fmt.Errorf("not a valid region (%q)", valid)
		}
	}

	body, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuild(guildID),
		params,
		discord.EndpointGuild(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var g guild.Guild
	return &g, unmarshal(body, &g)
}

// GuildDelete deletes a guild.Guild.
func (s *Session) GuildDelete(guildID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(http.MethodDelete, discord.EndpointGuild(guildID), nil, discord.EndpointGuild(guildID), options...)
	return err
}

// GuildLeave leaves a guild.Guild.
func (s *Session) GuildLeave(guildID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointUserGuild("@me", guildID),
		nil,
		discord.EndpointUserGuild("", guildID),
		options...,
	)
	return err
}

// GuildBans returns the guild.Ban in the given guild.
//
// limit is the limit of bans to return (max 1000).
// If not empty, all returned guild.Ban will be before the ID specified by beforeID.
// If not empty, all returned guild.Ban will be after the ID specified by afterID.
func (s *Session) GuildBans(guildID string, limit int, beforeID, afterID string, options ...discord.RequestOption) ([]*guild.Ban, error) {
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

	body, err := s.RequestWithBucketID(http.MethodGet, uri, nil, discord.EndpointGuildBans(guildID), options...)
	if err != nil {
		return nil, err
	}

	var b []*guild.Ban
	return b, unmarshal(body, &b)
}

// GuildBanCreate bans the given user.User from the given guild.Guild.
//
// days is the number of days of previous comments to delete.
//
// Note: See GuildBanCreate.
func (s *Session) GuildBanCreate(guildID, userID string, days int, options ...discord.RequestOption) error {
	return s.GuildBanCreateWithReason(guildID, userID, "", days, options...)
}

// GuildBan finds ban by given guild.Guild and user.User id and returns guild.Ban structure
func (s *Session) GuildBan(guildID, userID string, options ...discord.RequestOption) (*guild.Ban, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildBan(guildID, userID),
		nil,
		discord.EndpointGuildBan(guildID, userID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var b guild.Ban
	return &b, unmarshal(body, &b)
}

// GuildBanCreateWithReason bans the given user.User from the given guild.Guild also providing a reason.
//
// Note: See GuildBanCreate.
func (s *Session) GuildBanCreateWithReason(guildID, userID, reason string, days int, options ...discord.RequestOption) error {
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

	_, err := s.RequestWithBucketID(http.MethodPut, uri, nil, discord.EndpointGuildBan(guildID, ""), options...)
	return err
}

// GuildUnban unbans the given user.User from the given guild.Guild
func (s *Session) GuildUnban(guildID, userID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildBan(guildID, userID),
		nil,
		discord.EndpointGuildBan(guildID, ""),
		options...,
	)
	return err
}

// GuildMembers returns a list of members for a guild.Guild.
// If afterID is set, every member ID will be after this.
// limit is the maximum number of members to return (max 1000).
func (s *Session) GuildMembers(guildID string, afterID string, limit int, options ...discord.RequestOption) ([]*user.Member, error) {

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

	body, err := s.RequestWithBucketID(http.MethodGet, uri, nil, discord.EndpointGuildMembers(guildID), options...)
	if err != nil {
		return nil, err
	}

	var m []*user.Member
	return m, unmarshal(body, &m)
}

// GuildMembersSearch returns a list of user.Member whose username or nickname starts with a provided string.
// limit is the maximum number of members to return (min 1, max 1000).
func (s *Session) GuildMembersSearch(guildID, query string, limit int, options ...discord.RequestOption) ([]*user.Member, error) {

	uri := discord.EndpointGuildMembersSearch(guildID)

	queryParams := url.Values{}
	queryParams.Set("query", query)
	if limit > 1 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}

	body, err := s.RequestWithBucketID(http.MethodGet, uri+"?"+queryParams.Encode(), nil, uri, options...)
	if err != nil {
		return nil, err
	}

	var m []*user.Member
	return m, unmarshal(body, &m)
}

// GuildMember returns a user.Member of a guild.Guild.
func (s *Session) GuildMember(guildID, userID string, options ...discord.RequestOption) (*user.Member, error) {
	body, err := s.RequestWithBucketID(
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
	err = unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	// The returned object doesn't have the GuildID attribute so we will set it here.
	m.GuildID = guildID
	return &m, err
}

// GuildMemberAdd force joins a user.User to the guild.Guild with the given data.
func (s *Session) GuildMemberAdd(guildID, userID string, data *guild.MemberAddParams, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodPut,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// GuildMemberKick kicks the given user.User from the given guild.Guild.
func (s *Session) GuildMemberKick(guildID, userID string, options ...discord.RequestOption) error {
	return s.GuildMemberKickWithReason(guildID, userID, "", options...)
}

// GuildMemberKickWithReason removes the given user.User from the given guild.Guild with the given reason.
func (s *Session) GuildMemberKickWithReason(guildID, userID, reason string, options ...discord.RequestOption) error {
	uri := discord.EndpointGuildMember(guildID, userID)
	if reason != "" {
		uri += "?reason=" + url.QueryEscape(reason)
	}

	_, err := s.RequestWithBucketID(http.MethodDelete, uri, nil, discord.EndpointGuildMember(guildID, ""), options...)
	return err
}

// GuildMemberEdit edits a user.Member with the given data and returns them.
func (s *Session) GuildMemberEdit(guildID, userID string, data *guild.MemberParams, options ...discord.RequestOption) (*user.Member, error) {
	body, err := s.RequestWithBucketID(
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
	return &m, unmarshal(body, &m)
}

// GuildMemberMove moves a user.Member from one voice channel.Channel to another/none.
//
// Note: I am not entirely set on the name of this function, and it may change.
func (s *Session) GuildMemberMove(guildID string, userID string, channelID *string, options ...discord.RequestOption) error {
	data := struct {
		ChannelID *string `json:"channel_id"`
	}{channelID}

	_, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// GuildMemberNickname updates the nickname of a user.Member in a guild.Guild.
//
// Note: To reset the nickname, set it to an empty string.
func (s *Session) GuildMemberNickname(guildID, userID, nickname string, options ...discord.RequestOption) error {
	data := struct {
		Nick string `json:"nick"`
	}{nickname}

	if userID == "@me" {
		userID += "/nick"
	}

	_, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// GuildMemberMute (un)mutes a user.Member in a guild.Guild.
func (s *Session) GuildMemberMute(guildID string, userID string, mute bool, options ...discord.RequestOption) error {
	data := struct {
		Mute bool `json:"mute"`
	}{mute}

	_, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// GuildMemberTimeout times out a user.Member in a guild.Guild.
//
// Note: Set until to nil to remove timeout.
func (s *Session) GuildMemberTimeout(guildID string, userID string, until *time.Time, options ...discord.RequestOption) error {
	data := struct {
		CommunicationDisabledUntil *time.Time `json:"communication_disabled_until"`
	}{until}

	_, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// GuildMemberDeafen server deafens a user.Member in a guild.Guild.
func (s *Session) GuildMemberDeafen(guildID string, userID string, deaf bool, options ...discord.RequestOption) error {
	data := struct {
		Deaf bool `json:"deaf"`
	}{deaf}

	_, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildMember(guildID, userID),
		data,
		discord.EndpointGuildMember(guildID, ""),
		options...,
	)
	return err
}

// GuildMemberRoleAdd adds the specified guild.Role to a given user.Member.
func (s *Session) GuildMemberRoleAdd(guildID, userID, roleID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodPut,
		discord.EndpointGuildMemberRole(guildID, userID, roleID),
		nil,
		discord.EndpointGuildMemberRole(guildID, "", ""),
		options...,
	)
	return err
}

// GuildMemberRoleRemove removes the specified guild.Role to a given user.Member.
func (s *Session) GuildMemberRoleRemove(guildID, userID, roleID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildMemberRole(guildID, userID, roleID),
		nil,
		discord.EndpointGuildMemberRole(guildID, "", ""),
		options...,
	)
	return err
}

// GuildChannels returns the channel.Channel in the guild.Guild.
func (s *Session) GuildChannels(guildID string, options ...discord.RequestOption) ([]*channel.Channel, error) {
	body, err := s.RequestRaw(
		http.MethodGet,
		discord.EndpointGuildChannels(guildID),
		"",
		nil,
		discord.EndpointGuildChannels(guildID),
		0,
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st []*channel.Channel
	return st, unmarshal(body, &st)
}

// GuildChannelCreateData is provided to Session.GuildChannelCreateComplex
type GuildChannelCreateData struct {
	Name                 string                         `json:"name"`
	Type                 types.Channel                  `json:"type"`
	Topic                string                         `json:"topic,omitempty"`
	Bitrate              int                            `json:"bitrate,omitempty"`
	UserLimit            int                            `json:"user_limit,omitempty"`
	RateLimitPerUser     int                            `json:"rate_limit_per_user,omitempty"`
	Position             int                            `json:"position,omitempty"`
	PermissionOverwrites []*channel.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             string                         `json:"parent_id,omitempty"`
	NSFW                 bool                           `json:"nsfw,omitempty"`
}

// GuildChannelCreateComplex creates a new channel in the given guild.Guild
func (s *Session) GuildChannelCreateComplex(guildID string, data GuildChannelCreateData, options ...discord.RequestOption) (*channel.Channel, error) {
	body, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointGuildChannels(guildID),
		data,
		discord.EndpointGuildChannels(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st channel.Channel
	return &st, unmarshal(body, &st)
}

// GuildChannelCreate creates a new channel.Channel in the given guild.Guild.
func (s *Session) GuildChannelCreate(guildID, name string, ctype types.Channel, options ...discord.RequestOption) (st *channel.Channel, err error) {
	return s.GuildChannelCreateComplex(guildID, GuildChannelCreateData{
		Name: name,
		Type: ctype,
	}, options...)
}

// GuildChannelsReorder updates the order of channel.Channel in a guild.Guild.
func (s *Session) GuildChannelsReorder(guildID string, channels []*channel.Channel, options ...discord.RequestOption) error {
	data := make([]struct {
		ID       string `json:"id"`
		Position int    `json:"position"`
	}, len(channels))

	for i, c := range channels {
		data[i].ID = c.ID
		data[i].Position = c.Position
	}

	_, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildChannels(guildID),
		data,
		discord.EndpointGuildChannels(guildID),
		options...,
	)
	return err
}

// GuildInvites returns an array of invite.Invite for the given guild.Guild.
func (s *Session) GuildInvites(guildID string, options ...discord.RequestOption) ([]*invite.Invite, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildInvites(guildID),
		nil,
		discord.EndpointGuildInvites(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st []*invite.Invite
	return st, unmarshal(body, st)
}

// GuildRoleCreate creates a new guild.Role.
func (s *Session) GuildRoleCreate(guildID string, data *guild.RoleParams, options ...discord.RequestOption) (*guild.Role, error) {
	body, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointGuildRoles(guildID),
		data,
		discord.EndpointGuildRoles(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st guild.Role
	return &st, unmarshal(body, &st)
}

// GuildRoles returns all guild.Role for a given guild.Guild.
func (s *Session) GuildRoles(guildID string, options ...discord.RequestOption) ([]*guild.Role, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildRoles(guildID),
		nil,
		discord.EndpointGuildRoles(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st []*guild.Role
	return st, unmarshal(body, st)
}

// GuildRoleEdit updates an existing guild.Role and returns updated data.
func (s *Session) GuildRoleEdit(guildID, roleID string, data *guild.RoleParams, options ...discord.RequestOption) (*guild.Role, error) {
	// Prevent sending a color int that is too big.
	if data.Color != nil && *data.Color > 0xFFFFFF {
		return nil, fmt.Errorf("color value cannot be larger than 0xFFFFFF")
	}

	body, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildRole(guildID, roleID),
		data,
		discord.EndpointGuildRole(guildID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st guild.Role
	return &st, unmarshal(body, &st)
}

// GuildRoleReorder reoders guild.Role.
func (s *Session) GuildRoleReorder(guildID string, roles []*guild.Role, options ...discord.RequestOption) ([]*guild.Role, error) {
	body, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildRoles(guildID),
		roles,
		discord.EndpointGuildRoles(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st []*guild.Role
	return st, unmarshal(body, st)
}

// GuildRoleDelete deletes a guild.Role.
func (s *Session) GuildRoleDelete(guildID, roleID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildRole(guildID, roleID),
		nil,
		discord.EndpointGuildRole(guildID, ""),
		options...,
	)
	return err
}

// GuildPruneCount returns the number of user.Member that would be removed in a prune operation.
//
// Requires discord.PermissionKickMembers.
func (s *Session) GuildPruneCount(guildID string, days uint32, options ...discord.RequestOption) (uint32, error) {
	if days <= 0 {
		return 0, ErrPruneDaysBounds
	}

	p := struct {
		Pruned uint32 `json:"pruned"`
	}{}

	uri := discord.EndpointGuildPrune(guildID) + "?days=" + strconv.FormatUint(uint64(days), 10)
	body, err := s.RequestWithBucketID(http.MethodGet, uri, nil, discord.EndpointGuildPrune(guildID), options...)
	if err != nil {
		return 0, err
	}

	err = unmarshal(body, &p)
	if err != nil {
		return 0, err
	}

	return p.Pruned, err
}

// GuildPrune begins as prune operation.
// Returns the number of pruned members.
//
// Requires discord.PermissionKickMembers.
func (s *Session) GuildPrune(guildID string, days uint32, options ...discord.RequestOption) (uint32, error) {
	if days <= 0 {
		return 0, ErrPruneDaysBounds
	}

	data := struct {
		days uint32
	}{days}

	p := struct {
		Pruned uint32 `json:"pruned"`
	}{}

	body, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointGuildPrune(guildID),
		data,
		discord.EndpointGuildPrune(guildID),
		options...,
	)
	if err != nil {
		return 0, err
	}

	err = unmarshal(body, &p)
	if err != nil {
		return 0, err
	}

	return p.Pruned, nil
}

// GuildIntegrations returns user.Integration for a guild.Guild.
func (s *Session) GuildIntegrations(guildID string, options ...discord.RequestOption) ([]*user.Integration, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildIntegrations(guildID),
		nil,
		discord.EndpointGuildIntegrations(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st []*user.Integration
	return st, unmarshal(body, st)
}

// GuildIntegrationCreate creates a guild.Guild user.Integration.
func (s *Session) GuildIntegrationCreate(guildID, integrationType, integrationID string, options ...discord.RequestOption) error {
	data := struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}{integrationType, integrationID}

	_, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointGuildIntegrations(guildID),
		data,
		discord.EndpointGuildIntegrations(guildID),
		options...,
	)
	return err
}

// GuildIntegrationEdit edits a guild.Guild user.Integration.
//
// expireBehavior is the behavior when a user.Integration subscription lapses.
// expireGracePeriod is the period (in seconds) where the user.Integration will ignore lapsed subscriptions.
// enableEmoticons is true if emoticons should be synced for this user.Integration (twitch only currently).
func (s *Session) GuildIntegrationEdit(guildID, integrationID string, expireBehavior, expireGracePeriod int, enableEmoticons bool, options ...discord.RequestOption) error {
	data := struct {
		ExpireBehavior    int  `json:"expire_behavior"`
		ExpireGracePeriod int  `json:"expire_grace_period"`
		EnableEmoticons   bool `json:"enable_emoticons"`
	}{expireBehavior, expireGracePeriod, enableEmoticons}

	_, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildIntegration(guildID, integrationID),
		data,
		discord.EndpointGuildIntegration(guildID, ""),
		options...,
	)
	return err
}

// GuildIntegrationDelete removes the user.Integration from the guild.Guild.
func (s *Session) GuildIntegrationDelete(guildID, integrationID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildIntegration(guildID, integrationID),
		nil,
		discord.EndpointGuildIntegration(guildID, ""),
		options...,
	)
	return
}

// GuildIcon returns an image.Image of a guild.Guild icon.
func (s *Session) GuildIcon(guildID string, options ...discord.RequestOption) (image.Image, error) {
	g, err := s.Guild(guildID, options...)
	if err != nil {
		return nil, err
	}

	if g.Icon == "" {
		return nil, ErrGuildNoIcon
	}

	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildIcon(guildID, g.Icon),
		nil,
		discord.EndpointGuildIcon(guildID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	return img, err
}

// GuildSplash returns an image.Image of a guild.Guild splash image.
func (s *Session) GuildSplash(guildID string, options ...discord.RequestOption) (image.Image, error) {
	g, err := s.Guild(guildID, options...)
	if err != nil {
		return nil, err
	}

	if g.Splash == "" {
		return nil, err
	}

	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildSplash(guildID, g.Splash),
		nil,
		discord.EndpointGuildSplash(guildID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	return img, err
}

// GuildEmbed returns the guild.Embed for a guild.Guild.
func (s *Session) GuildEmbed(guildID string, options ...discord.RequestOption) (*guild.Embed, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildEmbed(guildID),
		nil,
		discord.EndpointGuildEmbed(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var em guild.Embed
	return &em, unmarshal(body, &em)
}

// GuildEmbedEdit edits the guild.Embed of a guild.Guild.
func (s *Session) GuildEmbedEdit(guildID string, data *guild.Embed, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildEmbed(guildID),
		data,
		discord.EndpointGuildEmbed(guildID),
		options...,
	)
	return err
}

// GuildAuditLog returns the guild.AuditLog for a guild.Guild.
//
// If provided, all returned guild.AuditLog will be filtered for the given userID.
// If provided all guild.AuditLog entries returned will be before the given beforeID.
// If provided the guild.AuditLog will be filtered for the given actionType.
// limit is the number of messages that can be returned (default 50, min 1, max 100).
func (s *Session) GuildAuditLog(guildID, userID, beforeID string, actionType, limit int, options ...discord.RequestOption) (*guild.AuditLog, error) {
	uri := discord.EndpointGuildAuditLogs(guildID)

	v := url.Values{}
	if userID != "" {
		v.Set("user_id", userID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if actionType > 0 {
		v.Set("action_type", strconv.Itoa(actionType))
	}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if len(v) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, v.Encode())
	}

	body, err := s.RequestWithBucketID(http.MethodGet, uri, nil, discord.EndpointGuildAuditLogs(guildID), options...)
	if err != nil {
		return nil, err
	}

	var al guild.AuditLog
	return &al, unmarshal(body, &al)
}

// GuildEmojis returns all emoji.Emoji.
func (s *Session) GuildEmojis(guildID string, options ...discord.RequestOption) ([]*emoji.Emoji, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildEmojis(guildID),
		nil,
		discord.EndpointGuildEmojis(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var em []*emoji.Emoji
	return em, unmarshal(body, em)
}

// GuildEmoji returns the emoji.Emoji in the given guild.Guild.
func (s *Session) GuildEmoji(guildID, emojiID string, options ...discord.RequestOption) (*emoji.Emoji, error) {
	var body []byte
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildEmoji(guildID, emojiID),
		nil,
		discord.EndpointGuildEmoji(guildID, emojiID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, unmarshal(body, &em)
}

// GuildEmojiCreate creates a new emoji.Emoji in the given guild.Guild.
func (s *Session) GuildEmojiCreate(guildID string, data *emoji.Params, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointGuildEmojis(guildID),
		data,
		discord.EndpointGuildEmojis(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, unmarshal(body, &em)
}

// GuildEmojiEdit modifies and returns updated emoji.Emoji in the given guild.Guild.
func (s *Session) GuildEmojiEdit(guildID, emojiID string, data *emoji.Params, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildEmoji(guildID, emojiID),
		data,
		discord.EndpointGuildEmojis(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, unmarshal(body, &em)
}

// GuildEmojiDelete deletes an emoji.Emoji in the given guild.Guild.
func (s *Session) GuildEmojiDelete(guildID, emojiID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildEmoji(guildID, emojiID),
		nil,
		discord.EndpointGuildEmojis(guildID),
		options...,
	)
	return err
}

// ApplicationEmojis returns all emoji.Emoji for the given application.Application
func (s *Session) ApplicationEmojis(appID string, options ...discord.RequestOption) (emojis []*emoji.Emoji, err error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointApplicationEmojis(appID),
		nil,
		discord.EndpointApplicationEmojis(appID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var data struct {
		Items []*emoji.Emoji `json:"items"`
	}

	emojis = data.Items
	return data.Items, unmarshal(body, &data)
}

// ApplicationEmoji returns the emoji.Emoji for the given application.Application.
func (s *Session) ApplicationEmoji(appID, emojiID string, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointApplicationEmoji(appID, emojiID),
		nil,
		discord.EndpointApplicationEmoji(appID, emojiID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, unmarshal(body, &em)
}

// ApplicationEmojiCreate creates a new emoji.Emoji for the given application.Application.
func (s *Session) ApplicationEmojiCreate(appID string, data *emoji.Params, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointApplicationEmojis(appID),
		data,
		discord.EndpointApplicationEmojis(appID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, unmarshal(body, &em)
}

// ApplicationEmojiEdit modifies and returns updated emoji.Emoji for the given application.Application.
func (s *Session) ApplicationEmojiEdit(appID string, emojiID string, data *emoji.Params, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointApplicationEmoji(appID, emojiID),
		data,
		discord.EndpointApplicationEmojis(appID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, unmarshal(body, &em)
}

// ApplicationEmojiDelete deletes an emoji.Emoji for the given application.Application.
func (s *Session) ApplicationEmojiDelete(appID, emojiID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointApplicationEmoji(appID, emojiID),
		nil,
		discord.EndpointApplicationEmojis(appID),
		options...,
	)
	return err
}

// GuildTemplate returns a guild.Template for the given code.
func (s *Session) GuildTemplate(code string, options ...discord.RequestOption) (*guild.Template, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildTemplate(code),
		nil,
		discord.EndpointGuildTemplate(code),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var t guild.Template
	return &t, unmarshal(body, &t)
}

// GuildCreateWithTemplate creates a guild.Guild based on a guild.Template.
//
// code is the Code of the guild.Template.
// name is the name of the guild.Guild (2-100 characters).
// icon is the base64 encoded 128x128 image for the guild.Guild icon.
func (s *Session) GuildCreateWithTemplate(templateCode, name, icon string, options ...discord.RequestOption) (*guild.Guild, error) {
	data := struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	}{name, icon}

	body, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointGuildTemplate(templateCode),
		data,
		discord.EndpointGuildTemplate(templateCode),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var g guild.Guild
	return &g, unmarshal(body, &g)
}

// GuildTemplates returns every guild.Template of the given guild.Guild.
func (s *Session) GuildTemplates(guildID string, options ...discord.RequestOption) ([]*guild.Template, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildTemplates(guildID),
		nil,
		discord.EndpointGuildTemplates(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var t []*guild.Template
	return t, unmarshal(body, &t)
}

// GuildTemplateCreate creates a guild.Template for the guild.Guild.
func (s *Session) GuildTemplateCreate(guildID string, data *guild.TemplateParams, options ...discord.RequestOption) (*guild.Template, error) {
	body, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointGuildTemplates(guildID),
		data,
		discord.EndpointGuildTemplates(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var t guild.Template
	return &t, unmarshal(body, &t)
}

// GuildTemplateSync syncs the guild.Template to the guild.Guild's current state
//
// code is the code of the guild.Template.
func (s *Session) GuildTemplateSync(guildID, code string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodPut,
		discord.EndpointGuildTemplateSync(guildID, code),
		nil,
		discord.EndpointGuildTemplateSync(guildID, ""),
		options...,
	)
	return err
}

// GuildTemplateEdit modifies the guild.Template's metadata of the given guild.Guild.
//
// code is the code of the guild.Template.
func (s *Session) GuildTemplateEdit(guildID, code string, data *guild.TemplateParams, options ...discord.RequestOption) (*guild.Template, error) {
	body, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildTemplateSync(guildID, code),
		data,
		discord.EndpointGuildTemplateSync(guildID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var t guild.Template
	return &t, unmarshal(body, &t)
}

// GuildTemplateDelete deletes the guild.Template of the given guild.Guild.
//
// code is the code of the guild.Template.
func (s *Session) GuildTemplateDelete(guildID, templateCode string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildTemplateSync(guildID, templateCode),
		nil,
		discord.EndpointGuildTemplateSync(guildID, ""),
		options...,
	)
	return err
}

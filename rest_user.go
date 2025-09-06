package gokord

import (
	"bytes"
	"image"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
)

// User returns the user.User details of the given userID (can be @me to be the current user.User ID).
func (s *Session) User(userID string, options ...RequestOption) (*user.User, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointUser(userID),
		nil,
		discord.EndpointUsers,
		options...,
	)
	if err != nil {
		return nil, err
	}

	var u user.User
	return &u, unmarshal(body, &u)
}

// UserAvatarDecode returns an image.Image of a user.User's Avatar.
func (s *Session) UserAvatarDecode(u *user.User, options ...RequestOption) (image.Image, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointUserAvatar(u.ID, u.Avatar),
		nil,
		discord.EndpointUserAvatar("", ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	return img, err
}

// UserUpdate updates current user settings.
//
// Note: Avatar must be either the hash/id of existing Avatar or
// data:image/png;base64,BASE64_STRING_OF_NEW_AVATAR_PNG to set a new avatar.
// If left blank, avatar will be set to null/blank.
func (s *Session) UserUpdate(username, avatar, banner string, options ...RequestOption) (*user.User, error) {
	data := struct {
		Username string `json:"username,omitempty"`
		Avatar   string `json:"avatar,omitempty"`
		Banner   string `json:"banner,omitempty"`
	}{username, avatar, banner}

	body, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointUser("@me"),
		data,
		discord.EndpointUsers,
		options...,
	)
	if err != nil {
		return nil, err
	}

	var u user.User
	return &u, unmarshal(body, &u)
}

// UserConnections returns the current user.Connection.
func (s *Session) UserConnections(options ...RequestOption) ([]*user.Connection, error) {
	response, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointUserConnections("@me"),
		nil,
		discord.EndpointUserConnections("@me"),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var conn []*user.Connection
	return conn, unmarshal(response, &conn)
}

// UserChannelCreate creates a new private channel.Channel (types.ChannelDM) with another user.User
func (s *Session) UserChannelCreate(userID string, options ...RequestOption) (*channel.Channel, error) {
	data := struct {
		RecipientID string `json:"recipient_id"`
	}{userID}

	body, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointUserChannels("@me"),
		data,
		discord.EndpointUserChannels(""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var c channel.Channel
	return &c, unmarshal(body, &c)
}

// UserGuildMember returns a user.Member for the current user.User in the given guild.Guild ID.
func (s *Session) UserGuildMember(guildID string, options ...RequestOption) (*user.Member, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointUserGuildMember("@me", guildID),
		nil,
		discord.EndpointUserGuildMember("@me", guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var m user.Member
	return &m, unmarshal(body, &m)
}

// UserGuilds returns an array of guild.UserGuild structures for all guilds.
//
// limit is the number of guilds that can be returned (max 200).
// If beforeID is set, it will return all guilds before this ID.
// If afterID is set, it will return all guilds after this ID.
// Set withCounts to true if you want to include approximate member and presence counts.
func (s *Session) UserGuilds(limit int, beforeID, afterID string, withCounts bool, options ...RequestOption) ([]*guild.UserGuild, error) {
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
	return ug, unmarshal(body, &ug)
}

// memberPermissions calculates the permissions for a user.Member.
// https://support.discord.com/hc/en-us/articles/206141927-How-is-the-permission-hierarchy-structured-
func memberPermissions(guild *guild.Guild, channel *channel.Channel, userID string, roles []string) int64 {
	if userID == guild.OwnerID {
		return discord.PermissionAll
	}

	var perms int64
	for _, role := range guild.Roles {
		if role.ID == guild.ID {
			perms |= role.Permissions
			break
		}
	}

	for _, role := range guild.Roles {
		for _, roleID := range roles {
			if role.ID == roleID {
				perms |= role.Permissions
				break
			}
		}
	}

	if perms&discord.PermissionAdministrator == discord.PermissionAdministrator {
		perms |= discord.PermissionAll
	}

	// Apply @everyone overrides from the channel.
	for _, overwrite := range channel.PermissionOverwrites {
		if guild.ID == overwrite.ID {
			perms &= ^overwrite.Deny
			perms |= overwrite.Allow
			break
		}
	}

	var denies, allows int64
	// Member overwrites can override role overrides, so do two passes
	for _, overwrite := range channel.PermissionOverwrites {
		for _, roleID := range roles {
			if overwrite.Type == types.PermissionOverwriteRole && roleID == overwrite.ID {
				denies |= overwrite.Deny
				allows |= overwrite.Allow
				break
			}
		}
	}

	perms &= ^denies
	perms |= allows

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.Type == types.PermissionOverwriteMember && overwrite.ID == userID {
			perms &= ^overwrite.Deny
			perms |= overwrite.Allow
			break
		}
	}

	if perms&discord.PermissionAdministrator == discord.PermissionAdministrator {
		perms |= discord.PermissionAllChannel
	}

	return perms
}

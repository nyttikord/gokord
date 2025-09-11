package userapi

import (
	"bytes"
	"image"
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
)

type Requester struct {
	discord.Requester
}

// User returns the user.User details of the given userID (can be @me to be the current user.User ID).
func (s Requester) User(userID string, options ...discord.RequestOption) (*user.User, error) {
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
	return &u, s.Unmarshal(body, &u)
}

// UserAvatarDecode returns an image.Image of a user.User's Avatar.
func (s Requester) UserAvatarDecode(u *user.User, options ...discord.RequestOption) (image.Image, error) {
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
func (s Requester) UserUpdate(username, avatar, banner string, options ...discord.RequestOption) (*user.User, error) {
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
	return &u, s.Unmarshal(body, &u)
}

// UserConnections returns the current user.Connection.
func (s Requester) UserConnections(options ...discord.RequestOption) ([]*user.Connection, error) {
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
	return conn, s.Unmarshal(response, &conn)
}

// UserChannelCreate creates a new private channel.Channel (types.ChannelDM) with another user.User
func (s Requester) UserChannelCreate(userID string, options ...discord.RequestOption) (*channel.Channel, error) {
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
	return &c, s.Unmarshal(body, &c)
}

// UserGuildMember returns a user.Member for the current user.User in the given guild.Guild ID.
func (s Requester) UserGuildMember(guildID string, options ...discord.RequestOption) (*user.Member, error) {
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
	return &m, s.Unmarshal(body, &m)
}

// MemberPermissions calculates the permissions for a user.Member.
// https://support.discord.com/hc/en-us/articles/206141927-How-is-the-permission-hierarchy-structured-
func MemberPermissions(guild *guild.Guild, channel *channel.Channel, userID string, roles []string) int64 {
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

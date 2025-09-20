// Package userapi contains everything to interact with everything located in the user package.
package userapi

import (
	"bytes"
	"image"
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/user"
)

// Requester handles everything inside the user package.
type Requester struct {
	discord.Requester
	State *State
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

// AvatarDecode returns an image.Image of a user.User Avatar.
func (s Requester) AvatarDecode(u *user.User, options ...discord.RequestOption) (image.Image, error) {
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

// Update updates current user.User settings.
//
// NOTE: Avatar must be either the hash/id of existing Avatar or
// data:image/png;base64,BASE64_STRING_OF_NEW_AVATAR_PNG to set a new avatar.
// If left blank, avatar will be set to null/blank.
func (s Requester) Update(username, avatar, banner string, options ...discord.RequestOption) (*user.User, error) {
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

// Connections returns the current user.Connection.
func (s Requester) Connections(options ...discord.RequestOption) ([]*user.Connection, error) {
	response, err := s.Request(http.MethodGet, discord.EndpointUserConnections("@me"), nil, options...)
	if err != nil {
		return nil, err
	}

	var conn []*user.Connection
	return conn, s.Unmarshal(response, &conn)
}

// ChannelCreate creates a new private channel.Channel (types.ChannelDM) with another user.User.
func (s Requester) ChannelCreate(userID string, options ...discord.RequestOption) (*channel.Channel, error) {
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

// GuildMember returns a user.Member for the current user.User in the given guild.Guild ID.
func (s Requester) GuildMember(guildID string, options ...discord.RequestOption) (*user.Member, error) {
	body, err := s.Request(http.MethodGet, discord.EndpointUserGuildMember("@me", guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var m user.Member
	return &m, s.Unmarshal(body, &m)
}

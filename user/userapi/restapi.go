// Package userapi contains everything to interact with everything located in the user package.
package userapi

import (
	"bytes"
	"image"
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/user"
)

// Requester handles everything inside the user package.
type Requester struct {
	RESTRequester
	State *State
}

// User returns the user.User details of the given userID (can be @me to be the current User ID).
func (r Requester) User(userID string) Request[*User] {
	return NewSimpleData[*User](
		r, http.MethodGet, discord.EndpointUser(userID),
	).WithBucketID(discord.EndpointUsers)
}

// AvatarDecode returns an image.Image of a user.User Avatar.
func (r Requester) AvatarDecode(u *User) Request[image.Image] {
	body, err := r.RequestWithBucketID(
		ctx,
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
func (r Requester) Update(username, avatar, banner string) Request[*User] {
	data := struct {
		Username string `json:"username,omitempty"`
		Avatar   string `json:"avatar,omitempty"`
		Banner   string `json:"banner,omitempty"`
	}{username, avatar, banner}

	return NewSimpleData[*User](
		r, http.MethodPatch, discord.EndpointUser("@me"),
	).WithBucketID(discord.EndpointUsers).WithData(data)
}

// Connections returns the current user.Connection.
func (r Requester) Connections() Request[[]*Connection] {
	return NewSimpleData[[]*Connection](
		r, http.MethodGet, discord.EndpointUserConnections("@me"),
	)
}

// ChannelCreate creates a new private channel.Channel (types.ChannelDM) with another user.User.
func (r Requester) ChannelCreate(userID string) Request[*channel.Channel] {
	data := struct {
		RecipientID string `json:"recipient_id"`
	}{userID}

	return NewSimpleData[*channel.Channel](
		r, http.MethodPost, discord.EndpointUserChannels("@me"),
	).WithBucketID(discord.EndpointUserChannels("")).WithData(data)
}

// GuildMember returns a user.Member for the current user.User in the given guild.Guild ID.
func (r Requester) GuildMember(guildID string) Request[*Member] {
	return NewSimpleData[*Member](
		r, http.MethodGet, discord.EndpointUserGuildMember("@me", guildID),
	).WithBucketID(discord.EndpointUserGuildMember("", guildID))
}

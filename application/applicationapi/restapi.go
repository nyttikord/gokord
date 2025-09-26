package applicationapi

import (
	"net/http"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/emoji"
)

// Emojis returns all emoji.Emoji for the given application.Application
func (r Requester) Emojis(appID string, options ...discord.RequestOption) (emojis []*emoji.Emoji, err error) {
	body, err := r.Request(http.MethodGet, discord.EndpointApplicationEmojis(appID), nil, options...)
	if err != nil {
		return nil, err
	}

	var data struct {
		Items []*emoji.Emoji `json:"items"`
	}

	emojis = data.Items
	return data.Items, r.Unmarshal(body, &data)
}

// Emoji returns the emoji.Emoji for the given application.Application.
func (r Requester) Emoji(appID, emojiID string, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointApplicationEmoji(appID, emojiID), nil, options...)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, r.Unmarshal(body, &em)
}

// EmojiCreate creates a new emoji.Emoji for the given application.Application.
func (r Requester) EmojiCreate(appID string, data *emoji.Params, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := r.Request(http.MethodPost, discord.EndpointApplicationEmojis(appID), data, options...)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, r.Unmarshal(body, &em)
}

// EmojiEdit modifies and returns updated emoji.Emoji for the given application.Application.
func (r Requester) EmojiEdit(appID string, emojiID string, data *emoji.Params, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := r.RequestWithBucketID(
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
	return &em, r.Unmarshal(body, &em)
}

// EmojiDelete deletes an emoji.Emoji for the given application.Application.
func (r Requester) EmojiDelete(appID, emojiID string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointApplicationEmoji(appID, emojiID),
		nil,
		discord.EndpointApplicationEmojis(appID),
		options...,
	)
	return err
}

// RoleConnectionMetadata returns application.RoleConnectionMetadata.
func (r Requester) RoleConnectionMetadata(appID string, options ...discord.RequestOption) ([]*application.RoleConnectionMetadata, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointApplicationRoleConnectionMetadata(appID), nil, options...)
	if err != nil {
		return nil, err
	}

	var m []*application.RoleConnectionMetadata
	return m, r.Unmarshal(body, &m)
}

// RoleConnectionMetadataUpdate updates and returns application.RoleConnectionMetadata.
func (r Requester) RoleConnectionMetadataUpdate(appID string, metadata []*application.RoleConnectionMetadata, options ...discord.RequestOption) ([]*application.RoleConnectionMetadata, error) {
	body, err := r.Request(http.MethodPut, discord.EndpointApplicationRoleConnectionMetadata(appID), metadata, options...)
	if err != nil {
		return nil, err
	}

	var m []*application.RoleConnectionMetadata
	return m, r.Unmarshal(body, &m)
}

// RoleConnection returns application.RoleConnection to the specified application.Application.
func (r Requester) RoleConnection(appID string, options ...discord.RequestOption) (*application.RoleConnection, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointUserApplicationRoleConnection(appID), nil, options...)
	if err != nil {
		return nil, err
	}

	var c application.RoleConnection
	return &c, r.Unmarshal(body, &c)
}

// RoleConnectionUpdate updates and returns application.RoleConnection to the specified application.Application.
func (r Requester) RoleConnectionUpdate(appID string, rconn *application.RoleConnection, options ...discord.RequestOption) (*application.RoleConnection, error) {
	body, err := r.Request(http.MethodPut, discord.EndpointUserApplicationRoleConnection(appID), rconn, options...)
	if err != nil {
		return nil, err
	}

	var c application.RoleConnection
	return &c, r.Unmarshal(body, &c)
}

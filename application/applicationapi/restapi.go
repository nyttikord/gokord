package applicationapi

import (
	"context"
	"net/http"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/emoji"
)

// Emojis returns all emoji.Emoji for the given application.Application
func (r Requester) Emojis(ctx context.Context, appID string, options ...discord.RequestOption) (emojis []*emoji.Emoji, err error) {
	body, err := r.Request(ctx, http.MethodGet, discord.EndpointApplicationEmojis(appID), nil, options...)
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
func (r Requester) Emoji(appID, emojiID string) request.Request[*emoji.Emoji] {
	return request.NewSimpleData[*emoji.Emoji](
		r, http.MethodGet, discord.EndpointApplicationEmoji(appID, emojiID),
	).WithBucketID(discord.EndpointApplicationEmojis(appID))
}

// EmojiCreate creates a new emoji.Emoji for the given application.Application.
func (r Requester) EmojiCreate(appID string, data *emoji.Params) request.Request[*emoji.Emoji] {
	return request.NewSimpleData[*emoji.Emoji](
		r, http.MethodPost, discord.EndpointApplicationEmojis(appID),
	).WithData(data)
}

// EmojiEdit modifies and returns updated emoji.Emoji for the given application.Application.
func (r Requester) EmojiEdit(appID string, emojiID string, data *emoji.Params) request.Request[*emoji.Emoji] {
	return request.NewSimpleData[*emoji.Emoji](
		r, http.MethodPatch, discord.EndpointApplicationEmoji(appID, emojiID),
	).WithData(data).WithBucketID(discord.EndpointApplicationEmojis(appID))
}

// EmojiDelete deletes an emoji.Emoji for the given application.Application.
func (r Requester) EmojiDelete(appID, emojiID string) request.EmptyRequest {
	req := request.NewSimple(
		r, http.MethodDelete, discord.EndpointApplicationEmoji(appID, emojiID),
	).WithBucketID(discord.EndpointApplicationEmojis(appID))
	return request.WrapAsEmpty(req)
}

// RoleConnectionMetadata returns application.RoleConnectionMetadata.
func (r Requester) RoleConnectionMetadata(appID string) request.Request[[]*application.RoleConnectionMetadata] {
	return request.NewSimpleData[[]*application.RoleConnectionMetadata](
		r, http.MethodGet, discord.EndpointApplicationRoleConnectionMetadata(appID),
	)
}

// RoleConnectionMetadataUpdate updates and returns application.RoleConnectionMetadata.
func (r Requester) RoleConnectionMetadataUpdate(appID string, metadata []*application.RoleConnectionMetadata) request.Request[[]*application.RoleConnectionMetadata] {
	return request.NewSimpleData[[]*application.RoleConnectionMetadata](
		r, http.MethodPut, discord.EndpointApplicationRoleConnectionMetadata(appID),
	).WithData(metadata)
}

// RoleConnection returns application.RoleConnection to the specified application.Application.
func (r Requester) RoleConnection(appID string) request.Request[*application.RoleConnection] {
	return request.NewSimpleData[*application.RoleConnection](
		r, http.MethodGet, discord.EndpointUserApplicationRoleConnection(appID),
	)
}

// RoleConnectionUpdate updates and returns application.RoleConnection to the specified application.Application.
func (r Requester) RoleConnectionUpdate(appID string, rconn *application.RoleConnection) request.Request[*application.RoleConnection] {
	return request.NewSimpleData[*application.RoleConnection](
		r, http.MethodPut, discord.EndpointUserApplicationRoleConnection(appID),
	).WithData(rconn)
}

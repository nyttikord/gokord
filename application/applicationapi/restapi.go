package applicationapi

import (
	"context"
	"net/http"

	. "github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/emoji"
)

// Emojis returns all emoji.Emoji for the given application.Application
func (r Requester) Emojis(appID string) Request[[]*emoji.Emoji] {
	return NewCustom[[]*emoji.Emoji](r, http.MethodGet, discord.EndpointApplicationEmojis(appID)).
		WithPost(func(ctx context.Context, b []byte) ([]*emoji.Emoji, error) {
			var data struct {
				Items []*emoji.Emoji `json:"items"`
			}
			return data.Items, r.Unmarshal(b, &data)
		})
}

// Emoji returns the emoji.Emoji for the given application.Application.
func (r Requester) Emoji(appID, emojiID string) Request[*emoji.Emoji] {
	return NewData[*emoji.Emoji](
		r, http.MethodGet, discord.EndpointApplicationEmoji(appID, emojiID),
	).WithBucketID(discord.EndpointApplicationEmojis(appID))
}

// EmojiCreate creates a new emoji.Emoji for the given application.Application.
func (r Requester) EmojiCreate(appID string, data *emoji.Params) Request[*emoji.Emoji] {
	return NewData[*emoji.Emoji](
		r, http.MethodPost, discord.EndpointApplicationEmojis(appID),
	).WithData(data)
}

// EmojiEdit modifies and returns updated emoji.Emoji for the given application.Application.
func (r Requester) EmojiEdit(appID string, emojiID string, data *emoji.Params) Request[*emoji.Emoji] {
	return NewData[*emoji.Emoji](
		r, http.MethodPatch, discord.EndpointApplicationEmoji(appID, emojiID),
	).WithData(data).WithBucketID(discord.EndpointApplicationEmojis(appID))
}

// EmojiDelete deletes an emoji.Emoji for the given application.Application.
func (r Requester) EmojiDelete(appID, emojiID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointApplicationEmoji(appID, emojiID),
	).WithBucketID(discord.EndpointApplicationEmojis(appID))
	return WrapAsEmpty(req)
}

// RoleConnectionMetadata returns application.RoleConnectionMetadata.
func (r Requester) RoleConnectionMetadata(appID string) Request[[]*RoleConnectionMetadata] {
	return NewData[[]*RoleConnectionMetadata](
		r, http.MethodGet, discord.EndpointApplicationRoleConnectionMetadata(appID),
	)
}

// RoleConnectionMetadataUpdate updates and returns application.RoleConnectionMetadata.
func (r Requester) RoleConnectionMetadataUpdate(appID string, metadata []*RoleConnectionMetadata) Request[[]*RoleConnectionMetadata] {
	return NewData[[]*RoleConnectionMetadata](
		r, http.MethodPut, discord.EndpointApplicationRoleConnectionMetadata(appID),
	).WithData(metadata)
}

// RoleConnection returns RoleConnection to the specified application.Application.
func (r Requester) RoleConnection(appID string) Request[*RoleConnection] {
	return NewData[*RoleConnection](
		r, http.MethodGet, discord.EndpointUserApplicationRoleConnection(appID),
	)
}

// RoleConnectionUpdate updates and returns RoleConnection to the specified application.Application.
func (r Requester) RoleConnectionUpdate(appID string, rconn *RoleConnection) Request[*RoleConnection] {
	return NewData[*RoleConnection](
		r, http.MethodPut, discord.EndpointUserApplicationRoleConnection(appID),
	).WithData(rconn)
}

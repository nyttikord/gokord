package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/emoji"
)

// Emojis returns all emoji.Emoji.
func (r Requester) Emojis(guildID string) Request[[]*emoji.Emoji] {
	return NewSimpleData[[]*emoji.Emoji](
		r, http.MethodGet, discord.EndpointGuildEmojis(guildID),
	)
}

// Emoji returns the emoji.Emoji in the given guild.Guild.
func (r Requester) Emoji(guildID, emojiID string) Request[*emoji.Emoji] {
	return NewSimpleData[*emoji.Emoji](
		r, http.MethodGet, discord.EndpointGuildEmoji(guildID, emojiID),
	).WithBucketID(discord.EndpointGuildEmojis(guildID))
}

// EmojiCreate creates a new emoji.Emoji in the given guild.Guild.
func (r Requester) EmojiCreate(guildID string, data *emoji.Params) Request[*emoji.Emoji] {
	return NewSimpleData[*emoji.Emoji](
		r, http.MethodPost, discord.EndpointGuildEmojis(guildID),
	).WithData(data)
}

// EmojiEdit modifies and returns updated emoji.Emoji in the given guild.Guild.
func (r Requester) EmojiEdit(guildID, emojiID string, data *emoji.Params) Request[*emoji.Emoji] {
	return NewSimpleData[*emoji.Emoji](
		r, http.MethodPatch, discord.EndpointGuildEmoji(guildID, emojiID),
	).WithBucketID(discord.EndpointGuildEmojis(guildID)).WithData(data)
}

// EmojiDelete deletes an emoji.Emoji in the given guild.Guild.
func (r Requester) EmojiDelete(guildID, emojiID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointGuildEmoji(guildID, emojiID),
	).WithBucketID(discord.EndpointGuildEmojis(guildID))
	return WrapAsEmpty(req)
}

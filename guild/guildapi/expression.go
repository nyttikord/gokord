package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/emoji"
)

// Emojis returns all emoji.Emoji.
func (r Requester) Emojis(guildID string, options ...discord.RequestOption) ([]*emoji.Emoji, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuildEmojis(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var em []*emoji.Emoji
	return em, r.Unmarshal(body, &em)
}

// Emoji returns the emoji.Emoji in the given guild.Guild.
func (r Requester) Emoji(guildID, emojiID string, options ...discord.RequestOption) (*emoji.Emoji, error) {
	var body []byte
	body, err := r.Request(http.MethodGet, discord.EndpointGuildEmoji(guildID, emojiID), nil, options...)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, r.Unmarshal(body, &em)
}

// EmojiCreate creates a new emoji.Emoji in the given guild.Guild.
func (r Requester) EmojiCreate(guildID string, data *emoji.Params, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := r.Request(http.MethodPost, discord.EndpointGuildEmojis(guildID), data, options...)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, r.Unmarshal(body, &em)
}

// EmojiEdit modifies and returns updated emoji.Emoji in the given guild.Guild.
func (r Requester) EmojiEdit(guildID, emojiID string, data *emoji.Params, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := r.RequestWithBucketID(
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
	return &em, r.Unmarshal(body, &em)
}

// EmojiDelete deletes an emoji.Emoji in the given guild.Guild.
func (r Requester) EmojiDelete(guildID, emojiID string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildEmoji(guildID, emojiID),
		nil,
		discord.EndpointGuildEmojis(guildID),
		options...,
	)
	return err
}

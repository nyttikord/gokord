package applicationapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/emoji"
)

// Emojis returns all emoji.Emoji for the given application.Application
func (s Requester) Emojis(appID string, options ...discord.RequestOption) (emojis []*emoji.Emoji, err error) {
	body, err := s.Request(http.MethodGet, discord.EndpointApplicationEmojis(appID), nil, options...)
	if err != nil {
		return nil, err
	}

	var data struct {
		Items []*emoji.Emoji `json:"items"`
	}

	emojis = data.Items
	return data.Items, s.Unmarshal(body, &data)
}

// Emoji returns the emoji.Emoji for the given application.Application.
func (s Requester) Emoji(appID, emojiID string, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := s.Request(http.MethodGet, discord.EndpointApplicationEmoji(appID, emojiID), nil, options...)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, s.Unmarshal(body, &em)
}

// EmojiCreate creates a new emoji.Emoji for the given application.Application.
func (s Requester) EmojiCreate(appID string, data *emoji.Params, options ...discord.RequestOption) (*emoji.Emoji, error) {
	body, err := s.Request(http.MethodPost, discord.EndpointApplicationEmojis(appID), data, options...)
	if err != nil {
		return nil, err
	}

	var em emoji.Emoji
	return &em, s.Unmarshal(body, &em)
}

// EmojiEdit modifies and returns updated emoji.Emoji for the given application.Application.
func (s Requester) EmojiEdit(appID string, emojiID string, data *emoji.Params, options ...discord.RequestOption) (*emoji.Emoji, error) {
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
	return &em, s.Unmarshal(body, &em)
}

// EmojiDelete deletes an emoji.Emoji for the given application.Application.
func (s Requester) EmojiDelete(appID, emojiID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointApplicationEmoji(appID, emojiID),
		nil,
		discord.EndpointApplicationEmojis(appID),
		options...,
	)
	return err
}

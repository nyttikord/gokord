package channel

import (
	"errors"

	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
)

// MessageEmbedFooter is a part of a [MessageEmbed] struct.
type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedImage is a part of a [MessageEmbed] struct.
type MessageEmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedThumbnail is a part of a [MessageEmbed] struct.
type MessageEmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedVideo is a part of a [MessageEmbed] struct.
type MessageEmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// MessageEmbedProvider is a part of a [MessageEmbed] struct.
type MessageEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

// MessageEmbedAuthor is a part of a [MessageEmbed] struct.
type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedField is a part of a [MessageEmbed] struct.
type MessageEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// An MessageEmbed stores data for [Message] embeds.
type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Type        types.Embed            `json:"type,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []*MessageEmbedField   `json:"fields,omitempty"`
}

var ErrReplyNilMessageRef = errors.New("reply attempted with nil message reference")

// SendEmbed to the given [Channel].
func SendEmbed(channelID string, embed *MessageEmbed) Request[*Message] {
	return SendEmbeds(channelID, []*MessageEmbed{embed})
}

// SendEmbeds sends multiple [MessageEmbed] to the given [Channel].
func SendEmbeds(channelID string, embeds []*MessageEmbed) Request[*Message] {
	return SendMessageComplex(channelID, &MessageSend{Embeds: embeds})
}

// SendEmbedReply sends a reply channel.MessageEmbed to the given channel.Channel.
//
// reference is the message reference to send containing the channel.Message to reply to.
func SendEmbedReply(channelID string, embed *MessageEmbed, reference *MessageReference) Request[*Message] {
	return SendEmbedsReply(channelID, []*MessageEmbed{embed}, reference)
}

// SendEmbedsReply sends a reply with multiple channel.MessageEmbed in the given channel.Channel.
//
// reference is the message reference to send containing the channel.Message to reply to.
func SendEmbedsReply(channelID string, embeds []*MessageEmbed, reference *MessageReference) Request[*Message] {
	if reference == nil {
		return NewError[*Message](ErrReplyNilMessageRef)
	}
	return SendMessageComplex(channelID, &MessageSend{
		Embeds:    embeds,
		Reference: reference,
	})
}

// EditEmbed, replacing it entirely with the given [MessageEmbed].
func EditEmbed(channelID, messageID string, embed *MessageEmbed) Request[*Message] {
	return EditEmbeds(channelID, messageID, []*MessageEmbed{embed})
}

// EditEmbeds, replacing it entirely with multiple [MessageEmbed]s.
func EditEmbeds(channelID, messageID string, embeds []*MessageEmbed) Request[*Message] {
	return EditMessageComplex(NewMessageEdit(channelID, messageID).SetEmbeds(embeds...))
}

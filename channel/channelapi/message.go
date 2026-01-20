package channelapi

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// Messages returns an array of channel.Message the given channel.Channel.
//
// limit is the number messages that can be returned (max 100).
// If provided all messages returned will be before beforeID.
// If provided all messages returned will be after afterID.
// If provided all messages returned will be around aroundID.
func (s Requester) Messages(channelID string, limit int, beforeID, afterID, aroundID string) Request[[]*channel.Message] {
	uri := discord.EndpointChannelMessages(channelID)

	v := url.Values{}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if aroundID != "" {
		v.Set("around", aroundID)
	}
	if len(v) > 0 {
		uri += "?" + v.Encode()
	}
	return NewSimpleData[[]*channel.Message](s, http.MethodGet, uri)
}

// Message gets channel.Message from a given channel.Channel.
func (s Requester) Message(channelID, messageID string) Request[*channel.Message] {
	return NewSimpleData[*channel.Message](
		s, http.MethodGet, discord.EndpointChannelMessage(channelID, messageID),
	).WithBucketID(discord.EndpointChannelMessage(channelID, ""))
}

// MessageSend sends a simple channel.Message to the given channel.Channel.
func (s Requester) MessageSend(channelID string, content string) Request[*channel.Message] {
	return s.MessageSendComplex(channelID, &channel.MessageSend{
		Content: content,
	})
}

// MessageSendComplex sends a channel.Message to the given channel.Channel.
func (s Requester) MessageSendComplex(channelID string, data *channel.MessageSend) Request[*channel.Message] {
	for _, embed := range data.Embeds {
		if embed.Type == "" {
			embed.Type = types.EmbedRich
		}
	}
	endpoint := discord.EndpointChannelMessages(channelID)

	if data.StickerIDs != nil {
		if len(data.StickerIDs) > 3 {
			return NewError[*channel.Message](ErrTooMuchStickers)
		}
	}

	files := data.Files
	if len(files) == 0 {
		return NewSimpleData[*channel.Message](s, http.MethodPost, endpoint).WithData(data)
	}
	contentType, body, encodeErr := channel.MultipartBodyWithJSON(data, files)
	if encodeErr != nil {
		return NewError[*channel.Message](encodeErr)
	}
	response, err := s.RequestRaw(ctx, http.MethodPost, endpoint, contentType, body, endpoint, 0, options...)
}

// MessageSendTTS sends a simple channel.Message to the given channel.Channel with Text to Speech.
func (s Requester) MessageSendTTS(channelID string, content string) Request[*channel.Message] {
	return s.MessageSendComplex(channelID, &channel.MessageSend{
		Content: content,
		TTS:     true,
	})
}

// MessageSendEmbed sends a channel.MessageEmbed to the given channel.Channel.
func (s Requester) MessageSendEmbed(channelID string, embed *channel.MessageEmbed) Request[*channel.Message] {
	return s.MessageSendEmbeds(channelID, []*channel.MessageEmbed{embed})
}

// MessageSendEmbeds sends multiple channel.MessageEmbed to the given channel.Channel.
func (s Requester) MessageSendEmbeds(channelID string, embeds []*channel.MessageEmbed) Request[*channel.Message] {
	return s.MessageSendComplex(channelID, &channel.MessageSend{
		Embeds: embeds,
	})
}

// MessageSendReply sends a reply channel.Message to the given channel.Channel.
//
// reference is the message reference to send containing the channel.Message to reply to.
func (s Requester) MessageSendReply(channelID string, content string, reference *channel.MessageReference) Request[*channel.Message] {
	if reference == nil {
		return NewError[*channel.Message](ErrReplyNilMessageRef)
	}
	return s.MessageSendComplex(channelID, &channel.MessageSend{
		Content:   content,
		Reference: reference,
	})
}

// MessageSendEmbedReply sends a reply channel.MessageEmbed to the given channel.Channel.
//
// reference is the message reference to send containing the channel.Message to reply to.
func (s Requester) MessageSendEmbedReply(channelID string, embed *channel.MessageEmbed, reference *channel.MessageReference) Request[*channel.Message] {
	return s.MessageSendEmbedsReply(channelID, []*channel.MessageEmbed{embed}, reference)
}

// MessageSendEmbedsReply sends a reply with multiple channel.MessageEmbed in the given channel.Channel.
//
// reference is the message reference to send containing the channel.Message to reply to.
func (s Requester) MessageSendEmbedsReply(channelID string, embeds []*channel.MessageEmbed, reference *channel.MessageReference) Request[*channel.Message] {
	if reference == nil {
		return NewError[*channel.Message](ErrReplyNilMessageRef)
	}
	return s.MessageSendComplex(channelID, &channel.MessageSend{
		Embeds:    embeds,
		Reference: reference,
	})
}

// MessageEdit edits an existing channel.Message, replacing it entirely with the given content.
func (s Requester) MessageEdit(channelID, messageID, content string) Request[*channel.Message] {
	return s.MessageEditComplex(channel.NewMessageEdit(channelID, messageID).SetContent(content))
}

// MessageEditComplex edits an existing channel.Message, replacing it entirely with the given channel.MessageEdit.
func (s Requester) MessageEditComplex(m *channel.MessageEdit) Request[*channel.Message] {
	if m.Embeds != nil {
		for _, embed := range *m.Embeds {
			if embed.Type == "" {
				embed.Type = types.EmbedRich
			}
		}
	}

	endpoint := discord.EndpointChannelMessage(m.Channel, m.ID)

	if len(m.Files) == 0 {
		return NewSimpleData[*channel.Message](
			s, http.MethodPatch, endpoint,
		).WithBucketID(discord.EndpointChannelMessage(m.Channel, "")).WithData(m)
	}
	contentType, body, encodeErr := channel.MultipartBodyWithJSON(m, m.Files)
	if encodeErr != nil {
		return NewError[*channel.Message](encodeErr)
	}
	response, err = s.RequestRaw(
		ctx,
		http.MethodPatch,
		endpoint,
		contentType,
		body,
		discord.EndpointChannelMessage(m.Channel, ""),
		0,
		options...,
	)
}

// MessageEditEmbed edits an existing channel.Message, replacing it entirely with the given channel.MessageEmbed.
func (s Requester) MessageEditEmbed(channelID, messageID string, embed *channel.MessageEmbed) Request[*channel.Message] {
	return s.MessageEditEmbeds(channelID, messageID, []*channel.MessageEmbed{embed})
}

// MessageEditEmbeds edits an existing channel.Message, replacing it entirely with the multiple channel.MessageEmbed.
func (s Requester) MessageEditEmbeds(channelID, messageID string, embeds []*channel.MessageEmbed) Request[*channel.Message] {
	return s.MessageEditComplex(channel.NewMessageEdit(channelID, messageID).SetEmbeds(embeds...))
}

// MessageDelete deletes a channel.Message from the given channel.Channel.
func (s Requester) MessageDelete(channelID, messageID string) Empty {
	req := NewSimple(
		s, http.MethodDelete, discord.EndpointChannelMessage(channelID, messageID),
	).WithBucketID(discord.EndpointChannelMessage(channelID, ""))
	return WrapAsEmpty(req)
}

// MessagesBulkDelete bulk deletes the channel.Message from the channel.Channel.
//
// messages contains the list of message's ID to delete (max 100).
//
// If only one messageID is in the slice, it calls ChannelMessageDelete.
// If the slice is empty, it does nothing.
func (s Requester) MessagesBulkDelete(channelID string, messages []string) Empty {
	if len(messages) == 0 {
		// to do nothing
		return WrapErrorAsEmpty(nil)
	}

	if len(messages) == 1 {
		return s.MessageDelete(channelID, messages[0])
	}

	if len(messages) > 100 {
		return WrapErrorAsEmpty(ErrTooMuchMessagesToDelete)
	}

	data := struct {
		Messages []string `json:"messages"`
	}{messages}

	req := NewSimple(s, http.MethodPost, discord.EndpointChannelMessagesBulkDelete(channelID)).WithData(data)
	return WrapAsEmpty(req)
}

// MessagePin pins a channel.Message within the given channel.Channel.
func (s Requester) MessagePin(channelID, messageID string) Empty {
	req := NewSimple(
		s, http.MethodPut, discord.EndpointChannelMessagePin(channelID, messageID),
	).WithBucketID(discord.EndpointChannelMessagePin(channelID, ""))
	return WrapAsEmpty(req)
}

// MessageUnpin unpins a channel.Message within the given channel.Channel.
func (s Requester) MessageUnpin(channelID, messageID string) Empty {
	req := NewSimple(
		s, http.MethodDelete, discord.EndpointChannelMessagePin(channelID, messageID),
	).WithBucketID(discord.EndpointChannelMessagePin(channelID, ""))
	return WrapAsEmpty(req)
}

// MessagesPinned returns channel.MessagesPinned within the given channel.Channel.
//
// limit is the max number of users to return (max 50).
// If provided all messages returned will be before the given time.
func (s Requester) MessagesPinned(channelID string, before *time.Time, limit int) Request[*channel.MessagesPinned] {
	uri := discord.EndpointChannelMessagesPins(channelID)

	v := url.Values{}
	if before != nil {
		v.Set("before", before.Format(time.RFC3339))
	}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	return NewSimpleData[*channel.MessagesPinned](s, http.MethodGet, uri)
}

// MessageCrosspost crossposts a channel.Message in a news channel.Channel to followers.
func (s Requester) MessageCrosspost(channelID, messageID string) Request[*channel.Message] {
	return NewSimpleData[*channel.Message](
		s, http.MethodPost, discord.EndpointChannelMessageCrosspost(channelID, messageID),
	).WithBucketID(discord.EndpointChannelMessageCrosspost(channelID, ""))
}

// MessageReactionAdd creates an emoji.Emoji reaction to a channel.Message.
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format (e.g. "hello:1234567654321").
func (s Requester) MessageReactionAdd(channelID, messageID, emojiID string) Empty {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.ReplaceAll(emojiID, "#", "%23")
	req := NewSimple(
		s, http.MethodPut, discord.EndpointMessageReaction(channelID, messageID, emojiID, "@me"),
	).WithBucketID(discord.EndpointMessageReaction(channelID, "", "", "@me"))
	return WrapAsEmpty(req)
}

// MessageReactionRemove deletes an emoji.Emoji reaction to a channel.Message.
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format (e.g. "hello:1234567654321").
func (s Requester) MessageReactionRemove(channelID, messageID, emojiID, userID string) Empty {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.ReplaceAll(emojiID, "#", "%23")
	req := NewSimple(
		s, http.MethodDelete, discord.EndpointMessageReaction(channelID, messageID, emojiID, userID),
	).WithBucketID(discord.EndpointMessageReaction(channelID, "", "", "@me"))
	return WrapAsEmpty(req)
}

// MessageReactionsRemoveAll deletes all reactions from a channel.Message.
func (s Requester) MessageReactionsRemoveAll(channelID, messageID string) Empty {
	req := NewSimple(
		s, http.MethodDelete, discord.EndpointMessageReactionsAll(channelID, messageID),
	).WithBucketID(discord.EndpointMessageReactionsAll(channelID, ""))
	return WrapAsEmpty(req)
}

// MessageReactionsRemoveEmoji deletes all reactions of a certain emoji.Emoji from a channel.Message.
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format (e.g. "hello:1234567654321").
func (s Requester) MessageReactionsRemoveEmoji(channelID, messageID, emojiID string) Empty {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.ReplaceAll(emojiID, "#", "%23")
	req := NewSimple(
		s, http.MethodDelete, discord.EndpointMessageReactions(channelID, messageID, emojiID),
	).WithBucketID(discord.EndpointMessageReactions(channelID, "", ""))
	return WrapAsEmpty(req)
}

// MessageReactions gets all the user.User reactions for a specific emoji.Emoji.
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format (e.g. "hello:1234567654321").
// limit is the max number of users to return (max 100).
// If provided all reactions returned will be before beforeID.
// If provided all reactions returned will be after afterID.
func (s Requester) MessageReactions(channelID, messageID, emojiID string, limit int, beforeID, afterID string) Request[[]*user.User] {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.ReplaceAll(emojiID, "#", "%23")
	uri := discord.EndpointMessageReactions(channelID, messageID, emojiID)

	v := url.Values{}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	return NewSimpleData[[]*user.User](
		s, http.MethodGet, uri,
	).WithBucketID(discord.EndpointMessageReaction(channelID, "", "", ""))
}

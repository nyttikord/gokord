package channelapi

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// Messages returns an array of channel.Message the given channel.Channel.
//
// limit is the number messages that can be returned (max 100).
// If provided all messages returned will be before beforeID.
// If provided all messages returned will be after afterID.
// If provided all messages returned will be around aroundID.
func (s Requester) Messages(channelID string, limit int, beforeID, afterID, aroundID string, options ...discord.RequestOption) ([]*channel.Message, error) {
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

	body, err := s.Request(http.MethodGet, uri, nil, options...)
	if err != nil {
		return nil, err
	}

	var c []*channel.Message
	return c, s.Unmarshal(body, &c)
}

// Message gets channel.Message from a given channel.Channel.
func (s Requester) Message(channelID, messageID string, options ...discord.RequestOption) (*channel.Message, error) {
	response, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointChannelMessage(channelID, messageID),
		nil,
		discord.EndpointChannelMessage(channelID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var m channel.Message
	return &m, s.Unmarshal(response, &m)
}

// MessageSend sends a simple channel.Message to the given channel.Channel.
func (s Requester) MessageSend(channelID string, content string, options ...discord.RequestOption) (*channel.Message, error) {
	return s.MessageSendComplex(channelID, &channel.MessageSend{
		Content: content,
	}, options...)
}

// MessageSendComplex sends a channel.Message to the given channel.Channel.
func (s Requester) MessageSendComplex(channelID string, data *channel.MessageSend, options ...discord.RequestOption) (*channel.Message, error) {
	for _, embed := range data.Embeds {
		if embed.Type == "" {
			embed.Type = types.EmbedRich
		}
	}
	endpoint := discord.EndpointChannelMessages(channelID)

	if data.StickerIDs != nil {
		if len(data.StickerIDs) > 3 {
			return nil, ErrTooMuchStickers
		}
	}

	files := data.Files
	var err error
	var response []byte
	if len(files) > 0 {
		contentType, body, encodeErr := channel.MultipartBodyWithJSON(data, files)
		if encodeErr != nil {
			return nil, encodeErr
		}
		response, err = s.RequestRaw(http.MethodPost, endpoint, contentType, body, endpoint, 0, options...)
	} else {
		response, err = s.Request(http.MethodPost, endpoint, data, options...)
	}
	if err != nil {
		return nil, err
	}

	var m channel.Message
	return &m, s.Unmarshal(response, &m)
}

// MessageSendTTS sends a simple channel.Message to the given channel.Channel with Text to Speech.
func (s Requester) MessageSendTTS(channelID string, content string, options ...discord.RequestOption) (*channel.Message, error) {
	return s.MessageSendComplex(channelID, &channel.MessageSend{
		Content: content,
		TTS:     true,
	}, options...)
}

// MessageSendEmbed sends a channel.MessageEmbed to the given channel.Channel.
func (s Requester) MessageSendEmbed(channelID string, embed *channel.MessageEmbed, options ...discord.RequestOption) (*channel.Message, error) {
	return s.MessageSendEmbeds(channelID, []*channel.MessageEmbed{embed}, options...)
}

// MessageSendEmbeds sends multiple channel.MessageEmbed to the given channel.Channel.
func (s Requester) MessageSendEmbeds(channelID string, embeds []*channel.MessageEmbed, options ...discord.RequestOption) (*channel.Message, error) {
	return s.MessageSendComplex(channelID, &channel.MessageSend{
		Embeds: embeds,
	}, options...)
}

// MessageSendReply sends a reply channel.Message to the given channel.Channel.
//
// reference is the message reference to send containing the channel.Message to reply to.
func (s Requester) MessageSendReply(channelID string, content string, reference *channel.MessageReference, options ...discord.RequestOption) (*channel.Message, error) {
	if reference == nil {
		return nil, ErrReplyNilMessageRef
	}
	return s.MessageSendComplex(channelID, &channel.MessageSend{
		Content:   content,
		Reference: reference,
	}, options...)
}

// MessageSendEmbedReply sends a reply channel.MessageEmbed to the given channel.Channel.
//
// reference is the message reference to send containing the channel.Message to reply to.
func (s Requester) MessageSendEmbedReply(channelID string, embed *channel.MessageEmbed, reference *channel.MessageReference, options ...discord.RequestOption) (*channel.Message, error) {
	return s.MessageSendEmbedsReply(channelID, []*channel.MessageEmbed{embed}, reference, options...)
}

// MessageSendEmbedsReply sends a reply with multiple channel.MessageEmbed in the given channel.Channel.
//
// reference is the message reference to send containing the channel.Message to reply to.
func (s Requester) MessageSendEmbedsReply(channelID string, embeds []*channel.MessageEmbed, reference *channel.MessageReference, options ...discord.RequestOption) (*channel.Message, error) {
	if reference == nil {
		return nil, ErrReplyNilMessageRef
	}
	return s.MessageSendComplex(channelID, &channel.MessageSend{
		Embeds:    embeds,
		Reference: reference,
	}, options...)
}

// MessageEdit edits an existing channel.Message, replacing it entirely with the given content.
func (s Requester) MessageEdit(channelID, messageID, content string, options ...discord.RequestOption) (*channel.Message, error) {
	return s.MessageEditComplex(channel.NewMessageEdit(channelID, messageID).SetContent(content), options...)
}

// MessageEditComplex edits an existing channel.Message, replacing it entirely with the given channel.MessageEdit.
func (s Requester) MessageEditComplex(m *channel.MessageEdit, options ...discord.RequestOption) (*channel.Message, error) {
	if m.Embeds != nil {
		for _, embed := range *m.Embeds {
			if embed.Type == "" {
				embed.Type = types.EmbedRich
			}
		}
	}

	endpoint := discord.EndpointChannelMessage(m.Channel, m.ID)

	var err error
	var response []byte
	if len(m.Files) > 0 {
		contentType, body, encodeErr := channel.MultipartBodyWithJSON(m, m.Files)
		if encodeErr != nil {
			return nil, encodeErr
		}
		response, err = s.RequestRaw(
			http.MethodPatch,
			endpoint,
			contentType,
			body,
			discord.EndpointChannelMessage(m.Channel, ""),
			0,
			options...,
		)
	} else {
		response, err = s.RequestWithBucketID(
			http.MethodPatch,
			endpoint,
			m,
			discord.EndpointChannelMessage(m.Channel, ""),
			options...,
		)
	}
	if err != nil {
		return nil, err
	}

	var msg channel.Message
	return &msg, s.Unmarshal(response, &msg)
}

// MessageEditEmbed edits an existing channel.Message, replacing it entirely with the given channel.MessageEmbed.
func (s Requester) MessageEditEmbed(channelID, messageID string, embed *channel.MessageEmbed, options ...discord.RequestOption) (*channel.Message, error) {
	return s.MessageEditEmbeds(channelID, messageID, []*channel.MessageEmbed{embed}, options...)
}

// MessageEditEmbeds edits an existing channel.Message, replacing it entirely with the multiple channel.MessageEmbed.
func (s Requester) MessageEditEmbeds(channelID, messageID string, embeds []*channel.MessageEmbed, options ...discord.RequestOption) (*channel.Message, error) {
	return s.MessageEditComplex(channel.NewMessageEdit(channelID, messageID).SetEmbeds(embeds), options...)
}

// MessageDelete deletes a channel.Message from the given channel.Channel.
func (s Requester) MessageDelete(channelID, messageID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointChannelMessage(channelID, messageID),
		nil,
		discord.EndpointChannelMessage(channelID, ""),
		options...,
	)
	return
}

// MessagesBulkDelete bulk deletes the channel.Message from the channel.Channel.
//
// messages contains the list of message's ID to delete (max 100).
//
// If only one messageID is in the slice, it calls ChannelMessageDelete.
// If the slice is empty, it does nothing.
func (s Requester) MessagesBulkDelete(channelID string, messages []string, options ...discord.RequestOption) error {
	if len(messages) == 0 {
		return nil
	}

	if len(messages) == 1 {
		return s.MessageDelete(channelID, messages[0], options...)
	}

	if len(messages) > 100 {
		return ErrTooMuchMessagesToDelete
	}

	data := struct {
		Messages []string `json:"messages"`
	}{messages}

	_, err := s.Request(http.MethodPost, discord.EndpointChannelMessagesBulkDelete(channelID), data, options...)
	return err
}

// MessagePin pins a channel.Message within a given channel.Channel.
func (s Requester) MessagePin(channelID, messageID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodPut,
		discord.EndpointChannelMessagePin(channelID, messageID),
		nil,
		discord.EndpointChannelMessagePin(channelID, ""),
		options...,
	)
	return err
}

// MessageUnpin unpins a channel.Message within a given channel.Channel.
func (s Requester) MessageUnpin(channelID, messageID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointChannelMessagePin(channelID, messageID),
		nil,
		discord.EndpointChannelMessagePin(channelID, ""),
		options...,
	)
	return
}

// MessagesPinned returns all pinned channel.Message within a given channel.Channel.
func (s Requester) MessagesPinned(channelID string, options ...discord.RequestOption) ([]*channel.Message, error) {
	body, err := s.Request(http.MethodGet, discord.EndpointChannelMessagesPins(channelID), nil, options...)
	if err != nil {
		return nil, err
	}

	var m []*channel.Message
	return m, json.Unmarshal(body, &m)
}

// MessageCrosspost crossposts a channel.Message in a news channel.Channel to followers.
func (s Requester) MessageCrosspost(channelID, messageID string, options ...discord.RequestOption) (*channel.Message, error) {
	body, err := s.Request(http.MethodPost, discord.EndpointChannelMessageCrosspost(channelID, messageID), nil, options...)
	if err != nil {
		return nil, err
	}

	var m channel.Message
	return &m, s.Unmarshal(body, &m)
}

// MessageReactionAdd creates an emoji.Emoji reaction to a channel.Message.
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format (e.g. "hello:1234567654321").
func (s Requester) MessageReactionAdd(channelID, messageID, emojiID string, options ...discord.RequestOption) error {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.Replace(emojiID, "#", "%23", -1)
	_, err := s.RequestWithBucketID(
		http.MethodPut,
		discord.EndpointMessageReaction(channelID, messageID, emojiID, "@me"),
		nil,
		discord.EndpointMessageReaction(channelID, "", "", ""),
		options...,
	)
	return err
}

// MessageReactionRemove deletes an emoji.Emoji reaction to a channel.Message.
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format (e.g. "hello:1234567654321").
func (s Requester) MessageReactionRemove(channelID, messageID, emojiID, userID string, options ...discord.RequestOption) error {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.Replace(emojiID, "#", "%23", -1)
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointMessageReaction(channelID, messageID, emojiID, userID),
		nil,
		discord.EndpointMessageReaction(channelID, "", "", ""),
		options...,
	)
	return err
}

// MessageReactionsRemoveAll deletes all reactions from a channel.Message.
func (s Requester) MessageReactionsRemoveAll(channelID, messageID string, options ...discord.RequestOption) error {
	_, err := s.Request(http.MethodDelete, discord.EndpointMessageReactionsAll(channelID, messageID), nil, options...)

	return err
}

// MessageReactionsRemoveEmoji deletes all reactions of a certain emoji.Emoji from a channel.Message.
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format (e.g. "hello:1234567654321").
func (s Requester) MessageReactionsRemoveEmoji(channelID, messageID, emojiID string, options ...discord.RequestOption) error {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.Replace(emojiID, "#", "%23", -1)
	_, err := s.Request(http.MethodDelete, discord.EndpointMessageReactions(channelID, messageID, emojiID), nil, options...)

	return err
}

// MessageReactions gets all the user.Get reactions for a specific emoji.Emoji.
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format (e.g. "hello:1234567654321").
// limit is the max number of users to return (max 100).
// If provided all reactions returned will be before beforeID.
// If provided all reactions returned will be after afterID.
func (s Requester) MessageReactions(channelID, messageID, emojiID string, limit int, beforeID, afterID string, options ...discord.RequestOption) ([]*user.User, error) {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.Replace(emojiID, "#", "%23", -1)
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

	body, err := s.RequestWithBucketID(
		http.MethodGet,
		uri,
		nil,
		discord.EndpointMessageReaction(channelID, "", "", ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var u []*user.User
	return u, s.Unmarshal(body, &u)
}

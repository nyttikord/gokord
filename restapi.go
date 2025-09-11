package gokord

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg" // For JPEG decoding
	_ "image/png"  // For PNG decoding
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/interactions"
	"github.com/nyttikord/gokord/premium"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/invite"
)

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Channels
// ------------------------------------------------------------------------------------------------

// Channel returns a Channel structure of a specific Channel.
// channelID  : The ID of the Channel you want returned.
func (s *Session) Channel(channelID string, options ...discord.RequestOption) (st *channel.Channel, err error) {
	body, err := s.RequestWithBucketID("GET", discord.EndpointChannel(channelID), nil, discord.EndpointChannel(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelEdit edits the given channel and returns the updated Channel data.
// channelID  : The ID of a Channel.
// data       : New Channel data.
func (s *Session) ChannelEdit(channelID string, data *channel.Edit, options ...discord.RequestOption) (st *channel.Channel, err error) {
	body, err := s.RequestWithBucketID("PATCH", discord.EndpointChannel(channelID), data, discord.EndpointChannel(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return

}

// ChannelEditComplex edits an existing channel, replacing the parameters entirely with ChannelEdit struct
// NOTE: deprecated, use ChannelEdit instead
// channelID     : The ID of a Channel
// data          : The channel struct to send
func (s *Session) ChannelEditComplex(channelID string, data *channel.Edit, options ...discord.RequestOption) (st *channel.Channel, err error) {
	return s.ChannelEdit(channelID, data, options...)
}

// ChannelDelete deletes the given channel
// channelID  : The ID of a Channel
func (s *Session) ChannelDelete(channelID string, options ...discord.RequestOption) (st *channel.Channel, err error) {

	body, err := s.RequestWithBucketID("DELETE", discord.EndpointChannel(channelID), nil, discord.EndpointChannel(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelTyping broadcasts to all members that authenticated user is typing in
// the given channel.
// channelID  : The ID of a Channel
func (s *Session) ChannelTyping(channelID string, options ...discord.RequestOption) (err error) {

	_, err = s.RequestWithBucketID("POST", discord.EndpointChannelTyping(channelID), nil, discord.EndpointChannelTyping(channelID), options...)
	return
}

// ChannelMessages returns an array of Message structures for messages within
// a given channel.
// channelID : The ID of a Channel.
// limit     : The number messages that can be returned. (max 100)
// beforeID  : If provided all messages returned will be before given ID.
// afterID   : If provided all messages returned will be after given ID.
// aroundID  : If provided all messages returned will be around given ID.
func (s *Session) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string, options ...discord.RequestOption) (st []*channel.Message, err error) {

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

	body, err := s.RequestWithBucketID("GET", uri, nil, discord.EndpointChannelMessages(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelMessage gets a single message by ID from a given channel.
// channeld  : The ID of a Channel
// messageID : the ID of a Message
func (s *Session) ChannelMessage(channelID, messageID string, options ...discord.RequestOption) (st *channel.Message, err error) {

	response, err := s.RequestWithBucketID("GET", discord.EndpointChannelMessage(channelID, messageID), nil, discord.EndpointChannelMessage(channelID, ""), options...)
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

// ChannelMessageSend sends a message to the given channel.
// channelID : The ID of a Channel.
// content   : The message to send.
func (s *Session) ChannelMessageSend(channelID string, content string, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelMessageSendComplex(channelID, &channel.MessageSend{
		Content: content,
	}, options...)
}

// ChannelMessageSendComplex sends a message to the given channel.
// channelID : The ID of a Channel.
// data      : The message struct to send.
func (s *Session) ChannelMessageSendComplex(channelID string, data *channel.MessageSend, options ...discord.RequestOption) (st *channel.Message, err error) {
	// TODO: Remove this when compatibility is not required.
	if data.Embed != nil {
		if data.Embeds == nil {
			data.Embeds = []*channel.MessageEmbed{data.Embed}
		} else {
			err = fmt.Errorf("cannot specify both Embed and Embeds")
			return
		}
	}

	for _, embed := range data.Embeds {
		if embed.Type == "" {
			embed.Type = "rich"
		}
	}
	endpoint := discord.EndpointChannelMessages(channelID)

	// TODO: Remove this when compatibility is not required.
	files := data.Files
	if data.File != nil {
		if files == nil {
			files = []*channel.File{data.File}
		} else {
			err = fmt.Errorf("cannot specify both File and Files")
			return
		}
	}

	if data.StickerIDs != nil {
		if len(data.StickerIDs) > 3 {
			err = fmt.Errorf("cannot send more than 3 stickers")
			return
		}
	}

	var response []byte
	if len(files) > 0 {
		contentType, body, encodeErr := channel.MultipartBodyWithJSON(data, files)
		if encodeErr != nil {
			return st, encodeErr
		}
		response, err = s.RequestRaw("POST", endpoint, contentType, body, endpoint, 0, options...)
	} else {
		response, err = s.RequestWithBucketID("POST", endpoint, data, endpoint, options...)
	}
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

// ChannelMessageSendTTS sends a message to the given channel with Text to Speech.
// channelID : The ID of a Channel.
// content   : The message to send.
func (s *Session) ChannelMessageSendTTS(channelID string, content string, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelMessageSendComplex(channelID, &channel.MessageSend{
		Content: content,
		TTS:     true,
	}, options...)
}

// ChannelMessageSendEmbed sends a message to the given channel with embedded data.
// channelID : The ID of a Channel.
// embed     : The embed data to send.
func (s *Session) ChannelMessageSendEmbed(channelID string, embed *channel.MessageEmbed, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelMessageSendEmbeds(channelID, []*channel.MessageEmbed{embed}, options...)
}

// ChannelMessageSendEmbeds sends a message to the given channel with multiple embedded data.
// channelID : The ID of a Channel.
// embeds    : The embeds data to send.
func (s *Session) ChannelMessageSendEmbeds(channelID string, embeds []*channel.MessageEmbed, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelMessageSendComplex(channelID, &channel.MessageSend{
		Embeds: embeds,
	}, options...)
}

// ChannelMessageSendReply sends a message to the given channel with reference data.
// channelID : The ID of a Channel.
// content   : The message to send.
// reference : The message reference to send.
func (s *Session) ChannelMessageSendReply(channelID string, content string, reference *channel.MessageReference, options ...discord.RequestOption) (*channel.Message, error) {
	if reference == nil {
		return nil, fmt.Errorf("reply attempted with nil message reference")
	}
	return s.ChannelMessageSendComplex(channelID, &channel.MessageSend{
		Content:   content,
		Reference: reference,
	}, options...)
}

// ChannelMessageSendEmbedReply sends a message to the given channel with reference data and embedded data.
// channelID : The ID of a Channel.
// embed   : The embed data to send.
// reference : The message reference to send.
func (s *Session) ChannelMessageSendEmbedReply(channelID string, embed *channel.MessageEmbed, reference *channel.MessageReference, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelMessageSendEmbedsReply(channelID, []*channel.MessageEmbed{embed}, reference, options...)
}

// ChannelMessageSendEmbedsReply sends a message to the given channel with reference data and multiple embedded data.
// channelID : The ID of a Channel.
// embeds    : The embeds data to send.
// reference : The message reference to send.
func (s *Session) ChannelMessageSendEmbedsReply(channelID string, embeds []*channel.MessageEmbed, reference *channel.MessageReference, options ...discord.RequestOption) (*channel.Message, error) {
	if reference == nil {
		return nil, fmt.Errorf("reply attempted with nil message reference")
	}
	return s.ChannelMessageSendComplex(channelID, &channel.MessageSend{
		Embeds:    embeds,
		Reference: reference,
	}, options...)
}

// ChannelMessageEdit edits an existing message, replacing it entirely with
// the given content.
// channelID  : The ID of a Channel
// messageID  : The ID of a Message
// content    : The contents of the message
func (s *Session) ChannelMessageEdit(channelID, messageID, content string, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelMessageEditComplex(channel.NewMessageEdit(channelID, messageID).SetContent(content), options...)
}

// ChannelMessageEditComplex edits an existing message, replacing it entirely with
// the given MessageEdit struct
func (s *Session) ChannelMessageEditComplex(m *channel.MessageEdit, options ...discord.RequestOption) (st *channel.Message, err error) {
	// TODO: Remove this when compatibility is not required.
	if m.Embed != nil {
		if m.Embeds == nil {
			m.Embeds = &[]*channel.MessageEmbed{m.Embed}
		} else {
			err = fmt.Errorf("cannot specify both Embed and Embeds")
			return
		}
	}

	if m.Embeds != nil {
		for _, embed := range *m.Embeds {
			if embed.Type == "" {
				embed.Type = "rich"
			}
		}
	}

	endpoint := discord.EndpointChannelMessage(m.Channel, m.ID)

	var response []byte
	if len(m.Files) > 0 {
		contentType, body, encodeErr := channel.MultipartBodyWithJSON(m, m.Files)
		if encodeErr != nil {
			return st, encodeErr
		}
		response, err = s.RequestRaw("PATCH", endpoint, contentType, body, discord.EndpointChannelMessage(m.Channel, ""), 0, options...)
	} else {
		response, err = s.RequestWithBucketID("PATCH", endpoint, m, discord.EndpointChannelMessage(m.Channel, ""), options...)
	}
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

// ChannelMessageEditEmbed edits an existing message with embedded data.
// channelID : The ID of a Channel
// messageID : The ID of a Message
// embed     : The embed data to send
func (s *Session) ChannelMessageEditEmbed(channelID, messageID string, embed *channel.MessageEmbed, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelMessageEditEmbeds(channelID, messageID, []*channel.MessageEmbed{embed}, options...)
}

// ChannelMessageEditEmbeds edits an existing message with multiple embedded data.
// channelID : The ID of a Channel
// messageID : The ID of a Message
// embeds    : The embeds data to send
func (s *Session) ChannelMessageEditEmbeds(channelID, messageID string, embeds []*channel.MessageEmbed, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelMessageEditComplex(channel.NewMessageEdit(channelID, messageID).SetEmbeds(embeds), options...)
}

// ChannelMessageDelete deletes a message from the Channel.
func (s *Session) ChannelMessageDelete(channelID, messageID string, options ...discord.RequestOption) (err error) {

	_, err = s.RequestWithBucketID("DELETE", discord.EndpointChannelMessage(channelID, messageID), nil, discord.EndpointChannelMessage(channelID, ""), options...)
	return
}

// ChannelMessagesBulkDelete bulk deletes the messages from the channel for the provided messageIDs.
// If only one messageID is in the slice call channelMessageDelete function.
// If the slice is empty do nothing.
// channelID : The ID of the channel for the messages to delete.
// messages  : The IDs of the messages to be deleted. A slice of string IDs. A maximum of 100 messages.
func (s *Session) ChannelMessagesBulkDelete(channelID string, messages []string, options ...discord.RequestOption) (err error) {

	if len(messages) == 0 {
		return
	}

	if len(messages) == 1 {
		err = s.ChannelMessageDelete(channelID, messages[0], options...)
		return
	}

	if len(messages) > 100 {
		messages = messages[:100]
	}

	data := struct {
		Messages []string `json:"messages"`
	}{messages}

	_, err = s.RequestWithBucketID("POST", discord.EndpointChannelMessagesBulkDelete(channelID), data, discord.EndpointChannelMessagesBulkDelete(channelID), options...)
	return
}

// ChannelMessagePin pins a message within a given channel.
// channelID: The ID of a channel.
// messageID: The ID of a message.
func (s *Session) ChannelMessagePin(channelID, messageID string, options ...discord.RequestOption) (err error) {

	_, err = s.RequestWithBucketID("PUT", discord.EndpointChannelMessagePin(channelID, messageID), nil, discord.EndpointChannelMessagePin(channelID, ""), options...)
	return
}

// ChannelMessageUnpin unpins a message within a given channel.
// channelID: The ID of a channel.
// messageID: The ID of a message.
func (s *Session) ChannelMessageUnpin(channelID, messageID string, options ...discord.RequestOption) (err error) {

	_, err = s.RequestWithBucketID("DELETE", discord.EndpointChannelMessagePin(channelID, messageID), nil, discord.EndpointChannelMessagePin(channelID, ""), options...)
	return
}

// ChannelMessagesPinned returns an array of Message structures for pinned messages
// within a given channel
// channelID : The ID of a Channel.
func (s *Session) ChannelMessagesPinned(channelID string, options ...discord.RequestOption) (st []*channel.Message, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointChannelMessagesPins(channelID), nil, discord.EndpointChannelMessagesPins(channelID), options...)

	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelFileSend sends a file to the given channel.
// channelID : The ID of a Channel.
// name: The name of the file.
// io.Reader : A reader for the file contents.
func (s *Session) ChannelFileSend(channelID, name string, r io.Reader, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelMessageSendComplex(channelID, &channel.MessageSend{File: &channel.File{Name: name, Reader: r}}, options...)
}

// ChannelFileSendWithMessage sends a file to the given channel with an message.
// DEPRECATED. Use ChannelMessageSendComplex instead.
// channelID : The ID of a Channel.
// content: Optional Message content.
// name: The name of the file.
// io.Reader : A reader for the file contents.
func (s *Session) ChannelFileSendWithMessage(channelID, content string, name string, r io.Reader, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelMessageSendComplex(channelID, &channel.MessageSend{File: &channel.File{Name: name, Reader: r}, Content: content}, options...)
}

// ChannelInvites returns an array of Invite structures for the given channel
// channelID   : The ID of a Channel
func (s *Session) ChannelInvites(channelID string, options ...discord.RequestOption) (st []*invite.Invite, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointChannelInvites(channelID), nil, discord.EndpointChannelInvites(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelInviteCreate creates a new invite for the given channel.
// channelID   : The ID of a Channel
// i           : An Invite struct with the values MaxAge, MaxUses and Temporary defined.
func (s *Session) ChannelInviteCreate(channelID string, i invite.Invite, options ...discord.RequestOption) (st *invite.Invite, err error) {

	data := struct {
		MaxAge    int  `json:"max_age"`
		MaxUses   int  `json:"max_uses"`
		Temporary bool `json:"temporary"`
		Unique    bool `json:"unique"`
	}{i.MaxAge, i.MaxUses, i.Temporary, i.Unique}

	body, err := s.RequestWithBucketID("POST", discord.EndpointChannelInvites(channelID), data, discord.EndpointChannelInvites(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelPermissionSet creates a Permission Override for the given channel.
// NOTE: This func name may changed.  Using Set instead of Create because
// you can both create a new override or update an override with this function.
func (s *Session) ChannelPermissionSet(channelID, targetID string, targetType types.PermissionOverwrite, allow, deny int64, options ...discord.RequestOption) (err error) {

	data := struct {
		ID    string                    `json:"id"`
		Type  types.PermissionOverwrite `json:"type"`
		Allow int64                     `json:"allow,string"`
		Deny  int64                     `json:"deny,string"`
	}{targetID, targetType, allow, deny}

	_, err = s.RequestWithBucketID("PUT", discord.EndpointChannelPermission(channelID, targetID), data, discord.EndpointChannelPermission(channelID, ""), options...)
	return
}

// ChannelPermissionDelete deletes a specific permission override for the given channel.
// NOTE: Name of this func may change.
func (s *Session) ChannelPermissionDelete(channelID, targetID string, options ...discord.RequestOption) (err error) {

	_, err = s.RequestWithBucketID("DELETE", discord.EndpointChannelPermission(channelID, targetID), nil, discord.EndpointChannelPermission(channelID, ""), options...)
	return
}

// ChannelMessageCrosspost cross posts a message in a news channel to followers
// of the channel
// channelID   : The ID of a Channel
// messageID   : The ID of a Message
func (s *Session) ChannelMessageCrosspost(channelID, messageID string, options ...discord.RequestOption) (st *channel.Message, err error) {

	endpoint := discord.EndpointChannelMessageCrosspost(channelID, messageID)

	body, err := s.RequestWithBucketID("POST", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelNewsFollow follows a news channel in the targetID
// channelID   : The ID of a News Channel
// targetID    : The ID of a Channel where the News Channel should post to
func (s *Session) ChannelNewsFollow(channelID, targetID string, options ...discord.RequestOption) (st *channel.Follow, err error) {

	endpoint := discord.EndpointChannelFollow(channelID)

	data := struct {
		WebhookChannelID string `json:"webhook_channel_id"`
	}{targetID}

	body, err := s.RequestWithBucketID("POST", endpoint, data, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Invites
// ------------------------------------------------------------------------------------------------

// Invite returns an Invite structure of the given invite
// inviteID : The invite code
func (s *Session) Invite(inviteID string, options ...discord.RequestOption) (st *invite.Invite, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointInvite(inviteID), nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// InviteWithCounts returns an Invite structure of the given invite including approximate member counts
// inviteID : The invite code
func (s *Session) InviteWithCounts(inviteID string, options ...discord.RequestOption) (st *invite.Invite, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointInvite(inviteID)+"?with_counts=true", nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// InviteComplex returns an Invite structure of the given invite including specified fields.
// inviteID                  : The invite code
// guildScheduledEventID     : If specified, includes specified guild scheduled event.
// withCounts                : Whether to include approximate member counts or not
// withExpiration            : Whether to include expiration time or not
func (s *Session) InviteComplex(inviteID, guildScheduledEventID string, withCounts, withExpiration bool, options ...discord.RequestOption) (st *invite.Invite, err error) {
	endpoint := discord.EndpointInvite(inviteID)
	v := url.Values{}
	if guildScheduledEventID != "" {
		v.Set("guild_scheduled_event_id", guildScheduledEventID)
	}
	if withCounts {
		v.Set("with_counts", "true")
	}
	if withExpiration {
		v.Set("with_expiration", "true")
	}

	if len(v) != 0 {
		endpoint += "?" + v.Encode()
	}

	body, err := s.RequestWithBucketID("GET", endpoint, nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// InviteDelete deletes an existing invite
// inviteID   : the code of an invite
func (s *Session) InviteDelete(inviteID string, options ...discord.RequestOption) (st *invite.Invite, err error) {

	body, err := s.RequestWithBucketID("DELETE", discord.EndpointInvite(inviteID), nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// InviteAccept accepts an Invite to a Guild or Channel
// inviteID : The invite code
func (s *Session) InviteAccept(inviteID string, options ...discord.RequestOption) (st *invite.Invite, err error) {

	body, err := s.RequestWithBucketID("POST", discord.EndpointInvite(inviteID), nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Voice
// ------------------------------------------------------------------------------------------------

// VoiceRegions returns the voice server regions
func (s *Session) VoiceRegions(options ...discord.RequestOption) (st []*discord.VoiceRegion, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointVoiceRegions, nil, discord.EndpointVoiceRegions, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Websockets
// ------------------------------------------------------------------------------------------------

// Gateway returns the websocket Gateway address
func (s *Session) Gateway(options ...discord.RequestOption) (gateway string, err error) {

	response, err := s.RequestWithBucketID("GET", discord.EndpointGateway, nil, discord.EndpointGateway, options...)
	if err != nil {
		return
	}

	temp := struct {
		URL string `json:"url"`
	}{}

	err = unmarshal(response, &temp)
	if err != nil {
		return
	}

	gateway = temp.URL

	// Ensure the gateway always has a trailing slash.
	// MacOS will fail to connect if we add query params without a trailing slash on the base domain.
	if !strings.HasSuffix(gateway, "/") {
		gateway += "/"
	}

	return
}

// GatewayBot returns the websocket Gateway address and the recommended number of shards
func (s *Session) GatewayBot(options ...discord.RequestOption) (st *GatewayBotResponse, err error) {

	response, err := s.RequestWithBucketID("GET", discord.EndpointGatewayBot, nil, discord.EndpointGatewayBot, options...)
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	if err != nil {
		return
	}

	// Ensure the gateway always has a trailing slash.
	// MacOS will fail to connect if we add query params without a trailing slash on the base domain.
	if !strings.HasSuffix(st.URL, "/") {
		st.URL += "/"
	}

	return
}

// Functions specific to Webhooks

// WebhookCreate returns a new Webhook.
// channelID: The ID of a Channel.
// name     : The name of the webhook.
// avatar   : The avatar of the webhook.
func (s *Session) WebhookCreate(channelID, name, avatar string, options ...discord.RequestOption) (st *channel.Webhook, err error) {

	data := struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	body, err := s.RequestWithBucketID("POST", discord.EndpointChannelWebhooks(channelID), data, discord.EndpointChannelWebhooks(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// ChannelWebhooks returns all webhooks for a given channel.
// channelID: The ID of a channel.
func (s *Session) ChannelWebhooks(channelID string, options ...discord.RequestOption) (st []*channel.Webhook, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointChannelWebhooks(channelID), nil, discord.EndpointChannelWebhooks(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildWebhooks returns all webhooks for a given guild.
// guildID: The ID of a Guild.
func (s *Session) GuildWebhooks(guildID string, options ...discord.RequestOption) (st []*channel.Webhook, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointGuildWebhooks(guildID), nil, discord.EndpointGuildWebhooks(guildID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// Webhook returns a webhook for a given ID
// webhookID: The ID of a webhook.
func (s *Session) Webhook(webhookID string, options ...discord.RequestOption) (st *channel.Webhook, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointWebhook(webhookID), nil, discord.EndpointWebhooks, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// WebhookWithToken returns a webhook for a given ID
// webhookID: The ID of a webhook.
// token    : The auth token for the webhook.
func (s *Session) WebhookWithToken(webhookID, token string, options ...discord.RequestOption) (st *channel.Webhook, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointWebhookToken(webhookID, token), nil, discord.EndpointWebhookToken("", ""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// WebhookEdit updates an existing Webhook.
// webhookID: The ID of a webhook.
// name     : The name of the webhook.
// avatar   : The avatar of the webhook.
func (s *Session) WebhookEdit(webhookID, name, avatar, channelID string, options ...discord.RequestOption) (st *channel.Webhook, err error) {

	data := struct {
		Name      string `json:"name,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
		ChannelID string `json:"channel_id,omitempty"`
	}{name, avatar, channelID}

	body, err := s.RequestWithBucketID("PATCH", discord.EndpointWebhook(webhookID), data, discord.EndpointWebhooks, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// WebhookEditWithToken updates an existing Webhook with an auth token.
// webhookID: The ID of a webhook.
// token    : The auth token for the webhook.
// name     : The name of the webhook.
// avatar   : The avatar of the webhook.
func (s *Session) WebhookEditWithToken(webhookID, token, name, avatar string, options ...discord.RequestOption) (st *channel.Webhook, err error) {

	data := struct {
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	var body []byte
	body, err = s.RequestWithBucketID("PATCH", discord.EndpointWebhookToken(webhookID, token), data, discord.EndpointWebhookToken("", ""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// WebhookDelete deletes a webhook for a given ID
// webhookID: The ID of a webhook.
func (s *Session) WebhookDelete(webhookID string, options ...discord.RequestOption) (err error) {

	_, err = s.RequestWithBucketID("DELETE", discord.EndpointWebhook(webhookID), nil, discord.EndpointWebhooks, options...)

	return
}

// WebhookDeleteWithToken deletes a webhook for a given ID with an auth token.
// webhookID: The ID of a webhook.
// token    : The auth token for the webhook.
func (s *Session) WebhookDeleteWithToken(webhookID, token string, options ...discord.RequestOption) (st *channel.Webhook, err error) {

	body, err := s.RequestWithBucketID("DELETE", discord.EndpointWebhookToken(webhookID, token), nil, discord.EndpointWebhookToken("", ""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) webhookExecute(webhookID, token string, wait bool, threadID string, data *channel.WebhookParams, options ...discord.RequestOption) (st *channel.Message, err error) {
	uri := discord.EndpointWebhookToken(webhookID, token)

	v := url.Values{}
	if wait {
		v.Set("wait", "true")
	}

	if threadID != "" {
		v.Set("thread_id", threadID)
	}
	if len(v) != 0 {
		uri += "?" + v.Encode()
	}

	var response []byte
	if len(data.Files) > 0 {
		contentType, body, encodeErr := channel.MultipartBodyWithJSON(data, data.Files)
		if encodeErr != nil {
			return st, encodeErr
		}

		response, err = s.RequestRaw("POST", uri, contentType, body, uri, 0, options...)
	} else {
		response, err = s.RequestWithBucketID("POST", uri, data, uri, options...)
	}
	if !wait || err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

// WebhookExecute executes a webhook.
// webhookID: The ID of a webhook.
// token    : The auth token for the webhook
// wait     : Waits for server confirmation of message send and ensures that the return struct is populated (it is nil otherwise)
func (s *Session) WebhookExecute(webhookID, token string, wait bool, data *channel.WebhookParams, options ...discord.RequestOption) (st *channel.Message, err error) {
	return s.webhookExecute(webhookID, token, wait, "", data, options...)
}

// WebhookThreadExecute executes a webhook in a thread.
// webhookID: The ID of a webhook.
// token    : The auth token for the webhook
// wait     : Waits for server confirmation of message send and ensures that the return struct is populated (it is nil otherwise)
// threadID :	Sends a message to the specified thread within a webhook's channel. The thread will automatically be unarchived.
func (s *Session) WebhookThreadExecute(webhookID, token string, wait bool, threadID string, data *channel.WebhookParams, options ...discord.RequestOption) (st *channel.Message, err error) {
	return s.webhookExecute(webhookID, token, wait, threadID, data, options...)
}

// WebhookMessage gets a webhook message.
// webhookID : The ID of a webhook
// token     : The auth token for the webhook
// messageID : The ID of message to get
func (s *Session) WebhookMessage(webhookID, token, messageID string, options ...discord.RequestOption) (message *channel.Message, err error) {
	uri := discord.EndpointWebhookMessage(webhookID, token, messageID)

	body, err := s.RequestWithBucketID("GET", uri, nil, discord.EndpointWebhookToken("", ""), options...)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &message)

	return
}

// WebhookMessageEdit edits a webhook message and returns a new one.
// webhookID : The ID of a webhook
// token     : The auth token for the webhook
// messageID : The ID of message to edit
func (s *Session) WebhookMessageEdit(webhookID, token, messageID string, data *channel.WebhookEdit, options ...discord.RequestOption) (st *channel.Message, err error) {
	uri := discord.EndpointWebhookMessage(webhookID, token, messageID)

	var response []byte
	if len(data.Files) > 0 {
		contentType, body, err := channel.MultipartBodyWithJSON(data, data.Files)
		if err != nil {
			return nil, err
		}

		response, err = s.RequestRaw("PATCH", uri, contentType, body, uri, 0, options...)
		if err != nil {
			return nil, err
		}
	} else {
		response, err = s.RequestWithBucketID("PATCH", uri, data, discord.EndpointWebhookToken("", ""), options...)

		if err != nil {
			return nil, err
		}
	}

	err = unmarshal(response, &st)
	return
}

// WebhookMessageDelete deletes a webhook message.
// webhookID : The ID of a webhook
// token     : The auth token for the webhook
// messageID : The ID of a message to edit
func (s *Session) WebhookMessageDelete(webhookID, token, messageID string, options ...discord.RequestOption) (err error) {
	uri := discord.EndpointWebhookMessage(webhookID, token, messageID)

	_, err = s.RequestWithBucketID("DELETE", uri, nil, discord.EndpointWebhookToken("", ""), options...)
	return
}

// MessageReactionAdd creates an emoji reaction to a message.
// channelID : The channel ID.
// messageID : The message ID.
// emojiID   : Either the unicode emoji for the reaction, or a guild emoji identifier in name:id format (e.g. "hello:1234567654321")
func (s *Session) MessageReactionAdd(channelID, messageID, emojiID string, options ...discord.RequestOption) error {

	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.Replace(emojiID, "#", "%23", -1)
	_, err := s.RequestWithBucketID("PUT", discord.EndpointMessageReaction(channelID, messageID, emojiID, "@me"), nil, discord.EndpointMessageReaction(channelID, "", "", ""), options...)

	return err
}

// MessageReactionRemove deletes an emoji reaction to a message.
// channelID : The channel ID.
// messageID : The message ID.
// emojiID   : Either the unicode emoji for the reaction, or a guild emoji identifier.
// userID	 : @me or ID of the user to delete the reaction for.
func (s *Session) MessageReactionRemove(channelID, messageID, emojiID, userID string, options ...discord.RequestOption) error {

	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.Replace(emojiID, "#", "%23", -1)
	_, err := s.RequestWithBucketID("DELETE", discord.EndpointMessageReaction(channelID, messageID, emojiID, userID), nil, discord.EndpointMessageReaction(channelID, "", "", ""), options...)

	return err
}

// MessageReactionsRemoveAll deletes all reactions from a message
// channelID : The channel ID
// messageID : The message ID.
func (s *Session) MessageReactionsRemoveAll(channelID, messageID string, options ...discord.RequestOption) error {

	_, err := s.RequestWithBucketID("DELETE", discord.EndpointMessageReactionsAll(channelID, messageID), nil, discord.EndpointMessageReactionsAll(channelID, messageID), options...)

	return err
}

// MessageReactionsRemoveEmoji deletes all reactions of a certain emoji from a message
// channelID : The channel ID
// messageID : The message ID
// emojiID   : The emoji ID
func (s *Session) MessageReactionsRemoveEmoji(channelID, messageID, emojiID string, options ...discord.RequestOption) error {

	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.Replace(emojiID, "#", "%23", -1)
	_, err := s.RequestWithBucketID("DELETE", discord.EndpointMessageReactions(channelID, messageID, emojiID), nil, discord.EndpointMessageReactions(channelID, messageID, emojiID), options...)

	return err
}

// MessageReactions gets all the users reactions for a specific emoji.
// channelID : The channel ID.
// messageID : The message ID.
// emojiID   : Either the unicode emoji for the reaction, or a guild emoji identifier.
// limit    : max number of users to return (max 100)
// beforeID  : If provided all reactions returned will be before given ID.
// afterID   : If provided all reactions returned will be after given ID.
func (s *Session) MessageReactions(channelID, messageID, emojiID string, limit int, beforeID, afterID string, options ...discord.RequestOption) (st []*user.User, err error) {
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

	body, err := s.RequestWithBucketID("GET", uri, nil, discord.EndpointMessageReaction(channelID, "", "", ""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to threads
// ------------------------------------------------------------------------------------------------

// MessageThreadStartComplex creates a new thread from an existing message.
// channelID : Channel to create thread in
// messageID : Message to start thread from
// data : Parameters of the thread
func (s *Session) MessageThreadStartComplex(channelID, messageID string, data *channel.ThreadStart, options ...discord.RequestOption) (ch *channel.Channel, err error) {
	endpoint := discord.EndpointChannelMessageThread(channelID, messageID)
	var body []byte
	body, err = s.RequestWithBucketID("POST", endpoint, data, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &ch)
	return
}

// MessageThreadStart creates a new thread from an existing message.
// channelID       : Channel to create thread in
// messageID       : Message to start thread from
// name            : Name of the thread
// archiveDuration : Auto archive duration (in minutes)
func (s *Session) MessageThreadStart(channelID, messageID string, name string, archiveDuration int, options ...discord.RequestOption) (ch *channel.Channel, err error) {
	return s.MessageThreadStartComplex(channelID, messageID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, options...)
}

// ThreadStartComplex creates a new thread.
// channelID : Channel to create thread in
// data : Parameters of the thread
func (s *Session) ThreadStartComplex(channelID string, data *channel.ThreadStart, options ...discord.RequestOption) (ch *channel.Channel, err error) {
	endpoint := discord.EndpointChannelThreads(channelID)
	var body []byte
	body, err = s.RequestWithBucketID("POST", endpoint, data, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &ch)
	return
}

// ThreadStart creates a new thread.
// channelID       : Channel to create thread in
// name            : Name of the thread
// archiveDuration : Auto archive duration (in minutes)
func (s *Session) ThreadStart(channelID, name string, typ types.Channel, archiveDuration int, options ...discord.RequestOption) (ch *channel.Channel, err error) {
	return s.ThreadStartComplex(channelID, &channel.ThreadStart{
		Name:                name,
		Type:                typ,
		AutoArchiveDuration: archiveDuration,
	}, options...)
}

// ForumThreadStartComplex starts a new thread (creates a post) in a forum channel.
// channelID   : Channel to create thread in.
// threadData  : Parameters of the thread.
// messageData : Parameters of the starting message.
func (s *Session) ForumThreadStartComplex(channelID string, threadData *channel.ThreadStart, messageData *channel.MessageSend, options ...discord.RequestOption) (th *channel.Channel, err error) {
	endpoint := discord.EndpointChannelThreads(channelID)

	// TODO: Remove this when compatibility is not required.
	if messageData.Embed != nil {
		if messageData.Embeds == nil {
			messageData.Embeds = []*channel.MessageEmbed{messageData.Embed}
		} else {
			err = fmt.Errorf("cannot specify both Embed and Embeds")
			return
		}
	}

	for _, embed := range messageData.Embeds {
		if embed.Type == "" {
			embed.Type = "rich"
		}
	}

	// TODO: Remove this when compatibility is not required.
	files := messageData.Files
	if messageData.File != nil {
		if files == nil {
			files = []*channel.File{messageData.File}
		} else {
			err = fmt.Errorf("cannot specify both File and Files")
			return
		}
	}

	data := struct {
		*channel.ThreadStart
		Message *channel.MessageSend `json:"message"`
	}{ThreadStart: threadData, Message: messageData}

	var response []byte
	if len(files) > 0 {
		contentType, body, encodeErr := channel.MultipartBodyWithJSON(data, files)
		if encodeErr != nil {
			return th, encodeErr
		}

		response, err = s.RequestRaw("POST", endpoint, contentType, body, endpoint, 0, options...)
	} else {
		response, err = s.RequestWithBucketID("POST", endpoint, data, endpoint, options...)
	}
	if err != nil {
		return
	}

	err = unmarshal(response, &th)
	return
}

// ForumThreadStart starts a new thread (post) in a forum channel.
// channelID       : Channel to create thread in.
// name            : Name of the thread.
// archiveDuration : Auto archive duration.
// content         : Content of the starting message.
func (s *Session) ForumThreadStart(channelID, name string, archiveDuration int, content string, options ...discord.RequestOption) (th *channel.Channel, err error) {
	return s.ForumThreadStartComplex(channelID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &channel.MessageSend{Content: content}, options...)
}

// ForumThreadStartEmbed starts a new thread (post) in a forum channel.
// channelID       : Channel to create thread in.
// name            : Name of the thread.
// archiveDuration : Auto archive duration.
// embed           : Embed data of the starting message.
func (s *Session) ForumThreadStartEmbed(channelID, name string, archiveDuration int, embed *channel.MessageEmbed, options ...discord.RequestOption) (th *channel.Channel, err error) {
	return s.ForumThreadStartComplex(channelID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &channel.MessageSend{Embeds: []*channel.MessageEmbed{embed}}, options...)
}

// ForumThreadStartEmbeds starts a new thread (post) in a forum channel.
// channelID       : Channel to create thread in.
// name            : Name of the thread.
// archiveDuration : Auto archive duration.
// embeds          : Embeds data of the starting message.
func (s *Session) ForumThreadStartEmbeds(channelID, name string, archiveDuration int, embeds []*channel.MessageEmbed, options ...discord.RequestOption) (th *channel.Channel, err error) {
	return s.ForumThreadStartComplex(channelID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &channel.MessageSend{Embeds: embeds}, options...)
}

// ThreadJoin adds current user to a thread
func (s *Session) ThreadJoin(id string, options ...discord.RequestOption) error {
	endpoint := discord.EndpointThreadMember(id, "@me")
	_, err := s.RequestWithBucketID("PUT", endpoint, nil, endpoint, options...)
	return err
}

// ThreadLeave removes current user to a thread
func (s *Session) ThreadLeave(id string, options ...discord.RequestOption) error {
	endpoint := discord.EndpointThreadMember(id, "@me")
	_, err := s.RequestWithBucketID("DELETE", endpoint, nil, endpoint, options...)
	return err
}

// ThreadMemberAdd adds another member to a thread
func (s *Session) ThreadMemberAdd(threadID, memberID string, options ...discord.RequestOption) error {
	endpoint := discord.EndpointThreadMember(threadID, memberID)
	_, err := s.RequestWithBucketID("PUT", endpoint, nil, endpoint, options...)
	return err
}

// ThreadMemberRemove removes another member from a thread
func (s *Session) ThreadMemberRemove(threadID, memberID string, options ...discord.RequestOption) error {
	endpoint := discord.EndpointThreadMember(threadID, memberID)
	_, err := s.RequestWithBucketID("DELETE", endpoint, nil, endpoint, options...)
	return err
}

// ThreadMember returns thread member object for the specified member of a thread.
// withMember : Whether to include a guild member object.
func (s *Session) ThreadMember(threadID, memberID string, withMember bool, options ...discord.RequestOption) (member *channel.ThreadMember, err error) {
	uri := discord.EndpointThreadMember(threadID, memberID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	var body []byte
	body, err = s.RequestWithBucketID("GET", uri, nil, uri, options...)

	if err != nil {
		return
	}

	err = unmarshal(body, &member)
	return
}

// ThreadMembers returns all members of specified thread.
// limit      : Max number of thread members to return (1-100). Defaults to 100.
// afterID    : Get thread members after this user ID.
// withMember : Whether to include a guild member object for each thread member.
func (s *Session) ThreadMembers(threadID string, limit int, withMember bool, afterID string, options ...discord.RequestOption) (members []*channel.ThreadMember, err error) {
	uri := discord.EndpointThreadMembers(threadID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		queryParams.Set("after", afterID)
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	var body []byte
	body, err = s.RequestWithBucketID("GET", uri, nil, uri, options...)

	if err != nil {
		return
	}

	err = unmarshal(body, &members)
	return
}

// ThreadsActive returns all active threads for specified channel.
func (s *Session) ThreadsActive(channelID string, options ...discord.RequestOption) (threads *channel.ThreadsList, err error) {
	var body []byte
	body, err = s.RequestWithBucketID("GET", discord.EndpointChannelActiveThreads(channelID), nil, discord.EndpointChannelActiveThreads(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &threads)
	return
}

// GuildThreadsActive returns all active threads for specified guild.
func (s *Session) GuildThreadsActive(guildID string, options ...discord.RequestOption) (threads *channel.ThreadsList, err error) {
	var body []byte
	body, err = s.RequestWithBucketID("GET", discord.EndpointGuildActiveThreads(guildID), nil, discord.EndpointGuildActiveThreads(guildID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &threads)
	return
}

// ThreadsArchived returns archived threads for specified channel.
// before : If specified returns only threads before the timestamp
// limit  : Optional maximum amount of threads to return.
func (s *Session) ThreadsArchived(channelID string, before *time.Time, limit int, options ...discord.RequestOption) (threads *channel.ThreadsList, err error) {
	endpoint := discord.EndpointChannelPublicArchivedThreads(channelID)
	v := url.Values{}
	if before != nil {
		v.Set("before", before.Format(time.RFC3339))
	}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if len(v) > 0 {
		endpoint += "?" + v.Encode()
	}

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &threads)
	return
}

// ThreadsPrivateArchived returns archived private threads for specified channel.
// before : If specified returns only threads before the timestamp
// limit  : Optional maximum amount of threads to return.
func (s *Session) ThreadsPrivateArchived(channelID string, before *time.Time, limit int, options ...discord.RequestOption) (threads *channel.ThreadsList, err error) {
	endpoint := discord.EndpointChannelPrivateArchivedThreads(channelID)
	v := url.Values{}
	if before != nil {
		v.Set("before", before.Format(time.RFC3339))
	}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if len(v) > 0 {
		endpoint += "?" + v.Encode()
	}
	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &threads)
	return
}

// ThreadsPrivateJoinedArchived returns archived joined private threads for specified channel.
// before : If specified returns only threads before the timestamp
// limit  : Optional maximum amount of threads to return.
func (s *Session) ThreadsPrivateJoinedArchived(channelID string, before *time.Time, limit int, options ...discord.RequestOption) (threads *channel.ThreadsList, err error) {
	endpoint := discord.EndpointChannelJoinedPrivateArchivedThreads(channelID)
	v := url.Values{}
	if before != nil {
		v.Set("before", before.Format(time.RFC3339))
	}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if len(v) > 0 {
		endpoint += "?" + v.Encode()
	}
	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &threads)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to application (slash) commands
// ------------------------------------------------------------------------------------------------

// ApplicationCommandCreate creates a global application command and returns it.
// appID       : The application ID.
// guildID     : Guild ID to create guild-specific application command. If empty - creates global application command.
// cmd         : New application command data.
func (s *Session) ApplicationCommandCreate(appID string, guildID string, cmd *interactions.Command, options ...discord.RequestOption) (ccmd *interactions.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	body, err := s.RequestWithBucketID("POST", endpoint, *cmd, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &ccmd)

	return
}

// ApplicationCommandEdit edits application command and returns new command data.
// appID       : The application ID.
// cmdID       : Application command ID to edit.
// guildID     : Guild ID to edit guild-specific application command. If empty - edits global application command.
// cmd         : Updated application command data.
func (s *Session) ApplicationCommandEdit(appID, guildID, cmdID string, cmd *interactions.Command, options ...discord.RequestOption) (updated *interactions.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	body, err := s.RequestWithBucketID("PATCH", endpoint, *cmd, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &updated)

	return
}

// ApplicationCommandBulkOverwrite Creates commands overwriting existing commands. Returns a list of commands.
// appID    : The application ID.
// commands : The commands to create.
func (s *Session) ApplicationCommandBulkOverwrite(appID string, guildID string, commands []*interactions.Command, options ...discord.RequestOption) (createdCommands []*interactions.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	body, err := s.RequestWithBucketID("PUT", endpoint, commands, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &createdCommands)

	return
}

// ApplicationCommandDelete deletes application command by ID.
// appID       : The application ID.
// cmdID       : Application command ID to delete.
// guildID     : Guild ID to delete guild-specific application command. If empty - deletes global application command.
func (s *Session) ApplicationCommandDelete(appID, guildID, cmdID string, options ...discord.RequestOption) error {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	_, err := s.RequestWithBucketID("DELETE", endpoint, nil, endpoint, options...)

	return err
}

// ApplicationCommand retrieves an application command by given ID.
// appID       : The application ID.
// cmdID       : Application command ID.
// guildID     : Guild ID to retrieve guild-specific application command. If empty - retrieves global application command.
func (s *Session) ApplicationCommand(appID, guildID, cmdID string, options ...discord.RequestOption) (cmd *interactions.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	body, err := s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &cmd)

	return
}

// ApplicationCommands retrieves all commands in application.
// appID       : The application ID.
// guildID     : Guild ID to retrieve all guild-specific application commands. If empty - retrieves global application commands.
func (s *Session) ApplicationCommands(appID, guildID string, options ...discord.RequestOption) (cmd []*interactions.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?with_localizations=true", nil, "GET "+endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &cmd)

	return
}

// GuildApplicationCommandsPermissions returns permissions for application commands in a guild.
// appID       : The application ID
// guildID     : Guild ID to retrieve application commands permissions for.
func (s *Session) GuildApplicationCommandsPermissions(appID, guildID string, options ...discord.RequestOption) (permissions []*interactions.GuildCommandPermissions, err error) {
	endpoint := discord.EndpointApplicationCommandsGuildPermissions(appID, guildID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &permissions)
	return
}

// ApplicationCommandPermissions returns all permissions of an application command
// appID       : The Application ID
// guildID     : The guild ID containing the application command
// cmdID       : The command ID to retrieve the permissions of
func (s *Session) ApplicationCommandPermissions(appID, guildID, cmdID string, options ...discord.RequestOption) (permissions *interactions.GuildCommandPermissions, err error) {
	endpoint := discord.EndpointApplicationCommandPermissions(appID, guildID, cmdID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &permissions)
	return
}

// ApplicationCommandPermissionsEdit edits the permissions of an application command
// appID       : The Application ID
// guildID     : The guild ID containing the application command
// cmdID       : The command ID to edit the permissions of
// permissions : An object containing a list of permissions for the application command
//
// NOTE: Requires OAuth2 token with applications.commands.permissions.update scope
func (s *Session) ApplicationCommandPermissionsEdit(appID, guildID, cmdID string, permissions *interactions.CommandPermissionsList, options ...discord.RequestOption) (err error) {
	endpoint := discord.EndpointApplicationCommandPermissions(appID, guildID, cmdID)

	_, err = s.RequestWithBucketID("PUT", endpoint, permissions, endpoint, options...)
	return
}

// ApplicationCommandPermissionsBatchEdit edits the permissions of a batch of commands
// appID       : The Application ID
// guildID     : The guild ID to batch edit commands of
// permissions : A list of permissions paired with a command ID, guild ID, and application ID per application command
//
// NOTE: This endpoint has been disabled with updates to command permissions (Permissions v2). Please use ApplicationCommandPermissionsEdit instead.
func (s *Session) ApplicationCommandPermissionsBatchEdit(appID, guildID string, permissions []*interactions.GuildCommandPermissions, options ...discord.RequestOption) (err error) {
	endpoint := discord.EndpointApplicationCommandsGuildPermissions(appID, guildID)

	_, err = s.RequestWithBucketID("PUT", endpoint, permissions, endpoint, options...)
	return
}

// InteractionRespond creates the response to an interaction.
// interaction : Interaction instance.
// resp        : Response message data.
func (s *Session) InteractionRespond(interaction *interactions.Interaction, resp *interactions.InteractionResponse, options ...discord.RequestOption) error {
	endpoint := discord.EndpointInteractionResponse(interaction.ID, interaction.Token)

	if resp.Data != nil && len(resp.Data.Files) > 0 {
		contentType, body, err := channel.MultipartBodyWithJSON(resp, resp.Data.Files)
		if err != nil {
			return err
		}

		_, err = s.RequestRaw("POST", endpoint, contentType, body, endpoint, 0, options...)
		return err
	}

	_, err := s.RequestWithBucketID("POST", endpoint, *resp, endpoint, options...)
	return err
}

// InteractionResponse gets the response to an interaction.
// interaction : Interaction instance.
func (s *Session) InteractionResponse(interaction *interactions.Interaction, options ...discord.RequestOption) (*channel.Message, error) {
	return s.WebhookMessage(interaction.AppID, interaction.Token, "@original", options...)
}

// InteractionResponseEdit edits the response to an interaction.
// interaction : Interaction instance.
// newresp     : Updated response message data.
func (s *Session) InteractionResponseEdit(interaction *interactions.Interaction, newresp *channel.WebhookEdit, options ...discord.RequestOption) (*channel.Message, error) {
	return s.WebhookMessageEdit(interaction.AppID, interaction.Token, "@original", newresp, options...)
}

// InteractionResponseDelete deletes the response to an interaction.
// interaction : Interaction instance.
func (s *Session) InteractionResponseDelete(interaction *interactions.Interaction, options ...discord.RequestOption) error {
	endpoint := discord.EndpointInteractionResponseActions(interaction.AppID, interaction.Token)

	_, err := s.RequestWithBucketID("DELETE", endpoint, nil, endpoint, options...)

	return err
}

// FollowupMessageCreate creates the followup message for an interaction.
// interaction : Interaction instance.
// wait        : Waits for server confirmation of message send and ensures that the return struct is populated (it is nil otherwise)
// data        : Data of the message to send.
func (s *Session) FollowupMessageCreate(interaction *interactions.Interaction, wait bool, data *channel.WebhookParams, options ...discord.RequestOption) (*channel.Message, error) {
	return s.WebhookExecute(interaction.AppID, interaction.Token, wait, data, options...)
}

// FollowupMessageEdit edits a followup message of an interaction.
// interaction : Interaction instance.
// messageID   : The followup message ID.
// data        : Data to update the message
func (s *Session) FollowupMessageEdit(interaction *interactions.Interaction, messageID string, data *channel.WebhookEdit, options ...discord.RequestOption) (*channel.Message, error) {
	return s.WebhookMessageEdit(interaction.AppID, interaction.Token, messageID, data, options...)
}

// FollowupMessageDelete deletes a followup message of an interaction.
// interaction : Interaction instance.
// messageID   : The followup message ID.
func (s *Session) FollowupMessageDelete(interaction *interactions.Interaction, messageID string, options ...discord.RequestOption) error {
	return s.WebhookMessageDelete(interaction.AppID, interaction.Token, messageID, options...)
}

// ------------------------------------------------------------------------------------------------
// Functions specific to stage instances
// ------------------------------------------------------------------------------------------------

// StageInstanceCreate creates and returns a new Stage instance associated to a Stage channel.
// data : Parameters needed to create a stage instance.
// data : The data of the Stage instance to create
func (s *Session) StageInstanceCreate(data *channel.StageInstanceParams, options ...discord.RequestOption) (si *channel.StageInstance, err error) {
	body, err := s.RequestWithBucketID("POST", discord.EndpointStageInstances, data, discord.EndpointStageInstances, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &si)
	return
}

// StageInstance will retrieve a Stage instance by ID of the Stage channel.
// channelID : The ID of the Stage channel
func (s *Session) StageInstance(channelID string, options ...discord.RequestOption) (si *channel.StageInstance, err error) {
	body, err := s.RequestWithBucketID("GET", discord.EndpointStageInstance(channelID), nil, discord.EndpointStageInstance(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &si)
	return
}

// StageInstanceEdit will edit a Stage instance by ID of the Stage channel.
// channelID : The ID of the Stage channel
// data : The data to edit the Stage instance
func (s *Session) StageInstanceEdit(channelID string, data *channel.StageInstanceParams, options ...discord.RequestOption) (si *channel.StageInstance, err error) {

	body, err := s.RequestWithBucketID("PATCH", discord.EndpointStageInstance(channelID), data, discord.EndpointStageInstance(channelID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &si)
	return
}

// StageInstanceDelete will delete a Stage instance by ID of the Stage channel.
// channelID : The ID of the Stage channel
func (s *Session) StageInstanceDelete(channelID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID("DELETE", discord.EndpointStageInstance(channelID), nil, discord.EndpointStageInstance(channelID), options...)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to guilds scheduled events
// ------------------------------------------------------------------------------------------------

// GuildScheduledEvents returns an array of GuildScheduledEvent for a guild
// guildID        : The ID of a Guild
// userCount      : Whether to include the user count in the response
func (s *Session) GuildScheduledEvents(guildID string, userCount bool, options ...discord.RequestOption) (st []*guild.ScheduledEvent, err error) {
	uri := discord.EndpointGuildScheduledEvents(guildID)
	if userCount {
		uri += "?with_user_count=true"
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, discord.EndpointGuildScheduledEvents(guildID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEvent returns a specific GuildScheduledEvent in a guild
// guildID        : The ID of a Guild
// eventID        : The ID of the event
// userCount      : Whether to include the user count in the response
func (s *Session) GuildScheduledEvent(guildID, eventID string, userCount bool, options ...discord.RequestOption) (st *guild.ScheduledEvent, err error) {
	uri := discord.EndpointGuildScheduledEvent(guildID, eventID)
	if userCount {
		uri += "?with_user_count=true"
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, discord.EndpointGuildScheduledEvent(guildID, eventID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEventCreate creates a GuildScheduledEvent for a guild and returns it
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventCreate(guildID string, event *guild.ScheduledEventParams, options ...discord.RequestOption) (st *guild.ScheduledEvent, err error) {
	body, err := s.RequestWithBucketID("POST", discord.EndpointGuildScheduledEvents(guildID), event, discord.EndpointGuildScheduledEvents(guildID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEventEdit updates a specific event for a guild and returns it.
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventEdit(guildID, eventID string, event *guild.ScheduledEventParams, options ...discord.RequestOption) (st *guild.ScheduledEvent, err error) {
	body, err := s.RequestWithBucketID("PATCH", discord.EndpointGuildScheduledEvent(guildID, eventID), event, discord.EndpointGuildScheduledEvent(guildID, eventID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEventDelete deletes a specific GuildScheduledEvent in a guild
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventDelete(guildID, eventID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID("DELETE", discord.EndpointGuildScheduledEvent(guildID, eventID), nil, discord.EndpointGuildScheduledEvent(guildID, eventID), options...)
	return
}

// GuildScheduledEventUsers returns an array of GuildScheduledEventUser for a particular event in a guild
// guildID    : The ID of a Guild
// eventID    : The ID of the event
// limit      : The maximum number of users to return (Max 100)
// withMember : Whether to include the member object in the response
// beforeID   : If is not empty all returned users entries will be before the given ID
// afterID    : If is not empty all returned users entries will be after the given ID
func (s *Session) GuildScheduledEventUsers(guildID, eventID string, limit int, withMember bool, beforeID, afterID string, options ...discord.RequestOption) (st []*guild.ScheduledEventUser, err error) {
	uri := discord.EndpointGuildScheduledEventUsers(guildID, eventID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}
	if beforeID != "" {
		queryParams.Set("before", beforeID)
	}
	if afterID != "" {
		queryParams.Set("after", afterID)
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, discord.EndpointGuildScheduledEventUsers(guildID, eventID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildOnboarding returns onboarding configuration of a guild.
// guildID   : The ID of the guild
func (s *Session) GuildOnboarding(guildID string, options ...discord.RequestOption) (onboarding *guild.Onboarding, err error) {
	endpoint := discord.EndpointGuildOnboarding(guildID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &onboarding)
	return
}

// GuildOnboardingEdit edits onboarding configuration of a guild.
// guildID   : The ID of the guild
// o         : New GuildOnboarding data
func (s *Session) GuildOnboardingEdit(guildID string, o *guild.Onboarding, options ...discord.RequestOption) (onboarding *guild.Onboarding, err error) {
	endpoint := discord.EndpointGuildOnboarding(guildID)

	var body []byte
	body, err = s.RequestWithBucketID("PUT", endpoint, o, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &onboarding)
	return
}

// ----------------------------------------------------------------------
// Functions specific to auto moderation
// ----------------------------------------------------------------------

// AutoModerationRules returns a list of auto moderation rules.
// guildID : ID of the guild
func (s *Session) AutoModerationRules(guildID string, options ...discord.RequestOption) (st []*guild.AutoModerationRule, err error) {
	endpoint := discord.EndpointGuildAutoModerationRules(guildID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// AutoModerationRule returns an auto moderation rule.
// guildID : ID of the guild
// ruleID  : ID of the auto moderation rule
func (s *Session) AutoModerationRule(guildID, ruleID string, options ...discord.RequestOption) (st *guild.AutoModerationRule, err error) {
	endpoint := discord.EndpointGuildAutoModerationRule(guildID, ruleID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// AutoModerationRuleCreate creates an auto moderation rule with the given data and returns it.
// guildID : ID of the guild
// rule    : Rule data
func (s *Session) AutoModerationRuleCreate(guildID string, rule *guild.AutoModerationRule, options ...discord.RequestOption) (st *guild.AutoModerationRule, err error) {
	endpoint := discord.EndpointGuildAutoModerationRules(guildID)

	var body []byte
	body, err = s.RequestWithBucketID("POST", endpoint, rule, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// AutoModerationRuleEdit edits and returns the updated auto moderation rule.
// guildID : ID of the guild
// ruleID  : ID of the auto moderation rule
// rule    : New rule data
func (s *Session) AutoModerationRuleEdit(guildID, ruleID string, rule *guild.AutoModerationRule, options ...discord.RequestOption) (st *guild.AutoModerationRule, err error) {
	endpoint := discord.EndpointGuildAutoModerationRule(guildID, ruleID)

	var body []byte
	body, err = s.RequestWithBucketID("PATCH", endpoint, rule, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// AutoModerationRuleDelete deletes an auto moderation rule.
// guildID : ID of the guild
// ruleID  : ID of the auto moderation rule
func (s *Session) AutoModerationRuleDelete(guildID, ruleID string, options ...discord.RequestOption) (err error) {
	endpoint := discord.EndpointGuildAutoModerationRule(guildID, ruleID)
	_, err = s.RequestWithBucketID("DELETE", endpoint, nil, endpoint, options...)
	return
}

// ApplicationRoleConnectionMetadata returns application role connection metadata.
// appID : ID of the application
func (s *Session) ApplicationRoleConnectionMetadata(appID string) (st []*application.RoleConnectionMetadata, err error) {
	endpoint := discord.EndpointApplicationRoleConnectionMetadata(appID)
	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ApplicationRoleConnectionMetadataUpdate updates and returns application role connection metadata.
// appID    : ID of the application
// metadata : New metadata
func (s *Session) ApplicationRoleConnectionMetadataUpdate(appID string, metadata []*application.RoleConnectionMetadata) (st []*application.RoleConnectionMetadata, err error) {
	endpoint := discord.EndpointApplicationRoleConnectionMetadata(appID)
	var body []byte
	body, err = s.RequestWithBucketID("PUT", endpoint, metadata, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserApplicationRoleConnection returns user role connection to the specified application.
// appID : ID of the application
func (s *Session) UserApplicationRoleConnection(appID string) (st *application.RoleConnection, err error) {
	endpoint := discord.EndpointUserApplicationRoleConnection(appID)
	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return

}

// UserApplicationRoleConnectionUpdate updates and returns user role connection to the specified application.
// appID      : ID of the application
// connection : New ApplicationRoleConnection data
func (s *Session) UserApplicationRoleConnectionUpdate(appID string, rconn *application.RoleConnection) (st *application.RoleConnection, err error) {
	endpoint := discord.EndpointUserApplicationRoleConnection(appID)
	var body []byte
	body, err = s.RequestWithBucketID("PUT", endpoint, rconn, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ----------------------------------------------------------------------
// Functions specific to polls
// ----------------------------------------------------------------------

// PollAnswerVoters returns users who voted for a particular answer in a poll on the specified message.
// channelID : ID of the channel.
// messageID : ID of the message.
// answerID  : ID of the answer.
func (s *Session) PollAnswerVoters(channelID, messageID string, answerID int) (voters []*user.User, err error) {
	endpoint := discord.EndpointPollAnswerVoters(channelID, messageID, answerID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	var r struct {
		Users []*user.User `json:"users"`
	}

	err = unmarshal(body, &r)
	if err != nil {
		return
	}

	voters = r.Users
	return
}

// PollExpire expires poll on the specified message.
// channelID : ID of the channel.
// messageID : ID of the message.
func (s *Session) PollExpire(channelID, messageID string) (msg *channel.Message, err error) {
	endpoint := discord.EndpointPollExpire(channelID, messageID)

	var body []byte
	body, err = s.RequestWithBucketID("POST", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &msg)
	return
}

// ----------------------------------------------------------------------
// Functions specific to monetization
// ----------------------------------------------------------------------

// SKUs returns all SKUs for a given application.
// appID : The ID of the application.
func (s *Session) SKUs(appID string) (skus []*premium.SKU, err error) {
	endpoint := discord.EndpointApplicationSKUs(appID)

	body, err := s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &skus)
	return
}

// Entitlements returns all Entitlements for a given app, active and expired.
// appID			: The ID of the application.
// filterOptions	: Optional filter options; otherwise set it to nil.
func (s *Session) Entitlements(appID string, filterOptions *premium.EntitlementFilterOptions, options ...discord.RequestOption) (entitlements []*premium.Entitlement, err error) {
	endpoint := discord.EndpointEntitlements(appID)

	queryParams := url.Values{}
	if filterOptions != nil {
		if filterOptions.UserID != "" {
			queryParams.Set("user_id", filterOptions.UserID)
		}
		if filterOptions.SkuIDs != nil && len(filterOptions.SkuIDs) > 0 {
			queryParams.Set("sku_ids", strings.Join(filterOptions.SkuIDs, ","))
		}
		if filterOptions.Before != nil {
			queryParams.Set("before", filterOptions.Before.Format(time.RFC3339))
		}
		if filterOptions.After != nil {
			queryParams.Set("after", filterOptions.After.Format(time.RFC3339))
		}
		if filterOptions.Limit > 0 {
			queryParams.Set("limit", strconv.Itoa(filterOptions.Limit))
		}
		if filterOptions.GuildID != "" {
			queryParams.Set("guild_id", filterOptions.GuildID)
		}
		if filterOptions.ExcludeEnded {
			queryParams.Set("exclude_ended", "true")
		}
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?"+queryParams.Encode(), nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &entitlements)
	return
}

// EntitlementConsume marks a given One-Time Purchase for the user as consumed.
func (s *Session) EntitlementConsume(appID, entitlementID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID("POST", discord.EndpointEntitlementConsume(appID, entitlementID), nil, discord.EndpointEntitlementConsume(appID, ""), options...)
	return
}

// EntitlementTestCreate creates a test entitlement to a given SKU for a given guild or user.
// Discord will act as though that user or guild has entitlement to your premium offering.
func (s *Session) EntitlementTestCreate(appID string, data *premium.EntitlementTest, options ...discord.RequestOption) (err error) {
	endpoint := discord.EndpointEntitlements(appID)

	_, err = s.RequestWithBucketID("POST", endpoint, data, endpoint, options...)
	return
}

// EntitlementTestDelete deletes a currently-active test entitlement. Discord will act as though
// that user or guild no longer has entitlement to your premium offering.
func (s *Session) EntitlementTestDelete(appID, entitlementID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID("DELETE", discord.EndpointEntitlement(appID, entitlementID), nil, discord.EndpointEntitlement(appID, ""), options...)
	return
}

// Subscriptions returns all subscriptions containing the SKU.
// skuID : The ID of the SKU.
// userID : User ID for which to return subscriptions. Required except for OAuth queries.
// before : Optional timestamp to retrieve subscriptions before this time.
// after : Optional timestamp to retrieve subscriptions after this time.
// limit : Optional maximum number of subscriptions to return (1-100, default 50).
func (s *Session) Subscriptions(skuID string, userID string, before, after *time.Time, limit int, options ...discord.RequestOption) (subscriptions []*premium.Subscription, err error) {
	endpoint := discord.EndpointSubscriptions(skuID)

	queryParams := url.Values{}
	if before != nil {
		queryParams.Set("before", before.Format(time.RFC3339))
	}
	if after != nil {
		queryParams.Set("after", after.Format(time.RFC3339))
	}
	if userID != "" {
		queryParams.Set("user_id", userID)
	}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?"+queryParams.Encode(), nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &subscriptions)
	return
}

// Subscription returns a subscription by its SKU and subscription ID.
// skuID : The ID of the SKU.
// subscriptionID : The ID of the subscription.
// userID : User ID for which to return the subscription. Required except for OAuth queries.
func (s *Session) Subscription(skuID, subscriptionID, userID string, options ...discord.RequestOption) (subscription *premium.Subscription, err error) {
	endpoint := discord.EndpointSubscription(skuID, subscriptionID)

	queryParams := url.Values{}
	if userID != "" {
		// Unlike stated in the documentation, the user_id parameter is required here.
		queryParams.Set("user_id", userID)
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?"+queryParams.Encode(), nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &subscription)
	return
}

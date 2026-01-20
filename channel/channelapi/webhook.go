package channelapi

import (
	"net/http"
	"net/url"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
)

// WebhookCreate creates a new channel.Webhook.
func (s Requester) WebhookCreate(channelID, name, avatar string) Request[*channel.Webhook] {
	data := struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	return NewData[*channel.Webhook](
		s, http.MethodPost, discord.EndpointChannelWebhooks(channelID),
	).WithData(data)
}

// Webhooks returns all channel.Webhook for a given channel.Channel.
func (s Requester) Webhooks(channelID string) Request[[]*channel.Webhook] {
	return NewData[[]*channel.Webhook](
		s, http.MethodGet, discord.EndpointChannelWebhooks(channelID),
	)
}

// Webhook returns the channel.Webhook.
func (s Requester) Webhook(webhookID string) Request[*channel.Webhook] {
	return NewData[*channel.Webhook](
		s, http.MethodGet, discord.EndpointWebhook(webhookID),
	).WithBucketID(discord.EndpointWebhooks)
}

// WebhookWithToken returns a channel.Webhook for a given ID with the given token.
func (s Requester) WebhookWithToken(webhookID, token string) Request[*channel.Webhook] {
	return NewData[*channel.Webhook](
		s, http.MethodGet, discord.EndpointWebhookToken(webhookID, token),
	).WithBucketID(discord.EndpointWebhooks)
}

// WebhookEdit updates an existing channel.Webhook.
func (s Requester) WebhookEdit(webhookID, name, avatar, channelID string) Request[*channel.Webhook] {
	data := struct {
		Name      string `json:"name,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
		ChannelID string `json:"channel_id,omitempty"`
	}{name, avatar, channelID}

	return NewData[*channel.Webhook](
		s, http.MethodPatch, discord.EndpointWebhook(webhookID),
	).WithBucketID(discord.EndpointWebhooks).WithData(data)
}

// WebhookEditWithToken updates an existing channel.Webhook with an auth token.
func (s Requester) WebhookEditWithToken(webhookID, token, name, avatar string) Request[*channel.Webhook] {
	data := struct {
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	return NewData[*channel.Webhook](
		s, http.MethodPatch, discord.EndpointWebhookToken(webhookID, token),
	).WithBucketID(discord.EndpointWebhooks).WithData(data)
}

// WebhookDelete deletes a channel.Webhook.
func (s Requester) WebhookDelete(webhookID string) Empty {
	req := NewSimple(
		s, http.MethodDelete, discord.EndpointWebhook(webhookID),
	).WithBucketID(discord.EndpointWebhooks)
	return WrapAsEmpty(req)
}

// WebhookDeleteWithToken deletes a channel.Webhook with an auth token.
func (s Requester) WebhookDeleteWithToken(webhookID, token string) Empty {
	req := NewSimple(
		s, http.MethodDelete, discord.EndpointWebhookToken(webhookID, token),
	).WithBucketID(discord.EndpointWebhooks)
	return WrapAsEmpty(req)
}

func (s Requester) webhookExecute(method, uri, bucket string, wait bool, threadID string, data *channel.WebhookParams) Request[*channel.Message] {
	v := url.Values{}
	if wait {
		v.Set("wait", "true")
	}

	if len(data.Components) > 0 {
		v.Set("with_components", "true")
	}

	if threadID != "" {
		v.Set("thread_id", threadID)
	}
	if len(v) != 0 {
		uri += "?" + v.Encode()
	}

	var err error
	var response []byte
	if len(data.Files) == 0 {
		return NewData[*channel.Message](
			s, method, uri,
		).WithBucketID(bucket).WithData(data)
	}
	contentType, body, encodeErr := channel.MultipartBodyWithJSON(data, data.Files)
	if encodeErr != nil {
		return NewError[*channel.Message](encodeErr)
	}

	response, err = s.RequestRaw(ctx, method, uri, contentType, body, bucket, 0, options...)
	if !wait || err != nil {
		return nil, err
	}

	var m channel.Message
	return &m, s.Unmarshal(response, &m)
}

// WebhookExecute executes a channel.Webhook.
//
// wait if must waits for server confirmation of message send and ensures that the return struct is populated (it is nil
// otherwise)
func (s Requester) WebhookExecute(webhookID, token string, wait bool, data *channel.WebhookParams) Request[*channel.Message] {
	return s.webhookExecute(
		http.MethodPost,
		discord.EndpointWebhookToken(webhookID, token),
		discord.EndpointWebhookToken("", ""),
		wait,
		"",
		data,
	)
}

// WebhookThreadExecute executes a channel.Webhook in a thread.
//
// wait if must waits for server confirmation of message send and ensures that the return struct is populated (it is nil
// otherwise)
//
// NOTE: The thread will automatically be unarchived.
func (s Requester) WebhookThreadExecute(webhookID, token string, wait bool, threadID string, data *channel.WebhookParams) Request[*channel.Message] {
	return s.webhookExecute(
		http.MethodPost,
		discord.EndpointWebhookToken(webhookID, token),
		discord.EndpointWebhookToken("", ""),
		wait,
		threadID,
		data,
	)
}

// WebhookMessage gets a channel.Webhook channel.Message.
func (s Requester) WebhookMessage(webhookID, token, messageID string) Request[*channel.Message] {
	return NewData[*channel.Message](
		s, http.MethodGet, discord.EndpointWebhookMessage(webhookID, token, messageID),
	).WithBucketID(discord.EndpointWebhooks)
}

// WebhookMessageEdit edits a channel.Webhook channel.Message and returns the updated channel.Message.
func (s Requester) WebhookMessageEdit(webhookID, token, messageID string, data *channel.WebhookEdit) Request[*channel.Message] {
	d := &channel.WebhookParams{
		Files:           data.Files,
		AllowedMentions: data.AllowedMentions,
	}
	if data.Content != nil {
		d.Content = *data.Content
	}
	if data.Components != nil {
		d.Components = *data.Components
	}
	if data.Embeds != nil {
		d.Embeds = *data.Embeds
	}
	if data.Attachments != nil {
		d.Attachments = *data.Attachments
	}
	return s.webhookExecute(
		http.MethodPatch,
		discord.EndpointWebhookMessage(webhookID, token, messageID),
		discord.EndpointWebhookToken("", ""),
		false,
		"",
		d,
	)
}

// WebhookMessageDelete deletes a channel.Webhook channel.Message.
func (s Requester) WebhookMessageDelete(webhookID, token, messageID string) Empty {
	req := NewSimple(
		s, http.MethodDelete, discord.EndpointWebhookMessage(webhookID, token, messageID),
	).WithBucketID(discord.EndpointWebhooks)
	return WrapAsEmpty(req)
}

package channelapi

import (
	"net/http"
	"net/url"

	. "github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
)

// WebhookCreate creates a new channel.Webhook.
func (r Requester) WebhookCreate(channelID, name, avatar string) Request[*Webhook] {
	data := struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	return NewData[*Webhook](
		r, http.MethodPost, discord.EndpointChannelWebhooks(channelID),
	).WithData(data)
}

// Webhooks returns all Webhook for a given channel.Channel.
func (r Requester) Webhooks(channelID string) Request[[]*Webhook] {
	return NewData[[]*Webhook](
		r, http.MethodGet, discord.EndpointChannelWebhooks(channelID),
	)
}

// Webhook returns the channel.Webhook.
func (r Requester) Webhook(webhookID string) Request[*Webhook] {
	return NewData[*Webhook](
		r, http.MethodGet, discord.EndpointWebhook(webhookID),
	).WithBucketID(discord.EndpointWebhooks)
}

// WebhookWithToken returns a channel.Webhook for a given ID with the given token.
func (r Requester) WebhookWithToken(webhookID, token string) Request[*Webhook] {
	return NewData[*Webhook](
		r, http.MethodGet, discord.EndpointWebhookToken(webhookID, token),
	).WithBucketID(discord.EndpointWebhooks)
}

// WebhookEdit updates an existing channel.Webhook.
func (r Requester) WebhookEdit(webhookID, name, avatar, channelID string) Request[*Webhook] {
	data := struct {
		Name      string `json:"name,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
		ChannelID string `json:"channel_id,omitempty"`
	}{name, avatar, channelID}

	return NewData[*Webhook](
		r, http.MethodPatch, discord.EndpointWebhook(webhookID),
	).WithBucketID(discord.EndpointWebhooks).WithData(data)
}

// WebhookEditWithToken updates an existing channel.Webhook with an auth token.
func (r Requester) WebhookEditWithToken(webhookID, token, name, avatar string) Request[*Webhook] {
	data := struct {
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	return NewData[*Webhook](
		r, http.MethodPatch, discord.EndpointWebhookToken(webhookID, token),
	).WithBucketID(discord.EndpointWebhooks).WithData(data)
}

// WebhookDelete deletes a channel.Webhook.
func (r Requester) WebhookDelete(webhookID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointWebhook(webhookID),
	).WithBucketID(discord.EndpointWebhooks)
	return WrapAsEmpty(req)
}

// WebhookDeleteWithToken deletes a channel.Webhook with an auth token.
func (r Requester) WebhookDeleteWithToken(webhookID, token string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointWebhookToken(webhookID, token),
	).WithBucketID(discord.EndpointWebhooks)
	return WrapAsEmpty(req)
}

func (r Requester) webhookExecute(method, uri, bucket string, wait bool, threadID string, data *WebhookParams) Request[*Message] {
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
		return NewData[*Message](
			r, method, uri,
		).WithBucketID(bucket).WithData(data)
	}
	contentType, body, encodeErr := MultipartBodyWithJSON(data, data.Files)
	if encodeErr != nil {
		return NewError[*Message](encodeErr)
	}

	response, err = r.RequestRaw(ctx, method, uri, contentType, body, bucket, 0, options...)
	if !wait || err != nil {
		return nil, err
	}

	var m Message
	return &m, r.Unmarshal(response, &m)
}

// WebhookExecute executes a channel.Webhook.
//
// wait if must waits for server confirmation of message send and ensures that the return struct is populated (it is nil
// otherwise)
func (r Requester) WebhookExecute(webhookID, token string, wait bool, data *WebhookParams) Request[*Message] {
	return r.webhookExecute(
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
func (r Requester) WebhookThreadExecute(webhookID, token string, wait bool, threadID string, data *WebhookParams) Request[*Message] {
	return r.webhookExecute(
		http.MethodPost,
		discord.EndpointWebhookToken(webhookID, token),
		discord.EndpointWebhookToken("", ""),
		wait,
		threadID,
		data,
	)
}

// WebhookMessage gets a channel.Webhook channel.Message.
func (r Requester) WebhookMessage(webhookID, token, messageID string) Request[*Message] {
	return NewData[*Message](
		r, http.MethodGet, discord.EndpointWebhookMessage(webhookID, token, messageID),
	).WithBucketID(discord.EndpointWebhooks)
}

// WebhookMessageEdit edits a channel.Webhook channel.Message and returns the updated channel.Message.
func (r Requester) WebhookMessageEdit(webhookID, token, messageID string, data *WebhookEdit) Request[*Message] {
	d := &WebhookParams{
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
	return r.webhookExecute(
		http.MethodPatch,
		discord.EndpointWebhookMessage(webhookID, token, messageID),
		discord.EndpointWebhookToken("", ""),
		false,
		"",
		d,
	)
}

// WebhookMessageDelete deletes a channel.Webhook channel.Message.
func (r Requester) WebhookMessageDelete(webhookID, token, messageID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointWebhookMessage(webhookID, token, messageID),
	).WithBucketID(discord.EndpointWebhooks)
	return WrapAsEmpty(req)
}

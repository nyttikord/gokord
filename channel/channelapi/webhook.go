package channelapi

import (
	"net/http"
	"net/url"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
)

// WebhookCreate creates a new channel.Webhook.
func (s Requester) WebhookCreate(channelID, name, avatar string, options ...discord.RequestOption) (*channel.Webhook, error) {
	data := struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	body, err := s.Request(http.MethodPost, discord.EndpointChannelWebhooks(channelID), data, options...)
	if err != nil {
		return nil, err
	}

	var w channel.Webhook
	return &w, s.Unmarshal(body, &w)
}

// Webhooks returns all channel.Webhook for a given channel.Channel.
func (s Requester) Webhooks(channelID string, options ...discord.RequestOption) ([]*channel.Webhook, error) {
	body, err := s.Request(http.MethodGet, discord.EndpointChannelWebhooks(channelID), nil, options...)
	if err != nil {
		return nil, err
	}

	var ws []*channel.Webhook
	return ws, s.Unmarshal(body, ws)
}

// GuildWebhooks returns all channel.Webhook for a given guild.Guild.
func (s Requester) GuildWebhooks(guildID string, options ...discord.RequestOption) ([]*channel.Webhook, error) {
	body, err := s.Request(http.MethodGet, discord.EndpointGuildWebhooks(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var ws []*channel.Webhook
	return ws, s.Unmarshal(body, &ws)
}

// Webhook returns the channel.Webhook.
func (s Requester) Webhook(webhookID string, options ...discord.RequestOption) (*channel.Webhook, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointWebhook(webhookID),
		nil,
		discord.EndpointWebhooks,
		options...,
	)
	if err != nil {
		return nil, err
	}

	var w channel.Webhook
	return &w, s.Unmarshal(body, &w)
}

// WebhookWithToken returns a channel.Webhook for a given ID with the given token.
func (s Requester) WebhookWithToken(webhookID, token string, options ...discord.RequestOption) (*channel.Webhook, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointWebhookToken(webhookID, token),
		nil,
		discord.EndpointWebhookToken("", ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var w channel.Webhook
	return &w, s.Unmarshal(body, &w)
}

// WebhookEdit updates an existing channel.Webhook.
func (s Requester) WebhookEdit(webhookID, name, avatar, channelID string, options ...discord.RequestOption) (*channel.Webhook, error) {
	data := struct {
		Name      string `json:"name,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
		ChannelID string `json:"channel_id,omitempty"`
	}{name, avatar, channelID}

	body, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointWebhook(webhookID),
		data,
		discord.EndpointWebhooks,
		options...,
	)
	if err != nil {
		return nil, err
	}

	var w channel.Webhook
	return &w, s.Unmarshal(body, &w)
}

// WebhookEditWithToken updates an existing channel.Webhook with an auth token.
func (s Requester) WebhookEditWithToken(webhookID, token, name, avatar string, options ...discord.RequestOption) (*channel.Webhook, error) {
	data := struct {
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	body, err := s.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointWebhookToken(webhookID, token),
		data,
		discord.EndpointWebhookToken("", ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var w channel.Webhook
	return &w, s.Unmarshal(body, &w)
}

// WebhookDelete deletes a channel.Webhook.
func (s Requester) WebhookDelete(webhookID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointWebhook(webhookID),
		nil,
		discord.EndpointWebhooks,
		options...,
	)
	return err
}

// WebhookDeleteWithToken deletes a channel.Webhook with an auth token.
func (s Requester) WebhookDeleteWithToken(webhookID, token string, options ...discord.RequestOption) (*channel.Webhook, error) {
	body, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointWebhookToken(webhookID, token),
		nil,
		discord.EndpointWebhookToken("", ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var w channel.Webhook
	return &w, s.Unmarshal(body, &w)
}

func (s Requester) webhookExecute(webhookID, token string, wait bool, threadID string, data *channel.WebhookParams, options ...discord.RequestOption) (*channel.Message, error) {
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

	var err error
	var response []byte
	if len(data.Files) > 0 {
		contentType, body, encodeErr := channel.MultipartBodyWithJSON(data, data.Files)
		if encodeErr != nil {
			return nil, encodeErr
		}

		response, err = s.RequestRaw("POST", uri, contentType, body, uri, 0, options...)
	} else {
		response, err = s.Request("POST", uri, data, options...)
	}
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
func (s Requester) WebhookExecute(webhookID, token string, wait bool, data *channel.WebhookParams, options ...discord.RequestOption) (*channel.Message, error) {
	return s.webhookExecute(webhookID, token, wait, "", data, options...)
}

// WebhookThreadExecute executes a channel.Webhook in a thread.
//
// wait if must waits for server confirmation of message send and ensures that the return struct is populated (it is nil
// otherwise)
//
// Note: The thread will automatically be unarchived.
func (s Requester) WebhookThreadExecute(webhookID, token string, wait bool, threadID string, data *channel.WebhookParams, options ...discord.RequestOption) (*channel.Message, error) {
	return s.webhookExecute(webhookID, token, wait, threadID, data, options...)
}

// WebhookMessage gets a channel.Webhook channel.Message.
func (s Requester) WebhookMessage(webhookID, token, messageID string, options ...discord.RequestOption) (*channel.Message, error) {
	body, err := s.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointWebhookMessage(webhookID, token, messageID),
		nil,
		discord.EndpointWebhookToken("", ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var m channel.Message
	return &m, s.Unmarshal(body, &m)
}

// WebhookMessageEdit edits a channel.Webhook channel.Message and returns the updated channel.Message.
func (s Requester) WebhookMessageEdit(webhookID, token, messageID string, data *channel.WebhookEdit, options ...discord.RequestOption) (*channel.Message, error) {
	uri := discord.EndpointWebhookMessage(webhookID, token, messageID)

	var err error
	var response []byte
	if len(data.Files) > 0 {
		var contentType string
		var body []byte
		contentType, body, err = channel.MultipartBodyWithJSON(data, data.Files)
		if err != nil {
			return nil, err
		}

		response, err = s.RequestRaw(http.MethodPatch, uri, contentType, body, uri, 0, options...)
	} else {
		response, err = s.RequestWithBucketID(
			http.MethodPatch,
			uri,
			data,
			discord.EndpointWebhookToken("", ""),
			options...,
		)
	}
	if err != nil {
		return nil, err
	}

	var m channel.Message
	return &m, s.Unmarshal(response, &m)
}

// WebhookMessageDelete deletes a channel.Webhook channel.Message.
func (s Requester) WebhookMessageDelete(webhookID, token, messageID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointWebhookMessage(webhookID, token, messageID),
		nil,
		discord.EndpointWebhookToken("", ""),
		options...,
	)
	return err
}

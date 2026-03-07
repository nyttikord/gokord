package channel

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/nyttikord/gokord/component"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// Webhook stores the data for a webhook.
type Webhook struct {
	ID        uint64        `json:"id,string"`
	Type      types.Webhook `json:"type"`
	GuildID   uint64        `json:"guild_id,string"`
	ChannelID uint64        `json:"channel_id,string"`
	User      *user.User    `json:"user"`
	Name      string        `json:"name"`
	Avatar    string        `json:"avatar"`
	Token     string        `json:"token"`

	// ApplicationID is the bot/OAuth2 application that created this [Webhook].
	ApplicationID uint64 `json:"application_id,omitempty,string"`
}

// WebhookParams is used in the [ExecuteWebhook].
type WebhookParams struct {
	Content         string                  `json:"content,omitempty"`
	Username        string                  `json:"username,omitempty"`
	AvatarURL       string                  `json:"avatar_url,omitempty"`
	TTS             bool                    `json:"tts,omitempty"`
	Files           []*File                 `json:"-"`
	Components      []component.Message     `json:"components"`
	Embeds          []*MessageEmbed         `json:"embeds,omitempty"`
	Attachments     []*MessageAttachment    `json:"attachments,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	// Only [MessageFlagsSuppressEmbeds] and [MessageFlagsEphemeral] can be set.
	// [MessageFlagsEphemeral] can only be set when using Followup Message Create endpoint.
	Flags MessageFlags `json:"flags,omitempty"`
	// Name of the thread to create.
	//
	// NOTE: can only be set if the [Channel] is a forum.
	ThreadName string `json:"thread_name,omitempty"`
}

// WebhookEdit stores data for editing of a [Webhook] [Message].
type WebhookEdit struct {
	Content         *string                 `json:"content,omitempty"`
	Components      *[]component.Message    `json:"components,omitempty"`
	Embeds          *[]*MessageEmbed        `json:"embeds,omitempty"`
	Files           []*File                 `json:"-"`
	Attachments     *[]*MessageAttachment   `json:"attachments,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlags            `json:"flags,omitempty"`
}

// CreateWebhook in the given [Channel].
func CreateWebhook(channelID uint64, name, avatar string) Request[*Webhook] {
	data := struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	return NewData[*Webhook](http.MethodPost, discord.EndpointChannelWebhooks(channelID)).
		WithData(data)
}

// ListWebhooks returns all [Webhook] for a given [Channel].
func ListWebhooks(channelID uint64) Request[[]*Webhook] {
	return NewData[[]*Webhook](http.MethodGet, discord.EndpointChannelWebhooks(channelID))
}

// GetWebhook returns the [Webhook].
func GetWebhook(webhookID uint64) Request[*Webhook] {
	return NewData[*Webhook](http.MethodGet, discord.EndpointWebhook(webhookID)).
		WithBucketID(discord.EndpointWebhooks)
}

// GetWebhookWithToken returns a [Webhook] for a given ID with the given token.
func GetWebhookWithToken(webhookID uint64, token string) Request[*Webhook] {
	return NewData[*Webhook](http.MethodGet, discord.EndpointWebhookToken(webhookID, token)).
		WithBucketID(discord.EndpointWebhooks)
}

// EditWebhook with the given data.
func EditWebhook(webhookID uint64, name, avatar string, channelID uint64) Request[*Webhook] {
	data := struct {
		Name      string `json:"name,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
		ChannelID uint64 `json:"channel_id,omitempty,string"`
	}{name, avatar, channelID}

	return NewData[*Webhook](http.MethodPatch, discord.EndpointWebhook(webhookID)).
		WithBucketID(discord.EndpointWebhooks).WithData(data)
}

// EditWebhookWithToken with an auth token.
func EditWebhookWithToken(webhookID uint64, token, name, avatar string) Request[*Webhook] {
	data := struct {
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	return NewData[*Webhook](http.MethodPatch, discord.EndpointWebhookToken(webhookID, token)).
		WithBucketID(discord.EndpointWebhooks).WithData(data)
}

// DeleteWebhook deletes a [Webhook].
func DeleteWebhook(webhookID uint64) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointWebhook(webhookID)).
		WithBucketID(discord.EndpointWebhooks)
	return WrapAsEmpty(req)
}

// DeleteWebhookWithToken with an auth token.
func DeleteWebhookWithToken(webhookID uint64, token string) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointWebhookToken(webhookID, token)).
		WithBucketID(discord.EndpointWebhooks)
	return WrapAsEmpty(req)
}

func webhookExecute(method, uri, bucket string, wait bool, threadID uint64, data *WebhookParams) Request[*Message] {
	v := url.Values{}
	if wait {
		v.Set("wait", "true")
	}

	if len(data.Components) > 0 {
		v.Set("with_components", "true")
	}

	if threadID != 0 {
		v.Set("thread_id", fmt.Sprintf("%d", threadID))
	}
	if len(v) != 0 {
		uri += "?" + v.Encode()
	}

	if len(data.Files) == 0 {
		return NewData[*Message](method, uri).WithBucketID(bucket).WithData(data)
	}
	return NewMultipart[*Message](http.MethodPost, uri, data, data.Files).WithBucketID(bucket)
}

// ExecuteWebhook sends a [Webhook] [Message] in a [Channel].
//
// wait if must waits for server confirmation of message send and ensures that the return struct is populated (it is nil
// otherwise)
func ExecuteWebhook(webhookID uint64, token string, wait bool, data *WebhookParams) Request[*Message] {
	return webhookExecute(
		http.MethodPost,
		discord.EndpointWebhookToken(webhookID, token),
		discord.EndpointWebhookToken(0, ""),
		wait,
		0,
		data,
	)
}

// ExecuteWebhookInThread sends a [Webhook] [Message] in a thread.
//
// wait if must waits for server confirmation of message send and ensures that the return struct is populated (it is nil
// otherwise)
//
// NOTE: The thread will automatically be unarchived.
func ExecuteWebhookInThread(webhookID uint64, token string, wait bool, threadID uint64, data *WebhookParams) Request[*Message] {
	return webhookExecute(
		http.MethodPost,
		discord.EndpointWebhookToken(webhookID, token),
		discord.EndpointWebhookToken(0, ""),
		wait,
		threadID,
		data,
	)
}

// GetWebhookMessage gets a [Message] sent by a [Webhook].
func GetWebhookMessage(webhookID uint64, token string, messageID uint64) Request[*Message] {
	return NewData[*Message](http.MethodGet, discord.EndpointWebhookMessage(webhookID, token, messageID)).
		WithBucketID(discord.EndpointWebhooks)
}

// EditWebhookMessage edits a [Message] sent by a [Webhook] and returns the updated [Message].
func EditWebhookMessage(webhookID uint64, token string, messageID uint64, data *WebhookEdit) Request[*Message] {
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
	return webhookExecute(
		http.MethodPatch,
		discord.EndpointWebhookMessage(webhookID, token, messageID),
		discord.EndpointWebhookToken(0, ""),
		false,
		0,
		d,
	)
}

// DeleteWebhookMessage deletes a [Message] sent by a [Webhook].
func DeleteWebhookMessage(webhookID uint64, token string, messageID uint64) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointWebhookMessage(webhookID, token, messageID)).
		WithBucketID(discord.EndpointWebhooks)
	return WrapAsEmpty(req)
}

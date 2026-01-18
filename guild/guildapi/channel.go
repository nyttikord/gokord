package guildapi

import (
	"context"
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
)

// Channels returns the list of channel.Channel in the guild.Guild.
func (r Requester) Channels(ctx context.Context, guildID string, options ...discord.RequestOption) ([]*channel.Channel, error) {
	body, err := r.RequestRaw(
		ctx,
		http.MethodGet,
		discord.EndpointGuildChannels(guildID),
		"",
		nil,
		discord.EndpointGuildChannels(guildID),
		0,
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st []*channel.Channel
	return st, r.Unmarshal(body, &st)
}

// ChannelCreateData is provided to Requester.ChannelCreateComplex
type ChannelCreateData struct {
	Name                 string                         `json:"name"`
	Type                 types.Channel                  `json:"type"`
	Topic                string                         `json:"topic,omitempty"`
	Bitrate              int                            `json:"bitrate,omitempty"`
	UserLimit            int                            `json:"user_limit,omitempty"`
	RateLimitPerUser     int                            `json:"rate_limit_per_user,omitempty"`
	Position             int                            `json:"position,omitempty"`
	PermissionOverwrites []*channel.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             string                         `json:"parent_id,omitempty"`
	NSFW                 bool                           `json:"nsfw,omitempty"`
}

// ChannelCreateComplex creates a new channel.Channel in the given guild.Guild.
func (r Requester) ChannelCreateComplex(ctx context.Context, guildID string, data ChannelCreateData, options ...discord.RequestOption) (*channel.Channel, error) {
	body, err := r.Request(ctx, http.MethodPost, discord.EndpointGuildChannels(guildID), data, options...)
	if err != nil {
		return nil, err
	}

	var st channel.Channel
	return &st, r.Unmarshal(body, &st)
}

// ChannelCreate creates a new channel.Channel in the given guild.Guild.
func (r Requester) ChannelCreate(ctx context.Context, guildID, name string, ctype types.Channel, options ...discord.RequestOption) (st *channel.Channel, err error) {
	return r.ChannelCreateComplex(ctx, guildID, ChannelCreateData{
		Name: name,
		Type: ctype,
	}, options...)
}

// ChannelsReorder updates the order of channel.Channel in a guild.Guild.
func (r Requester) ChannelsReorder(ctx context.Context, guildID string, channels []*channel.Channel, options ...discord.RequestOption) error {
	data := make([]struct {
		ID       string `json:"id"`
		Position int    `json:"position"`
	}, len(channels))

	for i, c := range channels {
		data[i].ID = c.ID
		data[i].Position = c.Position
	}

	_, err := r.Request(ctx, http.MethodPatch, discord.EndpointGuildChannels(guildID), data, options...)
	return err
}

// ThreadsActive returns all active threads in the given guild.Guild.
func (r Requester) ThreadsActive(ctx context.Context, guildID string, options ...discord.RequestOption) (*channel.ThreadsList, error) {
	body, err := r.Request(ctx, http.MethodGet, discord.EndpointGuildActiveThreads(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var tl channel.ThreadsList
	return &tl, r.Unmarshal(body, &tl)
}

// Webhooks returns all channel.Webhook for a given guild.Guild.
func (r Requester) Webhooks(ctx context.Context, guildID string, options ...discord.RequestOption) ([]*channel.Webhook, error) {
	body, err := r.Request(ctx, http.MethodGet, discord.EndpointGuildWebhooks(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var ws []*channel.Webhook
	return ws, r.Unmarshal(body, &ws)
}

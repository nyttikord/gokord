package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
)

// Channels returns the list of channel.Channel in the guild.Guild.
func (r Requester) Channels(guildID string) Request[[]*channel.Channel] {
	return NewData[[]*channel.Channel](
		r, http.MethodGet, discord.EndpointGuildChannels(guildID),
	)
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
func (r Requester) ChannelCreateComplex(guildID string, data ChannelCreateData) Request[*channel.Channel] {
	return NewData[*channel.Channel](
		r, http.MethodPost, discord.EndpointGuildChannels(guildID),
	).WithData(data)
}

// ChannelCreate creates a new channel.Channel in the given guild.Guild.
func (r Requester) ChannelCreate(guildID, name string, ctype types.Channel) Request[*channel.Channel] {
	return r.ChannelCreateComplex(guildID, ChannelCreateData{
		Name: name,
		Type: ctype,
	})
}

// ChannelsReorder updates the order of channel.Channel in a guild.Guild.
func (r Requester) ChannelsReorder(guildID string, channels []*channel.Channel) Empty {
	data := make([]struct {
		ID       string `json:"id"`
		Position int    `json:"position"`
	}, len(channels))

	for i, c := range channels {
		data[i].ID = c.ID
		data[i].Position = c.Position
	}

	req := NewSimple(r, http.MethodPatch, discord.EndpointGuildChannels(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// ThreadsActive returns all active threads in the given guild.Guild.
func (r Requester) ThreadsActive(guildID string) Request[*channel.ThreadsList] {
	return NewData[*channel.ThreadsList](
		r, http.MethodGet, discord.EndpointGuildActiveThreads(guildID),
	)
}

// Webhooks returns all channel.Webhook for a given guild.Guild.
func (r Requester) Webhooks(guildID string) Request[[]*channel.Webhook] {
	return NewData[[]*channel.Webhook](
		r, http.MethodGet, discord.EndpointGuildWebhooks(guildID),
	)
}

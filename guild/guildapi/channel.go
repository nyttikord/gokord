package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
)

// Channels returns the list of channel.Channel in the guild.Guild.
func (r Requester) Channels(guildID string, options ...discord.RequestOption) ([]*channel.Channel, error) {
	body, err := r.RequestRaw(
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

// ChannelCreateData is provided to Session.GuildChannelCreateComplex
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

// ChannelCreateComplex creates a new channel.Channel in the given guild.Guild
func (r Requester) ChannelCreateComplex(guildID string, data ChannelCreateData, options ...discord.RequestOption) (*channel.Channel, error) {
	body, err := r.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointGuildChannels(guildID),
		data,
		discord.EndpointGuildChannels(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st channel.Channel
	return &st, r.Unmarshal(body, &st)
}

// ChannelCreate creates a new channel.Channel in the given guild.Guild.
func (r Requester) ChannelCreate(guildID, name string, ctype types.Channel, options ...discord.RequestOption) (st *channel.Channel, err error) {
	return r.ChannelCreateComplex(guildID, ChannelCreateData{
		Name: name,
		Type: ctype,
	}, options...)
}

// ChannelsReorder updates the order of channel.Channel in a guild.Guild.
func (r Requester) ChannelsReorder(guildID string, channels []*channel.Channel, options ...discord.RequestOption) error {
	data := make([]struct {
		ID       string `json:"id"`
		Position int    `json:"position"`
	}, len(channels))

	for i, c := range channels {
		data[i].ID = c.ID
		data[i].Position = c.Position
	}

	_, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildChannels(guildID),
		data,
		discord.EndpointGuildChannels(guildID),
		options...,
	)
	return err
}

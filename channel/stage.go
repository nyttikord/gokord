package channel

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
)

// StageInstance holds information about a live stage.
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-resource
type StageInstance struct {
	ID        uint64 `json:"id,string"`
	GuildID   uint64 `json:"guild_id,string"`
	ChannelID uint64 `json:"channel_id,string"`
	Topic     string `json:"topic"`
	// PrivacyLevel of the [StageInstance].
	// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
	PrivacyLevel StageInstancePrivacyLevel `json:"privacy_level"`
	// GuildScheduledEventID linked with this [StageInstance].
	GuildScheduledEventID uint64 `json:"guild_scheduled_event_id,string"`
}

// StageInstanceParams represents the parameters needed to create or edit a [StageInstance].
type StageInstanceParams struct {
	ChannelID uint64 `json:"channel_id,omitempty,string"`
	Topic     string `json:"topic,omitempty"`
	// PrivacyLevel of the [StageInstance].
	//
	// Default: [StageInstancePrivacyLevelGuildOnly]
	PrivacyLevel StageInstancePrivacyLevel `json:"privacy_level,omitempty"`
	// SendStartNotification will notify @everyone that a [StageInstance] has started.
	SendStartNotification bool `json:"send_start_notification,omitempty"`
}

// StageInstancePrivacyLevel represents the privacy level of a [StageInstance].
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
type StageInstancePrivacyLevel int

const (
	// NOTE: the Level "1" is not used anymore, so it was deleted.

	// StageInstancePrivacyLevelGuildOnly is visible to only [user.Member]s.
	StageInstancePrivacyLevelGuildOnly StageInstancePrivacyLevel = 2
)

// CreateStageInstance returns a new [StageInstance] associated to a [types.ChannelGuildStageVoice].
func CreateStageInstance(data *StageInstanceParams) Request[*StageInstance] {
	return NewData[*StageInstance](http.MethodPost, discord.EndpointStageInstances).
		WithData(data)
}

// GetStageInstance will retrieve a [StageInstance] by the ID of the [types.ChannelGuildStageVoice].
func GetStageInstance(channelID uint64) Request[*StageInstance] {
	return NewData[*StageInstance](http.MethodGet, discord.EndpointStageInstance(channelID))
}

// EditStageInstance by ID of the [types.ChannelGuildStageVoice].
func EditStageInstance(channelID uint64, data *StageInstanceParams) Request[*StageInstance] {
	return NewData[*StageInstance](http.MethodPatch, discord.EndpointStageInstance(channelID)).
		WithData(data)
}

// DeleteStageInstance by ID of the [types.ChannelGuildStageVoice].
func DeleteStageInstance(channelID uint64) Empty {
	req := NewSimple(http.MethodGet, discord.EndpointStageInstance(channelID))
	return WrapAsEmpty(req)
}

package channel

// StageInstance holds information about a live stage.
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-resource
type StageInstance struct {
	// The id of this Stage instance.
	ID string `json:"id"`
	// The guild.Guild id of the associated Stage Channel.
	GuildID string `json:"guild_id"`
	// The id of the associated Stage Channel.
	ChannelID string `json:"channel_id"`
	// The topic of the Stage instance (1-120 characters).
	Topic string `json:"topic"`
	// The privacy level of the Stage instance.
	// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
	PrivacyLevel StageInstancePrivacyLevel `json:"privacy_level"`
	// The id of the guild.ScheduledEvent for this Stage instance.
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
}

// StageInstanceParams represents the parameters needed to create or edit a stage instance.
type StageInstanceParams struct {
	// ChannelID represents the id of the Stage Channel.
	ChannelID string `json:"channel_id,omitempty"`
	// Topic of the Stage instance (1-120 characters).
	Topic string `json:"topic,omitempty"`
	// PrivacyLevel of the Stage instance.
	//
	// Default: StageInstancePrivacyLevelGuildOnly
	PrivacyLevel StageInstancePrivacyLevel `json:"privacy_level,omitempty"`
	// SendStartNotification will notify @everyone that a Stage instance has started.
	SendStartNotification bool `json:"send_start_notification,omitempty"`
}

// StageInstancePrivacyLevel represents the privacy level of a Stage instance.
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
type StageInstancePrivacyLevel int

const (
	// NOTE: the Level "1" is not used anymore, so it was deleted.

	// StageInstancePrivacyLevelGuildOnly The Stage instance is visible to only guild members.
	StageInstancePrivacyLevelGuildOnly StageInstancePrivacyLevel = 2
)

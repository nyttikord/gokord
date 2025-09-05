package channel

// StageInstance holds information about a live stage.
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-resource
type StageInstance struct {
	// The id of this Stage instance
	ID string `json:"id"`
	// The guild id of the associated Stage channel
	GuildID string `json:"guild_id"`
	// The id of the associated Stage channel
	ChannelID string `json:"channel_id"`
	// The topic of the Stage instance (1-120 characters)
	Topic string `json:"topic"`
	// The privacy level of the Stage instance
	// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
	PrivacyLevel StageInstancePrivacyLevel `json:"privacy_level"`
	// Whether or not Stage Discovery is disabled (deprecated)
	DiscoverableDisabled bool `json:"discoverable_disabled"`
	// The id of the scheduled event for this Stage instance
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
}

// StageInstanceParams represents the parameters needed to create or edit a stage instance
type StageInstanceParams struct {
	// ChannelID represents the id of the Stage channel
	ChannelID string `json:"channel_id,omitempty"`
	// Topic of the Stage instance (1-120 characters)
	Topic string `json:"topic,omitempty"`
	// PrivacyLevel of the Stage instance (default GUILD_ONLY)
	PrivacyLevel StageInstancePrivacyLevel `json:"privacy_level,omitempty"`
	// SendStartNotification will notify @everyone that a Stage instance has started
	SendStartNotification bool `json:"send_start_notification,omitempty"`
}

// StageInstancePrivacyLevel represents the privacy level of a Stage instance
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
type StageInstancePrivacyLevel int

const (
	// StageInstancePrivacyLevelPublic The Stage instance is visible publicly. (deprecated)
	StageInstancePrivacyLevelPublic StageInstancePrivacyLevel = 1
	// StageInstancePrivacyLevelGuildOnly The Stage instance is visible to only guild members.
	StageInstancePrivacyLevelGuildOnly StageInstancePrivacyLevel = 2
)

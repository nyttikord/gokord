package guild

import (
	"encoding/json"
	"github.com/nyttikord/gokord/user"
	"time"
)

// ScheduledEvent is a representation of a scheduled event in a guild. Only for retrieval of the data.
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event
type ScheduledEvent struct {
	// The ID of the scheduled event
	ID string `json:"id"`
	// The guild id which the scheduled event belongs to
	GuildID string `json:"guild_id"`
	// The channel id in which the scheduled event will be hosted, or null if scheduled entity type is EXTERNAL
	ChannelID string `json:"channel_id"`
	// The id of the user that created the scheduled event
	CreatorID string `json:"creator_id"`
	// The name of the scheduled event (1-100 characters)
	Name string `json:"name"`
	// The description of the scheduled event (1-1000 characters)
	Description string `json:"description"`
	// The time the scheduled event will start
	ScheduledStartTime time.Time `json:"scheduled_start_time"`
	// The time the scheduled event will end, required only when entity_type is EXTERNAL
	ScheduledEndTime *time.Time `json:"scheduled_end_time"`
	// The privacy level of the scheduled event
	PrivacyLevel ScheduledEventPrivacyLevel `json:"privacy_level"`
	// The status of the scheduled event
	Status ScheduledEventStatus `json:"status"`
	// Type of the entity where event would be hosted
	// See field requirements
	// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-field-requirements-by-entity-type
	EntityType ScheduledEventEntityType `json:"entity_type"`
	// The id of an entity associated with a guild scheduled event
	EntityID string `json:"entity_id"`
	// Additional metadata for the guild scheduled event
	EntityMetadata ScheduledEventEntityMetadata `json:"entity_metadata"`
	// The user that created the scheduled event
	Creator *user.User `json:"creator"`
	// The number of users subscribed to the scheduled event
	UserCount int `json:"user_count"`
	// The cover image hash of the scheduled event
	// see https://discord.com/developers/docs/reference#image-formatting for more
	// information about image formatting
	Image string `json:"image"`
}

// ScheduledEventParams are the parameters allowed for creating or updating a scheduled event
// https://discord.com/developers/docs/resources/guild-scheduled-event#create-guild-scheduled-event
type ScheduledEventParams struct {
	// The channel id in which the scheduled event will be hosted, or null if scheduled entity type is EXTERNAL
	ChannelID string `json:"channel_id,omitempty"`
	// The name of the scheduled event (1-100 characters)
	Name string `json:"name,omitempty"`
	// The description of the scheduled event (1-1000 characters)
	Description string `json:"description,omitempty"`
	// The time the scheduled event will start
	ScheduledStartTime *time.Time `json:"scheduled_start_time,omitempty"`
	// The time the scheduled event will end, required only when entity_type is EXTERNAL
	ScheduledEndTime *time.Time `json:"scheduled_end_time,omitempty"`
	// The privacy level of the scheduled event
	PrivacyLevel ScheduledEventPrivacyLevel `json:"privacy_level,omitempty"`
	// The status of the scheduled event
	Status ScheduledEventStatus `json:"status,omitempty"`
	// Type of the entity where event would be hosted
	// See field requirements
	// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-field-requirements-by-entity-type
	EntityType ScheduledEventEntityType `json:"entity_type,omitempty"`
	// Additional metadata for the guild scheduled event
	EntityMetadata *ScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	// The cover image hash of the scheduled event
	// see https://discord.com/developers/docs/reference#image-formatting for more
	// information about image formatting
	Image string `json:"image,omitempty"`
}

// MarshalJSON is a helper function to marshal ScheduledEventParams
func (p ScheduledEventParams) MarshalJSON() ([]byte, error) {
	type guildScheduledEventParams ScheduledEventParams

	if p.EntityType == ScheduledEventEntityTypeExternal && p.ChannelID == "" {
		return json.Marshal(struct {
			guildScheduledEventParams
			ChannelID json.RawMessage `json:"channel_id"`
		}{
			guildScheduledEventParams: guildScheduledEventParams(p),
			ChannelID:                 json.RawMessage("null"),
		})
	}

	return json.Marshal(guildScheduledEventParams(p))
}

// ScheduledEventEntityMetadata holds additional metadata for guild scheduled event.
type ScheduledEventEntityMetadata struct {
	// location of the event (1-100 characters)
	// required for events with 'entity_type': EXTERNAL
	Location string `json:"location"`
}

// ScheduledEventPrivacyLevel is the privacy level of a scheduled event.
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
type ScheduledEventPrivacyLevel int

const (
	// ScheduledEventPrivacyLevelGuildOnly makes the scheduled
	// event is only accessible to guild members
	ScheduledEventPrivacyLevelGuildOnly ScheduledEventPrivacyLevel = 2
)

// ScheduledEventStatus is the status of a scheduled event
// Valid Guild Scheduled Event Status Transitions :
// SCHEDULED --> ACTIVE --> COMPLETED
// SCHEDULED --> CANCELED
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-status
type ScheduledEventStatus int

const (
	// ScheduledEventStatusScheduled represents the current event is in scheduled state
	ScheduledEventStatusScheduled ScheduledEventStatus = 1
	// ScheduledEventStatusActive represents the current event is in active state
	ScheduledEventStatusActive ScheduledEventStatus = 2
	// ScheduledEventStatusCompleted represents the current event is in completed state
	ScheduledEventStatusCompleted ScheduledEventStatus = 3
	// ScheduledEventStatusCanceled represents the current event is in canceled state
	ScheduledEventStatusCanceled ScheduledEventStatus = 4
)

// ScheduledEventEntityType is the type of entity associated with a guild scheduled event.
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types
type ScheduledEventEntityType int

const (
	// ScheduledEventEntityTypeStageInstance represents a stage channel
	ScheduledEventEntityTypeStageInstance ScheduledEventEntityType = 1
	// ScheduledEventEntityTypeVoice represents a voice channel
	ScheduledEventEntityTypeVoice ScheduledEventEntityType = 2
	// ScheduledEventEntityTypeExternal represents an external event
	ScheduledEventEntityTypeExternal ScheduledEventEntityType = 3
)

// ScheduledEventUser is a user subscribed to a scheduled event.
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-user-object
type ScheduledEventUser struct {
	GuildScheduledEventID string       `json:"guild_scheduled_event_id"`
	User                  *user.User   `json:"user"`
	Member                *user.Member `json:"member"`
}

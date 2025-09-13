package guild

import (
	"encoding/json"
	"time"

	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// ScheduledEvent is a representation of a scheduled event in a Guild.
// Only for retrieval of the data.
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event
type ScheduledEvent struct {
	// The ID of the ScheduledEvent.
	ID string `json:"id"`
	// The Guild.ID which the ScheduledEvent belongs to.
	GuildID string `json:"guild_id"`
	// The channel.Channel ID in which the ScheduledEvent will be hosted, or null if EntityType is
	// types.ScheduledEventEntityExternal.
	ChannelID string `json:"channel_id"`
	// The user.Get ID that created the ScheduledEvent.
	CreatorID string `json:"creator_id"`
	// The Name of the ScheduledEvent (1-100 characters)
	Name string `json:"name"`
	// The Description of the ScheduledEvent (1-1000 characters)
	Description string `json:"description"`
	// The time the ScheduledEvent will start.
	ScheduledStartTime time.Time `json:"scheduled_start_time"`
	// The time the ScheduledEvent will end, required only when EntityType is types.ScheduledEventEntityExternal.
	ScheduledEndTime *time.Time `json:"scheduled_end_time"`
	// The PrivacyLevel of the ScheduledEvent.
	PrivacyLevel ScheduledEventPrivacyLevel `json:"privacy_level"`
	// The Status of the ScheduledEvent.
	Status ScheduledEventStatus `json:"status"`
	// Type of the entity where ScheduledEvent would be hosted.
	//
	// Note: See field requirements.
	// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-field-requirements-by-entity-type
	EntityType types.ScheduledEventEntity `json:"entity_type"`
	// The ID of an entity associated with a ScheduledEvent.
	EntityID string `json:"entity_id"`
	// Additional metadata for the ScheduledEvent.
	EntityMetadata ScheduledEventEntityMetadata `json:"entity_metadata"`
	// The user.Get that created the ScheduledEvent.
	Creator *user.User `json:"creator"`
	// The number of users subscribed to the ScheduledEvent.
	UserCount int `json:"user_count"`
	// The cover Image hash of the ScheduledEvent.
	//
	// Note: See https://discord.com/developers/docs/reference#image-formatting for more information about image
	// formatting.
	Image string `json:"image"`
}

// ScheduledEventParams are the parameters allowed for creating or updating a ScheduledEvent.
// https://discord.com/developers/docs/resources/guild-scheduled-event#create-guild-scheduled-event
type ScheduledEventParams struct {
	// The channel.Channel ID in which the ScheduledEvent will be hosted, or null if EntityType is
	// types.ScheduledEventEntityExternal.
	ChannelID string `json:"channel_id,omitempty"`
	// The Name of the ScheduledEvent (1-100 characters)
	Name string `json:"name,omitempty"`
	// The Description of the ScheduledEvent (1-1000 characters)
	Description string `json:"description,omitempty"`
	// The time the ScheduledEvent will start.
	ScheduledStartTime *time.Time `json:"scheduled_start_time,omitempty"`
	// The time the ScheduledEvent will end, required only when EntityType is types.ScheduledEventEntityExternal.
	ScheduledEndTime *time.Time `json:"scheduled_end_time,omitempty"`
	// The PrivacyLevel of the ScheduledEvent.
	PrivacyLevel ScheduledEventPrivacyLevel `json:"privacy_level,omitempty"`
	// The Status of the ScheduledEvent.
	Status ScheduledEventStatus `json:"status,omitempty"`
	// Type of the entity where ScheduledEvent would be hosted.
	//
	// Note: See field requirements
	// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-field-requirements-by-entity-type
	EntityType types.ScheduledEventEntity `json:"entity_type,omitempty"`
	// Additional metadata for the ScheduledEvent.
	EntityMetadata *ScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	// The cover Image hash of the ScheduledEvent.
	//
	// Note: See https://discord.com/developers/docs/reference#image-formatting for more information about image
	// formatting.
	Image string `json:"image,omitempty"`
}

// MarshalJSON is a helper function to marshal ScheduledEventParams
func (p ScheduledEventParams) MarshalJSON() ([]byte, error) {
	type guildScheduledEventParams ScheduledEventParams

	if p.EntityType == types.ScheduledEventEntityExternal && p.ChannelID == "" {
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

// ScheduledEventEntityMetadata holds additional metadata for ScheduledEvent.
type ScheduledEventEntityMetadata struct {
	// Location of the ScheduledEvent (1-100 characters)
	//
	// Required for types.ScheduledEventEntityExternal
	Location string `json:"location"`
}

// ScheduledEventPrivacyLevel is the privacy level of a ScheduledEvent.
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
type ScheduledEventPrivacyLevel int

const (
	// ScheduledEventPrivacyLevelGuildOnly makes the scheduled
	// event is only accessible to guild members
	ScheduledEventPrivacyLevelGuildOnly ScheduledEventPrivacyLevel = 2
)

// ScheduledEventStatus is the status of a ScheduledEvent.
//
// Valid ScheduledEvent Status Transitions :
//
//	SCHEDULED --> ACTIVE --> COMPLETED
//	SCHEDULED --> CANCELED
//
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

// ScheduledEventUser is a user.Get subscribed to a ScheduledEvent.
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-user-object
type ScheduledEventUser struct {
	GuildScheduledEventID string       `json:"guild_scheduled_event_id"`
	User                  *user.User   `json:"user"`
	Member                *user.Member `json:"member"`
}

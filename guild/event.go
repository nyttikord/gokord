package guild

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// ScheduledEvent is a representation of a scheduled event in a [Guild].
// Only for retrieval of the data.
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event
type ScheduledEvent struct {
	ID uint64 `json:"id,string"`
	// The [Guild.ID] which the [ScheduledEvent] belongs to.
	GuildID uint64 `json:"guild_id,string"`
	// The [channel.Channel] ID in which the [ScheduledEvent] will be hosted, or null if [EntityType] is
	// [types.ScheduledEventEntityExternal].
	ChannelID uint64 `json:"channel_id,string"`
	// The [user.User] ID that created the ScheduledEvent.
	CreatorID   uint64 `json:"creator_id,string"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// The time the [ScheduledEvent] will start.
	ScheduledStartTime time.Time `json:"scheduled_start_time"`
	// The time the [ScheduledEvent] will end, required only when [EntityType] is [types.ScheduledEventEntityExternal].
	ScheduledEndTime *time.Time `json:"scheduled_end_time"`
	// The PrivacyLevel of the [ScheduledEvent].
	PrivacyLevel ScheduledEventPrivacyLevel `json:"privacy_level"`
	// Status of the [ScheduledEvent].
	Status ScheduledEventStatus `json:"status"`
	// Type of the entity where [ScheduledEvent] would be hosted.
	//
	// NOTE: See field requirements.
	// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-field-requirements-by-entity-type
	EntityType types.ScheduledEventEntity `json:"entity_type"`
	// The ID of an entity associated with a [ScheduledEvent].
	EntityID uint64 `json:"entity_id,string"`
	// Additional metadata for the [ScheduledEvent].
	EntityMetadata ScheduledEventEntityMetadata `json:"entity_metadata"`
	// The [user.User] that created the [ScheduledEvent].
	Creator *user.User `json:"creator"`
	// The number of [user.User]s subscribed to the [ScheduledEvent].
	UserCount int `json:"user_count"`
	// The cover Image hash of the [ScheduledEvent].
	//
	// NOTE: See https://discord.com/developers/docs/reference#image-formatting for more information about image
	// formatting.
	Image string `json:"image"`
}

// ScheduledEventParams are the parameters allowed for creating or updating a [ScheduledEvent].
// https://discord.com/developers/docs/resources/guild-scheduled-event#create-guild-scheduled-event
//
// See [ScheduledEvent] for the documentations of the fields.
type ScheduledEventParams struct {
	ChannelID          uint64                        `json:"channel_id,omitempty,string"`
	Name               string                        `json:"name,omitempty"`
	Description        string                        `json:"description,omitempty"`
	ScheduledStartTime *time.Time                    `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   *time.Time                    `json:"scheduled_end_time,omitempty"`
	PrivacyLevel       ScheduledEventPrivacyLevel    `json:"privacy_level,omitempty"`
	Status             ScheduledEventStatus          `json:"status,omitempty"`
	EntityType         types.ScheduledEventEntity    `json:"entity_type,omitempty"`
	EntityMetadata     *ScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	// NOTE: See https://discord.com/developers/docs/reference#image-formatting for more information about image
	// formatting.
	Image string `json:"image,omitempty"`
}

// I have commented this, because it could be an infinite recursive function.
/*
func (p ScheduledEventParams) MarshalJSON() ([]byte, error) {
	if p.EntityType == types.ScheduledEventEntityExternal && p.ChannelID == "" {
		return json.Marshal(struct {
			ScheduledEventParams
			ChannelID json.RawMessage `json:"channel_id"`
		}{
			ScheduledEventParams: p,
			ChannelID:            json.RawMessage("null"),
		})
	}
	return json.Marshal(p)
}
*/

// ScheduledEventEntityMetadata holds additional metadata for [ScheduledEvent].
type ScheduledEventEntityMetadata struct {
	// Location of the ScheduledEvent (1-100 characters).
	//
	// Required for [types.ScheduledEventEntityExternal].
	Location string `json:"location"`
}

// ScheduledEventPrivacyLevel is the privacy level of a [ScheduledEvent].
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
type ScheduledEventPrivacyLevel int

const (
	// ScheduledEventPrivacyLevelGuildOnly makes the scheduled event is only accessible to [Guild] [user.Member]s.
	ScheduledEventPrivacyLevelGuildOnly ScheduledEventPrivacyLevel = 2
)

// ScheduledEventStatus is the status of a [ScheduledEvent].
//
// Valid [ScheduledEventStatus] Transitions :
//
//	SCHEDULED --> ACTIVE --> COMPLETED
//	SCHEDULED --> CANCELED
//
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-status
type ScheduledEventStatus int

const (
	// ScheduledEventStatusScheduled represents the current event is in scheduled state.
	ScheduledEventStatusScheduled ScheduledEventStatus = 1
	// ScheduledEventStatusActive represents the current event is in active state.
	ScheduledEventStatusActive ScheduledEventStatus = 2
	// ScheduledEventStatusCompleted represents the current event is in completed state.
	ScheduledEventStatusCompleted ScheduledEventStatus = 3
	// ScheduledEventStatusCanceled represents the current event is in canceled state.
	ScheduledEventStatusCanceled ScheduledEventStatus = 4
)

// ScheduledEventUser is a [user.User] subscribed to a [ScheduledEvent].
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-user-object
type ScheduledEventUser struct {
	GuildScheduledEventID uint64       `json:"guild_scheduled_event_id,string"`
	User                  *user.User   `json:"user"`
	Member                *user.Member `json:"member"`
}

// ListScheduledEvents returns all [ScheduledEvent]s for the [Guild].
//
// userCount indicates whether to include the user count in the response.
func ListScheduledEvents(guildID uint64, userCount bool) Request[[]*ScheduledEvent] {
	uri := discord.EndpointGuildScheduledEvents(guildID)
	if userCount {
		uri += "?with_user_count=true"
	}

	return NewData[[]*ScheduledEvent](http.MethodGet, uri)
}

// GetScheduledEvent returns a specific [ScheduledEvent] in a [Guild].
//
// userCount indicates whether to include the user count in the response.
func GetScheduledEvent(guildID, eventID uint64, userCount bool) Request[*ScheduledEvent] {
	uri := discord.EndpointGuildScheduledEvent(guildID, eventID)
	if userCount {
		uri += "?with_user_count=true"
	}

	return NewData[*ScheduledEvent](http.MethodGet, uri)
}

// CreateScheduledEvent for a [Guild] and returns it.
func CreateScheduledEvent(guildID uint64, event *ScheduledEventParams) Request[*ScheduledEvent] {
	return NewData[*ScheduledEvent](http.MethodPost, discord.EndpointGuildScheduledEvents(guildID)).
		WithData(event)
}

// EditScheduledEvent for a [Guild] and returns it.
func EditScheduledEvent(guildID, eventID uint64, event *ScheduledEventParams) Request[*ScheduledEvent] {
	return NewData[*ScheduledEvent](http.MethodPatch, discord.EndpointGuildScheduledEvent(guildID, eventID)).
		WithBucketID(discord.EndpointGuildScheduledEvent(guildID, 0)).WithData(event)
}

// DeleteScheduledEvent in a [Guild].
func DeleteScheduledEvent(guildID, eventID uint64) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointGuildScheduledEvent(guildID, eventID)).
		WithBucketID(discord.EndpointGuildScheduledEvent(guildID, 0))
	return WrapAsEmpty(req)
}

// ListScheduledEventUsers returns an array of [ScheduledEventUser] for a particular [ScheduledEvent] in a [Guild].
//
// limit is the maximum number of users to return (max 100).
// withMember indicates whether to include the member object in the response.
// If is not empty all returned users entries will be before beforeID.
// If is not empty all returned users entries will be after afterID.
func ListScheduledEventUsers(guildID, eventID uint64, limit int, withMember bool, beforeID, afterID string) Request[[]*ScheduledEventUser] {
	uri := discord.EndpointGuildScheduledEventUsers(guildID, eventID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}
	if beforeID != "" {
		queryParams.Set("before", beforeID)
	}
	if afterID != "" {
		queryParams.Set("after", afterID)
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	return NewData[[]*ScheduledEventUser](http.MethodPost, uri).
		WithBucketID(discord.EndpointGuildScheduledEventUsers(guildID, 0))
}

package guildapi

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/guild"
)

// ScheduledEvents returns an array of guild.ScheduledEvent for a guild.Guild.
//
// userCount indicates whether to include the user count in the response.
func (r Requester) ScheduledEvents(guildID string, userCount bool) Request[[]*ScheduledEvent] {
	uri := discord.EndpointGuildScheduledEvents(guildID)
	if userCount {
		uri += "?with_user_count=true"
	}

	return NewSimpleData[[]*ScheduledEvent](r, http.MethodGet, uri)
}

// ScheduledEvent returns a specific guild.ScheduledEvent in a guild.Guild.
//
// userCount indicates whether to include the user count in the response.
func (r Requester) ScheduledEvent(guildID, eventID string, userCount bool) Request[*ScheduledEvent] {
	uri := discord.EndpointGuildScheduledEvent(guildID, eventID)
	if userCount {
		uri += "?with_user_count=true"
	}

	return NewSimpleData[*ScheduledEvent](r, http.MethodGet, uri)
}

// ScheduledEventCreate creates a guild.ScheduledEvent for a guild.Guild and returns it.
func (r Requester) ScheduledEventCreate(guildID string, event *ScheduledEventParams) Request[*ScheduledEvent] {
	return NewSimpleData[*ScheduledEvent](
		r, http.MethodPost, discord.EndpointGuildScheduledEvents(guildID),
	).WithData(event)
}

// ScheduledEventEdit updates a guild.ScheduledEvent for a guild.Guild and returns it.
func (r Requester) ScheduledEventEdit(guildID, eventID string, event *ScheduledEventParams) Request[*ScheduledEvent] {
	return NewSimpleData[*ScheduledEvent](
		r, http.MethodPatch, discord.EndpointGuildScheduledEvent(guildID, eventID),
	).WithBucketID(discord.EndpointGuildScheduledEvent(guildID, "")).WithData(event)
}

// ScheduledEventDelete deletes a specific guild.ScheduledEvent in a guild.Guild.
func (r Requester) ScheduledEventDelete(guildID, eventID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointGuildScheduledEvent(guildID, eventID),
	).WithBucketID(discord.EndpointGuildScheduledEvent(guildID, eventID))
	return WrapAsEmpty(req)
}

// ScheduledEventUsers returns an array of guild.ScheduledEventUser for a particular event in a guild.Guild.
//
// limit is the maximum number of users to return (max 100).
// withMember indicates whether to include the member object in the response.
// If is not empty all returned users entries will be before beforeID.
// If is not empty all returned users entries will be after afterID.
func (r Requester) ScheduledEventUsers(guildID, eventID string, limit int, withMember bool, beforeID, afterID string) Request[[]*ScheduledEventUser] {
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

	return NewSimpleData[[]*ScheduledEventUser](
		r, http.MethodPost, uri,
	).WithBucketID(discord.EndpointGuildScheduledEventUsers(guildID, ""))
}

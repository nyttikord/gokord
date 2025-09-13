package guildapi

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/guild"
)

// ScheduledEvents returns an array of guild.ScheduledEvent for a guild.Guild.
//
// userCount indicates whether to include the user count in the response.
func (r Requester) ScheduledEvents(guildID string, userCount bool, options ...discord.RequestOption) ([]*guild.ScheduledEvent, error) {
	uri := discord.EndpointGuildScheduledEvents(guildID)
	if userCount {
		uri += "?with_user_count=true"
	}

	body, err := r.Request(http.MethodGet, uri, nil, options...)
	if err != nil {
		return nil, err
	}

	var se []*guild.ScheduledEvent
	return se, r.Unmarshal(body, &se)
}

// ScheduledEvent returns a specific guild.ScheduledEvent in a guild.Guild.
//
// userCount indicates whether to include the user count in the response.
func (r Requester) ScheduledEvent(guildID, eventID string, userCount bool, options ...discord.RequestOption) (*guild.ScheduledEvent, error) {
	uri := discord.EndpointGuildScheduledEvent(guildID, eventID)
	if userCount {
		uri += "?with_user_count=true"
	}

	body, err := r.Request(http.MethodGet, uri, nil, options...)
	if err != nil {
		return nil, err
	}

	var se *guild.ScheduledEvent
	return se, r.Unmarshal(body, &se)
}

// ScheduledEventCreate creates a guild.ScheduledEvent for a guild.Guild and returns it.
func (r Requester) ScheduledEventCreate(guildID string, event *guild.ScheduledEventParams, options ...discord.RequestOption) (*guild.ScheduledEvent, error) {
	body, err := r.Request(http.MethodPost, discord.EndpointGuildScheduledEvents(guildID), event, options...)
	if err != nil {
		return nil, err
	}

	var se *guild.ScheduledEvent
	return se, r.Unmarshal(body, &se)
}

// ScheduledEventEdit updates a guild.ScheduledEvent for a guild.Guild and returns it.
func (r Requester) ScheduledEventEdit(guildID, eventID string, event *guild.ScheduledEventParams, options ...discord.RequestOption) (*guild.ScheduledEvent, error) {
	body, err := r.Request(http.MethodPatch, discord.EndpointGuildScheduledEvent(guildID, eventID), event, options...)
	if err != nil {
		return nil, err
	}

	var se *guild.ScheduledEvent
	return se, r.Unmarshal(body, &se)
}

// ScheduledEventDelete deletes a specific guild.ScheduledEvent in a guild.Guild.
func (r Requester) ScheduledEventDelete(guildID, eventID string, options ...discord.RequestOption) error {
	_, err := r.Request(http.MethodDelete, discord.EndpointGuildScheduledEvent(guildID, eventID), nil, options...)
	return err
}

// ScheduledEventUsers returns an array of guild.ScheduledEventUser for a particular event in a guild.Guild.
//
// limit is the maximum number of users to return (max 100).
// withMember indicates whether to include the member object in the response.
// If is not empty all returned users entries will be before beforeID.
// If is not empty all returned users entries will be after afterID.
func (r Requester) ScheduledEventUsers(guildID, eventID string, limit int, withMember bool, beforeID, afterID string, options ...discord.RequestOption) ([]*guild.ScheduledEventUser, error) {
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

	body, err := r.Request("GET", uri, nil, options...)
	if err != nil {
		return nil, err
	}

	var seu []*guild.ScheduledEventUser
	return seu, r.Unmarshal(body, &seu)
}

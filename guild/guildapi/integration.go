package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/user"
)

// Integrations returns the list of user.Integration for a guild.Guild.
func (r Requester) Integrations(guildID string) Request[[]*user.Integration] {
	return NewData[[]*user.Integration](
		r, http.MethodGet, discord.EndpointGuildIntegrations(guildID),
	)
}

// IntegrationCreate creates a guild.Guild user.Integration.
func (r Requester) IntegrationCreate(guildID, integrationType, integrationID string) Request[*user.Integration] {
	data := struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}{integrationType, integrationID}

	return NewData[*user.Integration](
		r, http.MethodPost, discord.EndpointGuildIntegrations(guildID),
	).WithData(data)
}

// IntegrationEdit edits a guild.Guild user.Integration.
//
// expireBehavior is the behavior when a user.Integration subscription lapses.
// expireGracePeriod is the period (in seconds) where the user.Integration will ignore lapsed subscriptions.
// enableEmoticons is true if emoticons should be synced for this user.Integration (twitch only currently).
func (r Requester) IntegrationEdit(guildID, integrationID string, expireBehavior, expireGracePeriod int, enableEmoticons bool) Request[*user.Integration] {
	data := struct {
		ExpireBehavior    int  `json:"expire_behavior"`
		ExpireGracePeriod int  `json:"expire_grace_period"`
		EnableEmoticons   bool `json:"enable_emoticons"`
	}{expireBehavior, expireGracePeriod, enableEmoticons}

	return NewData[*user.Integration](
		r, http.MethodPatch, discord.EndpointGuildIntegration(guildID, integrationID),
	).WithBucketID(discord.EndpointGuildIntegrations(guildID)).WithData(data)
}

// IntegrationDelete removes the user.Integration from the guild.Guild.
func (r Requester) IntegrationDelete(guildID, integrationID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointGuildIntegration(guildID, integrationID),
	).WithBucketID(discord.EndpointGuildIntegrations(guildID))
	return WrapAsEmpty(req)
}

package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/user"
)

// Integrations returns the list of user.Integration for a guild.Guild.
func (r Requester) Integrations(guildID string, options ...discord.RequestOption) ([]*user.Integration, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuildIntegrations(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var st []*user.Integration
	return st, r.Unmarshal(body, &st)
}

// IntegrationCreate creates a guild.Guild user.Integration.
func (r Requester) IntegrationCreate(guildID, integrationType, integrationID string, options ...discord.RequestOption) error {
	data := struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}{integrationType, integrationID}

	_, err := r.Request(http.MethodPost, discord.EndpointGuildIntegrations(guildID), data, options...)
	return err
}

// IntegrationEdit edits a guild.Guild user.Integration.
//
// expireBehavior is the behavior when a user.Integration subscription lapses.
// expireGracePeriod is the period (in seconds) where the user.Integration will ignore lapsed subscriptions.
// enableEmoticons is true if emoticons should be synced for this user.Integration (twitch only currently).
func (r Requester) IntegrationEdit(guildID, integrationID string, expireBehavior, expireGracePeriod int, enableEmoticons bool, options ...discord.RequestOption) error {
	data := struct {
		ExpireBehavior    int  `json:"expire_behavior"`
		ExpireGracePeriod int  `json:"expire_grace_period"`
		EnableEmoticons   bool `json:"enable_emoticons"`
	}{expireBehavior, expireGracePeriod, enableEmoticons}

	_, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildIntegration(guildID, integrationID),
		data,
		discord.EndpointGuildIntegration(guildID, ""),
		options...,
	)
	return err
}

// IntegrationDelete removes the user.Integration from the guild.Guild.
func (r Requester) IntegrationDelete(guildID, integrationID string, options ...discord.RequestOption) (err error) {
	_, err = r.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildIntegration(guildID, integrationID),
		nil,
		discord.EndpointGuildIntegration(guildID, ""),
		options...,
	)
	return
}

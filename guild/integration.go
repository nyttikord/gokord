package guild

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/user"
)

// ListIntegrations returns the list of user.Integration for a guild.Guild.
func ListIntegrations(guildID uint64) Request[[]*user.Integration] {
	return NewData[[]*user.Integration](http.MethodGet, discord.EndpointGuildIntegrations(guildID))
}

// CreateIntegration in the given [Guild].
func CreateIntegration(guildID uint64, integrationType string, integrationID uint64) Request[*user.Integration] {
	data := struct {
		Type string `json:"type"`
		ID   uint64 `json:"id,string"`
	}{integrationType, integrationID}

	return NewData[*user.Integration](http.MethodPost, discord.EndpointGuildIntegrations(guildID)).
		WithData(data)
}

// EditIntegration in the given [Guild].
//
// expireBehavior is the behavior when a user.Integration subscription lapses.
// expireGracePeriod is the period (in seconds) where the user.Integration will ignore lapsed subscriptions.
// enableEmoticons is true if emoticons should be synced for this user.Integration (twitch only currently).
func EditIntegration(guildID, integrationID uint64, expireBehavior, expireGracePeriod int, enableEmoticons bool) Request[*user.Integration] {
	data := struct {
		ExpireBehavior    int  `json:"expire_behavior"`
		ExpireGracePeriod int  `json:"expire_grace_period"`
		EnableEmoticons   bool `json:"enable_emoticons"`
	}{expireBehavior, expireGracePeriod, enableEmoticons}

	return NewData[*user.Integration](http.MethodPatch, discord.EndpointGuildIntegration(guildID, integrationID)).
		WithBucketID(discord.EndpointGuildIntegrations(guildID)).WithData(data)
}

// DeleteIntegration from the guild.Guild.
func DeleteIntegration(guildID, integrationID uint64) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointGuildIntegration(guildID, integrationID)).
		WithBucketID(discord.EndpointGuildIntegrations(guildID))
	return WrapAsEmpty(req)
}

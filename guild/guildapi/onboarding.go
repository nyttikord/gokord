package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/guild"
)

// Onboarding returns guild.Onboarding configuration of a guild.Guild.
func (r Requester) Onboarding(guildID string, options ...discord.RequestOption) (*guild.Onboarding, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuildOnboarding(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var ob guild.Onboarding
	return &ob, r.Unmarshal(body, &ob)
}

// OnboardingEdit edits guild.Onboarding configuration of a guild.Guild.
func (r Requester) OnboardingEdit(guildID string, o *guild.Onboarding, options ...discord.RequestOption) (*guild.Onboarding, error) {
	body, err := r.Request("PUT", discord.EndpointGuildOnboarding(guildID), o, options...)
	if err != nil {
		return nil, err
	}

	var ob guild.Onboarding
	return &ob, r.Unmarshal(body, &ob)
}

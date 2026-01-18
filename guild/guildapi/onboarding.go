package guildapi

import (
	"context"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/guild"
)

// Onboarding returns guild.Onboarding configuration of a guild.Guild.
func (r Requester) Onboarding(ctx context.Context, guildID string, options ...discord.RequestOption) (*guild.Onboarding, error) {
	body, err := r.Request(ctx, http.MethodGet, discord.EndpointGuildOnboarding(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var ob guild.Onboarding
	return &ob, r.Unmarshal(body, &ob)
}

// OnboardingEdit edits guild.Onboarding configuration of a guild.Guild.
func (r Requester) OnboardingEdit(ctx context.Context, guildID string, o *guild.Onboarding, options ...discord.RequestOption) (*guild.Onboarding, error) {
	body, err := r.Request(ctx, http.MethodPut, discord.EndpointGuildOnboarding(guildID), o, options...)
	if err != nil {
		return nil, err
	}

	var ob guild.Onboarding
	return &ob, r.Unmarshal(body, &ob)
}

package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/guild"
)

// Onboarding returns guild.Onboarding configuration of a guild.Guild.
func (r Requester) Onboarding(guildID string) Request[*Onboarding] {
	return NewSimpleData[*Onboarding](
		r, http.MethodGet, discord.EndpointGuildOnboarding(guildID),
	)
}

// OnboardingEdit edits guild.Onboarding configuration of a guild.Guild.
func (r Requester) OnboardingEdit(guildID string, o *Onboarding) Request[*Onboarding] {
	return NewSimpleData[*Onboarding](
		r, http.MethodPut, discord.EndpointGuildOnboarding(guildID),
	).WithData(o)
}

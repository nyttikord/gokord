// Package guildapi contains everything to interact with everything located in the guild package.
package guildapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/guild"
)

var (
	ErrVerificationLevelBounds = errors.New("VerificationLevel out of bounds, should be between 0 and 3")
)

// Guild returns the guild.Guild with the given guildID.
func (r Requester) Guild(guildID string) Request[*Guild] {
	return NewSimpleData[*Guild](
		r, http.MethodGet, discord.EndpointGuild(guildID),
	)
}

// GuildWithCounts returns the guild.Guild with the given guildID with approximate user.Member and status.Presence counts.
func (r Requester) GuildWithCounts(guildID string) Request[*Guild] {
	return NewSimpleData[*Guild](
		r, http.MethodGet, discord.EndpointGuild(guildID)+"?with_counts=true",
	)
}

// GuildPreview returns the Preview for the given public Guild guildID.
func (r Requester) GuildPreview(guildID string) Request[*Preview] {
	return NewSimpleData[*Preview](
		r, http.MethodGet, discord.EndpointGuildPreview(guildID),
	)
}

// GuildEdit edits a guild.Guild with the given params.
func (r Requester) GuildEdit(guildID string, params *Params) Request[*Guild] {
	// Bounds checking for regions
	if params.Region != "" {
		isValid := false
		regions, _ := r.VoiceRegions()
		for _, r := range regions {
			if params.Region == r.ID {
				isValid = true
			}
		}
		if !isValid {
			var valid []string
			for _, r := range regions {
				valid = append(valid, r.ID)
			}
			return nil, fmt.Errorf("not a valid region (%q)", valid)
		}
	}

	body, err := r.Request(ctx, http.MethodPatch, discord.EndpointGuild(guildID), params, options...)
	if err != nil {
		return nil, err
	}

	var g Guild
	return &g, r.Unmarshal(body, &g)
}

// GuildDelete deletes a guild.Guild.
func (r Requester) GuildDelete(guildID string) EmptyRequest {
	req := NewSimple(r, http.MethodDelete, discord.EndpointGuild(guildID))
	return WrapAsEmpty(req)
}

// GuildLeave leaves a guild.Guild.
func (r Requester) GuildLeave(guildID string) EmptyRequest {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointUserGuild("@e", guildID),
	).WithBucketID(discord.EndpointUserGuild("", guildID))
	return WrapAsEmpty(req)
}

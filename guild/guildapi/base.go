// Package guildapi contains everything to interact with everything located in the guild package.
package guildapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/guild"
)

var (
	ErrVerificationLevelBounds = errors.New("VerificationLevel out of bounds, should be between 0 and 3")
	ErrInvalidVoiceRegions     = errors.New("invalid voice regions")
)

// Guild returns the guild.Guild with the given guildID.
func (r Requester) Guild(guildID string) Request[*Guild] {
	return NewData[*Guild](
		r, http.MethodGet, discord.EndpointGuild(guildID),
	)
}

// GuildWithCounts returns the guild.Guild with the given guildID with approximate user.Member and status.Presence counts.
func (r Requester) GuildWithCounts(guildID string) Request[*Guild] {
	return NewData[*Guild](
		r, http.MethodGet, discord.EndpointGuild(guildID)+"?with_counts=true",
	)
}

// GuildPreview returns the Preview for the given public Guild guildID.
func (r Requester) GuildPreview(guildID string) Request[*Preview] {
	return NewData[*Preview](
		r, http.MethodGet, discord.EndpointGuildPreview(guildID),
	)
}

// GuildEdit edits a guild.Guild with the given params.
func (r Requester) GuildEdit(guildID string, params *Params) Request[*Guild] {
	return NewData[*Guild](r, http.MethodPatch, discord.EndpointGuild(guildID)).
		WithData(params).
		WithPre(func(ctx context.Context, do *Do) error {
			if params.Region == "" {
				return nil
			}
			valid := false
			regions, err := r.VoiceRegions().Do(ctx)
			if err != nil {
				return err
			}
			for _, r := range regions {
				if params.Region == r.ID {
					valid = true
				}
			}
			if valid {
				return nil
			}
			var validRegions []string
			for _, r := range regions {
				validRegions = append(validRegions, r.ID)
			}
			return errors.Join(
				ErrInvalidVoiceRegions, fmt.Errorf("%s is not a voice region (%q)", params.Region, validRegions),
			)
		})
}

// GuildDelete deletes a guild.Guild.
func (r Requester) GuildDelete(guildID string) Empty {
	req := NewSimple(r, http.MethodDelete, discord.EndpointGuild(guildID))
	return WrapAsEmpty(req)
}

// GuildLeave leaves a guild.Guild.
func (r Requester) GuildLeave(guildID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointUserGuild("@e", guildID),
	).WithBucketID(discord.EndpointUserGuild("", guildID))
	return WrapAsEmpty(req)
}

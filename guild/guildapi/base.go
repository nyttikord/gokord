// Package guildapi contains everything to interact with everything located in the guild package.
package guildapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/guild"
)

var (
	ErrVerificationLevelBounds = errors.New("VerificationLevel out of bounds, should be between 0 and 3")
)

// Guild returns the guild.Guild with the given guildID.
func (r Requester) Guild(guildID string, options ...discord.RequestOption) (*guild.Guild, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuild(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var g guild.Guild
	return &g, r.Unmarshal(body, &g)
}

// GuildWithCounts returns the guild.Guild with the given guildID with approximate user.Member and status.Presence counts.
func (r Requester) GuildWithCounts(guildID string, options ...discord.RequestOption) (*guild.Guild, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuild(guildID)+"?with_counts=true", nil, options...)
	if err != nil {
		return nil, err
	}

	var g guild.Guild
	return &g, r.Unmarshal(body, &g)
}

// GuildPreview returns the guild.Preview for the given public guild.Guild guildID.
func (r Requester) GuildPreview(guildID string, options ...discord.RequestOption) (*guild.Preview, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuildPreview(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var gp guild.Preview
	return &gp, r.Unmarshal(body, &gp)
}

// GuildEdit edits a guild.Guild with the given params.
func (r Requester) GuildEdit(guildID string, params *guild.Params, options ...discord.RequestOption) (*guild.Guild, error) {
	// Bounds checking for VerificationLevel, interval: [0, 4]
	if params.VerificationLevel != nil {
		val := *params.VerificationLevel
		if val < 0 || val > 4 {
			return nil, ErrVerificationLevelBounds
		}
	}

	// Bounds checking for regions
	if params.Region != "" {
		isValid := false
		regions, _ := r.VoiceAPI().VoiceRegions(options...)
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

	body, err := r.Request(http.MethodPatch, discord.EndpointGuild(guildID), params, options...)
	if err != nil {
		return nil, err
	}

	var g guild.Guild
	return &g, r.Unmarshal(body, &g)
}

// GuildDelete deletes a guild.Guild.
func (r Requester) GuildDelete(guildID string, options ...discord.RequestOption) error {
	_, err := r.Request(http.MethodDelete, discord.EndpointGuild(guildID), nil, options...)
	return err
}

// GuildLeave leaves a guild.Guild.
func (r Requester) GuildLeave(guildID string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointUserGuild("@me", guildID),
		nil,
		discord.EndpointUserGuild("", guildID),
		options...,
	)
	return err
}

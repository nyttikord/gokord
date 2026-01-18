package guildapi

import (
	"context"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/guild"
)

// AutoModerationRules returns a list of guild.AutoModerationRule in the given guild.Guild.
func (r Requester) AutoModerationRules(ctx context.Context, guildID string, options ...discord.RequestOption) ([]*guild.AutoModerationRule, error) {
	body, err := r.Request(ctx, http.MethodGet, discord.EndpointGuildAutoModerationRules(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var rules []*guild.AutoModerationRule
	return rules, r.Unmarshal(body, &rules)
}

// AutoModerationRule returns a guild.AutoModerationRule.
func (r Requester) AutoModerationRule(ctx context.Context, guildID, ruleID string, options ...discord.RequestOption) (*guild.AutoModerationRule, error) {
	body, err := r.Request(ctx, http.MethodGet, discord.EndpointGuildAutoModerationRule(guildID, ruleID), nil, options...)
	if err != nil {
		return nil, err
	}

	var rule *guild.AutoModerationRule
	return rule, r.Unmarshal(body, &rule)
}

// AutoModerationRuleCreate creates a guild.AutoModerationRule and returns it.
func (r Requester) AutoModerationRuleCreate(ctx context.Context, guildID string, rule *guild.AutoModerationRule, options ...discord.RequestOption) (*guild.AutoModerationRule, error) {
	body, err := r.Request(ctx, http.MethodPost, discord.EndpointGuildAutoModerationRules(guildID), rule, options...)
	if err != nil {
		return nil, err
	}

	var rl *guild.AutoModerationRule
	return rl, r.Unmarshal(body, &rl)
}

// AutoModerationRuleEdit edits and returns the updated guild.AutoModerationRule.
func (r Requester) AutoModerationRuleEdit(ctx context.Context, guildID, ruleID string, rule *guild.AutoModerationRule, options ...discord.RequestOption) (*guild.AutoModerationRule, error) {
	body, err := r.Request(ctx, http.MethodPatch, discord.EndpointGuildAutoModerationRule(guildID, ruleID), rule, options...)
	if err != nil {
		return nil, err
	}

	var rl *guild.AutoModerationRule
	return rl, r.Unmarshal(body, &rl)
}

// AutoModerationRuleDelete deletes a guild.AutoModerationRule.
func (r Requester) AutoModerationRuleDelete(ctx context.Context, guildID, ruleID string, options ...discord.RequestOption) error {
	_, err := r.Request(ctx, http.MethodDelete, discord.EndpointGuildAutoModerationRule(guildID, ruleID), nil, options...)
	return err
}

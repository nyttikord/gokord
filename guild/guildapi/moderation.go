package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/guild"
)

// AutoModerationRules returns a list of guild.AutoModerationRule in the given guild.Guild.
func (r Requester) AutoModerationRules(guildID string) Request[[]*AutoModerationRule] {
	return NewData[[]*AutoModerationRule](
		r, http.MethodGet, discord.EndpointGuildAutoModerationRules(guildID),
	)
}

// AutoModerationRule returns a guild.AutoModerationRule in the guild.Guild.
func (r Requester) AutoModerationRule(guildID, ruleID string) Request[*AutoModerationRule] {
	return NewData[*AutoModerationRule](
		r, http.MethodGet, discord.EndpointGuildAutoModerationRule(guildID, ruleID),
	).WithBucketID(discord.EndpointGuildAutoModerationRules(guildID))
}

// AutoModerationRuleCreate creates a AutoModerationRule and returns it.
func (r Requester) AutoModerationRuleCreate(guildID string, rule *AutoModerationRule) Request[*AutoModerationRule] {
	return NewData[*AutoModerationRule](
		r, http.MethodGet, discord.EndpointGuildAutoModerationRules(guildID),
	).WithData(rule)
}

// AutoModerationRuleEdit edits and returns the updated AutoModerationRule.
func (r Requester) AutoModerationRuleEdit(guildID, ruleID string, rule *AutoModerationRule) Request[*AutoModerationRule] {
	return NewData[*AutoModerationRule](
		r, http.MethodPatch, discord.EndpointGuildAutoModerationRule(guildID, ruleID),
	).WithBucketID(discord.EndpointGuildAutoModerationRules(guildID)).WithData(rule)
}

// AutoModerationRuleDelete deletes a AutoModerationRule.
func (r Requester) AutoModerationRuleDelete(guildID, ruleID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointGuildAutoModerationRule(guildID, ruleID),
	).WithBucketID(discord.EndpointGuildAutoModerationRules(guildID))
	return WrapAsEmpty(req)
}

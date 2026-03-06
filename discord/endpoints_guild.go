package discord

import "fmt"

func EndpointGuild(gID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointGuilds, gID)
}

func EndpointGuildAutoModeration(gID uint64) string {
	return EndpointGuild(gID) + "/auto-moderation"
}
func EndpointGuildAutoModerationRules(gID uint64) string {
	return EndpointGuildAutoModeration(gID) + "/rules"
}
func EndpointGuildAutoModerationRule(gID, rID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointGuildAutoModerationRules(gID), rID)
}

func EndpointGuildThreads(gID uint64) string {
	return EndpointGuild(gID) + "/threads"
}
func EndpointGuildActiveThreads(gID uint64) string {
	return EndpointGuildThreads(gID) + "/active"
}

func EndpointGuildPreview(gID uint64) string {
	return EndpointGuild(gID) + "/preview"
}

func EndpointGuildChannels(gID uint64) string {
	return EndpointGuild(gID) + "/channels"
}

func EndpointGuildMembers(gID uint64) string {
	return EndpointGuild(gID) + "/members"
}
func EndpointGuildMembersSearch(gID uint64) string {
	return EndpointGuildMembers(gID) + "/search"
}
func EndpointGuildMember(gID, uID uint64) string {
	return fmt.Sprintf("%s/members/%d", EndpointGuild(gID), uID)
}
func EndpointGuildMemberRole(gID, uID, rID uint64) string {
	return fmt.Sprintf("%s/roles/%d", EndpointGuildMember(gID, uID), rID)
}
func EndpointGuildBans(gID uint64) string {
	return EndpointGuild(gID) + "/bans"
}
func EndpointGuildBan(gID, uID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointGuildBans(gID), uID)
}

func EndpointGuildIntegrations(gID uint64) string {
	return EndpointGuild(gID) + "/integrations"
}
func EndpointGuildIntegration(gID, iID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointGuildIntegrations(gID), iID)
}

func EndpointGuildRoles(gID uint64) string {
	return EndpointGuild(gID) + "/roles"
}
func EndpointGuildRole(gID, rID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointGuildRoles(gID), rID)
}
func EndpointGuildRoleMemberCounts(gID uint64) string {
	return EndpointGuildRoles(gID) + "/member-counts"
}

func EndpointGuildInvites(gID uint64) string {
	return EndpointGuild(gID) + "/invites"
}

func EndpointGuildWidget(gID uint64) string {
	return EndpointGuild(gID) + "/widget"
}

var EndpointGuildEmbed = EndpointGuildWidget

func EndpointGuildPrune(gID uint64) string {
	return EndpointGuild(gID) + "/prune"
}

func EndpointGuildWebhooks(gID uint64) string {
	return EndpointGuild(gID) + "/webhooks"
}
func EndpointGuildAuditLogs(gID uint64) string {
	return EndpointGuild(gID) + "/audit-logs"
}
func EndpointGuildEmojis(gID uint64) string {
	return EndpointGuild(gID) + "/emojis"
}
func EndpointGuildEmoji(gID, eID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointGuildEmojis(gID), eID)
}
func EndpointGuildStickers(gID uint64) string {
	return EndpointGuild(gID) + "/stickers"
}
func EndpointGuildSticker(gID, sID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointGuildStickers(gID), sID)
}

func EndpointStageInstance(cID string) string {
	return EndpointStageInstances + "/" + cID
}
func EndpointGuildScheduledEvents(gID uint64) string {
	return EndpointGuild(gID) + "/scheduled-events"
}
func EndpointGuildScheduledEvent(gID, eID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointGuildScheduledEvents(gID), eID)
}
func EndpointGuildScheduledEventUsers(gID, eID uint64) string {
	return EndpointGuildScheduledEvent(gID, eID) + "/users"
}

func EndpointGuildOnboarding(gID uint64) string {
	return EndpointGuild(gID) + "/onboarding"
}

func EndpointGuildTemplate(tID string) string {
	return fmt.Sprintf("%s/templates/%d", EndpointGuilds, tID)
}
func EndpointGuildTemplates(gID uint64) string {
	return EndpointGuild(gID) + "/templates"
}
func EndpointGuildTemplateSync(gID uint64, tID string) string {
	return fmt.Sprintf("%s/%s", EndpointGuildTemplates(gID), tID)
}

func EndpointGuildSoundboardSounds(gId uint64) string {
	return EndpointGuild(gId) + "/soundboard-sounds"
}
func EndpointGuildSoundboardSound(gId, sId uint64) string {
	return fmt.Sprintf("%s/%d", EndpointGuildSoundboardSounds(gId), sId)
}

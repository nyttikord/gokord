package guild

import (
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// A Ban stores data for a Guild ban.
type Ban struct {
	Reason string     `json:"reason"`
	User   *user.User `json:"user"`
}

// AutoModerationRuleEvent indicates in what event context a guild.AutoModerationRule should be checked.
type AutoModerationRuleEvent int

const (
	// AutoModerationRuleEventMessageSend is checked when a member sends or edits a message in the guild
	AutoModerationRuleEventMessageSend AutoModerationRuleEvent = 1
)

// AutoModerationRuleTrigger represents the type of content which can trigger the guild.AutoModerationRule.
type AutoModerationRuleTrigger int

const (
	AutoModerationRuleTriggerKeyword       AutoModerationRuleTrigger = 1
	AutoModerationRuleTriggerHarmfulLink   AutoModerationRuleTrigger = 2
	AutoModerationRuleTriggerSpam          AutoModerationRuleTrigger = 3
	AutoModerationRuleTriggerKeywordPreset AutoModerationRuleTrigger = 4
)

// AutoModerationRule stores data for an auto moderation rule.
type AutoModerationRule struct {
	ID              string                         `json:"id,omitempty"`
	GuildID         string                         `json:"guild_id,omitempty"`
	Name            string                         `json:"name,omitempty"`
	CreatorID       string                         `json:"creator_id,omitempty"`
	EventType       AutoModerationRuleEvent        `json:"event_type,omitempty"`
	TriggerType     AutoModerationRuleTrigger      `json:"trigger_type,omitempty"`
	TriggerMetadata *AutoModerationTriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []AutoModerationAction         `json:"actions,omitempty"`
	Enabled         *bool                          `json:"enabled,omitempty"`
	ExemptRoles     *[]string                      `json:"exempt_roles,omitempty"`
	ExemptChannels  *[]string                      `json:"exempt_channels,omitempty"`
}

// AutoModerationKeywordPreset represents an internally pre-defined word set.
type AutoModerationKeywordPreset uint

// Auto moderation keyword presets.
const (
	AutoModerationKeywordPresetProfanity     AutoModerationKeywordPreset = 1
	AutoModerationKeywordPresetSexualContent AutoModerationKeywordPreset = 2
	AutoModerationKeywordPresetSlurs         AutoModerationKeywordPreset = 3
)

// AutoModerationTriggerMetadata represents additional metadata used to determine whether AutoModerationRule should be
// triggered.
type AutoModerationTriggerMetadata struct {
	// Substrings which will be searched for in content.
	// NOTE: should be only used with keyword trigger type.
	KeywordFilter []string `json:"keyword_filter,omitempty"`
	// Regular expression patterns which will be matched against content (maximum of 10).
	// NOTE: should be only used with keyword trigger type.
	RegexPatterns []string `json:"regex_patterns,omitempty"`

	// Internally pre-defined wordsets which will be searched for in content.
	// NOTE: should be only used with keyword preset trigger type.
	Presets []AutoModerationKeywordPreset `json:"presets,omitempty"`

	// Substrings which should not trigger the rule.
	// NOTE: should be only used with keyword or keyword preset trigger type.
	AllowList *[]string `json:"allow_list,omitempty"`

	// Total number of unique role and user mentions allowed per message.
	// NOTE: should be only used with mention spam trigger type.
	MentionTotalLimit int `json:"mention_total_limit,omitempty"`
}

// AutoModerationActionMetadata represents additional metadata needed during execution for a specific
// AutoModerationAction type.
type AutoModerationActionMetadata struct {
	// Channel to which user content should be logged.
	// NOTE: should be only used with send alert message action type.
	ChannelID string `json:"channel_id,omitempty"`

	// Timeout duration in seconds (maximum of 2419200 - 4 weeks).
	// NOTE: should be only used with timeout action type.
	Duration int `json:"duration_seconds,omitempty"`

	// Additional explanation that will be shown to members whenever their message is blocked (maximum of 150 characters).
	// NOTE: should be only used with block message action type.
	CustomMessage string `json:"custom_message,omitempty"`
}

// AutoModerationAction stores data for an auto moderation action.
type AutoModerationAction struct {
	Type     types.AutoModerationAction    `json:"type"`
	Metadata *AutoModerationActionMetadata `json:"metadata,omitempty"`
}

// A AuditLog stores data for a Guild audit Log.
// https://discord.com/developers/docs/resources/audit-log#audit-log-object-audit-log-structure
type AuditLog struct {
	Webhooks        []*channel.Webhook  `json:"webhooks,omitempty"`
	Users           []*user.User        `json:"users,omitempty"`
	AuditLogEntries []*AuditLogEntry    `json:"audit_log_entries"`
	Integrations    []*user.Integration `json:"integrations"`
}

// AuditLogEntry for a AuditLog.
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-entry-structure
type AuditLogEntry struct {
	TargetID   string            `json:"target_id"`
	Changes    []*AuditLogChange `json:"changes"`
	UserID     string            `json:"user_id"`
	ID         string            `json:"id"`
	ActionType *AuditLogAction   `json:"action_type"`
	Options    *AuditLogOptions  `json:"options"`
	Reason     string            `json:"reason"`
}

// AuditLogChange for an AuditLogEntry.
type AuditLogChange struct {
	NewValue interface{}        `json:"new_value"`
	OldValue interface{}        `json:"old_value"`
	Key      *AuditLogChangeKey `json:"key"`
}

// AuditLogChangeKey value for AuditLogChange.
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-key
type AuditLogChangeKey string

const (
	// AuditLogChangeKeyAfkChannelID is sent when afk channel changed (snowflake) - Guild.
	AuditLogChangeKeyAfkChannelID AuditLogChangeKey = "afk_channel_id"
	// AuditLogChangeKeyAfkTimeout is sent when afk timeout duration changed (int) - Guild.
	AuditLogChangeKeyAfkTimeout AuditLogChangeKey = "afk_timeout"
	// AuditLogChangeKeyAllow is sent when a permission on a text or voice channel was allowed for a role (string) - Role.
	AuditLogChangeKeyAllow AuditLogChangeKey = "allow"
	// AuditLogChangeKeyApplicationID is sent when application id of the added or removed webhook or bot (snowflake) -
	// channel.Channel.
	AuditLogChangeKeyApplicationID AuditLogChangeKey = "application_id"
	// AuditLogChangeKeyArchived is sent when thread was archived/unarchived (bool) - thread
	AuditLogChangeKeyArchived AuditLogChangeKey = "archived"
	// AuditLogChangeKeyAsset is sent when asset is changed (string) - emoji.Sticker
	AuditLogChangeKeyAsset AuditLogChangeKey = "asset"
	// AuditLogChangeKeyAutoArchiveDuration is sent when auto archive duration changed (int) - thread
	AuditLogChangeKeyAutoArchiveDuration AuditLogChangeKey = "auto_archive_duration"
	// AuditLogChangeKeyAvailable is sent when availability of sticker changed (bool) - emoji.Sticker
	AuditLogChangeKeyAvailable AuditLogChangeKey = "available"
	// AuditLogChangeKeyAvatarHash is sent when user avatar changed (string) - user.User
	AuditLogChangeKeyAvatarHash AuditLogChangeKey = "avatar_hash"
	// AuditLogChangeKeyBannerHash is sent when guild banner changed (string) - Guild
	AuditLogChangeKeyBannerHash AuditLogChangeKey = "banner_hash"
	// AuditLogChangeKeyBitrate is sent when voice channel bitrate changed (int) - channel.Channel
	AuditLogChangeKeyBitrate AuditLogChangeKey = "bitrate"
	// AuditLogChangeKeyChannelID is sent when channel for invite code or guild scheduled event changed (snowflake) -
	// invite.Invite or ScheduledEvent
	AuditLogChangeKeyChannelID AuditLogChangeKey = "channel_id"
	// AuditLogChangeKeyCode is sent when invite code changed (string) - invite.Invite
	AuditLogChangeKeyCode AuditLogChangeKey = "code"
	// AuditLogChangeKeyColor is sent when role color changed (int) - Role
	AuditLogChangeKeyColor AuditLogChangeKey = "color"
	// AuditLogChangeKeyCommunicationDisabledUntil is sent when member timeout state changed (ISO8601 timestamp) -
	// user.Member
	AuditLogChangeKeyCommunicationDisabledUntil AuditLogChangeKey = "communication_disabled_until"
	// AuditLogChangeKeyDeaf is sent when user server deafened/undeafened (bool) - user.Member
	AuditLogChangeKeyDeaf AuditLogChangeKey = "deaf"
	// AuditLogChangeKeyDefaultAutoArchiveDuration is sent when default auto archive duration for newly created threads changed (int) -
	// channel.Channel
	AuditLogChangeKeyDefaultAutoArchiveDuration AuditLogChangeKey = "default_auto_archive_duration"
	// AuditLogChangeKeyDefaultMessageNotification is sent when default message notification level changed (int) - Guild
	AuditLogChangeKeyDefaultMessageNotification AuditLogChangeKey = "default_message_notifications"
	// AuditLogChangeKeyDeny is sent when a permission on a text or voice channel was denied for a role (string) - Role
	AuditLogChangeKeyDeny AuditLogChangeKey = "deny"
	// AuditLogChangeKeyDescription is sent when description changed (string) - Guild, emoji.Sticker, or ScheduledEvent
	AuditLogChangeKeyDescription AuditLogChangeKey = "description"
	// AuditLogChangeKeyDiscoverySplashHash is sent when discovery splash changed (string) - Guild
	AuditLogChangeKeyDiscoverySplashHash AuditLogChangeKey = "discovery_splash_hash"
	// AuditLogChangeKeyEnableEmoticons is sent when integration emoticons enabled/disabled (bool) - types.IntegrationInstall
	AuditLogChangeKeyEnableEmoticons AuditLogChangeKey = "enable_emoticons"
	// AuditLogChangeKeyEntityType is sent when entity type of guild scheduled event was changed (int) - ScheduledEvent
	AuditLogChangeKeyEntityType AuditLogChangeKey = "entity_type"
	// AuditLogChangeKeyExpireBehavior is sent when integration expiring subscriber behavior changed (int) - types.IntegrationInstall
	AuditLogChangeKeyExpireBehavior AuditLogChangeKey = "expire_behavior"
	// AuditLogChangeKeyExpireGracePeriod is sent when integration expire grace period changed (int) - types.IntegrationInstall
	AuditLogChangeKeyExpireGracePeriod AuditLogChangeKey = "expire_grace_period"
	// AuditLogChangeKeyExplicitContentFilter is sent when change in whose messages are scanned and deleted for explicit
	// content in the server is made (int) - Guild
	AuditLogChangeKeyExplicitContentFilter AuditLogChangeKey = "explicit_content_filter"
	// AuditLogChangeKeyFormatType is sent when format type of sticker changed (int - sticker format type) - emoji.Sticker
	AuditLogChangeKeyFormatType AuditLogChangeKey = "format_type"
	// AuditLogChangeKeyGuildID is sent when guild sticker is in changed (snowflake) - emoji.Sticker
	AuditLogChangeKeyGuildID AuditLogChangeKey = "guild_id"
	// AuditLogChangeKeyHoist is sent when role is now displayed/no longer displayed separate from online users (bool) - Role
	AuditLogChangeKeyHoist AuditLogChangeKey = "hoist"
	// AuditLogChangeKeyIconHash is sent when icon changed (string) - Guild or Role
	AuditLogChangeKeyIconHash AuditLogChangeKey = "icon_hash"
	// AuditLogChangeKeyID is sent when the id of the changed entity - sometimes used in conjunction with other keys (snowflake) - any
	AuditLogChangeKeyID AuditLogChangeKey = "id"
	// AuditLogChangeKeyInvitable is sent when private thread is now invitable/uninvitable (bool) - thread
	AuditLogChangeKeyInvitable AuditLogChangeKey = "invitable"
	// AuditLogChangeKeyInviterID is sent when person who created invite code changed (snowflake) - invite.Invite
	AuditLogChangeKeyInviterID AuditLogChangeKey = "inviter_id"
	// AuditLogChangeKeyLocation is sent when channel id for guild scheduled event changed (string) - ScheduledEvent
	AuditLogChangeKeyLocation AuditLogChangeKey = "location"
	// AuditLogChangeKeyLocked is sent when thread was locked/unlocked (bool) - thread
	AuditLogChangeKeyLocked AuditLogChangeKey = "locked"
	// AuditLogChangeKeyMaxAge is sent when invite code expiration time changed (int) - invite.Invite
	AuditLogChangeKeyMaxAge AuditLogChangeKey = "max_age"
	// AuditLogChangeKeyMaxUses is sent when max number of times invite code can be used changed (int) - invite.Invite
	AuditLogChangeKeyMaxUses AuditLogChangeKey = "max_uses"
	// AuditLogChangeKeyMentionable is sent when role is now mentionable/unmentionable (bool) - Role
	AuditLogChangeKeyMentionable AuditLogChangeKey = "mentionable"
	// AuditLogChangeKeyMfaLevel is sent when two-factor auth requirement changed (int - mfa level) - Guild
	AuditLogChangeKeyMfaLevel AuditLogChangeKey = "mfa_level"
	// AuditLogChangeKeyMute is sent when user server muted/unmuted (bool) - user.Member
	AuditLogChangeKeyMute AuditLogChangeKey = "mute"
	// AuditLogChangeKeyName is sent when name changed (string) - any
	AuditLogChangeKeyName AuditLogChangeKey = "name"
	// AuditLogChangeKeyNick is sent when user nickname changed (string) - user.Member
	AuditLogChangeKeyNick AuditLogChangeKey = "nick"
	// AuditLogChangeKeyNSFW is sent when channel nsfw restriction changed (bool) - channel.Channel
	AuditLogChangeKeyNSFW AuditLogChangeKey = "nsfw"
	// AuditLogChangeKeyOwnerID is sent when owner changed (snowflake) - Guild
	AuditLogChangeKeyOwnerID AuditLogChangeKey = "owner_id"
	// AuditLogChangeKeyPermissionOverwrite is sent when permissions on a channel changed (array of channel overwrite
	// objects) - channel.Channel
	AuditLogChangeKeyPermissionOverwrite AuditLogChangeKey = "permission_overwrites"
	// AuditLogChangeKeyPermissions is sent when permissions for a role changed (string) - Role
	AuditLogChangeKeyPermissions AuditLogChangeKey = "permissions"
	// AuditLogChangeKeyPosition is sent when text or voice channel position changed (int) - channel.Channel
	AuditLogChangeKeyPosition AuditLogChangeKey = "position"
	// AuditLogChangeKeyPreferredLocale is sent when preferred locale changed (string) - Guild
	AuditLogChangeKeyPreferredLocale AuditLogChangeKey = "preferred_locale"
	// AuditLogChangeKeyPrivacylevel is sent when privacy level of the stage instance changed (integer - privacy level) -
	// channel.StageInstance or ScheduledEvent
	AuditLogChangeKeyPrivacylevel AuditLogChangeKey = "privacy_level"
	// AuditLogChangeKeyPruneDeleteDays is sent when number of days after which inactive and role-unassigned members are
	// kicked changed (int) - Guild
	AuditLogChangeKeyPruneDeleteDays AuditLogChangeKey = "prune_delete_days"
	// AuditLogChangeKeyPublicUpdatesChannelID is sent when id of the public updates channel changed (snowflake) - Guild
	AuditLogChangeKeyPublicUpdatesChannelID AuditLogChangeKey = "public_updates_channel_id"
	// AuditLogChangeKeyRateLimitPerUser is sent when amount of seconds a user has to wait before sending another message
	// changed (int) - channel.Channel
	AuditLogChangeKeyRateLimitPerUser AuditLogChangeKey = "rate_limit_per_user"
	// AuditLogChangeKeyRegion is sent when region changed (string) - Guild
	AuditLogChangeKeyRegion AuditLogChangeKey = "region"
	// AuditLogChangeKeyRulesChannelID is sent when id of the rules channel changed (snowflake) - Guild
	AuditLogChangeKeyRulesChannelID AuditLogChangeKey = "rules_channel_id"
	// AuditLogChangeKeySplashHash is sent when invite splash page artwork changed (string) - Guild
	AuditLogChangeKeySplashHash AuditLogChangeKey = "splash_hash"
	// AuditLogChangeKeyStatus is sent when status of guild scheduled event was changed (int - guild scheduled event
	// status) - ScheduledEvent
	AuditLogChangeKeyStatus AuditLogChangeKey = "status"
	// AuditLogChangeKeySystemChannelID is sent when id of the system channel changed (snowflake) - Guild
	AuditLogChangeKeySystemChannelID AuditLogChangeKey = "system_channel_id"
	// AuditLogChangeKeyTags is sent when related emoji of sticker changed (string) - emoji.Sticker
	AuditLogChangeKeyTags AuditLogChangeKey = "tags"
	// AuditLogChangeKeyTemporary is sent when invite code is now temporary or never expires (bool) - invite.Invite
	AuditLogChangeKeyTemporary AuditLogChangeKey = "temporary"
	// TODO: remove when compatibility is not required
	AuditLogChangeKeyTempoary = AuditLogChangeKeyTemporary
	// AuditLogChangeKeyTopic is sent when text channel topic or stage instance topic changed (string) - channel.Channel
	// or channel.StageInstance
	AuditLogChangeKeyTopic AuditLogChangeKey = "topic"
	// AuditLogChangeKeyType is sent when type of entity created (int or string) - any
	AuditLogChangeKeyType AuditLogChangeKey = "type"
	// AuditLogChangeKeyUnicodeEmoji is sent when role unicode emoji changed (string) - Role
	AuditLogChangeKeyUnicodeEmoji AuditLogChangeKey = "unicode_emoji"
	// AuditLogChangeKeyUserLimit is sent when new user limit in a voice channel set (int) - voice channel.Channel
	AuditLogChangeKeyUserLimit AuditLogChangeKey = "user_limit"
	// AuditLogChangeKeyUses is sent when number of times invite code used changed (int) - invite.Invite
	AuditLogChangeKeyUses AuditLogChangeKey = "uses"
	// AuditLogChangeKeyVanityURLCode is sent when guild invite vanity url changed (string) - Guild
	AuditLogChangeKeyVanityURLCode AuditLogChangeKey = "vanity_url_code"
	// AuditLogChangeKeyVerificationLevel is sent when required verification level changed (int - verification level) - Guild
	AuditLogChangeKeyVerificationLevel AuditLogChangeKey = "verification_level"
	// AuditLogChangeKeyWidgetChannelID is sent when channel id of the server widget changed (snowflake) - Guild
	AuditLogChangeKeyWidgetChannelID AuditLogChangeKey = "widget_channel_id"
	// AuditLogChangeKeyWidgetEnabled is sent when server widget enabled/disabled (bool) - Guild
	AuditLogChangeKeyWidgetEnabled AuditLogChangeKey = "widget_enabled"
	// AuditLogChangeKeyRoleAdd is sent when new role added (array of partial role objects) - Guild
	AuditLogChangeKeyRoleAdd AuditLogChangeKey = "$add"
	// AuditLogChangeKeyRoleRemove is sent when role removed (array of partial role objects) - Guild
	AuditLogChangeKeyRoleRemove AuditLogChangeKey = "$remove"
)

// AuditLogOptions optional data for the AuditLog
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditLogOptions struct {
	DeleteMemberDays              string                 `json:"delete_member_days"`
	MembersRemoved                string                 `json:"members_removed"`
	ChannelID                     string                 `json:"channel_id"`
	MessageID                     string                 `json:"message_id"`
	Count                         string                 `json:"count"`
	ID                            string                 `json:"id"`
	Type                          *types.AuditLogOptions `json:"type"`
	RoleName                      string                 `json:"role_name"`
	ApplicationID                 string                 `json:"application_id"`
	AutoModerationRuleName        string                 `json:"auto_moderation_rule_name"`
	AutoModerationRuleTriggerType string                 `json:"auto_moderation_rule_trigger_type"`
	IntegrationType               string                 `json:"integration_type"`
}

// AuditLogAction is the Action of the AuditLog
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-events
type AuditLogAction int

const (
	AuditLogActionGuildUpdate AuditLogAction = 1

	AuditLogActionChannelCreate          AuditLogAction = 10
	AuditLogActionChannelUpdate          AuditLogAction = 11
	AuditLogActionChannelDelete          AuditLogAction = 12
	AuditLogActionChannelOverwriteCreate AuditLogAction = 13
	AuditLogActionChannelOverwriteUpdate AuditLogAction = 14
	AuditLogActionChannelOverwriteDelete AuditLogAction = 15

	AuditLogActionMemberKick       AuditLogAction = 20
	AuditLogActionMemberPrune      AuditLogAction = 21
	AuditLogActionMemberBanAdd     AuditLogAction = 22
	AuditLogActionMemberBanRemove  AuditLogAction = 23
	AuditLogActionMemberUpdate     AuditLogAction = 24
	AuditLogActionMemberRoleUpdate AuditLogAction = 25
	AuditLogActionMemberMove       AuditLogAction = 26
	AuditLogActionMemberDisconnect AuditLogAction = 27
	AuditLogActionBotAdd           AuditLogAction = 28

	AuditLogActionRoleCreate AuditLogAction = 30
	AuditLogActionRoleUpdate AuditLogAction = 31
	AuditLogActionRoleDelete AuditLogAction = 32

	AuditLogActionInviteCreate AuditLogAction = 40
	AuditLogActionInviteUpdate AuditLogAction = 41
	AuditLogActionInviteDelete AuditLogAction = 42

	AuditLogActionWebhookCreate AuditLogAction = 50
	AuditLogActionWebhookUpdate AuditLogAction = 51
	AuditLogActionWebhookDelete AuditLogAction = 52

	AuditLogActionEmojiCreate AuditLogAction = 60
	AuditLogActionEmojiUpdate AuditLogAction = 61
	AuditLogActionEmojiDelete AuditLogAction = 62

	AuditLogActionMessageDelete     AuditLogAction = 72
	AuditLogActionMessageBulkDelete AuditLogAction = 73
	AuditLogActionMessagePin        AuditLogAction = 74
	AuditLogActionMessageUnpin      AuditLogAction = 75

	AuditLogActionIntegrationCreate   AuditLogAction = 80
	AuditLogActionIntegrationUpdate   AuditLogAction = 81
	AuditLogActionIntegrationDelete   AuditLogAction = 82
	AuditLogActionStageInstanceCreate AuditLogAction = 83
	AuditLogActionStageInstanceUpdate AuditLogAction = 84
	AuditLogActionStageInstanceDelete AuditLogAction = 85

	AuditLogActionStickerCreate AuditLogAction = 90
	AuditLogActionStickerUpdate AuditLogAction = 91
	AuditLogActionStickerDelete AuditLogAction = 92

	AuditLogGuildScheduledEventCreate AuditLogAction = 100
	AuditLogGuildScheduledEventUpdate AuditLogAction = 101
	AuditLogGuildScheduledEventDelete AuditLogAction = 102

	AuditLogActionThreadCreate AuditLogAction = 110
	AuditLogActionThreadUpdate AuditLogAction = 111
	AuditLogActionThreadDelete AuditLogAction = 112

	AuditLogActionApplicationCommandPermissionUpdate AuditLogAction = 121

	AuditLogActionAutoModerationRuleCreate                AuditLogAction = 140
	AuditLogActionAutoModerationRuleUpdate                AuditLogAction = 141
	AuditLogActionAutoModerationRuleDelete                AuditLogAction = 142
	AuditLogActionAutoModerationBlockMessage              AuditLogAction = 143
	AuditLogActionAutoModerationFlagToChannel             AuditLogAction = 144
	AuditLogActionAutoModerationUserCommunicationDisabled AuditLogAction = 145

	AuditLogActionCreatorMonetizationRequestCreated AuditLogAction = 150
	AuditLogActionCreatorMonetizationTermsAccepted  AuditLogAction = 151

	AuditLogActionOnboardingPromptCreate AuditLogAction = 163
	AuditLogActionOnboardingPromptUpdate AuditLogAction = 164
	AuditLogActionOnboardingPromptDelete AuditLogAction = 165
	AuditLogActionOnboardingCreate       AuditLogAction = 166
	AuditLogActionOnboardingUpdate       AuditLogAction = 167

	AuditLogActionHomeSettingsCreate = 190
	AuditLogActionHomeSettingsUpdate = 191
)

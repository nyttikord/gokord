package types

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

// AutoModerationAction represents an action which will execute whenever a guild.AutoModerationAction is triggered.
type AutoModerationAction int

const (
	AutoModerationActionBlockMessage     AutoModerationAction = 1
	AutoModerationActionSendAlertMessage AutoModerationAction = 2
	AutoModerationActionTimeout          AutoModerationAction = 3
)

// OnboardingPrompt is the type of guild.OnboardingPrompt.
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-prompt-types
type OnboardingPrompt int

const (
	OnboardingPromptMultipleChoice OnboardingPrompt = 0
	OnboardingPromptDropdown       OnboardingPrompt = 1
)

// ScheduledEventEntity is the type of entity associated with guild.ScheduledEvent.
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types
type ScheduledEventEntity int

const (
	// ScheduledEventEntityStageInstance represents a stage channel
	ScheduledEventEntityStageInstance ScheduledEventEntity = 1
	// ScheduledEventEntityVoice represents a voice channel
	ScheduledEventEntityVoice ScheduledEventEntity = 2
	// ScheduledEventEntityExternal represents an external event
	ScheduledEventEntityExternal ScheduledEventEntity = 3
)

// AuditLogOptions is the type of guild.AuditLogOptions.
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditLogOptions string

const (
	AuditLogOptionsRole   AuditLogOptions = "0"
	AuditLogOptionsMember AuditLogOptions = "1"
)

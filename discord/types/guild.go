package types

// AutoModerationRuleEvent indicates in what event context a rule should be checked.
type AutoModerationRuleEvent int

// Auto moderation rule event types.
const (
	// AutoModerationRuleEventMessageSend is checked when a member sends or edits a message in the guild
	AutoModerationRuleEventMessageSend AutoModerationRuleEvent = 1
)

// AutoModerationRuleTrigger represents the type of content which can trigger the rule.
type AutoModerationRuleTrigger int

// Auto moderation rule trigger types.
const (
	AutoModerationRuleTriggerKeyword       AutoModerationRuleTrigger = 1
	AutoModerationRuleTriggerHarmfulLink   AutoModerationRuleTrigger = 2
	AutoModerationRuleTriggerSpam          AutoModerationRuleTrigger = 3
	AutoModerationRuleTriggerKeywordPreset AutoModerationRuleTrigger = 4
)

// AutoModerationAction represents an action which will execute whenever a rule is triggered.
type AutoModerationAction int

// Auto moderation actions types.
const (
	AutoModerationActionBlockMessage     AutoModerationAction = 1
	AutoModerationActionSendAlertMessage AutoModerationAction = 2
	AutoModerationActionTimeout          AutoModerationAction = 3
)

// OnboardingPrompt is the type of an onboarding prompt.
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-prompt-types
type OnboardingPrompt int

// Block containing known OnboardingPrompt values.
const (
	OnboardingPromptMultipleChoice OnboardingPrompt = 0
	OnboardingPromptDropdown       OnboardingPrompt = 1
)

// ScheduledEventEntity is the type of entity associated with a guild scheduled event.
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

// AuditLogOptions of the AuditLogOption
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditLogOptions string

// Valid Types for AuditLogOptions
const (
	AuditLogOptionsRole   AuditLogOptions = "0"
	AuditLogOptionsMember AuditLogOptions = "1"
)

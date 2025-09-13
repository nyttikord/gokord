package types

// AutoModerationAction if the type of action which will execute whenever a guild.AutoModerationAction is triggered.
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

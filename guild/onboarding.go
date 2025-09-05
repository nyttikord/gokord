package guild

import "github.com/nyttikord/gokord/user"

// OnboardingMode defines the criteria used to satisfy constraints that are required for enabling onboarding.
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-onboarding-mode
type OnboardingMode int

// Block containing known OnboardingMode values.
const (
	// OnboardingModeDefault counts default channels towards constraints.
	OnboardingModeDefault OnboardingMode = 0
	// OnboardingModeAdvanced counts default channels and questions towards constraints.
	OnboardingModeAdvanced OnboardingMode = 1
)

// Onboarding represents the onboarding flow for a guild.
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object
type Onboarding struct {
	// ID of the guild this onboarding flow is part of.
	GuildID string `json:"guild_id,omitempty"`

	// Prompts shown during onboarding and in the customize community (Channels & Roles) tab.
	Prompts *[]OnboardingPrompt `json:"prompts,omitempty"`

	// Channel IDs that members get opted into automatically.
	DefaultChannelIDs []string `json:"default_channel_ids,omitempty"`

	// Whether onboarding is enabled in the guild.
	Enabled *bool `json:"enabled,omitempty"`

	// Mode of onboarding.
	Mode *OnboardingMode `json:"mode,omitempty"`
}

// OnboardingPromptType is the type of an onboarding prompt.
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-prompt-types
type OnboardingPromptType int

// Block containing known OnboardingPromptType values.
const (
	OnboardingPromptTypeMultipleChoice OnboardingPromptType = 0
	OnboardingPromptTypeDropdown       OnboardingPromptType = 1
)

// OnboardingPrompt is a prompt shown during onboarding and in the customize community (Channels & Roles) tab.
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-onboarding-prompt-structure
type OnboardingPrompt struct {
	// ID of the prompt.
	// NOTE: always requires to be a valid snowflake (e.g. "0"), see
	// https://github.com/discord/discord-api-docs/issues/6320 for more information.
	ID string `json:"id,omitempty"`

	// Type of the prompt.
	Type OnboardingPromptType `json:"type"`

	// Options available within the prompt.
	Options []OnboardingPromptOption `json:"options"`

	// Title of the prompt.
	Title string `json:"title"`

	// Indicates whether users are limited to selecting one option for the prompt.
	SingleSelect bool `json:"single_select"`

	// Indicates whether the prompt is required before a user completes the onboarding flow.
	Required bool `json:"required"`

	// Indicates whether the prompt is present in the onboarding flow.
	// If false, the prompt will only appear in the customize community (Channels & Roles) tab.
	InOnboarding bool `json:"in_onboarding"`
}

// OnboardingPromptOption is an option available within an onboarding prompt.
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-prompt-option-structure
type OnboardingPromptOption struct {
	// ID of the prompt option.
	ID string `json:"id,omitempty"`

	// IDs for channels a member is added to when the option is selected.
	ChannelIDs []string `json:"channel_ids"`

	// IDs for roles assigned to a member when the option is selected.
	RoleIDs []string `json:"role_ids"`

	// Emoji of the option.
	// NOTE: when creating or updating a prompt option
	// EmojiID, EmojiName and EmojiAnimated should be used instead.
	Emoji *user.Emoji `json:"emoji,omitempty"`

	// Title of the option.
	Title string `json:"title"`

	// Description of the option.
	Description string `json:"description"`

	// ID of the option's emoji.
	// NOTE: only used when creating or updating a prompt option.
	EmojiID string `json:"emoji_id,omitempty"`
	// Name of the option's emoji.
	// NOTE: only used when creating or updating a prompt option.
	EmojiName string `json:"emoji_name,omitempty"`
	// Whether the option's emoji is animated.
	// NOTE: only used when creating or updating a prompt option.
	EmojiAnimated *bool `json:"emoji_animated,omitempty"`
}

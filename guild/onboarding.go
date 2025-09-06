package guild

import (
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
)

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
	// Guild.ID this onboarding flow is part of.
	GuildID string `json:"guild_id,omitempty"`

	// Prompts shown during onboarding and in the customize community (Channels & Roles) tab.
	Prompts *[]OnboardingPrompt `json:"prompts,omitempty"`

	// channel.Channel IDs that members get opted into automatically.
	DefaultChannelIDs []string `json:"default_channel_ids,omitempty"`

	// Whether onboarding is enabled in the Guild.
	Enabled *bool `json:"enabled,omitempty"`

	// Mode of onboarding.
	Mode *OnboardingMode `json:"mode,omitempty"`
}

// OnboardingPrompt is a prompt shown during Onboarding and in the customize community (Channels & Roles) tab.
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-onboarding-prompt-structure
type OnboardingPrompt struct {
	// ID of the prompt.
	//
	// Note: always requires to be a valid snowflake (e.g. "0"), see
	// https://github.com/discord/discord-api-docs/issues/6320 for more information.
	ID string `json:"id,omitempty"`

	// Type of the OnboardingPrompt.
	Type types.OnboardingPrompt `json:"type"`

	// Options available within the OnboardingPrompt.
	Options []OnboardingPromptOption `json:"options"`

	// Title of the OnboardingPrompt.
	Title string `json:"title"`

	// Indicates whether users are limited to selecting one option for the OnboardingPrompt.
	SingleSelect bool `json:"single_select"`

	// Indicates whether the OnboardingPrompt is required before a user completes the onboarding flow.
	Required bool `json:"required"`

	// Indicates whether the OnboardingPrompt is present in the onboarding flow.
	// If false, the OnboardingPrompt will only appear in the customize community (Channels & Roles) tab.
	InOnboarding bool `json:"in_onboarding"`
}

// OnboardingPromptOption is an option available within an OnboardingPrompt.
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-prompt-option-structure
type OnboardingPromptOption struct {
	// ID of the OnboardingPromptOption.
	ID string `json:"id,omitempty"`

	// IDs for channels a user.Member is added to when the OnboardingPromptOption is selected.
	ChannelIDs []string `json:"channel_ids"`

	// IDs for Roles assigned to a user.Member when the OnboardingPromptOption is selected.
	RoleIDs []string `json:"role_ids"`

	// Emoji of the option.
	//
	// Note: when creating or updating a OnboardingPromptOption
	// EmojiID, EmojiName and EmojiAnimated should be used instead.
	Emoji *emoji.Emoji `json:"emoji,omitempty"`

	// Title of the OnboardingPromptOption.
	Title string `json:"title"`

	// Description of the OnboardingPromptOption.
	Description string `json:"description"`

	// ID of the option's emoji.Emoji.
	//
	// Note: only used when creating or updating a OnboardingPromptOption.
	EmojiID string `json:"emoji_id,omitempty"`
	// Name of the option's emoji.Emoji.
	//
	// Note: only used when creating or updating a OnboardingPromptOption.
	EmojiName string `json:"emoji_name,omitempty"`
	// Whether the option's emoji.Emoji is animated.
	//
	// Note: only used when creating or updating a OnboardingPromptOption.
	EmojiAnimated *bool `json:"emoji_animated,omitempty"`
}

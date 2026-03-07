package guild

import (
	"encoding/json"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/internal/structs"
)

// OnboardingMode defines the criteria used to satisfy constraints that are required for enabling [Onboarding].
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-onboarding-mode
type OnboardingMode int

const (
	// OnboardingModeDefault counts default channels towards constraints.
	OnboardingModeDefault OnboardingMode = 0
	// OnboardingModeAdvanced counts default channels and questions towards constraints.
	OnboardingModeAdvanced OnboardingMode = 1
)

// Onboarding represents the onboarding flow for a [Guild].
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object
type Onboarding struct {
	GuildID uint64 `json:"guild_id,omitempty,string"`
	// Prompts shown during [Onboarding] and in the customize community ([channel.Channel]s & [Role]s) tab.
	Prompts *[]OnboardingPrompt `json:"prompts,omitempty"`
	// [channel.Channel] IDs that members get opted into automatically.
	DefaultChannelIDs []string `json:"default_channel_ids,omitempty"`
	// Whether [Onboarding] is enabled in the [Guild].
	Enabled *bool `json:"enabled,omitempty"`
	// Mode of [Onboarding].
	Mode *OnboardingMode `json:"mode,omitempty"`
}

// OnboardingPrompt is a prompt shown during [Onboarding] and in the customize community ([channel.Channel]s and
// [Role]s) tab.
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-onboarding-prompt-structure
type OnboardingPrompt struct {
	// ID of the prompt.
	//
	// Note: always requires to be a valid snowflake (e.g. "0"), see
	// https://github.com/discord/discord-api-docs/issues/6320 for more information.
	ID uint64 `json:"id,omitempty,string"`
	// Type of the [OnboardingPrompt].
	Type types.OnboardingPrompt `json:"type"`
	// Options available within the [OnboardingPrompt].
	Options []OnboardingPromptOption `json:"options"`
	Title   string                   `json:"title"`
	// Indicates whether [user.User]s are limited to selecting one option for the [OnboardingPrompt].
	SingleSelect bool `json:"single_select"`
	// Indicates whether the [OnboardingPrompt] is required before a user completes the [Onboarding] flow.
	Required bool `json:"required"`
	// Indicates whether the [OnboardingPrompt] is present in the onboarding flow.
	// If false, the [OnboardingPrompt] will only appear in the customize community ([channel.Channel]s & [Role]s) tab.
	InOnboarding bool `json:"in_onboarding"`
}

// OnboardingPromptOption is an option available within an [OnboardingPrompt].
// https://discord.com/developers/docs/resources/guild#guild-onboarding-object-prompt-option-structure
type OnboardingPromptOption struct {
	ID uint64 `json:"id,omitempty,string"`
	// IDs for channels a [user.Member] is added to when the [OnboardingPromptOption] is selected.
	ChannelIDs []uint64 `json:"-"`
	// IDs for [Role]s assigned to a [user.Member] when the [OnboardingPromptOption] is selected.
	RoleIDs []uint64 `json:"-"`
	// [emoji.Emoji] of the option.
	//
	// NOTE: when creating or updating a [OnboardingPromptOption] EmojiID, EmojiName and EmojiAnimated should be used
	// instead.
	Emoji       *emoji.Emoji `json:"emoji,omitempty"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	// See [OnboardingPromptOption.Emoji].
	EmojiID uint64 `json:"emoji_id,omitempty,string"`
	// See [OnboardingPromptOption.Emoji].
	EmojiName string `json:"emoji_name,omitempty"`
	// See [OnboardingPromptOption.Emoji].
	EmojiAnimated *bool `json:"emoji_animated,omitempty"`
}

func (o *OnboardingPromptOption) MarshalJSON() ([]byte, error) {
	type t OnboardingPromptOption
	v := struct {
		t
		ChannelIDs []string `json:"channel_ids"`
		RoleIDs    []string `json:"role_ids"`
	}{t(*o), structs.UintsToSnowflakes(o.ChannelIDs), structs.UintsToSnowflakes(o.RoleIDs)}
	return json.Marshal(v)
}

func (o *OnboardingPromptOption) UnmarshalJSON(data []byte) error {
	type t OnboardingPromptOption
	var v struct {
		t
		ChannelIDs []string `json:"channel_ids"`
		RoleIDs    []string `json:"role_ids"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*o = OnboardingPromptOption(v.t)
	o.ChannelIDs = structs.SnowflakesToUints(v.ChannelIDs)
	o.RoleIDs = structs.SnowflakesToUints(v.RoleIDs)
	return nil
}

// GetOnboarding returns [Onboarding] configuration of a [Guild].
func GetOnboarding(guildID uint64) Request[*Onboarding] {
	return NewData[*Onboarding](http.MethodGet, discord.EndpointGuildOnboarding(guildID))
}

// EditOnboarding configuration of a [Guild].
func EditOnboarding(guildID uint64, o *Onboarding) Request[*Onboarding] {
	return NewData[*Onboarding](http.MethodPut, discord.EndpointGuildOnboarding(guildID)).
		WithData(o)
}

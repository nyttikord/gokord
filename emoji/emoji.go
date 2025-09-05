package emoji

import (
	"github.com/nyttikord/gokord/user"
	"regexp"
)

// Emoji struct holds data related to Emoji's
type Emoji struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Roles         []string   `json:"roles"`
	User          *user.User `json:"user"`
	RequireColons bool       `json:"require_colons"`
	Managed       bool       `json:"managed"`
	Animated      bool       `json:"animated"`
	Available     bool       `json:"available"`
}

// EmojiRegex is the regex used to find and identify emojis in messages
var (
	EmojiRegex = regexp.MustCompile(`<(a|):[A-Za-z0-9_~]+:[0-9]{18,20}>`)
)

// MessageFormat returns a correctly formatted Emoji for use in Message content and embeds
func (e *Emoji) MessageFormat() string {
	if e.ID != "" && e.Name != "" {
		if e.Animated {
			return "<a:" + e.APIName() + ">"
		}

		return "<:" + e.APIName() + ">"
	}

	return e.APIName()
}

// APIName returns an correctly formatted API name for use in the MessageReactions endpoints.
func (e *Emoji) APIName() string {
	if e.ID != "" && e.Name != "" {
		return e.Name + ":" + e.ID
	}
	if e.Name != "" {
		return e.Name
	}
	return e.ID
}

// EmojiParams represents parameters needed to create or update an Emoji.
type EmojiParams struct {
	// Name of the emoji
	Name string `json:"name,omitempty"`
	// A base64 encoded emoji image, has to be smaller than 256KB.
	// NOTE: can be only set on creation.
	Image string `json:"image,omitempty"`
	// Roles for which this emoji will be available.
	// NOTE: can not be used with application emoji endpoints.
	Roles []string `json:"roles,omitempty"`
}

// Component represents button Emoji, if it does have one.
type Component struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

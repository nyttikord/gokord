// Package emoji contains every thing related to emoji and stickers.
package emoji

import (
	"net/http"
	"regexp"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/user"
)

// Emoji struct holds data related to emoji's
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

var (
	// Regex is the regexp.Regexp used to find and identify emojis in messages
	Regex = regexp.MustCompile(`<(a|):[A-Za-z0-9_~]+:[0-9]{18,20}>`)
)

// MessageFormat returns a correctly formatted Emoji for use in channel.Message content and channel.MessageEmbed.
func (e *Emoji) MessageFormat() string {
	if e.ID != "" && e.Name != "" {
		if e.Animated {
			return "<a:" + e.APIName() + ">"
		}

		return "<:" + e.APIName() + ">"
	}

	return e.APIName()
}

// APIName returns a correctly formatted API name for use in the channel.MessageReactions endpoints.
func (e *Emoji) APIName() string {
	if e.ID != "" && e.Name != "" {
		return e.Name + ":" + e.ID
	}
	if e.Name != "" {
		return e.Name
	}
	return e.ID
}

// Params represents parameters needed to create or update an Emoji.
type Params struct {
	// Name of the Emoji.
	Name string `json:"name,omitempty"`
	// A base64 encoded emoji image, has to be smaller than 256KB.
	//
	// NOTE: can be only set on creation.
	Image string `json:"image,omitempty"`
	// Roles for which this Emoji will be available.
	//
	// NOTE: can not be used with application emoji endpoints.
	Roles []string `json:"roles,omitempty"`
}

// Component represents component.Button's Emoji, if it does have one.
// Also used by channel.Poll.
type Component struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

// List returns all [Emoji] in the given [guild.Guild].
func List(guildID string) Request[[]*Emoji] {
	return NewData[[]*Emoji](http.MethodGet, discord.EndpointGuildEmojis(guildID))
}

// Get returns the [Emoji] in the given [guild.Guild].
func Get(guildID, emojiID string) Request[*Emoji] {
	return NewData[*Emoji](http.MethodGet, discord.EndpointGuildEmoji(guildID, emojiID)).
		WithBucketID(discord.EndpointGuildEmojis(guildID))
}

// Create a new [Emoji] in the given [guild.Guild].
func Create(guildID string, data *Params) Request[*Emoji] {
	return NewData[*Emoji](http.MethodPost, discord.EndpointGuildEmojis(guildID)).
		WithData(data)
}

// Update and returns the updated [Emoji] in the given [guild.Guild].
func Update(guildID, emojiID string, data *Params) Request[*Emoji] {
	return NewData[*Emoji](http.MethodPatch, discord.EndpointGuildEmoji(guildID, emojiID)).
		WithBucketID(discord.EndpointGuildEmojis(guildID)).WithData(data)
}

// Delete an [Emoji] in the given [guild.Guild].
func Delete(guildID, emojiID string) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointGuildEmoji(guildID, emojiID)).
		WithBucketID(discord.EndpointGuildEmojis(guildID))
	return WrapAsEmpty(req)
}

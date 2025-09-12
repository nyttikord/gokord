package channel

import (
	"github.com/nyttikord/gokord/component"
	"github.com/nyttikord/gokord/discord/types"
)

// MessageSend stores all parameters you can send with ChannelMessageSendComplex.
type MessageSend struct {
	Content         string                  `json:"content,omitempty"`
	Embeds          []*MessageEmbed         `json:"embeds"`
	TTS             bool                    `json:"tts"`
	Components      []component.Message     `json:"components"`
	Files           []*File                 `json:"-"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Reference       *MessageReference       `json:"message_reference,omitempty"`
	StickerIDs      []string                `json:"sticker_ids"`
	Flags           MessageFlags            `json:"flags,omitempty"`
	Poll            *Poll                   `json:"poll,omitempty"`
}

// MessageEdit is used to chain parameters via ChannelMessageEditComplex, which
// is also where you should get the instance from.
type MessageEdit struct {
	Content         *string                 `json:"content,omitempty"`
	Components      *[]component.Message    `json:"components,omitempty"`
	Embeds          *[]*MessageEmbed        `json:"embeds,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlags            `json:"flags,omitempty"`
	// Files to append to the message
	Files []*File `json:"-"`
	// Overwrite existing attachments
	Attachments *[]*MessageAttachment `json:"attachments,omitempty"`

	ID      string
	Channel string
}

// NewMessageEdit returns a MessageEdit struct, initialized with the Channel and ID.
func NewMessageEdit(channelID string, messageID string) *MessageEdit {
	return &MessageEdit{
		Channel: channelID,
		ID:      messageID,
	}
}

// SetContent is the same as setting the variable Content, except it doesn't take a pointer.
func (m *MessageEdit) SetContent(str string) *MessageEdit {
	m.Content = &str
	return m
}

// SetEmbed is a convenience function for setting the embed, so you can chain commands.
func (m *MessageEdit) SetEmbed(embed *MessageEmbed) *MessageEdit {
	m.Embeds = &[]*MessageEmbed{embed}
	return m
}

// SetEmbeds is a convenience function for setting the embeds, so you can chain commands.
func (m *MessageEdit) SetEmbeds(embeds []*MessageEmbed) *MessageEdit {
	m.Embeds = &embeds
	return m
}

// MessageAllowedMentions allows the user to specify which mentions Discord is allowed to parse in this message.
// This is useful when sending user input as a message, as it prevents unwanted mentions.
// If this type is used, all mentions must be explicitly whitelisted, either by putting an AllowedMentionType in the
// Parse slice (allowing all mentions of that type) or, in the case of roles and users, explicitly allowing those
// mentions on an ID-by-ID basis.
// For more information on this functionality, see:
// https://discordapp.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-reference
type MessageAllowedMentions struct {
	// The mention types that are allowed to be parsed in this message.
	// Please note that this is purposely **not** marked as omitempty, so if a zero-value MessageAllowedMentions object
	// is provided no mentions will be allowed.
	Parse []types.AllowedMention `json:"parse"`

	// A list of guild.Role IDs to allow.
	// This cannot be used when specifying types.AllowedMentionRoles in the Parse slice.
	Roles []string `json:"roles,omitempty"`

	// A list of user.User IDs to allow.
	// This cannot be used when specifying types.AllowedMentionUsers in the Parse slice.
	Users []string `json:"users,omitempty"`

	// For replies, whether to mention the author of the message being replied to
	RepliedUser bool `json:"replied_user"`
}

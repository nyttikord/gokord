package user

import "regexp"

// Emoji struct holds data related to Emoji's
type Emoji struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles"`
	User          *User    `json:"user"`
	RequireColons bool     `json:"require_colons"`
	Managed       bool     `json:"managed"`
	Animated      bool     `json:"animated"`
	Available     bool     `json:"available"`
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

// StickerFormat is the file format of the Sticker.
type StickerFormat int

// Defines all known Sticker types.
const (
	StickerFormatTypePNG    StickerFormat = 1
	StickerFormatTypeAPNG   StickerFormat = 2
	StickerFormatTypeLottie StickerFormat = 3
	StickerFormatTypeGIF    StickerFormat = 4
)

// StickerType is the type of sticker.
type StickerType int

// Defines Sticker types.
const (
	StickerTypeStandard StickerType = 1
	StickerTypeGuild    StickerType = 2
)

// Sticker represents a sticker object that can be sent in a Message.
type Sticker struct {
	ID          string        `json:"id"`
	PackID      string        `json:"pack_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Tags        string        `json:"tags"`
	Type        StickerType   `json:"type"`
	FormatType  StickerFormat `json:"format_type"`
	Available   bool          `json:"available"`
	GuildID     string        `json:"guild_id"`
	User        *User         `json:"user"`
	SortValue   int           `json:"sort_value"`
}

// StickerItem represents the smallest amount of data required to render a sticker. A partial sticker object.
type StickerItem struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	FormatType StickerFormat `json:"format_type"`
}

// StickerPack represents a pack of standard stickers.
type StickerPack struct {
	ID             string     `json:"id"`
	Stickers       []*Sticker `json:"stickers"`
	Name           string     `json:"name"`
	SKUID          string     `json:"sku_id"`
	CoverStickerID string     `json:"cover_sticker_id"`
	Description    string     `json:"description"`
	BannerAssetID  string     `json:"banner_asset_id"`
}

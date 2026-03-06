package emoji

import (
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// StickerFormat is the file format of the Sticker.
type StickerFormat int

// Defines all known Sticker format.
const (
	StickerFormatPNG    StickerFormat = 1
	StickerFormatAPNG   StickerFormat = 2
	StickerFormatLottie StickerFormat = 3
	StickerFormatGIF    StickerFormat = 4
)

// Sticker represents a sticker object that can be sent in a [channel.Message].
type Sticker struct {
	ID          uint64        `json:"id,string"`
	PackID      uint64        `json:"pack_id,string"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Tags        string        `json:"tags"`
	Type        types.Sticker `json:"type"`
	FormatType  StickerFormat `json:"format_type"`
	Available   bool          `json:"available"`
	GuildID     uint64        `json:"guild_id,string"`
	User        *user.User    `json:"user"`
	SortValue   int           `json:"sort_value"`
}

// StickerItem represents the smallest amount of data required to render a [Sticker].
// A partial [Sticker] object.
type StickerItem struct {
	ID         uint64        `json:"id,string"`
	Name       string        `json:"name"`
	FormatType StickerFormat `json:"format_type"`
}

// StickerPack represents a pack of standard [Sticker]s.
type StickerPack struct {
	ID             uint64     `json:"id,string"`
	Stickers       []*Sticker `json:"stickers"`
	Name           string     `json:"name"`
	SKUID          uint64     `json:"sku_id,string"`
	CoverStickerID uint64     `json:"cover_sticker_id,string"`
	Description    string     `json:"description"`
	BannerAssetID  uint64     `json:"banner_asset_id,string"`
}

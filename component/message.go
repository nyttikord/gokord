package component

import (
	"encoding/json"

	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
)

// Message is implemented by all message components.
type Message interface {
	json.Marshaler
	Type() types.Component
	message()
}

// ActionsRow is a top-level container Component for displaying a row of interactive components.
type ActionsRow struct {
	// Can contain Button, SelectMenu and TextInput.
	//
	// NOTE: maximum of 5.
	Components []Message `json:"components"`
	// Unique identifier for the Component; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`
}

func (r *ActionsRow) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ActionsRow
		Type types.Component `json:"type"`
	}{
		ActionsRow: *r,
		Type:       r.Type(),
	})
}

func (r *ActionsRow) Type() types.Component {
	return types.ComponentActionsRow
}

func (r *ActionsRow) message() {}

// ButtonStyle is style of Button.
type ButtonStyle uint

const (
	// ButtonStylePrimary is a button with blurple color.
	ButtonStylePrimary ButtonStyle = 1
	// ButtonStyleSecondary is a button with grey color.
	ButtonStyleSecondary ButtonStyle = 2
	// ButtonStyleSuccess is a button with green color.
	ButtonStyleSuccess ButtonStyle = 3
	// ButtonStyleDanger is a button with red color.
	ButtonStyleDanger ButtonStyle = 4
	// ButtonStyleLink is a special type of button which navigates to a URL. Has grey color.
	ButtonStyleLink ButtonStyle = 5
	// ButtonStylePremium is a special type of button with a blurple color that links to a SKU.
	ButtonStylePremium ButtonStyle = 6
)

// Button represents button Component.
type Button struct {
	Label    string           `json:"label"`
	Style    ButtonStyle      `json:"style"`
	Disabled bool             `json:"disabled"`
	Emoji    *emoji.Component `json:"emoji,omitempty"`

	// NOTE: Only button with ButtonStyleLink style can have link. Also, URL is mutually exclusive with CustomID.
	URL      string `json:"url,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
	// Identifier for a purchasable premium.SKU. Only available when using ButtonStylePremium.
	SKUID string `json:"sku_id,omitempty"`
	// Unique identifier for the Component; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`
}

func (b *Button) MarshalJSON() ([]byte, error) {
	if b.Style == 0 {
		b.Style = ButtonStylePrimary
	}
	return json.Marshal(struct {
		Button
		Type types.Component `json:"type"`
	}{
		Button: *b,
		Type:   b.Type(),
	})
}

func (*Button) Type() types.Component {
	return types.ComponentButton
}

func (b *Button) message() {}

// Section is a top-level layout Component that allows you to join Message contextually with an Accessory.
type Section struct {
	// Unique identifier for the Component; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`
	// Array of text display components; max of 3.
	Components []Message `json:"components"`
	// Can be Button or Thumbnail
	Accessory Message `json:"accessory"`
}

func (*Section) Type() types.Component {
	return types.ComponentSection
}

func (s *Section) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Section
		Type types.Component `json:"type"`
	}{
		Section: *s,
		Type:    s.Type(),
	})
}

func (*Section) message() {}

// TextDisplay is a top-level Component that allows you to add markdown-formatted text to the Message.
type TextDisplay struct {
	Content string `json:"content"`
}

func (*TextDisplay) Type() types.Component {
	return types.ComponentTextDisplay
}

func (t *TextDisplay) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TextDisplay
		Type types.Component `json:"type"`
	}{
		TextDisplay: *t,
		Type:        t.Type(),
	})
}

func (*TextDisplay) message() {}

// Thumbnail can be used as an accessory for a Section component.
type Thumbnail struct {
	// Unique identifier for the Component; autopopulated through increment if not provided.
	ID          int               `json:"id,omitempty"`
	Media       UnfurledMediaItem `json:"media"`
	Description *string           `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler,omitempty"`
}

func (*Thumbnail) Type() types.Component {
	return types.ComponentThumbnail
}

func (t *Thumbnail) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Thumbnail
		Type types.Component `json:"type"`
	}{
		Thumbnail: *t,
		Type:      t.Type(),
	})
}

func (*Thumbnail) message() {}

// MediaGallery is a top-level Component allows you to group images, videos or gifs into a gallery grid.
type MediaGallery struct {
	// Unique identifier for the Component; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`
	// Array of media gallery items; max of 10.
	Items []MediaGalleryItem `json:"items"`
}

func (*MediaGallery) Type() types.Component {
	return types.ComponentMediaGallery
}

func (m *MediaGallery) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		MediaGallery
		Type types.Component `json:"type"`
	}{
		MediaGallery: *m,
		Type:         m.Type(),
	})
}

func (*MediaGallery) message() {}

// MediaGalleryItem represents an item used in MediaGallery.
type MediaGalleryItem struct {
	Media       UnfurledMediaItem `json:"media"`
	Description *string           `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler"`
}

// File is a top-level Component that allows you to display an uploaded file as an attachment to the message and
// reference it in the Component.
type File struct {
	// Unique identifier for the Component; autopopulated through increment if not provided.
	ID      int               `json:"id,omitempty"`
	File    UnfurledMediaItem `json:"file"`
	Spoiler bool              `json:"spoiler"`
}

func (*File) Type() types.Component {
	return types.ComponentFile
}

func (f *File) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		File
		Type types.Component `json:"type"`
	}{
		File: *f,
		Type: f.Type(),
	})
}

func (*File) message() {}

// SeparatorSpacingSize represents spacing size around the Separator.
type SeparatorSpacingSize uint

const (
	SeparatorSpacingSizeSmall SeparatorSpacingSize = 1
	SeparatorSpacingSizeLarge SeparatorSpacingSize = 2
)

// Separator is a top-level layout Component that adds vertical padding and visual division between other components.
type Separator struct {
	// Unique identifier for the component; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`

	Divider *bool                 `json:"divider,omitempty"`
	Spacing *SeparatorSpacingSize `json:"spacing,omitempty"`
}

func (*Separator) Type() types.Component {
	return types.ComponentSeparator
}

func (s *Separator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Separator
		Type types.Component `json:"type"`
	}{
		Separator: *s,
		Type:      s.Type(),
	})
}

func (*Separator) message() {}

// Container is a top-level layout Component.
// Containers are visually distinct from surrounding components and have an optional customizable color bar (similar to
// channel.MessageEmbed).
type Container struct {
	// Unique identifier for the Component; autopopulated through increment if not provided.
	ID          int       `json:"id,omitempty"`
	AccentColor *int      `json:"accent_color,omitempty"`
	Spoiler     bool      `json:"spoiler"`
	Components  []Message `json:"components"`
}

func (*Container) Type() types.Component {
	return types.ComponentContainer
}

func (c *Container) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Container
		Type types.Component `json:"type"`
	}{
		Container: *c,
		Type:      c.Type(),
	})
}

func (*Container) message() {}

// UnfurledMediaItem represents an unfurled media item.
type UnfurledMediaItem struct {
	URL string `json:"url"`
}

// UnfurledMediaItemLoadingState is the loading state of the unfurled media item.
type UnfurledMediaItemLoadingState uint

const (
	UnfurledMediaItemLoadingStateUnknown        UnfurledMediaItemLoadingState = 0
	UnfurledMediaItemLoadingStateLoading        UnfurledMediaItemLoadingState = 1
	UnfurledMediaItemLoadingStateLoadingSuccess UnfurledMediaItemLoadingState = 2
	UnfurledMediaItemLoadingStateLoadedNotFound UnfurledMediaItemLoadingState = 3
)

// ResolvedUnfurledMediaItem represents a resolved unfurled media item.
type ResolvedUnfurledMediaItem struct {
	URL         string `json:"url"`
	ProxyURL    string `json:"proxy_url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	ContentType string `json:"content_type"`
}

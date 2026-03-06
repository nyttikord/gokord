package component

import (
	"encoding/json"

	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
)

// Message is implemented by every message [Component].
type Message interface {
	Component
	message()
}

// ActionsRow is a top-level container [Component] for displaying a row of interactive [Component]s.
type ActionsRow struct {
	// Components holds [Message] [Component] to display.
	Components []Message `json:"components"`
	// Unique identifier for the [Component]; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`
}

func (r *ActionsRow) MarshalJSON() ([]byte, error) {
	return marshalJSON(r)
}

func (r *ActionsRow) UnmarshalJSON(data []byte) error {
	type t ActionsRow
	var v struct {
		t
		RawComponents []Unmarshaler `json:"component"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*r = ActionsRow(v.t)
	r.Components = make([]Message, len(v.RawComponents))
	for i, vr := range v.RawComponents {
		r.Components[i] = vr.Component.(Message)
	}
	return nil
}

func (r *ActionsRow) Type() types.Component {
	return types.ComponentActionsRow
}

func (r *ActionsRow) message() {}

// ButtonStyle is style of [Button].
type ButtonStyle uint

const (
	// ButtonStylePrimary is a [Button] with blurple color.
	ButtonStylePrimary ButtonStyle = 1
	// ButtonStyleSecondary is a [Button] with grey color.
	ButtonStyleSecondary ButtonStyle = 2
	// ButtonStyleSuccess is a [Button] with green color.
	ButtonStyleSuccess ButtonStyle = 3
	// ButtonStyleDanger is a [Button] with red color.
	ButtonStyleDanger ButtonStyle = 4
	// ButtonStyleLink is a special type of [Button] which navigates to a URL. Has grey color.
	ButtonStyleLink ButtonStyle = 5
	// ButtonStylePremium is a special type of [Button] with a blurple color that links to a [premium.SKU].
	ButtonStylePremium ButtonStyle = 6
)

// Button represents button [Component].
type Button struct {
	Label    string           `json:"label"`
	Style    ButtonStyle      `json:"style"`
	Disabled bool             `json:"disabled"`
	Emoji    *emoji.Component `json:"emoji,omitempty"`

	// NOTE: Only button with ButtonStyleLink style can have link.
	// Also, URL is mutually exclusive with CustomID.
	URL      string `json:"url,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
	// Identifier for a purchasable [premium.SKU].
	// Only available when using [ButtonStylePremium].
	SKUID uint64 `json:"sku_id,omitempty,string"`
	// Unique identifier for the [Component]; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`
}

func (b *Button) MarshalJSON() ([]byte, error) {
	if b.Style == 0 {
		b.Style = ButtonStylePrimary
	}
	return marshalJSON(b)
}

func (*Button) Type() types.Component {
	return types.ComponentButton
}

func (b *Button) message() {}

// Section is a top-level layout [Component] that allows you to join [Message] contextually with an accessory.
type Section struct {
	// Unique identifier for the [Component]; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`
	// Array of [TextDisplay] [Component]s.
	//
	// NOTE: maximum of 3.
	Components []Message `json:"components"`
	// Can be a [Button] or a [Thumbnail].
	Accessory Message `json:"accessory"`
}

func (*Section) Type() types.Component {
	return types.ComponentSection
}

func (s *Section) MarshalJSON() ([]byte, error) {
	return marshalJSON(s)
}

func (s *Section) UnmarshalJSON(data []byte) error {
	type t Section
	var v struct {
		t
		RawComponents []Unmarshaler `json:"components"`
		RawAccessory  Unmarshaler   `json:"accessory"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*s = Section(v.t)
	s.Components = make([]Message, len(v.RawComponents))
	for i, vr := range v.RawComponents {
		s.Components[i] = vr.Component.(Message)
	}
	s.Accessory = v.RawAccessory.Component.(Message)
	return nil
}

func (*Section) message() {}

// Thumbnail can be used as an accessory for a [Section] [Component].
type Thumbnail struct {
	// Unique identifier for the [Component]; autopopulated through increment if not provided.
	ID          int               `json:"id,omitempty"`
	Media       UnfurledMediaItem `json:"media"`
	Description *string           `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler,omitempty"`
}

func (*Thumbnail) Type() types.Component {
	return types.ComponentThumbnail
}

func (t *Thumbnail) MarshalJSON() ([]byte, error) {
	return marshalJSON(t)
}

func (*Thumbnail) message() {}

// MediaGallery is a top-level [Component] allows you to group images, videos or gifs into a gallery grid.
type MediaGallery struct {
	// Unique identifier for the [Component]; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`
	// Array of [MediaGalleryItem]s.
	//
	// NOTE: maximum of 10.
	Items []MediaGalleryItem `json:"items"`
}

func (*MediaGallery) Type() types.Component {
	return types.ComponentMediaGallery
}

func (m *MediaGallery) MarshalJSON() ([]byte, error) {
	return marshalJSON(m)
}

func (*MediaGallery) message() {}

// MediaGalleryItem represents an item used in [MediaGallery].
type MediaGalleryItem struct {
	Media       UnfurledMediaItem `json:"media"`
	Description *string           `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler"`
}

// File is a top-level [Component] that allows you to display an uploaded file as an attachment to the message and
// reference it in the [Component].
type File struct {
	// Unique identifier for the [Component]; autopopulated through increment if not provided.
	ID      int               `json:"id,omitempty"`
	File    UnfurledMediaItem `json:"file"`
	Spoiler bool              `json:"spoiler"`
}

func (*File) Type() types.Component {
	return types.ComponentFile
}

func (f *File) MarshalJSON() ([]byte, error) {
	return marshalJSON(f)
}

func (*File) message() {}

// SeparatorSpacingSize represents spacing size around the [Separator].
type SeparatorSpacingSize uint

const (
	SeparatorSpacingSizeSmall SeparatorSpacingSize = 1
	SeparatorSpacingSizeLarge SeparatorSpacingSize = 2
)

// Separator is a top-level layout [Component] that adds vertical padding and visual division between other
// [Component]s.
type Separator struct {
	// Unique identifier for the [Component]; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`

	Divider *bool                 `json:"divider,omitempty"`
	Spacing *SeparatorSpacingSize `json:"spacing,omitempty"`
}

func (*Separator) Type() types.Component {
	return types.ComponentSeparator
}

func (s *Separator) MarshalJSON() ([]byte, error) {
	return marshalJSON(s)
}

func (*Separator) message() {}

// Container is a top-level layout [Component].
// Containers are visually distinct from surrounding [Component]s and have an optional customizable color bar (similar
// to [channel.MessageEmbed]).
type Container struct {
	// Unique identifier for the [Component]; autopopulated through increment if not provided.
	ID          int       `json:"id,omitempty"`
	AccentColor *int      `json:"accent_color,omitempty"`
	Spoiler     bool      `json:"spoiler"`
	Components  []Message `json:"components"`
}

func (*Container) Type() types.Component {
	return types.ComponentContainer
}

func (c *Container) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *Container) UnmarshalJSON(data []byte) error {
	type t Container
	var v struct {
		t
		RawComponents []Unmarshaler `json:"component"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*c = Container(v.t)
	c.Components = make([]Message, len(v.RawComponents))
	for i, vr := range v.RawComponents {
		c.Components[i] = vr.Component.(Message)
	}
	return nil
}

func (*Container) message() {}

// UnfurledMediaItem represents an unfurled media item.
type UnfurledMediaItem struct {
	URL string `json:"url"`
}

// UnfurledMediaItemLoadingState is the loading state of the [UnfurledMediaItem].
type UnfurledMediaItemLoadingState uint

const (
	UnfurledMediaItemLoadingStateUnknown        UnfurledMediaItemLoadingState = 0
	UnfurledMediaItemLoadingStateLoading        UnfurledMediaItemLoadingState = 1
	UnfurledMediaItemLoadingStateLoadingSuccess UnfurledMediaItemLoadingState = 2
	UnfurledMediaItemLoadingStateLoadedNotFound UnfurledMediaItemLoadingState = 3
)

// ResolvedUnfurledMediaItem represents a resolved [UnfurledMediaItem].
type ResolvedUnfurledMediaItem struct {
	URL         string `json:"url"`
	ProxyURL    string `json:"proxy_url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	ContentType string `json:"content_type"`
}

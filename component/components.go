package component

import (
	"encoding/json"

	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
)

// Component represents every component.
type Component interface {
	json.Marshaler
	Type() types.Component
}

// Message is implemented by all message components.
type Message interface {
	json.Marshaler
	Type() types.Component
	message()
}

// Modal is implemented by all modal components.
type Modal interface {
	json.Marshaler
	Type() types.Component
	modal()
}

func toJson(m Component) ([]byte, error) {
	return json.Marshal(struct {
		Component
		Type types.Component `json:"type"`
	}{m, m.Type()})
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
	return toJson(r)
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
	return toJson(b)
}

func (*Button) Type() types.Component {
	return types.ComponentButton
}

func (b *Button) message() {}

// SelectMenuOption represents an option for a SelectMenu.
type SelectMenuOption struct {
	Label       string           `json:"label,omitempty"`
	Value       string           `json:"value"`
	Description string           `json:"description"`
	Emoji       *emoji.Component `json:"emoji,omitempty"`
	// Determines whenever option is selected by default or not.
	Default bool `json:"default"`
}

// SelectMenuDefaultValue represents an entity selected by default in autopopulated select menus.
type SelectMenuDefaultValue struct {
	// ID of the entity.
	ID string `json:"id"`
	// Type of the entity.
	Type types.SelectMenuDefaultValue `json:"type"`
}

// SelectMenu represents select menu Component.
type SelectMenu struct {
	// Type of the SelectMenu.
	MenuType types.SelectMenu `json:"type,omitempty"`
	// CustomID is a developer-defined identifier for the SelectMenu.
	CustomID string `json:"custom_id,omitempty"`
	// The text which will be shown in the menu if there's no default options or all options was deselected and Component was closed.
	Placeholder string `json:"placeholder"`
	// This value determines the minimal amount of selected items in the menu.
	MinValues *int `json:"min_values,omitempty"`
	// This value determines the maximal amount of selected items in the menu.
	// If MaxValues or MinValues are greater than one then the user can select multiple items in the Component.
	MaxValues int `json:"max_values,omitempty"`
	// List of default values for autopopulated select menus.
	//
	// NOTE: Number of entries should be in the range defined by MinValues and MaxValues.
	DefaultValues []SelectMenuDefaultValue `json:"default_values,omitempty"`

	Options []SelectMenuOption `json:"options,omitempty"`
	// The list of value(s) selected from the predefined options.
	//
	// NOTE: This will only exist if the Interaction was a ModalSubmit otherwise you should (still) be using
	// gokord.InteractionResponse.MessageComponentData()
	Values   []string `json:"values,omitempty"`
	Disabled bool     `json:"disabled"`

	// NOTE: Can only be used in SelectMenu with types.SelectMenuChannel.
	ChannelTypes []types.Channel `json:"channel_types,omitempty"`

	// Unique identifier for the Component; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`
}

func (s *SelectMenu) Type() types.Component {
	if s.MenuType != 0 {
		return types.Component(s.MenuType)
	}
	return types.ComponentSelectMenu
}

func (s *SelectMenu) MarshalJSON() ([]byte, error) {
	return toJson(s)
}

func (s *SelectMenu) message() {}

func (s *SelectMenu) modal() {}

// TextInput represents text input Component.
type TextInput struct {
	CustomID    string         `json:"custom_id"`
	Style       TextInputStyle `json:"style"`
	Placeholder string         `json:"placeholder,omitempty"`
	Value       string         `json:"value,omitempty"`
	Required    bool           `json:"required"`
	MinLength   int            `json:"min_length,omitempty"`
	MaxLength   int            `json:"max_length,omitempty"`

	// Unique identifier for the Component; autopopulated through increment if not provided.
	ID int `json:"id,omitempty"`
}

func (*TextInput) Type() types.Component {
	return types.ComponentTextInput
}

func (m *TextInput) MarshalJSON() ([]byte, error) {
	return toJson(m)
}

func (*TextInput) modal() {}

// TextInputStyle is style of text in TextInput Component.
type TextInputStyle uint

// Text styles
const (
	TextInputShort     TextInputStyle = 1
	TextInputParagraph TextInputStyle = 2
)

// Section is a top-level layout Component that allows you to join Message Component contextually with an Accessory.
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
	return toJson(s)
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
	return toJson(t)
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
	return toJson(t)
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
	return toJson(m)
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
	return toJson(f)
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
	return toJson(s)
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
	return toJson(c)
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

// Label is a top-level layout Component.
// It wraps modal components with text as a label and optional description.
type Label struct {
	// Unique identifier for the Component; autopopulated through increment if not provided.
	ID          int    `json:"id,omitempty"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Component   Modal  `json:"component"`
}

func (*Label) Type() types.Component {
	return types.ComponentLabel
}

func (l *Label) MarshalJSON() ([]byte, error) {
	return toJson(l)
}

func (*Label) modal() {}

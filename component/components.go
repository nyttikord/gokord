package component

import (
	"encoding/json"

	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
)

type Component interface {
	json.Marshaler
	Type() types.Component
}

// Message is a base interface for all message components.
type Message interface {
	json.Marshaler
	Type() types.Component
	message()
}

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

// ActionsRow is a top-level container component for displaying a row of interactive components.
type ActionsRow struct {
	// Can contain Button, SelectMenu and TextInput.
	// NOTE: maximum of 5.
	Components []Message `json:"components"`
	// Unique identifier for the component; auto populated through increment if not provided.
	ID int `json:"id,omitempty"`
}

// MarshalJSON is a method for marshaling ActionsRow to a JSON object.
func (r *ActionsRow) MarshalJSON() ([]byte, error) {
	return toJson(r)
}

// Type is a method to get the type component.
func (r *ActionsRow) Type() types.Component {
	return types.ComponentActionsRow
}

func (r *ActionsRow) message() {}

// ButtonStyle is style of button.
type ButtonStyle uint

// Button styles.
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

// Button represents button component.
type Button struct {
	Label    string           `json:"label"`
	Style    ButtonStyle      `json:"style"`
	Disabled bool             `json:"disabled"`
	Emoji    *emoji.Component `json:"emoji,omitempty"`

	// NOTE: Only button with ButtonStyleLink style can have link. Also, URL is mutually exclusive with CustomID.
	URL      string `json:"url,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
	// Identifier for a purchasable SKU. Only available when using premium-style buttons.
	SKUID string `json:"sku_id,omitempty"`
	// Unique identifier for the component; auto populated through increment if not provided.
	ID int `json:"id,omitempty"`
}

// MarshalJSON is a method for marshaling Button to a JSON object.
func (b Button) MarshalJSON() ([]byte, error) {
	if b.Style == 0 {
		b.Style = ButtonStylePrimary
	}
	return toJson(b)
}

// Type is a method to get the type of a component.
func (Button) Type() types.Component {
	return types.ComponentButton
}

func (b Button) message() {}

// SelectMenuOption represents an option for a select menu.
type SelectMenuOption struct {
	Label       string           `json:"label,omitempty"`
	Value       string           `json:"value"`
	Description string           `json:"description"`
	Emoji       *emoji.Component `json:"emoji,omitempty"`
	// Determines whenever option is selected by default or not.
	Default bool `json:"default"`
}

// SelectMenuDefaultValueType represents the type of an entity selected by default in auto-populated select menus.
type SelectMenuDefaultValueType string

// SelectMenuDefaultValue types.
const (
	SelectMenuDefaultValueUser    SelectMenuDefaultValueType = "user"
	SelectMenuDefaultValueRole    SelectMenuDefaultValueType = "role"
	SelectMenuDefaultValueChannel SelectMenuDefaultValueType = "channel"
)

// SelectMenuDefaultValue represents an entity selected by default in auto-populated select menus.
type SelectMenuDefaultValue struct {
	// ID of the entity.
	ID string `json:"id"`
	// Type of the entity.
	Type SelectMenuDefaultValueType `json:"type"`
}

// SelectMenuType represents select menu type.
type SelectMenuType types.Component

// SelectMenu types.
const (
	StringSelectMenu      = SelectMenuType(types.ComponentSelectMenu)
	UserSelectMenu        = SelectMenuType(types.ComponentUserSelectMenu)
	RoleSelectMenu        = SelectMenuType(types.ComponentRoleSelectMenu)
	MentionableSelectMenu = SelectMenuType(types.ComponentMentionableSelectMenu)
	ChannelSelectMenu     = SelectMenuType(types.ComponentChannelSelectMenu)
)

// SelectMenu represents select menu component.
type SelectMenu struct {
	// Type of the select menu.
	MenuType SelectMenuType `json:"type,omitempty"`
	// CustomID is a developer-defined identifier for the select menu.
	CustomID string `json:"custom_id,omitempty"`
	// The text which will be shown in the menu if there's no default options or all options was deselected and component was closed.
	Placeholder string `json:"placeholder"`
	// This value determines the minimal amount of selected items in the menu.
	MinValues *int `json:"min_values,omitempty"`
	// This value determines the maximal amount of selected items in the menu.
	// If MaxValues or MinValues are greater than one then the user can select multiple items in the component.
	MaxValues int `json:"max_values,omitempty"`
	// List of default values for auto-populated select menus.
	// NOTE: Number of entries should be in the range defined by MinValues and MaxValues.
	DefaultValues []SelectMenuDefaultValue `json:"default_values,omitempty"`

	Options []SelectMenuOption `json:"options,omitempty"`
	// The list of value(s) selected from the predefined options.
	// NOTE: This will only exist if the InteractionType was a ModalSubmit
	// otherwise you should (still) be using `MessageComponentData`
	Values   []string `json:"values,omitempty"`
	Disabled bool     `json:"disabled"`

	// NOTE: Can only be used in SelectMenu with Channel menu type.
	ChannelTypes []types.Channel `json:"channel_types,omitempty"`

	// Unique identifier for the component; auto populated through increment if not provided.
	ID int `json:"id,omitempty"`
}

// Type is a method to get the type of a component.
func (s SelectMenu) Type() types.Component {
	if s.MenuType != 0 {
		return types.Component(s.MenuType)
	}
	return types.ComponentSelectMenu
}

// MarshalJSON is a method for marshaling SelectMenu to a JSON object.
func (s SelectMenu) MarshalJSON() ([]byte, error) {
	return toJson(s)
}

func (s SelectMenu) message() {}

func (s SelectMenu) modal() {} // for StringSelectMenu

// TextInput represents text input component.
type TextInput struct {
	CustomID    string         `json:"custom_id"`
	Label       string         `json:"label"`
	Style       TextInputStyle `json:"style"`
	Placeholder string         `json:"placeholder,omitempty"`
	Value       string         `json:"value,omitempty"`
	Required    bool           `json:"required"`
	MinLength   int            `json:"min_length,omitempty"`
	MaxLength   int            `json:"max_length,omitempty"`

	// Unique identifier for the component; auto populated through increment if not provided.
	ID int `json:"id,omitempty"`
}

// Type is a method to get the type of a component.
func (TextInput) Type() types.Component {
	return types.ComponentTextInput
}

// MarshalJSON is a method for marshaling TextInput to a JSON object.
func (m TextInput) MarshalJSON() ([]byte, error) {
	return toJson(m)
}

func (TextInput) modal() {}

// TextInputStyle is style of text in TextInput component.
type TextInputStyle uint

// Text styles
const (
	TextInputShort     TextInputStyle = 1
	TextInputParagraph TextInputStyle = 2
)

// Section is a top-level layout component that allows you to join text contextually with an accessory.
type Section struct {
	// Unique identifier for the component; auto populated through increment if not provided.
	ID int `json:"id,omitempty"`
	// Array of text display components; max of 3.
	Components []Message `json:"components"`
	// Can be Button or Thumbnail
	Accessory Message `json:"accessory"`
}

// Type is a method to get the type of a component.
func (*Section) Type() types.Component {
	return types.ComponentSection
}

// MarshalJSON is a method for marshaling Section to a JSON object.
func (s *Section) MarshalJSON() ([]byte, error) {
	return toJson(s)
}

func (*Section) message() {}

// TextDisplay is a top-level component that allows you to add markdown-formatted text to the message.
type TextDisplay struct {
	Content string `json:"content"`
}

// Type is a method to get the type of a component.
func (TextDisplay) Type() types.Component {
	return types.ComponentTextDisplay
}

// MarshalJSON is a method for marshaling TextDisplay to a JSON object.
func (t TextDisplay) MarshalJSON() ([]byte, error) {
	return toJson(t)
}

func (TextDisplay) message() {}

// Thumbnail component can be used as an accessory for a section component.
type Thumbnail struct {
	// Unique identifier for the component; auto populated through increment if not provided.
	ID          int               `json:"id,omitempty"`
	Media       UnfurledMediaItem `json:"media"`
	Description *string           `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler,omitempty"`
}

// Type is a method to get the type of a component.
func (Thumbnail) Type() types.Component {
	return types.ComponentThumbnail
}

// MarshalJSON is a method for marshaling Thumbnail to a JSON object.
func (t Thumbnail) MarshalJSON() ([]byte, error) {
	return toJson(t)
}

func (Thumbnail) message() {}

// MediaGallery is a top-level component allows you to group images, videos or gifs into a gallery grid.
type MediaGallery struct {
	// Unique identifier for the component; auto populated through increment if not provided.
	ID int `json:"id,omitempty"`
	// Array of media gallery items; max of 10.
	Items []MediaGalleryItem `json:"items"`
}

// Type is a method to get the type of a component.
func (MediaGallery) Type() types.Component {
	return types.ComponentMediaGallery
}

// MarshalJSON is a method for marshaling MediaGallery to a JSON object.
func (m MediaGallery) MarshalJSON() ([]byte, error) {
	return toJson(m)
}

func (MediaGallery) message() {}

// MediaGalleryItem represents an item used in MediaGallery.
type MediaGalleryItem struct {
	Media       UnfurledMediaItem `json:"media"`
	Description *string           `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler"`
}

// FileComponent is a top-level component that allows you to display an uploaded file as an attachment to the message and reference it in the component.
type FileComponent struct {
	// Unique identifier for the component; auto populated through increment if not provided.
	ID      int               `json:"id,omitempty"`
	File    UnfurledMediaItem `json:"file"`
	Spoiler bool              `json:"spoiler"`
}

// Type is a method to get the type of a component.
func (FileComponent) Type() types.Component {
	return types.ComponentFile
}

// MarshalJSON is a method for marshaling FileComponent to a JSON object.
func (f FileComponent) MarshalJSON() ([]byte, error) {
	return toJson(f)
}

func (FileComponent) message() {}

// SeparatorSpacingSize represents spacing size around the separator.
type SeparatorSpacingSize uint

// Separator spacing sizes.
const (
	SeparatorSpacingSizeSmall SeparatorSpacingSize = 1
	SeparatorSpacingSizeLarge SeparatorSpacingSize = 2
)

// Separator is a top-level layout component that adds vertical padding and visual division between other components.
type Separator struct {
	// Unique identifier for the component; auto populated through increment if not provided.
	ID int `json:"id,omitempty"`

	Divider *bool                 `json:"divider,omitempty"`
	Spacing *SeparatorSpacingSize `json:"spacing,omitempty"`
}

// Type is a method to get the type of a component.
func (Separator) Type() types.Component {
	return types.ComponentSeparator
}

// MarshalJSON is a method for marshaling Separator to a JSON object.
func (s Separator) MarshalJSON() ([]byte, error) {
	return toJson(s)
}

func (Separator) message() {}

// Container is a top-level layout component.
// Containers are visually distinct from surrounding components and have an optional customizable color bar (similar to embeds).
type Container struct {
	// Unique identifier for the component; auto populated through increment if not provided.
	ID          int       `json:"id,omitempty"`
	AccentColor *int      `json:"accent_color,omitempty"`
	Spoiler     bool      `json:"spoiler"`
	Components  []Message `json:"components"`
}

// Type is a method to get the type of a component.
func (*Container) Type() types.Component {
	return types.ComponentContainer
}

// MarshalJSON is a method for marshaling Container to a JSON object.
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

// Unfurled media item loading states.
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

// Label is a top-level layout component.
// Labels wrap modal components with text as a label and optional description.
type Label struct {
	// Unique identifier for the component; auto populated through increment if not provided.
	ID          int     `json:"id,omitempty"`
	Label       string  `json:"label"`
	Description string  `json:"description,omitempty"`
	Component   Message `json:"component"`
}

// Type is a method to get the type of a component.
func (*Label) Type() types.Component {
	return types.ComponentLabel
}

// MarshalJSON is a method for marshaling Label to a JSON object.
func (l *Label) MarshalJSON() ([]byte, error) {
	return toJson(l)
}

func (*Label) modal() {}

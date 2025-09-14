package component

import (
	"encoding/json"

	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
	//"github.com/nyttikord/gokord/logger"
)

// Modal is implemented by all modal components.
type Modal interface {
	Component
	modal()
}

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
	return json.Marshal(struct {
		SelectMenu
		Type types.Component `json:"type"`
	}{
		SelectMenu: *s,
		Type:       s.Type(),
	})
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

func (t *TextInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TextInput
		Type types.Component `json:"type"`
	}{
		TextInput: *t,
		Type:      t.Type(),
	})
}

func (*TextInput) modal() {}

// TextInputStyle is style of text in TextInput Component.
type TextInputStyle uint

// Text styles
const (
	TextInputShort     TextInputStyle = 1
	TextInputParagraph TextInputStyle = 2
)

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
	return json.Marshal(struct {
		Label
		Type types.Component `json:"type"`
	}{
		Label: *l,
		Type:  l.Type(),
	})
}

func (l *Label) UnmarshalJSON(data []byte) error {
	println("label before")
	type t Label 
	var v struct {
		t 
		RawComponent Unmarshalable `json:"component"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	println("label after first unmarshalling")
	*l = Label(v.t)
	l.Component = v.RawComponent.Component.(Modal)
	println("label finished")
	return nil
}

func (*Label) modal() {}

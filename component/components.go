// Package component contains everything related to Component including Modal component and Message component.
package component

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nyttikord/gokord/discord/types"
)

var (
	ErrUnknownComponents = errors.New("unknown component")
)

// Component represents every component.
//
// NOTE to gokord contributors: when you are implementing a new component, don't forget to create a custom UnmarshalJSON
// if you are using an interface as a value type.
type Component interface {
	json.Marshaler
	Type() types.Component
}

// Unmarshaler is used to convert raw json bytes into a valid Component.
type Unmarshaler struct {
	Component
}

// UnmarshalJSON converts json bytes into a valid Component.
func (un *Unmarshaler) UnmarshalJSON(data []byte) error {
	var v struct {
		Type types.Component `json:"type"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch v.Type {
	case types.ComponentActionsRow:
		un.Component = &ActionsRow{}
	case types.ComponentButton:
		un.Component = &Button{}
	case types.ComponentSelectMenu, types.ComponentChannelSelectMenu, types.ComponentUserSelectMenu,
		types.ComponentRoleSelectMenu, types.ComponentMentionableSelectMenu:
		un.Component = &SelectMenu{}
	case types.ComponentTextInput:
		un.Component = &TextInput{}
	case types.ComponentSection:
		un.Component = &Section{}
	case types.ComponentTextDisplay:
		un.Component = &TextDisplay{}
	case types.ComponentThumbnail:
		un.Component = &Thumbnail{}
	case types.ComponentMediaGallery:
		un.Component = &MediaGallery{}
	case types.ComponentFile:
		un.Component = &File{}
	case types.ComponentSeparator:
		un.Component = &Separator{}
	case types.ComponentContainer:
		un.Component = &Container{}
	case types.ComponentLabel:
		un.Component = &Label{}
	default:
		return fmt.Errorf("unknown component type: %d", v.Type)
	}
	return json.Unmarshal(data, un.Component)
}

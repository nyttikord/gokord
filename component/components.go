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

// Unmarshaler is used to convert raw json bytes into a valid Component
type Unmarshaler struct {
	Component
}

// UnmarshalJSON converts json bytes into a valid Component
func (un *Unmarshaler) UnmarshalJSON(data []byte) error {
	var err error
	un.Component, err = unmarshalComponent(data)
	return err
}

func unmarshalComponent(data []byte) (Component, error) {
	var v struct {
		Type types.Component `json:"type"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	var c Component
	switch v.Type {
	case types.ComponentActionsRow:
		c = &ActionsRow{}
	case types.ComponentButton:
		c = &Button{}
	case types.ComponentSelectMenu, types.ComponentChannelSelectMenu, types.ComponentUserSelectMenu,
		types.ComponentRoleSelectMenu, types.ComponentMentionableSelectMenu:
		c = &SelectMenu{}
	case types.ComponentTextInput:
		c = &TextInput{}
	case types.ComponentSection:
		c = &Section{}
	case types.ComponentTextDisplay:
		c = &TextDisplay{}
	case types.ComponentThumbnail:
		c = &Thumbnail{}
	case types.ComponentMediaGallery:
		c = &MediaGallery{}
	case types.ComponentFile:
		c = &File{}
	case types.ComponentSeparator:
		c = &Separator{}
	case types.ComponentContainer:
		c = &Container{}
	case types.ComponentLabel:
		c = &Label{}
	default:
		return nil, fmt.Errorf("unknown component type: %d", v.Type)
	}
	return c, json.Unmarshal(data, c)
}

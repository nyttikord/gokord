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
type Component interface {
	json.Marshaler
	Type() types.Component
}

type Unmarshalable struct {
	Component
}

// UnmarshalJSON converts json bytes to a valid Component 
func (un *Unmarshalable) UnmarshalJSON(data []byte) error {
	var err error
	un.Component, err = unmarshalComponent(data)
	return err
}

func unmarshalComponent(data []byte) (Component, error) {
	var v struct {
		Type types.Component `json:"type"`
	}
	println("in before")
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	println("in after")

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
		return nil, fmt.Errorf("c.nown component type: %d", v.Type)
	}
	println("in last marshal")
	return c, json.Unmarshal(data, c)
}


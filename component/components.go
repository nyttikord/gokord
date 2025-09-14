package component

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/logger"
)

var (
	ErrUnknownComponents = errors.New("unknown component")
)

// Component represents every component.
type Component interface {
	json.Marshaler
	json.Unmarshaler
	Type() types.Component
}

func unmarshalComponent(data []byte) (Component, error) {
	logger.Log(logger.LevelDebug, 0, "called for %s", data)
	var t struct {
		Type types.Component `json:"type"`
	}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	var c Component
	switch t.Type {
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
		return nil, errors.Join(ErrUnknownComponents, fmt.Errorf("uknown type %d", t.Type))
	}
	err = json.Unmarshal(data, c)
	return c, err
}

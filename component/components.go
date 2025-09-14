package component

import (
	"encoding/json"

	"github.com/nyttikord/gokord/discord/types"
)

// Component represents every component.
type Component interface {
	json.Marshaler
	Type() types.Component
}

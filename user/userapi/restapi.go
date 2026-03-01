// Package userapi contains everything to interact with everything located in the user package.
package userapi

import (
	. "github.com/nyttikord/gokord/discord/request"
)

// Requester handles everything inside the user package.
type Requester struct {
	REST
	State *State
}

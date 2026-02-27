package guildapi

import (
	. "github.com/nyttikord/gokord/discord/request"
)

// Requester handles everything inside the guild package.
type Requester struct {
	REST
	Websocket
	State *State
}

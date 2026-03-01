// Package channelapi contains everything to interact with everything located in the channel package.
package channelapi

import (
	. "github.com/nyttikord/gokord/discord/request"
)

var (
// ErrReplyNilMessageRef = errors.New("reply attempted with nil message reference")
)

// Requester handles everything inside the channel package.
type Requester struct {
	REST
	State *State
}

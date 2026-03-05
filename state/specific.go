// Package state contains interfaces and variables used by every [State].
//
// You can get a state with [gokord.Session]:
//
//	var s *gokord.Session
//	s.GuildAPI().State // state related to guilds
package state

import (
	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/user"
)

// Bot represents the state related to gokord.Session (including if the session is not a bot).
type Bot interface {
	User() *user.User
	SessionID() string
	Shard() *[2]int
	Application() *application.Application
}

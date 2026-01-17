package state

import (
	"errors"
)

// ErrStateNotFound is returned when the state cache requested is not found
var ErrStateNotFound = errors.New("state cache not found")

// State represents the cache to prevent using too much requests.
type State interface {
	// MemberState returns the state of user.Member.
	MemberState() Member
	// ChannelState returns the state for channel package.
	ChannelState() Channel
	// GuildState returns the state for guild package and emoji package.
	GuildState() Guild
	// BotState returns the state of the Bot.
	BotState() Bot

	// Params returns the state.Params of the State
	Params() Params
}

// Params describes the parameters of the State.
type Params struct {
	// MaxMessageCount represents how many messages per channel the state will store.
	MaxMessageCount    int
	TrackChannels      bool
	TrackThreads       bool
	TrackEmojis        bool
	TrackStickers      bool
	TrackMembers       bool
	TrackThreadMembers bool
	TrackRoles         bool
	TrackVoice         bool
	TrackPresences     bool
}

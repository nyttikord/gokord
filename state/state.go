package state

import (
	"errors"
	"sync"
)

// ErrStateNotFound is returned when the state cache requested is not found
var ErrStateNotFound = errors.New("state cache not found")

// State represents the cache to prevent using too much requests.
type State interface {
	// GetMutex returns sync.RWMutex associated with the State.
	GetMutex() *sync.RWMutex

	// MemberState returns the State for user.Member.
	MemberState() Member
	// ChannelState returns the State for channel package.
	ChannelState() Channel
	// GuildState returns the State for guild package and emoji package.
	GuildState() Guild

	// GetMaxMessageCount returns how many messages per channel the State will store.
	GetMaxMessageCount() int
	// AreChannelsTracked returns true if the State must track channel.Channel.
	AreChannelsTracked() bool
	// AreThreadsTracked returns true if the State must track threads.
	AreThreadsTracked() bool
	// AreEmojisTracked returns true if the State must track emoji.Emoji.
	AreEmojisTracked() bool
	// AreStickersTracked returns true if the State must track emoji.Sticker.
	AreStickersTracked() bool
	// AreMembersTracked returns true if the State must track user.Member.
	AreMembersTracked() bool
	// AreThreadMembersTracked returns true if the State must track channel.ThreadMember.
	AreThreadMembersTracked() bool
	// AreRolesTracked returns true if the State must track guild.Role.
	AreRolesTracked() bool
	// AreVoiceTracked returns true if the State must track voice related things.
	AreVoiceTracked() bool
	// ArePresencesTracked returns true if the State must track status.Presence.
	ArePresencesTracked() bool
}

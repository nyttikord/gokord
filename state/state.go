package state

import (
	"errors"
	"sync"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
)

// ErrStateNotFound is returned when the state cache requested is not found
var ErrStateNotFound = errors.New("state cache not found")

// State represents the cache to prevent using too much requests.
type State interface {
	// GetMutex returns sync.RWMutex associated with the State.
	GetMutex() *sync.RWMutex

	// GetGuilds returns the guild.Guild cached.
	GetGuilds() []*guild.Guild

	// MemberAdd adds a user.Member to the State.
	MemberAdd(*user.Member) error
	// ChannelAdd adds a channel.Channel to the State.
	ChannelAdd(*channel.Channel) error
	// Guild returns the guild.Guild cached.
	Guild(string) (*guild.Guild, error)

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

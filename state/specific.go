// Package state contains interfaces and variables used by every State.
//
// You can get a state with gokord.Session:
//
//	var s *gokord.Session
//	s.GuildAPI().State // state related to guilds
package state

import (
	"iter"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
)

// Member represents the State related to user.Member (including status.Presence).
type Member interface {
	MemberAdd(*user.Member) error
	MemberRemove(*user.Member) error
	Member(string, string) (*user.Member, error)

	PresenceAdd(string, *status.Presence) error
	PresenceRemove(string, *status.Presence) error
	Presence(string, string) (*status.Presence, error)
}

// Channel represents the State related to channel.Channel (including channel.Message and threads).
type Channel interface {
	ChannelAdd(*channel.Channel) error
	ChannelRemove(*channel.Channel) error
	Channel(string) (*channel.Channel, error)
	PrivateChannels() []*channel.Channel

	MessageAdd(*channel.Message) error
	MessageRemove(*channel.Message) error
	MessageRemoveByID(string, string) error
	Message(string, string) (*channel.Message, error)

	ThreadListSync(string, []string, []*channel.Channel, []*channel.ThreadMember) error
	ThreadMembersUpdate(string, string, int, []channel.AddedThreadMember, []string) error
	ThreadMemberUpdate(*channel.ThreadMember) error

	// AppendGuildChannel is for internal use only.
	// Use ChannelAdd instead.
	AppendGuildChannel(c *channel.Channel)
}

// Guild represents the State related to guild.Guild (including guild.Role) and emoji.Emoji.
type Guild interface {
	GuildAdd(*guild.Guild)
	GuildRemove(*guild.Guild) error
	Guild(string) (*guild.Guild, error)
	Guilds() iter.Seq[*guild.Guild]

	RoleAdd(string, *guild.Role) error
	RoleRemove(string, string) error
	Role(string, string) (*guild.Role, error)
}

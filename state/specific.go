package state

import (
	"iter"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
)

type Member interface {
	MemberAdd(*user.Member) error
	MemberRemove(*user.Member) error
	Member(string, string) (*user.Member, error)

	PresenceAdd(string, *status.Presence) error
	PresenceRemove(string, *status.Presence) error
	Presence(string, string) (*status.Presence, error)
}

type Channel interface {
	ChannelAdd(*channel.Channel) error
	ChannelRemove(*channel.Channel) error
	Channel(string) (*channel.Channel, error)
	PrivateChannels() []*channel.Channel

	MessageAdd(*channel.Message) error
	MessageRemove(*channel.Message) error
	Message(string, string) (*channel.Message, error)
}

type Guild interface {
	GuildAdd(*guild.Guild)
	GuildRemove(*guild.Guild) error
	Guild(string) (*guild.Guild, error)
	Guilds() iter.Seq[*guild.Guild]

	RoleAdd(string, *guild.Role) error
	RoleRemove(string, string) error
	Role(string, string) (*guild.Role, error)
}

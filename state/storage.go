package state

import (
	"errors"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
)

// Key is the unique key used to store data in the Storage.
type Key string

const (
	KeyMemberPrefix  = "m:" // KeyMemberPrefix is the prefix used before each user.Member Key
	KeyGuildPrefix   = "g:" // KeyGuildPrefix is the prefix used before each guild.Guild Key
	KeyChannelPrefix = "c:" // KeyChannelPrefix is the prefix used before each Channel.Channel Key
)

var (
	ErrItemNotFound    = errors.New("item not found in storage")
	ErrInvalidDataType = errors.New("invalid data type")
)

// Storage represents a storage used to cache information.
// This is typically used by a State.
type Storage interface {
	// Get returns the data attached with the key in the Storage.
	//
	// Returns nil if the data was not found and throw the error.
	Get(key Key) (any, error)
	// Write the data in the Storage at the key location.
	Write(key Key, data any) error
}

// KeyMember returns the unique Key linked with the given user.Member.
func KeyMember(m *user.Member) Key {
	return KeyMemberRaw(m.GuildID, m.User.ID)
}

// KeyMemberRaw returns the unique Key linked with the user.Member described by the given parameters.
func KeyMemberRaw(guildID, userID string) Key {
	return KeyMemberPrefix + Key(guildID+":"+userID)
}

// KeyGuild returns the unique Key linked with the given guild.Guild.
func KeyGuild(g *guild.Guild) Key {
	return KeyGuildRaw(g.ID)
}

// KeyGuildRaw returns the unique Key linked with the guild.Guild described by the given parameter.
func KeyGuildRaw(guildID string) Key {
	return KeyGuildPrefix + Key(guildID)
}

// KeyChannel returns the unique Key linked with the given channel.Channel.
func KeyChannel(c *channel.Channel) Key {
	return KeyChannelRaw(c.GuildID, c.ID)
}

// KeyChannelRaw returns the unique Key linked with the channel.Channel described by the given parameters.
func KeyChannelRaw(guildID, channelID string) Key {
	return KeyChannelPrefix + Key(guildID+":"+channelID)
}

type MapStorage[T any] struct {
	storage map[Key]T
}

func (m *MapStorage[T]) Get(key Key) (any, error) {
	v, ok := m.storage[key]
	if !ok {
		return nil, ErrItemNotFound
	}
	return v, nil
}

func (m *MapStorage[T]) Write(key Key, data any) error {
	v, ok := data.(T)
	if !ok {
		return ErrInvalidDataType
	}
	m.storage[key] = v
	return nil
}

package state

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/nyttikord/avl"
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
	ErrInvalidDataType = errors.New("invalid data type")
)

// Storage represents a storage used to cache information.
// This is typically used by a State.
//
// When a data is saved in the Storage, it cannot be modified without calling Write.
// The content of the Storage must be immutable.
// Thus, do not store pointers!
type Storage[T any] interface {
	// Get returns the data attached with the key in the Storage.
	// It should never return a pointer to a struct.
	//
	// Returns nil if the data was not found and throw the error.
	Get(key Key) (T, error)
	// Write the data in the Storage at the key location.
	Write(key Key, data T) error
	// Delete a value associated with the key.
	//
	// Does not return an error if the value was not present.
	Delete(key Key) error
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
	return KeyChannelRaw(c.ID)
}

// KeyChannelRaw returns the unique Key linked with the channel.Channel described by the given parameter.
func KeyChannelRaw(channelID string) Key {
	return KeyChannelPrefix + Key(channelID)
}

// MapStorage is a standard implementation of Storage used.
// It uses a Go map to store data.
//
// See AVLStorage for the default implementation used.
type MapStorage[T any] map[Key]T

// deepCopy is an ugly code performing a deep copy.
// It works, but feel free to refactor it if you have any better idea.
// (The proposition https://github.com/golang-design/reflect using reflect package does not work sadly...)
func deepCopy[T any](t T) (T, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return t, err
	}
	var copied T
	err = json.Unmarshal(b, &copied)
	return copied, err
}

func (m MapStorage[T]) Get(key Key) (T, error) {
	v, ok := m[key]
	if !ok {
		return v, ErrNotFound
	}
	return deepCopy(v)
}

func (m MapStorage[T]) Write(key Key, data T) error {
	var err error
	m[key], err = deepCopy(data)
	return err
}

func (m MapStorage[T]) Delete(key Key) error {
	delete(m, key)
	return nil
}

// AVLStorage is a standard implementation of Storage used if no implementation is given.
// It uses an AVL (self-balancing binary search tree) to store data.
//
// See MapStorage for another standard implementation of Storage.
type AVLStorage[T any] struct {
	tree *avl.KeyAVL[Key, T]
}

// NewAVLStorage creates a new AVLStorage.
func NewAVLStorage[T any]() *AVLStorage[T] {
	tree := avl.NewKeyImmutable(func(k1, k2 Key) int {
		return strings.Compare(string(k1), string(k2))
	}, func(v T) T {
		cp, err := deepCopy(v)
		if err != nil {
			panic(err)
		}
		return cp
	})
	return &AVLStorage[T]{tree: tree}
}

func (a *AVLStorage[T]) Get(key Key) (v T, err error) {
	tv := a.tree.Get(key)
	if tv == nil {
		err = ErrNotFound
	} else {
		v = *tv
	}
	return
}

func (a *AVLStorage[T]) Write(key Key, data T) error {
	a.tree.Insert(key, data)
	return nil
}

func (a *AVLStorage[T]) Delete(key Key) error {
	a.tree.Delete(key)
	return nil
}

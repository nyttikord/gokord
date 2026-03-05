package state

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/nyttikord/avl"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
)

var (
	ErrInvalidDataType = errors.New("invalid data type")
)

// Storage represents a storage used to cache information.
// This is typically used by a [State].
//
// When a data is saved in the Storage, it cannot be modified without calling [Storage.Write].
// The content of the Storage must be immutable.
// Thus, do not store pointers!
type Storage[K, V any] interface {
	// Get returns the data attached with the key in the [Storage].
	// It should never return a pointer to a struct.
	//
	// Returns nil if the data was not found and throw an [ErrNotFound].
	Get(key K) (V, error)
	// Write the data in the [Storage] at the key location.
	Write(key K, data V) error
	// Delete a value associated with the key.
	//
	// Does not return an error if the value was not present.
	Delete(key K) error
}

func stringToUint(s string) uint64 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// KeyMember returns the unique key linked with the given [user.Member].
func KeyMember(m *user.Member) uint64 {
	return stringToUint(m.User.ID)
}

// KeyGuild returns the unique key linked with the given [guild.Guild].
func KeyGuild(g *guild.Guild) uint64 {
	return stringToUint(g.ID)
}

// KeyChannel returns the unique key linked with the given [channel.Channel].
func KeyChannel(c *channel.Channel) uint64 {
	return stringToUint(c.ID)
}

// MapStorage is a standard implementation of [Storage] used.
// It uses a Go map to store data.
//
// See [AVLStorage] for the default implementation used.
type MapStorage[K comparable, V any] map[K]V

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

func (m MapStorage[K, V]) Get(key K) (V, error) {
	v, ok := m[key]
	if !ok {
		return v, ErrNotFound
	}
	return deepCopy(v)
}

func (m MapStorage[K, V]) Write(key K, data V) error {
	var err error
	m[key], err = deepCopy(data)
	return err
}

func (m MapStorage[K, V]) Delete(key K) error {
	delete(m, key)
	return nil
}

// AVLStorage is a standard implementation of [Storage] used if no implementation is given.
// It uses an AVL (self-balancing binary search tree) to store data.
//
// See [MapStorage] for another standard implementation of [Storage].
type AVLStorage[K, V any] struct {
	tree *avl.KeyAVL[K, V]
}

// NewAVLStorage creates a new AVLStorage.
func NewAVLStorage[K, V any](cmp avl.CompareFunc[K]) *AVLStorage[K, V] {
	tree := avl.NewKeyImmutable(cmp, func(v V) V {
		cp, err := deepCopy(v)
		if err != nil {
			panic(err)
		}
		return cp
	})
	return &AVLStorage[K, V]{tree: tree}
}

func (a *AVLStorage[K, V]) Get(key K) (v V, err error) {
	tv := a.tree.Get(key)
	if tv == nil {
		err = ErrNotFound
	} else {
		v = *tv
	}
	return
}

func (a *AVLStorage[K, V]) Write(key K, data V) error {
	a.tree.Insert(key, data)
	return nil
}

func (a *AVLStorage[K, V]) Delete(key K) error {
	a.tree.Delete(key)
	return nil
}

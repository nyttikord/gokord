package state

import (
	"slices"
	"sync"

	"github.com/nyttikord/avl"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/guild"
)

type GuildStorage Storage[uint64, guild.Guild]

type Guild struct {
	State
	mu      sync.RWMutex
	storage GuildStorage
	guilds  *avl.SimpleAVL[uint64]
	params  *Params
}

func NewGuild(state State, storage GuildStorage, params *Params) *Guild {
	return &Guild{
		State:   state,
		storage: storage,
		guilds:  avl.NewSimple[uint64](),
		params:  params,
	}
}

// KeyGuild returns the unique key linked with the given [guild.Guild].
func KeyGuild(g *guild.Guild) uint64 {
	return KeyGuildReverse(g.ID)
}

// KeyGuildReverse returns the key linked with the requested [guild.Guild].
func KeyGuildReverse(guildID uint64) uint64 {
	return guildID
}

// AddGuild adds a [guild.Guild] to the current [Guild] state, or updates it if it already exists.
func (s *Guild) AddGuild(g *guild.Guild) error {
	if gl, err := s.GetGuild(g.ID); err == nil {
		if g.MemberCount == 0 {
			g.MemberCount = gl.MemberCount
		}
		if g.Roles == nil {
			g.Roles = gl.Roles
		}
		if g.Emojis == nil {
			g.Emojis = gl.Emojis
		}
		if g.Members == nil {
			g.Members = gl.Members
		}
		if g.Presences == nil {
			g.Presences = gl.Presences
		}
		if g.Channels == nil {
			g.Channels = gl.Channels
		}
		if g.Threads == nil {
			g.Threads = gl.Threads
		}
		if g.VoiceStates == nil {
			g.VoiceStates = gl.VoiceStates
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Update the channels to point to the right gd
	for _, c := range g.Channels {
		if err := s.ChannelState().storage.Write(KeyChannel(c), *c); err != nil {
			return err
		}
	}

	// Add all the threads to the state in case of thread sync list.
	for _, t := range g.Threads {
		if err := s.ChannelState().storage.Write(KeyChannel(t), *t); err != nil {
			return err
		}
	}

	err := s.storage.Write(KeyGuild(g), *g)
	if err != nil {
		return err
	}
	if g.ID != 0 {
		s.guilds.Insert(g.ID)
	}
	return nil
}

// RemoveGuild removes a [guild.Guild] from current [Guild] state.
func (s *Guild) RemoveGuild(guild *guild.Guild) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.storage.Delete(KeyGuild(guild))
	if err != nil {
		return err
	}
	s.guilds.Delete(guild.ID)
	return nil
}

// GetGuild returns the [guild.Guild].
//
// Useful for querying if @me is in a [guild.Guild]:
//
//	_, err := s.GuildState().GetGuild(guildID)
//	isInGuild := !errors.Is(err, state.ErrStateNotFound)
func (s *Guild) GetGuild(guildID uint64) (*guild.Guild, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	g, err := s.storage.Get(KeyGuildReverse(guildID))
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// ListGuilds returns the sorted list of [guild.Guild]s ID.
func (s *Guild) ListGuilds() []uint64 {
	return s.guilds.Sort()
}

// AddRole adds a [guild.Role] to the current [Guild] state, or updates it if it already exists.
func (s *Guild) AddRole(guildID uint64, role *guild.Role) error {
	g, err := s.GetGuild(guildID)
	if err != nil {
		return err
	}

	if _, err = s.GetRole(guildID, role.ID); err == nil {
		id := slices.IndexFunc(g.Roles, func(r *guild.Role) bool { return r.ID == role.ID })
		g.Roles[id] = role
	} else {
		g.Roles = append(g.Roles, role)
	}

	return s.AddGuild(g)
}

// RemoveRole removes a [guild.Role] from current [Role] state.
func (s *Guild) RemoveRole(guildID, roleID uint64) error {
	g, err := s.GetGuild(guildID)
	if err != nil {
		return err
	}

	g.Roles = slices.DeleteFunc(g.Roles, func(r *guild.Role) bool { return r.ID == roleID })

	return s.AddGuild(g)
}

// GetRole returns the [guild.GetRole] from a [guild.Guild].
func (s *Guild) GetRole(guildID, roleID uint64) (*guild.Role, error) {
	g, err := s.GetGuild(guildID)
	if err != nil {
		return nil, err
	}

	for _, r := range g.Roles {
		if r.ID == roleID {
			return r, nil
		}
	}

	return nil, ErrNotFound
}

// GetEmoji returns an [emoji.Emoji] in the [guild.Guild].
func (s *Guild) GetEmoji(guildID, emojiID uint64) (*emoji.Emoji, error) {
	g, err := s.GetGuild(guildID)
	if err != nil {
		return nil, err
	}

	for _, e := range g.Emojis {
		if e.ID == emojiID {
			return e, nil
		}
	}

	return nil, ErrNotFound
}

// AddEmoji adds an [emoji.Emoji] to the current [Guild] state.
func (s *Guild) AddEmoji(guildID uint64, em *emoji.Emoji) error {
	g, err := s.GetGuild(guildID)
	if err != nil {
		return err
	}

	if _, err = s.GetEmoji(guildID, em.ID); err == nil {
		id := slices.IndexFunc(g.Emojis, func(e *emoji.Emoji) bool { return e.ID == em.ID })
		g.Emojis[id] = em
	} else {
		g.Emojis = append(g.Emojis, em)
	}

	return s.AddGuild(g)
}

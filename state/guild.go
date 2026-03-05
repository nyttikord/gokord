package state

import (
	"slices"
	"sync"

	"github.com/nyttikord/avl"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/guild"
)

type Guild struct {
	State
	mu      sync.RWMutex
	storage Storage[uint64, guild.Guild]
	guilds  *avl.SimpleAVL[string]
	params  *Params
}

func NewGuild(state State, storage Storage[uint64, guild.Guild], params *Params) *Guild {
	return &Guild{
		State:   state,
		storage: storage,
		guilds:  avl.NewString(),
		params:  params,
	}
}

// GuildAdd adds a guild.Guild to the current State, or updates it if it already exists.
func (s *Guild) GuildAdd(g *guild.Guild) error {
	if gl, err := s.Guild(g.ID); err == nil {
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
		if err := s.ChannelState().AppendGuildChannel(c); err != nil {
			return err
		}
	}

	// Add all the threads to the state in case of thread sync list.
	for _, t := range g.Threads {
		if err := s.ChannelState().AppendGuildChannel(t); err != nil {
			return err
		}
	}

	err := s.storage.Write(KeyGuild(g), *g)
	if err != nil {
		return err
	}
	if len(g.ID) > 0 {
		s.guilds.Insert(g.ID)
	}
	return nil
}

// GuildRemove removes a guild.Guild from current State.
func (s *Guild) GuildRemove(guild *guild.Guild) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.storage.Delete(KeyGuild(guild))
	if err != nil {
		return err
	}
	s.guilds.Delete(guild.ID)
	return nil
}

// Guild returns the guild.Guild.
//
// Useful for querying if @me is in a guild:
//
//	_, err := s.GuildState().Guild(guildID)
//	isInGuild := !errors.Is(err, state.ErrStateNotFound)
func (s *Guild) Guild(guildID string) (*guild.Guild, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	g, err := s.storage.Get(stringToUint(guildID))
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// Guilds returns the sorted list of guilds ID.
func (s *Guild) Guilds() []string {
	return s.guilds.Sort()
}

// RoleAdd adds a guild.Role to the current State, or updates it if it already exists.
func (s *Guild) RoleAdd(guildID string, role *guild.Role) error {
	g, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	if _, err = s.Role(guildID, role.ID); err == nil {
		id := slices.IndexFunc(g.Roles, func(r *guild.Role) bool { return r.ID == role.ID })
		g.Roles[id] = role
	} else {
		g.Roles = append(g.Roles, role)
	}

	return s.GuildAdd(g)
}

// RoleRemove removes a guild.Role from current State.
func (s *Guild) RoleRemove(guildID, roleID string) error {
	g, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	g.Roles = slices.DeleteFunc(g.Roles, func(r *guild.Role) bool { return r.ID == roleID })

	return s.GuildAdd(g)
}

// Role returns the guild.Role from a guild.Guild.
func (s *Guild) Role(guildID, roleID string) (*guild.Role, error) {
	g, err := s.Guild(guildID)
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

// Emoji returns an emoji.Emoji in the guild.Guild.
func (s *Guild) Emoji(guildID, emojiID string) (*emoji.Emoji, error) {
	g, err := s.Guild(guildID)
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

// EmojiAdd adds an emoji.Emoji to the current State.
func (s *Guild) EmojiAdd(guildID string, em *emoji.Emoji) error {
	g, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	if _, err = s.Emoji(guildID, em.ID); err == nil {
		id := slices.IndexFunc(g.Emojis, func(e *emoji.Emoji) bool { return e.ID == em.ID })
		g.Emojis[id] = em
	} else {
		g.Emojis = append(g.Emojis, em)
	}

	return s.GuildAdd(g)
}

// EmojisAdd adds multiple emoji.Emoji to the current State.
func (s *Guild) EmojisAdd(guildID string, emojis []*emoji.Emoji) error {
	for _, e := range emojis {
		if err := s.EmojiAdd(guildID, e); err != nil {
			return err
		}
	}
	return nil
}

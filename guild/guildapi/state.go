package guildapi

import (
	"slices"

	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/state"
)

type State struct {
	state.State
	storage state.Storage
	guilds  []string
}

func NewState(state state.State, storage state.Storage) *State {
	return &State{
		State:   state,
		storage: storage,
		guilds:  make([]string, 0),
	}
}

// GuildAdd adds a guild.Guild to the current State, or updates it if it already exists.
func (s *State) GuildAdd(g *guild.Guild) error {
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	// Update the channels to point to the right gd
	for _, c := range g.Channels {
		// no need to unlock here
		if err := s.ChannelState().AppendGuildChannel(c); err != nil {
			return err
		}
	}

	// Add all the threads to the state in case of thread sync list.
	for _, t := range g.Threads {
		// no need to unlock here
		if err := s.ChannelState().AppendGuildChannel(t); err != nil {
			return err
		}
	}

	err := s.storage.Write(state.KeyGuild(g), g)
	if err != nil {
		return err
	}
	s.guilds = append(s.guilds, g.ID)
	return nil
}

// GuildRemove removes a guild.Guild from current State.
func (s *State) GuildRemove(guild *guild.Guild) error {
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	err := s.storage.Delete(state.KeyGuild(guild))
	if err != nil {
		return err
	}
	slices.DeleteFunc(s.guilds, func(s string) bool { return s == guild.ID })
	return nil
}

// Guild returns the guild.Guild.
//
// Useful for querying if @me is in a guild:
//
//	   _, err := discordgo.Session.State.Application(guildID)
//		  isInGuild := errors.Is(err, state.ErrStateNotFound)
func (s *State) Guild(guildID string) (*guild.Guild, error) {
	s.GetMutex().RLock()
	defer s.GetMutex().RUnlock()

	gRaw, err := s.storage.Get(state.KeyGuildRaw(guildID))
	if err != nil {
		return nil, err
	}
	g := gRaw.(guild.Guild)

	return &g, nil
}

func (s *State) Guilds() []string {
	return s.guilds
}

// RoleAdd adds a guild.Role to the current State, or updates it if it already exists.
func (s *State) RoleAdd(guildID string, role *guild.Role) error {
	g, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	for i, r := range g.Roles {
		if r.ID == role.ID {
			g.Roles[i] = role
			return s.storage.Write(state.KeyGuild(g), g)
		}
	}

	g.Roles = append(g.Roles, role)

	return s.storage.Write(state.KeyGuild(g), g)
}

// RoleRemove removes a guild.Role from current State.
func (s *State) RoleRemove(guildID, roleID string) error {
	g, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	slices.DeleteFunc(g.Roles, func(r *guild.Role) bool { return r.ID == roleID })

	return s.storage.Write(state.KeyGuild(g), g)
}

// Role returns the guild.Role from a guild.Guild.
func (s *State) Role(guildID, roleID string) (*guild.Role, error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	s.GetMutex().RLock()
	defer s.GetMutex().RUnlock()

	for _, r := range g.Roles {
		if r.ID == roleID {
			return r, nil
		}
	}

	return nil, state.ErrStateNotFound
}

// Emoji returns an emoji.Emoji in the guild.Guild.
func (s *State) Emoji(guildID, emojiID string) (*emoji.Emoji, error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	s.GetMutex().RLock()
	defer s.GetMutex().RUnlock()

	for _, e := range g.Emojis {
		if e.ID == emojiID {
			return e, nil
		}
	}

	return nil, state.ErrStateNotFound
}

// EmojiAdd adds an emoji.Emoji to the current State.
func (s *State) EmojiAdd(guildID string, emoji *emoji.Emoji) error {
	g, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	for i, e := range g.Emojis {
		if e.ID == emoji.ID {
			g.Emojis[i] = emoji
			return s.storage.Write(state.KeyGuild(g), g)
		}
	}

	g.Emojis = append(g.Emojis, emoji)
	return s.storage.Write(state.KeyGuild(g), g)
}

// EmojisAdd adds multiple emoji.Emoji to the current State.
func (s *State) EmojisAdd(guildID string, emojis []*emoji.Emoji) error {
	for _, e := range emojis {
		if err := s.EmojiAdd(guildID, e); err != nil {
			return err
		}
	}
	return nil
}

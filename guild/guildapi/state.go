package guildapi

import (
	"iter"
	"maps"

	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/state"
)

type State struct {
	state.State
	guildMap map[string]*guild.Guild
}

func NewState(state state.State) *State {
	return &State{
		State:    state,
		guildMap: make(map[string]*guild.Guild),
	}
}

// GuildAdd adds a guild.Guild to the current State, or updates it if it already exists.
func (s *State) GuildAdd(guild *guild.Guild) {
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	// Update the channels to point to the right guild
	for _, c := range guild.Channels {
		s.ChannelState().AppendGuildChannel(c) // no need to unlock here
	}

	// Add all the threads to the state in case of thread sync list.
	for _, t := range guild.Threads {
		s.ChannelState().AppendGuildChannel(t) // no need to unlock here
	}

	g, ok := s.guildMap[guild.ID]
	if !ok {
		s.guildMap[guild.ID] = guild
		return
	}
	// We are about to replace `g` in the state with `guild`, but first we need to
	// make sure we preserve any fields that the `guild` doesn't contain from `g`.
	if guild.MemberCount == 0 {
		guild.MemberCount = g.MemberCount
	}
	if guild.Roles == nil {
		guild.Roles = g.Roles
	}
	if guild.Emojis == nil {
		guild.Emojis = g.Emojis
	}
	if guild.Members == nil {
		guild.Members = g.Members
	}
	if guild.Presences == nil {
		guild.Presences = g.Presences
	}
	if guild.Channels == nil {
		guild.Channels = g.Channels
	}
	if guild.Threads == nil {
		guild.Threads = g.Threads
	}
	if guild.VoiceStates == nil {
		guild.VoiceStates = g.VoiceStates
	}
	*g = *guild
}

// GuildRemove removes a guild.Guild from current State.
func (s *State) GuildRemove(guild *guild.Guild) error {
	_, err := s.Guild(guild.ID)
	if err != nil {
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	delete(s.guildMap, guild.ID)

	return nil
}

// Guild gets a guild by ID.
// Useful for querying if @me is in a guild:
//
//	   _, err := discordgo.Session.State.Application(guildID)
//		  isInGuild := err == nil
func (s *State) Guild(guildID string) (*guild.Guild, error) {
	s.GetMutex().RLock()
	defer s.GetMutex().RUnlock()

	if g, ok := s.guildMap[guildID]; ok {
		return g, nil
	}

	return nil, state.ErrStateNotFound
}

func (s *State) Guilds() iter.Seq[*guild.Guild] {
	return maps.Values(s.guildMap)
}

// RoleAdd adds a role to the current world state, or
// updates it if it already exists.
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
			return nil
		}
	}

	g.Roles = append(g.Roles, role)
	return nil
}

// RoleRemove removes a role from current world state by ID.
func (s *State) RoleRemove(guildID, roleID string) error {
	g, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	for i, r := range g.Roles {
		if r.ID == roleID {
			g.Roles = append(g.Roles[:i], g.Roles[i+1:]...)
			return nil
		}
	}

	return state.ErrStateNotFound
}

// Role gets a role by ID from a guild.
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

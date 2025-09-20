package userapi

import (
	"errors"

	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
)

type State struct {
	state.State
	memberMap map[string]map[string]*user.Member
}

var ErrGuildNotCached = errors.New("member's guild not cached")

func NewState(state state.State) *State {
	return &State{
		State:     state,
		memberMap: make(map[string]map[string]*user.Member),
	}
}

func (s *State) createMemberMap(g *guild.Guild) {
	members := make(map[string]*user.Member)
	for _, m := range g.Members {
		members[m.User.ID] = m
	}
	s.memberMap[g.ID] = members
}

// MemberAdd adds a member to the current world state, or updates it if it already exists.
func (s *State) MemberAdd(member *user.Member) error {
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	g, err := s.GuildState().Guild(member.GuildID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	members, ok := s.memberMap[member.GuildID]
	if !ok {
		s.createMemberMap(g)
	}

	m, ok := members[member.User.ID]
	if !ok {
		members[member.User.ID] = member
		g.Members = append(g.Members, member)
	} else {
		// We are about to replace `m` in the state with `member`, but first we need to
		// make sure we preserve any fields that the `member` doesn't contain from `m`.
		if member.JoinedAt.IsZero() {
			member.JoinedAt = m.JoinedAt
		}
		*m = *member
	}
	return nil
}

// MemberRemove removes a member from current world state.
func (s *State) MemberRemove(member *user.Member) error {
	g, err := s.GuildState().Guild(member.GuildID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	members, ok := s.memberMap[member.GuildID]
	if !ok {
		return state.ErrStateNotFound
	}

	_, ok = members[member.User.ID]
	if !ok {
		return state.ErrStateNotFound
	}
	delete(members, member.User.ID)

	for i, m := range g.Members {
		if m.User.ID == member.User.ID {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			return nil
		}
	}
	// this is technically not reachable
	return state.ErrStateNotFound
}

// Member gets a member by ID from a guild.
func (s *State) Member(guildID, userID string) (*user.Member, error) {
	s.GetMutex().RLock()
	defer s.GetMutex().RUnlock()

	members, ok := s.memberMap[guildID]
	if !ok {
		return nil, state.ErrStateNotFound
	}

	m, ok := members[userID]
	if ok {
		return m, nil
	}

	return nil, state.ErrStateNotFound
}

// PresenceAdd adds a presence to the current world state, or
// updates it if it already existuserapis.
func (s *State) PresenceAdd(guildID string, presence *status.Presence) error {
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	g, err := s.GuildState().Guild(guildID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	for i, p := range g.Presences {
		if p.User.ID == presence.User.ID {
			//g.Presences[i] = presence

			//Update status
			g.Presences[i].Activities = presence.Activities
			if presence.Status != "" {
				g.Presences[i].Status = presence.Status
			}
			if presence.ClientStatus.Desktop != "" {
				g.Presences[i].ClientStatus.Desktop = presence.ClientStatus.Desktop
			}
			if presence.ClientStatus.Mobile != "" {
				g.Presences[i].ClientStatus.Mobile = presence.ClientStatus.Mobile
			}
			if presence.ClientStatus.Web != "" {
				g.Presences[i].ClientStatus.Web = presence.ClientStatus.Web
			}

			//Update the optionally sent user information
			//ID Is a mandatory field so you should not need to check if it is empty
			g.Presences[i].User.ID = presence.User.ID

			if presence.User.Avatar != "" {
				g.Presences[i].User.Avatar = presence.User.Avatar
			}
			if presence.User.Discriminator != "" {
				g.Presences[i].User.Discriminator = presence.User.Discriminator
			}
			if presence.User.Email != "" {
				g.Presences[i].User.Email = presence.User.Email
			}
			if presence.User.Token != "" {
				g.Presences[i].User.Token = presence.User.Token
			}
			if presence.User.Username != "" {
				g.Presences[i].User.Username = presence.User.Username
			}

			return nil
		}
	}

	g.Presences = append(g.Presences, presence)
	return nil
}

// PresenceRemove removes a presence from the current world state.
func (s *State) PresenceRemove(guildID string, presence *status.Presence) error {
	g, err := s.GuildState().Guild(guildID)
	if err != nil {
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	for i, p := range g.Presences {
		if p.User.ID == presence.User.ID {
			g.Presences = append(g.Presences[:i], g.Presences[i+1:]...)
			return nil
		}
	}

	return state.ErrStateNotFound
}

// Presence gets a presence by ID from a guild.
func (s *State) Presence(guildID, userID string) (*status.Presence, error) {
	g, err := s.GuildState().Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, p := range g.Presences {
		if p.User.ID == userID {
			return p, nil
		}
	}

	return nil, state.ErrStateNotFound
}

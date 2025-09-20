package userapi

import (
	"errors"

	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user"
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

func (s *State) memberAdd(member *user.Member) error {
	g, err := s.Guild(member.GuildID)
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

// MemberAdd adds a member to the current world state, or updates it if it already exists.
func (s *State) MemberAdd(member *user.Member) error {
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	return s.memberAdd(member)
}

// MemberRemove removes a member from current world state.
func (s *State) MemberRemove(member *user.Member) error {
	g, err := s.Guild(member.GuildID)
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

package userapi

import (
	"errors"
	"slices"

	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
)

type State struct {
	state.State
	storage   state.Storage
	memberMap map[string]map[string]*user.Member
}

var ErrGuildNotCached = errors.New("member's guild not cached")

func NewState(state state.State, storage state.Storage) *State {
	return &State{
		State:     state,
		storage:   storage,
		memberMap: make(map[string]map[string]*user.Member),
	}
}

// MemberAdd adds a user.Member to the current State, or updates it if it already exists.
func (s *State) MemberAdd(member *user.Member) error {
	g, err := s.GuildState().Guild(member.GuildID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	if m, err := s.Member(member.GuildID, member.User.ID); err == nil {
		if member.JoinedAt.IsZero() {
			member.JoinedAt = m.JoinedAt
		}
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	id := slices.IndexFunc(g.Members, func(m *user.Member) bool { return m.User.ID == member.User.ID })
	if id == -1 {
		g.Members = append(g.Members, member)
	} else {
		g.Members[id] = member
	}
	err = s.GuildState().GuildAdd(g)
	if err != nil {
		return err
	}

	return s.storage.Write(state.KeyMember(member), member)
}

// MemberRemove removes a user.Member from current State.
func (s *State) MemberRemove(member *user.Member) error {
	_, err := s.Member(member.GuildID, member.User.ID)
	if err != nil {
		return err
	}
	g, err := s.GuildState().Guild(member.GuildID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	g.Members = slices.DeleteFunc(g.Members, func(m *user.Member) bool { return m.User.ID == member.User.ID })
	err = s.GuildState().GuildAdd(g)
	if err != nil {
		return err
	}

	return s.storage.Delete(state.KeyMember(member))
}

// Member returns the user.Member from a guild.Guild.
func (s *State) Member(guildID, userID string) (*user.Member, error) {
	s.GetMutex().RLock()
	defer s.GetMutex().RUnlock()

	mRaw, err := s.storage.Get(state.KeyMemberRaw(guildID, userID))
	if err != nil {
		return nil, err
	}
	m := mRaw.(user.Member)
	return &m, nil
}

// PresenceAdd adds a status.Presence to the current State, or updates it if it already exists.
func (s *State) PresenceAdd(guildID string, presence *status.Presence) error {
	g, err := s.GuildState().Guild(guildID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	if p, err := s.Presence(guildID, presence.User.ID); err == nil {
		if presence.Status == "" {
			presence.Status = p.Status
		}
		if presence.ClientStatus.Desktop == "" {
			presence.ClientStatus.Desktop = p.ClientStatus.Desktop
		}
		if presence.ClientStatus.Mobile == "" {
			presence.ClientStatus.Mobile = p.ClientStatus.Mobile
		}
		if presence.ClientStatus.Web == "" {
			presence.ClientStatus.Web = p.ClientStatus.Web
		}
		if presence.User.Avatar == "" {
			presence.User.Avatar = p.User.Avatar
		}
		if presence.User.Discriminator == "" {
			presence.User.Discriminator = p.User.Discriminator
		}
		if presence.User.Email == "" {
			presence.User.Email = p.User.Email
		}
		if presence.User.Token == "" {
			presence.User.Token = p.User.Token
		}
		if presence.User.Username == "" {
			presence.User.Username = p.User.Username
		}
		id := slices.IndexFunc(g.Presences, func(p *status.Presence) bool { return p.User.ID == presence.User.ID })
		g.Presences[id] = presence
	} else {
		g.Presences = append(g.Presences, presence)
	}
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	return nil
}

// PresenceRemove removes a status.Presence from the current State.
func (s *State) PresenceRemove(guildID string, presence *status.Presence) error {
	g, err := s.GuildState().Guild(guildID)
	if err != nil {
		return err
	}

	_, err = s.Presence(guildID, presence.User.ID)
	if err != nil {
		return err
	}

	g.Presences = slices.DeleteFunc(g.Presences, func(p *status.Presence) bool { return p.User.ID == presence.User.ID })

	err = s.GuildState().GuildAdd(g)
	if err != nil {
		return err
	}

	return nil
}

// Presence returns the status.Presence from a guild.Guild.
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

// UserColor returns the color of a user.User in a channel.Channel.
// While colors are defined at a guild.Guild level, determining for a channel.Channel is more useful in message handlers.
// Returns 0 in cases of error, which is the color of @everyone.
func (s *State) UserColor(userID, channelID string) int {
	c, err := s.ChannelState().Channel(channelID)
	if err != nil {
		return 0
	}

	g, err := s.GuildState().Guild(c.GuildID)
	if err != nil {
		return 0
	}

	member, err := s.Member(g.ID, userID)
	if err != nil {
		return 0
	}

	return guild.FirstRoleColor(g, member.Roles)
}

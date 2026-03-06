package state

import (
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
)

type MemberStorage Storage[string, user.Member]

type Member struct {
	State
	mu      sync.RWMutex
	storage MemberStorage
	params  *Params
}

var ErrMemberGuildNotCached = errors.New("member's guild not cached")

func NewMember(state State, storage MemberStorage, params *Params) *Member {
	return &Member{
		State:   state,
		storage: storage,
		params:  params,
	}
}

// KeyMember returns the unique key linked with the given [user.Member].
func KeyMember(m *user.Member) string {
	return KeyMemberReverse(m.GuildID, m.User.ID)
}

// KeyMemberReverse returns the key linked with the requested [user.Member].
func KeyMemberReverse(guildID, userID uint64) string {
	return fmt.Sprintf("%d:%d", guildID, userID)
}

// AddMember adds a [user.Member] to the current [Member] state, or updates it if it already exists.
func (s *Member) AddMember(member *user.Member) error {
	g, err := s.GuildState().GetGuild(member.GuildID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.Join(err, ErrMemberGuildNotCached)
		}
		return err
	}

	if m, err := s.GetMember(member.GuildID, member.User.ID); err == nil {
		if member.JoinedAt.IsZero() {
			member.JoinedAt = m.JoinedAt
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := slices.IndexFunc(g.Members, func(m *user.Member) bool { return m.User.ID == member.User.ID })
	if id == -1 {
		g.Members = append(g.Members, member)
	} else {
		g.Members[id] = member
	}
	err = s.GuildState().AddGuild(g)
	if err != nil {
		return err
	}

	return s.storage.Write(KeyMember(member), *member)
}

// RemoveMember removes a [user.Member] from current [Member] state.
func (s *Member) RemoveMember(member *user.Member) error {
	_, err := s.GetMember(member.GuildID, member.User.ID)
	if err != nil {
		return err
	}
	g, err := s.GuildState().GetGuild(member.GuildID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.Join(err, ErrMemberGuildNotCached)
		}
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	g.Members = slices.DeleteFunc(g.Members, func(m *user.Member) bool { return m.User.ID == member.User.ID })
	err = s.GuildState().AddGuild(g)
	if err != nil {
		return err
	}

	return s.storage.Delete(KeyMember(member))
}

// GetMember returns the [user.Member] from a [guild.Guild].
func (s *Member) GetMember(guildID, userID uint64) (*user.Member, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, err := s.storage.Get(KeyMemberReverse(guildID, userID))
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// AddPresence adds a [status.Presence] to the current [Member] state, or updates it if it already exists.
func (s *Member) AddPresence(guildID uint64, presence *status.Presence) error {
	g, err := s.GuildState().GetGuild(guildID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.Join(err, ErrMemberGuildNotCached)
		}
		return err
	}

	if p, err := s.GetPresence(guildID, presence.User.ID); err == nil {
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

	return s.GuildState().AddGuild(g)
}

// RemovePresence removes a [status.Presence] from the current [Member] state.
func (s *Member) RemovePresence(guildID uint64, presence *status.Presence) error {
	g, err := s.GuildState().GetGuild(guildID)
	if err != nil {
		return err
	}

	_, err = s.GetPresence(guildID, presence.User.ID)
	if err != nil {
		return err
	}

	g.Presences = slices.DeleteFunc(g.Presences, func(p *status.Presence) bool { return p.User.ID == presence.User.ID })

	return s.GuildState().AddGuild(g)
}

// GetPresence returns the [status.Presence] from a [guild.Guild].
func (s *Member) GetPresence(guildID, userID uint64) (*status.Presence, error) {
	g, err := s.GuildState().GetGuild(guildID)
	if err != nil {
		return nil, err
	}

	for _, p := range g.Presences {
		if p.User.ID == userID {
			return p, nil
		}
	}

	return nil, ErrNotFound
}

// GetUserColor returns the color of a [user.User] in a [channel.Channel].
// While colors are defined at a [guild.Guild] level, determining for a [channel.Channel] is more useful in message
// handlers.
// Returns 0 in cases of error, which is the color of @everyone.
func (s *Member) GetUserColor(userID, channelID uint64) int {
	c, err := s.ChannelState().GetChannel(channelID)
	if err != nil {
		return 0
	}

	g, err := s.GuildState().GetGuild(c.GuildID)
	if err != nil {
		return 0
	}

	member, err := s.GetMember(g.ID, userID)
	if err != nil {
		return 0
	}

	return guild.FirstRoleColor(g, member.Roles)
}

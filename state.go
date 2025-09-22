package gokord

import (
	"sync"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
)

// State contains the current known state.
type State struct {
	sync.RWMutex
	session     *Session
	User        *user.User
	SessionID   string
	Shard       *[2]int
	Application *application.Application

	// MaxMessageCount represents how many messages per channel the state will store.
	MaxMessageCount    int
	TrackChannels      bool
	TrackThreads       bool
	TrackEmojis        bool
	TrackStickers      bool
	TrackMembers       bool
	TrackThreadMembers bool
	TrackRoles         bool
	TrackVoice         bool
	TrackPresences     bool
}

func (s *State) GetMaxMessageCount() int {
	return s.MaxMessageCount
}

func (s *State) AreChannelsTracked() bool {
	return s.TrackChannels
}

func (s *State) AreThreadsTracked() bool {
	return s.TrackThreads
}

func (s *State) AreEmojisTracked() bool {
	return s.TrackEmojis
}

func (s *State) AreStickersTracked() bool {
	return s.TrackStickers
}

func (s *State) AreMembersTracked() bool {
	return s.TrackMembers
}

func (s *State) AreThreadMembersTracked() bool {
	return s.TrackThreadMembers
}

func (s *State) AreRolesTracked() bool {
	return s.TrackRoles
}

func (s *State) AreVoiceTracked() bool {
	return s.TrackVoice
}

func (s *State) ArePresencesTracked() bool {
	return s.TrackPresences
}

// GetMutex returns the sync.RWMutex of the State.
// You do not have to modify this.
func (s *State) GetMutex() *sync.RWMutex {
	return &s.RWMutex
}

// MemberState returns the state.State related to user.Member.
// Use Session.UserAPI().State instead.
func (s *State) MemberState() state.Member {
	return s.session.UserAPI().State
}

// ChannelState returns the state.State related to channel.Channel.
// Use Session.ChannelAPI().State instead.
func (s *State) ChannelState() state.Channel {
	return s.session.ChannelAPI().State
}

// GuildState returns the state.State related to guild.Guild.
// Use Session.GuildAPI().State instead.
func (s *State) GuildState() state.Guild {
	return s.session.GuildAPI().State
}

// NewState creates an empty state.
func NewState(s *Session) *State {
	return &State{
		session:            s,
		TrackChannels:      true,
		TrackThreads:       true,
		TrackEmojis:        true,
		TrackStickers:      true,
		TrackMembers:       true,
		TrackThreadMembers: true,
		TrackRoles:         true,
		TrackVoice:         true,
		TrackPresences:     true,
	}
}

func (s *State) voiceStateUpdate(update *event.VoiceStateUpdate) error {
	g, err := s.GuildState().Guild(update.GuildID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	// Handle Leaving Application
	if update.ChannelID == "" {
		for i, st := range g.VoiceStates {
			if st.UserID == update.UserID {
				g.VoiceStates = append(g.VoiceStates[:i], g.VoiceStates[i+1:]...)
				return nil
			}
		}
	} else {
		for i, st := range g.VoiceStates {
			if st.UserID == update.UserID {
				g.VoiceStates[i] = update.VoiceState
				return nil
			}
		}

		g.VoiceStates = append(g.VoiceStates, update.VoiceState)
	}

	return nil
}

// VoiceState gets a VoiceState by guild and user ID.
func (s *State) VoiceState(guildID, userID string) (*user.VoiceState, error) {
	g, err := s.GuildState().Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, st := range g.VoiceStates {
		if st.UserID == userID {
			return st, nil
		}
	}

	return nil, state.ErrStateNotFound
}

// OnReady takes a Ready event and updates all internal state.
func (s *State) onReady(se *Session, r *event.Ready) error {
	s.Lock()
	defer s.Unlock()

	// We must store the bare essentials like the current user.User or the SessionID.
	if !se.StateEnabled {
		s.SessionID = r.SessionID
		s.User = r.User
		s.Shard = r.Shard
		s.Application = r.Application

		return nil
	}

	for _, g := range r.Guilds {
		s.GuildState().GuildAdd(g)
	}

	for _, c := range r.PrivateChannels {
		if err := s.ChannelState().ChannelAdd(c); err != nil {
			return err
		}
	}

	return nil
}

// onInterface handles all events related to State.
func (s *State) onInterface(se *Session, i interface{}) error {
	r, ok := i.(*event.Ready)
	if ok {
		return s.onReady(se, r)
	}

	if !se.StateEnabled {
		return nil
	}

	var err error
	switch t := i.(type) {
	case *event.GuildCreate:
		s.GuildState().GuildAdd(t.Guild)
	case *event.GuildUpdate:
		s.GuildState().GuildAdd(t.Guild)
	case *event.GuildDelete:
		var old *guild.Guild
		old, err = s.GuildState().Guild(t.ID)
		if err == nil {
			oldCopy := *old
			t.BeforeDelete = &oldCopy
		}

		err = s.GuildState().GuildRemove(t.Guild)
	case *event.GuildMemberAdd:
		var g *guild.Guild
		// Updates the MemberCount of the guild.
		g, err = s.GuildState().Guild(t.Member.GuildID)
		if err != nil {
			return err
		}
		g.MemberCount++

		// Caches member if tracking is enabled.
		if s.TrackMembers {
			err = s.MemberState().MemberAdd(t.Member)
		}
	case *event.GuildMemberUpdate:
		if s.TrackMembers {
			var old *user.Member
			old, err = s.MemberState().Member(t.GuildID, t.User.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.MemberState().MemberAdd(t.Member)
		}
	case *event.GuildMemberRemove:
		var g *guild.Guild
		// Updates the MemberCount of the g.
		g, err = s.GuildState().Guild(t.Member.GuildID)
		if err != nil {
			return err
		}
		g.MemberCount--

		// Removes member from the cache if tracking is enabled.
		if s.TrackMembers {
			var old *user.Member
			old, err = s.MemberState().Member(t.Member.GuildID, t.Member.User.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}

			err = s.MemberState().MemberRemove(t.Member)
		}
	case *event.GuildMembersChunk:
		if s.TrackMembers {
			for i := range t.Members {
				t.Members[i].GuildID = t.GuildID
				err = s.MemberState().MemberAdd(t.Members[i])
			}
		}

		if s.TrackPresences {
			for _, p := range t.Presences {
				err = s.MemberState().PresenceAdd(t.GuildID, p)
			}
		}
	case *event.GuildRoleCreate:
		if s.TrackRoles {
			err = s.GuildState().RoleAdd(t.GuildID, t.Role)
		}
	case *event.GuildRoleUpdate:
		if s.TrackRoles {
			var old *guild.Role
			old, err = s.GuildState().Role(t.GuildID, t.Role.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.GuildState().RoleAdd(t.GuildID, t.Role)
		}
	case *event.GuildRoleDelete:
		if s.TrackRoles {
			var old *guild.Role
			old, err = s.GuildState().Role(t.GuildID, t.RoleID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}

			err = s.GuildState().RoleRemove(t.GuildID, t.RoleID)
		}
	case *event.GuildEmojisUpdate:
		if s.TrackEmojis {
			var g *guild.Guild
			g, err = s.GuildState().Guild(t.GuildID)
			if err != nil {
				return err
			}
			s.Lock()
			defer s.Unlock()
			g.Emojis = t.Emojis
		}
	case *event.GuildStickersUpdate:
		if s.TrackStickers {
			var g *guild.Guild
			g, err = s.GuildState().Guild(t.GuildID)
			if err != nil {
				return err
			}
			s.Lock()
			defer s.Unlock()
			g.Stickers = t.Stickers
		}
	case *event.ChannelCreate:
		if s.TrackChannels {
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *event.ChannelUpdate:
		if s.TrackChannels {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *event.ChannelDelete:
		if s.TrackChannels {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			err = s.ChannelState().ChannelRemove(t.Channel)
		}
	case *event.ThreadCreate:
		if s.TrackThreads {
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *event.ThreadUpdate:
		if s.TrackThreads {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *event.ThreadDelete:
		if s.TrackThreads {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			err = s.ChannelState().ChannelRemove(t.Channel)
		}
	case *event.ThreadMemberUpdate:
		if s.TrackThreads {
			err = s.ChannelState().ThreadMemberUpdate(t.ThreadMember)
		}
	case *event.ThreadMembersUpdate:
		if s.TrackThreadMembers {
			err = s.ChannelState().ThreadMembersUpdate(t.ID, t.GuildID, t.MemberCount, t.AddedMembers, t.RemovedMembers)
		}
	case *event.ThreadListSync:
		if s.TrackThreads {
			err = s.ChannelState().ThreadListSync(t.GuildID, t.ChannelIDs, t.Threads, t.Members)
		}
	case *event.MessageCreate:
		if s.MaxMessageCount != 0 {
			err = s.ChannelState().MessageAdd(t.Message)
		}
	case *event.MessageUpdate:
		if s.MaxMessageCount != 0 {
			var old *channel.Message
			old, err = s.ChannelState().Message(t.ChannelID, t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.ChannelState().MessageAdd(t.Message)
		}
	case *event.MessageDelete:
		if s.MaxMessageCount != 0 {
			var old *channel.Message
			old, err = s.ChannelState().Message(t.ChannelID, t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}

			err = s.ChannelState().MessageRemove(t.Message)
		}
	case *event.MessageDeleteBulk:
		if s.MaxMessageCount != 0 {
			for _, mID := range t.Messages {
				err = s.ChannelState().MessageRemoveByID(t.ChannelID, mID)
				if err != nil {
					return err
				}
			}
		}
	case *event.VoiceStateUpdate:
		if s.TrackVoice {
			var old *user.VoiceState
			old, err = s.VoiceState(t.GuildID, t.UserID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.voiceStateUpdate(t)
		}
	case *event.PresenceUpdate:
		if s.TrackPresences {
			err = s.MemberState().PresenceAdd(t.GuildID, &t.Presence)
		}
		if s.TrackMembers {
			if t.Status == status.Offline {
				return err
			}

			var m *user.Member
			m, err = s.MemberState().Member(t.GuildID, t.User.ID)

			if err != nil {
				// Member not found; this is a user coming online
				m = &user.Member{
					GuildID: t.GuildID,
					User:    t.User,
				}
			} else {
				if t.User.Username != "" {
					m.User.Username = t.User.Username
				}
			}

			err = s.MemberState().MemberAdd(m)
		}

	}

	return nil
}

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

// sessionState contains the current known state.
type sessionState struct {
	sync.RWMutex
	session     *Session
	user        *user.User
	sessionID   string
	shard       *[2]int
	application *application.Application

	params state.Params
}

func (s *sessionState) User() *user.User {
	return s.user
}

func (s *sessionState) SessionID() string {
	return s.sessionID
}

func (s *sessionState) Shard() *[2]int {
	return s.shard
}

func (s *sessionState) Application() *application.Application {
	return s.application
}

func (s *sessionState) GetMutex() *sync.RWMutex {
	return &s.RWMutex
}

func (s *sessionState) MemberState() state.Member {
	return s.session.UserAPI().State
}

func (s *sessionState) ChannelState() state.Channel {
	return s.session.ChannelAPI().State
}

func (s *sessionState) GuildState() state.Guild {
	return s.session.GuildAPI().State
}

func (s *sessionState) BotState() state.Bot {
	return s
}

func (s *sessionState) Params() state.Params {
	return s.params
}

// NewState creates an empty state.State.
func NewState(s *Session) state.State {
	return &sessionState{
		session: s,
		params: state.Params{
			TrackChannels:      true,
			TrackThreads:       true,
			TrackEmojis:        true,
			TrackStickers:      true,
			TrackMembers:       true,
			TrackThreadMembers: true,
			TrackRoles:         true,
			TrackVoice:         true,
			TrackPresences:     true,
		},
	}
}

func (s *sessionState) voiceStateUpdate(update *event.VoiceStateUpdate) error {
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
func (s *sessionState) VoiceState(guildID, userID string) (*user.VoiceState, error) {
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
func (s *sessionState) onReady(se *Session, r *event.Ready) error {
	s.Lock()
	defer s.Unlock()

	// We must store the bare essentials like the current user.User or the SessionID.
	if !se.StateEnabled {
		s.sessionID = r.SessionID
		s.user = r.User
		s.shard = r.Shard
		s.application = r.Application

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
func (s *sessionState) onInterface(se *Session, i interface{}) error {
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
		if s.params.TrackMembers {
			err = s.MemberState().MemberAdd(t.Member)
		}
	case *event.GuildMemberUpdate:
		if s.params.TrackMembers {
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
		if s.params.TrackMembers {
			var old *user.Member
			old, err = s.MemberState().Member(t.Member.GuildID, t.Member.User.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}

			err = s.MemberState().MemberRemove(t.Member)
		}
	case *event.GuildMembersChunk:
		if s.params.TrackMembers {
			for i := range t.Members {
				t.Members[i].GuildID = t.GuildID
				err = s.MemberState().MemberAdd(t.Members[i])
			}
		}

		if s.params.TrackPresences {
			for _, p := range t.Presences {
				err = s.MemberState().PresenceAdd(t.GuildID, p)
			}
		}
	case *event.GuildRoleCreate:
		if s.params.TrackRoles {
			err = s.GuildState().RoleAdd(t.GuildID, t.Role)
		}
	case *event.GuildRoleUpdate:
		if s.params.TrackRoles {
			var old *guild.Role
			old, err = s.GuildState().Role(t.GuildID, t.Role.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.GuildState().RoleAdd(t.GuildID, t.Role)
		}
	case *event.GuildRoleDelete:
		if s.params.TrackRoles {
			var old *guild.Role
			old, err = s.GuildState().Role(t.GuildID, t.RoleID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}

			err = s.GuildState().RoleRemove(t.GuildID, t.RoleID)
		}
	case *event.GuildEmojisUpdate:
		if s.params.TrackEmojis {
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
		if s.params.TrackStickers {
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
		if s.params.TrackChannels {
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *event.ChannelUpdate:
		if s.params.TrackChannels {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *event.ChannelDelete:
		if s.params.TrackChannels {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			err = s.ChannelState().ChannelRemove(t.Channel)
		}
	case *event.ThreadCreate:
		if s.params.TrackThreads {
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *event.ThreadUpdate:
		if s.params.TrackThreads {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *event.ThreadDelete:
		if s.params.TrackThreads {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			err = s.ChannelState().ChannelRemove(t.Channel)
		}
	case *event.ThreadMemberUpdate:
		if s.params.TrackThreads {
			err = s.ChannelState().ThreadMemberUpdate(t.ThreadMember)
		}
	case *event.ThreadMembersUpdate:
		if s.params.TrackThreadMembers {
			err = s.ChannelState().ThreadMembersUpdate(t.ID, t.GuildID, t.MemberCount, t.AddedMembers, t.RemovedMembers)
		}
	case *event.ThreadListSync:
		if s.params.TrackThreads {
			err = s.ChannelState().ThreadListSync(t.GuildID, t.ChannelIDs, t.Threads, t.Members)
		}
	case *event.MessageCreate:
		if s.params.MaxMessageCount != 0 {
			err = s.ChannelState().MessageAdd(t.Message)
		}
	case *event.MessageUpdate:
		if s.params.MaxMessageCount != 0 {
			var old *channel.Message
			old, err = s.ChannelState().Message(t.ChannelID, t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.ChannelState().MessageAdd(t.Message)
		}
	case *event.MessageDelete:
		if s.params.MaxMessageCount != 0 {
			var old *channel.Message
			old, err = s.ChannelState().Message(t.ChannelID, t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}

			err = s.ChannelState().MessageRemove(t.Message)
		}
	case *event.MessageDeleteBulk:
		if s.params.MaxMessageCount != 0 {
			for _, mID := range t.Messages {
				err = s.ChannelState().MessageRemoveByID(t.ChannelID, mID)
				if err != nil {
					return err
				}
			}
		}
	case *event.VoiceStateUpdate:
		if s.params.TrackVoice {
			var old *user.VoiceState
			old, err = s.VoiceState(t.GuildID, t.UserID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.voiceStateUpdate(t)
		}
	case *event.PresenceUpdate:
		if s.params.TrackPresences {
			err = s.MemberState().PresenceAdd(t.GuildID, &t.Presence)
		}
		if s.params.TrackMembers {
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

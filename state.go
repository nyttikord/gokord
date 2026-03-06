package gokord

import (
	"errors"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user"
)

// sessionState contains the current known state.
type sessionState struct {
	*Session
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

func (s *sessionState) SessionState() state.Bot {
	return s
}

// newState creates an empty [state.State].
func NewState(s *Session) state.State {
	return &sessionState{
		Session: s,
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

func (s *sessionState) voiceStateUpdate(update *event.VoiceStateUpdate) (err error) {
	var g *guild.Guild
	g, err = s.GuildState().GetGuild(update.GuildID)
	if err != nil {
		return
	}

	defer func() {
		if err == nil {
			err = s.GuildState().AddGuild(g)
		}
	}()

	// Handle Leaving Application
	if update.ChannelID == 0 {
		for i, st := range g.VoiceStates {
			if st.UserID == update.UserID {
				g.VoiceStates = append(g.VoiceStates[:i], g.VoiceStates[i+1:]...)
				return
			}
		}
	} else {
		for i, st := range g.VoiceStates {
			if st.UserID == update.UserID {
				g.VoiceStates[i] = update.VoiceState
				return
			}
		}

		g.VoiceStates = append(g.VoiceStates, update.VoiceState)
	}
	return
}

// VoiceState gets a VoiceState by guild and user ID.
func (s *sessionState) VoiceState(guildID, userID uint64) (*user.VoiceState, error) {
	g, err := s.GuildState().GetGuild(guildID)
	if err != nil {
		return nil, err
	}

	for _, st := range g.VoiceStates {
		if st.UserID == userID {
			return st, nil
		}
	}

	return nil, state.ErrNotFound
}

// OnReady takes a Ready event and updates all internal state.
func (s *sessionState) onReady(se *Session, r *event.Ready) error {
	// We must store the bare essentials like the current user.User or the SessionID.
	// assuming that a mutex is not needed for this because it is always the first handled called
	s.sessionID = r.SessionID
	s.user = r.User
	s.shard = r.Shard
	s.application = r.Application

	if !se.Options.StateEnabled {
		return nil
	}

	for _, g := range r.Guilds {
		err := s.GuildState().AddGuild(g)
		if err != nil {
			return err
		}
	}

	for _, c := range r.PrivateChannels {
		if err := s.ChannelState().AddChannel(c); err != nil {
			return err
		}
	}

	return nil
}

// onInterface handles all events related to State.
func (s *sessionState) onInterface(se *Session, i any) error {
	r, ok := i.(*event.Ready)
	if ok {
		return s.onReady(se, r)
	}

	if !se.Options.StateEnabled {
		return nil
	}

	switch t := i.(type) {
	case *event.GuildCreate:
		return s.GuildState().AddGuild(t.Guild)
	case *event.GuildUpdate:
		return s.GuildState().AddGuild(t.Guild)
	case *event.GuildDelete:
		old, err := s.GuildState().GetGuild(t.ID)
		if err == nil {
			oldCopy := *old
			t.BeforeDelete = &oldCopy
		}
		return s.GuildState().RemoveGuild(t.Guild)
	case *event.GuildMemberAdd:
		// Updates the MemberCount of the guild.
		g, err := s.GuildState().GetGuild(t.Member.GuildID)
		if err != nil {
			return err
		}
		g.MemberCount++

		// Caches member if tracking is enabled.
		if s.params.TrackMembers {
			err = s.MemberState().AddMember(t.Member)
			if err != nil {
				return err
			}
		}
		return s.GuildState().AddGuild(g)
	case *event.GuildMemberUpdate:
		if s.params.TrackMembers {
			old, err := s.MemberState().GetMember(t.GuildID, t.User.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			return s.MemberState().AddMember(t.Member)
		}
	case *event.GuildMemberRemove:
		// Updates the MemberCount of the g.
		g, err := s.GuildState().GetGuild(t.Member.GuildID)
		if err != nil {
			return err
		}
		g.MemberCount--

		// Removes member from the cache if tracking is enabled.
		if s.params.TrackMembers {
			var old *user.Member
			old, err = s.MemberState().GetMember(t.Member.GuildID, t.Member.User.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			err = s.MemberState().RemoveMember(t.Member)
			if err != nil {
				return err
			}
		}
		return s.GuildState().AddGuild(g)
	case *event.GuildMembersChunk:
		if s.params.TrackMembers {
			for i := range t.Members {
				t.Members[i].GuildID = t.GuildID
				return s.MemberState().AddMember(t.Members[i])
			}
		}

		if s.params.TrackPresences {
			for _, p := range t.Presences {
				err := s.MemberState().AddPresence(t.GuildID, p)
				if err != nil {
					return err
				}
			}
		}
	case *event.GuildRoleCreate:
		if s.params.TrackRoles {
			return s.GuildState().AddRole(t.GuildID, t.Role)
		}
	case *event.GuildRoleUpdate:
		if s.params.TrackRoles {
			old, err := s.GuildState().GetRole(t.GuildID, t.Role.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			return s.GuildState().AddRole(t.GuildID, t.Role)
		}
	case *event.GuildRoleDelete:
		if s.params.TrackRoles {
			old, err := s.GuildState().GetRole(t.GuildID, t.RoleID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			return s.GuildState().RemoveRole(t.GuildID, t.RoleID)
		}
	case *event.GuildEmojisUpdate:
		if s.params.TrackEmojis {
			g, err := s.GuildState().GetGuild(t.GuildID)
			if err != nil {
				return err
			}
			g.Emojis = t.Emojis
			return s.GuildState().AddGuild(g)
		}
	case *event.GuildStickersUpdate:
		if s.params.TrackStickers {
			g, err := s.GuildState().GetGuild(t.GuildID)
			if err != nil {
				return err
			}
			g.Stickers = t.Stickers
			return s.GuildState().AddGuild(g)
		}
	case *event.ChannelCreate:
		if s.params.TrackChannels {
			return s.ChannelState().AddChannel(t.Channel)
		}
	case *event.ChannelUpdate:
		if s.params.TrackChannels {
			old, err := s.ChannelState().GetChannel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			return s.ChannelState().AddChannel(t.Channel)
		}
	case *event.ChannelDelete:
		if s.params.TrackChannels {
			old, err := s.ChannelState().GetChannel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			return s.ChannelState().RemoveChannel(t.Channel)
		}
	case *event.ThreadCreate:
		if s.params.TrackThreads {
			return s.ChannelState().AddChannel(t.Channel)
		}
	case *event.ThreadUpdate:
		if s.params.TrackThreads {
			old, err := s.ChannelState().GetChannel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			return s.ChannelState().AddChannel(t.Channel)
		}
	case *event.ThreadDelete:
		if s.params.TrackThreads {
			old, err := s.ChannelState().GetChannel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			return s.ChannelState().RemoveChannel(t.Channel)
		}
	case *event.ThreadMemberUpdate:
		if s.params.TrackThreads {
			return s.ChannelState().ThreadMemberUpdate(t.ThreadMember)
		}
	case *event.ThreadMembersUpdate:
		if s.params.TrackThreadMembers {
			return s.ChannelState().OnThreadMembersUpdate(t.ID, t.GuildID, t.MemberCount, t.AddedMembers, t.RemovedMembers)
		}
	case *event.ThreadListSync:
		if s.params.TrackThreads {
			return s.ChannelState().OnThreadListSync(t.GuildID, t.ChannelIDs, t.Threads, t.Members)
		}
	case *event.MessageCreate:
		if s.params.MaxMessageCount != 0 {
			return s.ChannelState().AddMessage(t.Message)
		}
	case *event.MessageUpdate:
		if s.params.MaxMessageCount != 0 {
			old, err := s.ChannelState().GetMessage(t.ChannelID, t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			return s.ChannelState().AddMessage(t.Message)
		}
	case *event.MessageDelete:
		if s.params.MaxMessageCount != 0 {
			old, err := s.ChannelState().GetMessage(t.ChannelID, t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			return s.ChannelState().RemoveMessage(t.Message)
		}
	case *event.MessageDeleteBulk:
		if s.params.MaxMessageCount != 0 {
			for _, mID := range t.Messages {
				err := s.ChannelState().RemoveMessageByID(t.ChannelID, mID)
				if err != nil {
					return err
				}
			}
		}
	case *event.VoiceStateUpdate:
		if s.params.TrackVoice {
			old, err := s.VoiceState(t.GuildID, t.UserID)
			if err == nil {
				t.BeforeUpdate = old
			} else if !errors.Is(err, state.ErrNotFound) {
				s.logger.Error("fetching before state", "error", err)
			}
			return s.voiceStateUpdate(t)
		}
	case *event.PresenceUpdate:
		if s.params.TrackPresences {
			err := s.MemberState().AddPresence(t.GuildID, &t.Presence)
			if err != nil {
				return err
			}
		}
		if s.params.TrackMembers {
			m, err := s.MemberState().GetMember(t.GuildID, t.User.ID)
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
			return s.MemberState().AddMember(m)
		}
	}
	return nil
}

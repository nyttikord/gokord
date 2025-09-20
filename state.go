package gokord

import (
	"errors"
	"sort"
	"sync"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
)

// ErrNilState is returned when the state is nil.
var ErrNilState = errors.New("state not instantiated")

// ErrStateNotFound is returned when the state cache requested is not found
var ErrStateNotFound = errors.New("state cache not found")

// ErrMessageIncompletePermissions is returned when the message requested for permissions does not contain enough data to
// generate the permissions.
var ErrMessageIncompletePermissions = errors.New("message incomplete, unable to determine permissions")

// A State contains the current known state.
// As discord sends this in a READY blob, it seems reasonable to simply use that struct as the data store.
type State struct {
	sync.RWMutex
	Ready
	session Session

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

func (s *State) GetMutex() *sync.RWMutex {
	return &s.RWMutex
}

func (s *State) MemberState() state.Member {
	return s.session.UserAPI().State
}

func (s *State) ChannelState() state.Channel {
	return s.session.ChannelAPI().State
}

func (s *State) GuildState() state.Guild {
	return s.session.GuildAPI().State
}

// NewState creates an empty state.
func NewState() *State {
	return &State{
		Ready: Ready{
			PrivateChannels: []*channel.Channel{},
			Guilds:          []*guild.Guild{},
		},
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

// Emoji returns an emoji for a guild and emoji id.
func (s *State) Emoji(guildID, emojiID string) (*emoji.Emoji, error) {
	if s == nil {
		return nil, ErrNilState
	}

	g, err := s.GuildState().Guild(guildID)
	if err != nil {
		return nil, err
	}

	s.RLock()
	defer s.RUnlock()

	for _, e := range g.Emojis {
		if e.ID == emojiID {
			return e, nil
		}
	}

	return nil, ErrStateNotFound
}

// EmojiAdd adds an emoji to the current world state.
func (s *State) EmojiAdd(guildID string, emoji *emoji.Emoji) error {
	if s == nil {
		return ErrNilState
	}

	g, err := s.GuildState().Guild(guildID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	for i, e := range g.Emojis {
		if e.ID == emoji.ID {
			g.Emojis[i] = emoji
			return nil
		}
	}

	g.Emojis = append(g.Emojis, emoji)
	return nil
}

// EmojisAdd adds multiple emojis to the world state.
func (s *State) EmojisAdd(guildID string, emojis []*emoji.Emoji) error {
	for _, e := range emojis {
		if err := s.EmojiAdd(guildID, e); err != nil {
			return err
		}
	}
	return nil
}

func (s *State) voiceStateUpdate(update *VoiceStateUpdate) error {
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
	if s == nil {
		return nil, ErrNilState
	}

	g, err := s.GuildState().Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, st := range g.VoiceStates {
		if st.UserID == userID {
			return st, nil
		}
	}

	return nil, ErrStateNotFound
}

// OnReady takes a Ready event and updates all internal state.
func (s *State) onReady(se *Session, r *Ready) (err error) {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	// We must track at least the current user for Voice, even
	// if state is disabled, store the bare essentials.
	if !se.StateEnabled {
		ready := Ready{
			Version:     r.Version,
			SessionID:   r.SessionID,
			User:        r.User,
			Shard:       r.Shard,
			Application: r.Application,
		}

		s.Ready = ready

		return nil
	}

	s.Ready = *r

	for _, g := range s.Guilds {
		s.guildMap[g.ID] = g
		s.createMemberMap(g)

		for _, c := range g.Channels {
			s.channelMap[c.ID] = c
		}
	}

	for _, c := range s.PrivateChannels {
		s.channelMap[c.ID] = c
	}

	return nil
}

// OnInterface handles all events related to states.
func (s *State) OnInterface(se *Session, i interface{}) error {
	if s == nil {
		return ErrNilState
	}

	r, ok := i.(*Ready)
	if ok {
		return s.onReady(se, r)
	}

	if !se.StateEnabled {
		return nil
	}

	var err error
	switch t := i.(type) {
	case *GuildCreate:
		s.GuildState().GuildAdd(t.Guild)
	case *GuildUpdate:
		s.GuildState().GuildAdd(t.Guild)
	case *GuildDelete:
		var old *guild.Guild
		old, err = s.GuildState().Guild(t.ID)
		if err == nil {
			oldCopy := *old
			t.BeforeDelete = &oldCopy
		}

		err = s.GuildState().GuildRemove(t.Guild)
	case *GuildMemberAdd:
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
	case *GuildMemberUpdate:
		if s.TrackMembers {
			var old *user.Member
			old, err = s.MemberState().Member(t.GuildID, t.User.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.MemberState().MemberAdd(t.Member)
		}
	case *GuildMemberRemove:
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
	case *GuildMembersChunk:
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
	case *GuildRoleCreate:
		if s.TrackRoles {
			err = s.GuildState().RoleAdd(t.GuildID, t.Role)
		}
	case *GuildRoleUpdate:
		if s.TrackRoles {
			var old *guild.Role
			old, err = s.GuildState().Role(t.GuildID, t.Role.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.GuildState().RoleAdd(t.GuildID, t.Role)
		}
	case *GuildRoleDelete:
		if s.TrackRoles {
			var old *guild.Role
			old, err = s.GuildState().Role(t.GuildID, t.RoleID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}

			err = s.GuildState().RoleRemove(t.GuildID, t.RoleID)
		}
	case *GuildEmojisUpdate:
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
	case *GuildStickersUpdate:
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
	case *ChannelCreate:
		if s.TrackChannels {
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *ChannelUpdate:
		if s.TrackChannels {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *ChannelDelete:
		if s.TrackChannels {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			err = s.ChannelState().ChannelRemove(t.Channel)
		}
	case *ThreadCreate:
		if s.TrackThreads {
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *ThreadUpdate:
		if s.TrackThreads {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}
			err = s.ChannelState().ChannelAdd(t.Channel)
		}
	case *ThreadDelete:
		if s.TrackThreads {
			var old *channel.Channel
			old, err = s.ChannelState().Channel(t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}
			err = s.ChannelState().ChannelRemove(t.Channel)
		}
	case *ThreadMemberUpdate:
		if s.TrackThreads {
			err = s.ThreadMemberUpdate(t)
		}
	case *ThreadMembersUpdate:
		if s.TrackThreadMembers {
			err = s.ThreadMembersUpdate(t)
		}
	case *ThreadListSync:
		if s.TrackThreads {
			err = s.ThreadListSync(t)
		}
	case *MessageCreate:
		if s.MaxMessageCount != 0 {
			err = s.ChannelState().MessageAdd(t.Message)
		}
	case *MessageUpdate:
		if s.MaxMessageCount != 0 {
			var old *channel.Message
			old, err = s.ChannelState().Message(t.ChannelID, t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.ChannelState().MessageAdd(t.Message)
		}
	case *MessageDelete:
		if s.MaxMessageCount != 0 {
			var old *channel.Message
			old, err = s.ChannelState().Message(t.ChannelID, t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}

			err = s.ChannelState().MessageRemove(t.Message)
		}
	case *MessageDeleteBulk:
		if s.MaxMessageCount != 0 {
			for _, mID := range t.Messages {
				s.ChannelState().messageRemoveByID(t.ChannelID, mID)
			}
		}
	case *VoiceStateUpdate:
		if s.TrackVoice {
			var old *user.VoiceState
			old, err = s.VoiceState(t.GuildID, t.UserID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.voiceStateUpdate(t)
		}
	case *PresenceUpdate:
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

// UserChannelPermissions returns the permission of a user in a channel.
// userID    : The ID of the user to calculate permissions for.
// channelID : The ID of the channel to calculate permission for.
func (s *State) UserChannelPermissions(userID, channelID string) (apermissions int64, err error) {
	if s == nil {
		return 0, ErrNilState
	}

	c, err := s.ChannelState().Channel(channelID)
	if err != nil {
		return
	}

	g, err := s.GuildState().Guild(c.GuildID)
	if err != nil {
		return
	}

	member, err := s.MemberState().Member(g.ID, userID)
	if err != nil {
		return
	}

	return MemberPermissions(g, c, userID, member.Roles), nil
}

// MessagePermissions returns the permissions of the author of the message
// in the channel in which it was sent.
func (s *State) MessagePermissions(message *channel.Message) (apermissions int64, err error) {
	if s == nil {
		return 0, ErrNilState
	}

	if message.Author == nil || message.Member == nil {
		return 0, ErrMessageIncompletePermissions
	}

	c, err := s.ChannelState().Channel(message.ChannelID)
	if err != nil {
		return
	}

	g, err := s.GuildState().Guild(c.GuildID)
	if err != nil {
		return
	}

	return MemberPermissions(g, c, message.Author.ID, message.Member.Roles), nil
}

// UserColor returns the color of a user in a channel.
// While colors are defined at a Guild level, determining for a channel is more useful in message handlers.
// 0 is returned in cases of error, which is the color of @everyone.
// userID    : The ID of the user to calculate the color for.
// channelID   : The ID of the channel to calculate the color for.
func (s *State) UserColor(userID, channelID string) int {
	if s == nil {
		return 0
	}

	c, err := s.ChannelState().Channel(channelID)
	if err != nil {
		return 0
	}

	g, err := s.GuildState().Guild(c.GuildID)
	if err != nil {
		return 0
	}

	member, err := s.MemberState().Member(g.ID, userID)
	if err != nil {
		return 0
	}

	return firstRoleColorColor(g, member.Roles)
}

// MessageColor returns the color of the author's name as displayed
// in the client associated with this message.
func (s *State) MessageColor(message *channel.Message) int {
	if s == nil {
		return 0
	}

	if message.Member == nil || message.Member.Roles == nil {
		return 0
	}

	c, err := s.ChannelState().Channel(message.ChannelID)
	if err != nil {
		return 0
	}

	g, err := s.GuildState().Guild(c.GuildID)
	if err != nil {
		return 0
	}

	return firstRoleColorColor(g, message.Member.Roles)
}

func firstRoleColorColor(g *guild.Guild, memberRoles []string) int {
	roles := guild.Roles(g.Roles)
	sort.Sort(roles)

	for _, role := range roles {
		for _, roleID := range memberRoles {
			if role.ID == roleID {
				if role.Color != 0 {
					return role.Color
				}
			}
		}
	}

	for _, role := range roles {
		if role.ID == g.ID {
			return role.Color
		}
	}

	return 0
}

// ThreadListSync syncs guild threads with provided ones.
func (s *State) ThreadListSync(tls *ThreadListSync) error {
	g, err := s.GuildState().Guild(tls.GuildID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	// This algorithm filters out archived or
	// threads which are children of channels in channelIDs
	// and then it adds all synced threads to guild threads and cache
	index := 0
outer:
	for _, t := range g.Threads {
		if !t.ThreadMetadata.Archived && tls.ChannelIDs != nil {
			for _, v := range tls.ChannelIDs {
				if t.ParentID == v {
					delete(s.channelMap, t.ID)
					continue outer
				}
			}
			g.Threads[index] = t
			index++
		} else {
			delete(s.channelMap, t.ID)
		}
	}
	g.Threads = g.Threads[:index]
	for _, t := range tls.Threads {
		s.channelMap[t.ID] = t
		g.Threads = append(g.Threads, t)
	}

	for _, m := range tls.Members {
		if c, ok := s.channelMap[m.ID]; ok {
			c.Member = m
		}
	}

	return nil
}

// ThreadMembersUpdate updates thread members list
func (s *State) ThreadMembersUpdate(tmu *ThreadMembersUpdate) error {
	thread, err := s.Channel(tmu.ID)
	if err != nil {
		return err
	}
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	for idx, member := range thread.Members {
		for _, removedMember := range tmu.RemovedMembers {
			if member.ID == removedMember {
				thread.Members = append(thread.Members[:idx], thread.Members[idx+1:]...)
				break
			}
		}
	}

	for _, addedMember := range tmu.AddedMembers {
		thread.Members = append(thread.Members, addedMember.ThreadMember)
		if addedMember.Member != nil {
			err = s.memberAdd(addedMember.Member)
			if err != nil {
				return err
			}
		}
		if addedMember.Presence != nil {
			err = s.presenceAdd(tmu.GuildID, addedMember.Presence)
			if err != nil {
				return err
			}
		}
	}
	thread.MemberCount = tmu.MemberCount

	return nil
}

// ThreadMemberUpdate sets or updates member data for the current user.
func (s *State) ThreadMemberUpdate(mu *ThreadMemberUpdate) error {
	thread, err := s.Channel(mu.ID)
	if err != nil {
		return err
	}

	thread.Member = mu.ThreadMember
	return nil
}

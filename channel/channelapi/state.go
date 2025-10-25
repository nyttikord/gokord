package channelapi

import (
	"errors"
	"slices"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/state"
)

type State struct {
	state.State
	storage         state.Storage
	privateChannels []*channel.Channel
}

var (
	ErrGuildNotCached   = errors.New("channel's guild not cached")
	ErrChannelNotCached = errors.New("message or thread's channel not cached")
	ErrMemberNotCached  = errors.New("member not cached")
	// ErrMessageIncompletePermissions is returned when the message requested for permissions does not contain enough data to
	// generate the permissions.
	ErrMessageIncompletePermissions = errors.New("message incomplete, unable to determine permissions")
)

func NewState(state state.State, storage state.Storage) *State {
	return &State{
		State:           state,
		storage:         storage,
		privateChannels: make([]*channel.Channel, 0),
	}
}

// AppendGuildChannel is for internal use only.
// Use ChannelAdd instead.
func (s *State) AppendGuildChannel(c *channel.Channel) error {
	return s.storage.Write(state.KeyChannel(c), c)
}

// ChannelAdd adds a channel.Channel to the current State, or updates it if it already exists.
// Channels may exist either as PrivateChannels or inside a guild.Guild.
func (s *State) ChannelAdd(chann *channel.Channel) error {
	g, err := s.GuildState().Guild(chann.GuildID)

	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	if c, err := s.Channel(chann.ID); err == nil {
		if chann.Messages == nil {
			chann.Messages = c.Messages
		}
		if chann.PermissionOverwrites == nil {
			chann.PermissionOverwrites = c.PermissionOverwrites
		}
		if chann.ThreadMetadata == nil {
			chann.ThreadMetadata = c.ThreadMetadata
		}
	}

	fn := func(sl []*channel.Channel) {
		id := slices.IndexFunc(sl, func(c *channel.Channel) bool { return c.ID == chann.ID })
		if id == -1 {
			sl = append(sl, chann)
		} else {
			sl[id] = chann
		}
	}

	if chann.Type == types.ChannelDM || chann.Type == types.ChannelGroupDM {
		fn(s.privateChannels)
	} else {
		if chann.IsThread() {
			fn(g.Threads)
		} else {
			fn(g.Channels)
		}
		err = s.GuildState().GuildAdd(g)
		if err != nil {
			return err
		}
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	return s.AppendGuildChannel(chann)
}

// ChannelRemove removes a channel.Channel from current State.
func (s *State) ChannelRemove(chann *channel.Channel) error {
	_, err := s.Channel(chann.ID)
	if err != nil {
		return err
	}

	if chann.Type == types.ChannelDM || chann.Type == types.ChannelGroupDM {
		s.GetMutex().Lock()
		defer s.GetMutex().Unlock()

		s.privateChannels = slices.DeleteFunc(s.privateChannels, func(c *channel.Channel) bool { return c.ID == chann.ID })
		return s.storage.Delete(state.KeyChannel(chann))
	}

	g, err := s.GuildState().Guild(chann.GuildID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	if chann.IsThread() {
		g.Threads = slices.DeleteFunc(g.Threads, func(c *channel.Channel) bool { return c.ID == chann.ID })
	} else {
		g.Channels = slices.DeleteFunc(g.Channels, func(c *channel.Channel) bool { return c.ID == chann.ID })
	}

	err = s.GuildState().GuildAdd(g)
	if err != nil {
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	return s.storage.Delete(state.KeyChannel(chann))
}

// Channel returns the channel.Channel.
func (s *State) Channel(channelID string) (*channel.Channel, error) {
	s.GetMutex().RLock()
	defer s.GetMutex().RUnlock()

	cRaw, err := s.storage.Get(state.KeyChannelRaw(channelID))
	if err != nil {
		return nil, err
	}
	c := cRaw.(channel.Channel)
	return &c, nil
}

// PrivateChannels returns all private channels.
func (s *State) PrivateChannels() []*channel.Channel {
	return s.privateChannels
}

// MessageAdd adds a channel.Message to the current State, or updates it if it exists.
// If the channel cannot be found, the message is discarded.
// Messages are kept in state up to state.State GetMaxMessageCount per channel.
func (s *State) MessageAdd(message *channel.Message) error {
	c, err := s.Channel(message.ChannelID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrChannelNotCached)
		}
		return err
	}

	if m, err := s.Message(message.ChannelID, message.ID); err == nil {
		if message.Content == "" {
			message.Content = m.Content
		}
		if message.EditedTimestamp == nil {
			message.EditedTimestamp = m.EditedTimestamp
		}
		if message.Mentions == nil {
			message.Mentions = m.Mentions
		}
		if message.Embeds == nil {
			message.Embeds = m.Embeds
		}
		if message.Attachments == nil {
			message.Attachments = m.Attachments
		}
		if message.Timestamp.IsZero() {
			message.Timestamp = m.Timestamp
		}
		if message.Author == nil {
			message.Author = m.Author
		}
		if message.Components == nil {
			message.Components = m.Components
		}
		id := slices.IndexFunc(c.Messages, func(m *channel.Message) bool { return m.ID == message.ID })
		c.Messages[id] = message
	} else {
		c.Messages = append(c.Messages, message)
	}

	if len(c.Messages) > s.Params().MaxMessageCount {
		c.Messages = c.Messages[len(c.Messages)-s.Params().MaxMessageCount:]
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	return s.storage.Write(state.KeyChannel(c), c)
}

// MessageRemove removes a channel.Message from the current State.
func (s *State) MessageRemove(message *channel.Message) error {
	return s.MessageRemoveByID(message.ChannelID, message.ID)
}

// MessageRemoveByID removes a channel.Message by channelID and messageID from the current State.
func (s *State) MessageRemoveByID(channelID, messageID string) error {
	c, err := s.Channel(channelID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrChannelNotCached)
		}
		return err
	}

	c.Messages = slices.DeleteFunc(c.Messages, func(m *channel.Message) bool { return m.ID == messageID })

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	return s.storage.Write(state.KeyChannel(c), c)
}

// Message gets a message by channel and message ID.
func (s *State) Message(channelID, messageID string) (*channel.Message, error) {
	c, err := s.Channel(channelID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return nil, errors.Join(err, ErrChannelNotCached)
		}
		return nil, err
	}

	for _, m := range c.Messages {
		if m.ID == messageID {
			return m, nil
		}
	}

	return nil, state.ErrStateNotFound
}

// ThreadListSync syncs guild threads with provided ones.
// TODO: use gokord.ThreadListSync when event will be remade
func (s *State) ThreadListSync(guildID string, channelIDs []string, threads []*channel.Channel, members []*channel.ThreadMember) error {
	g, err := s.GuildState().Guild(guildID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	ths := make([]*channel.Channel, 0, len(g.Threads))
	messages := make(map[string][]*channel.Message, len(g.Threads))
	// converting channelIDs to map to have better perf
	var channels map[string]struct{} // stored value is never used, we use it like a set
	if len(channelIDs) > 0 {
		channels = make(map[string]struct{}, len(channelIDs))
		for _, id := range channelIDs {
			channels[id] = struct{}{}
		}
	}
	// removing from map archived/deleted thread and saving untouched threads
	for i, c := range g.Channels {
		if c.IsThread() {
			// if thread is in sync list
			ok := true
			if channels != nil {
				_, ok = channels[c.ID]
			}
			if ok {
				// cleaning the map from old thread
				// if the thread continue to exist, it will be added later
				// we just save cached messages before
				messages[c.ID] = c.Messages
				g.Channels = slices.Delete(g.Channels, i, i+1)
			} else {
				// saved because we don't want to touch it
				ths = append(ths, c)
			}
		}
	}
	// updating guild threads and channel map with touched thread
	for _, c := range threads {
		// we add cached messages if we have deleted the thread previously
		c.Messages = messages[c.ID]
		for _, m := range members {
			if m.ID == c.ID {
				c.Member = m
			}
		}
		ths = append(ths, c)
		err = s.ChannelAdd(c)
		if err != nil {
			return err
		}
	}
	g.Threads = ths
	return s.GuildState().GuildAdd(g)
}

// ThreadMembersUpdate updates thread members list.
// TODO: use gokord.ThreadMembersUpdate when event will be remade
func (s *State) ThreadMembersUpdate(id string, guildID string, count int, addedMembers []channel.AddedThreadMember, removedMembers []string) error {
	thread, err := s.Channel(id)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrChannelNotCached)
		}
		return err
	}
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	for _, removedMember := range removedMembers {
		thread.Members = slices.DeleteFunc(thread.Members, func(m *channel.ThreadMember) bool { return m.ID == removedMember })
	}

	for _, addedMember := range addedMembers {
		thread.Members = append(thread.Members, addedMember.ThreadMember)
		if addedMember.Member != nil {
			s.GetMutex().Unlock() // unlock to add the member
			err = s.MemberState().MemberAdd(addedMember.Member)
			s.GetMutex().Lock()
			if err != nil {
				return err
			}
		}
		if addedMember.Presence != nil {
			s.GetMutex().Unlock() // unlock to add the presence
			err = s.MemberState().PresenceAdd(guildID, addedMember.Presence)
			s.GetMutex().Lock()
			if err != nil {
				return err
			}
		}
	}
	thread.MemberCount = count

	return s.ChannelAdd(thread)
}

// ThreadMemberUpdate sets or updates member data for the current user.
func (s *State) ThreadMemberUpdate(tm *channel.ThreadMember) error {
	thread, err := s.Channel(tm.ID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrChannelNotCached)
		}
		return err
	}

	thread.Member = tm
	return nil
}

// UserChannelPermissions returns the permission of a user in a channel.
func (s *State) UserChannelPermissions(userID, channelID string) (int64, error) {
	c, err := s.Channel(channelID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return 0, errors.Join(err, ErrChannelNotCached)
		}
		return 0, err
	}

	g, err := s.GuildState().Guild(c.GuildID)
	if err != nil {
		// checking for state.ErrStateNotFound is useless because it is already checked by Channel
		return 0, err
	}

	member, err := s.MemberState().Member(g.ID, userID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return 0, errors.Join(err, ErrMemberNotCached)
		}
		return 0, err
	}

	return guild.MemberPermissions(g, c, userID, member.Roles), nil
}

// MessagePermissions returns the permissions of the author of the channel.Message in the channel.Channel in which it
// was sent.
func (s *State) MessagePermissions(message *channel.Message) (int64, error) {
	if message.Author == nil || message.Member == nil {
		return 0, ErrMessageIncompletePermissions
	}

	c, err := s.Channel(message.ChannelID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return 0, errors.Join(err, ErrChannelNotCached)
		}
		return 0, err
	}

	g, err := s.GuildState().Guild(c.GuildID)
	if err != nil {
		// checking for state.ErrStateNotFound is useless because it is already checked by Channel
		return 0, err
	}

	return guild.MemberPermissions(g, c, message.Author.ID, message.Member.Roles), nil
}

// MessageColor returns the color of the author's name as displayed in the client associated with this channel.Message.
// Returns 0 in cases of error, which is the color of @everyone.
func (s *State) MessageColor(message *channel.Message) int {
	if message.Member == nil || message.Member.Roles == nil {
		return 0
	}

	c, err := s.Channel(message.ChannelID)
	if err != nil {
		return 0
	}

	g, err := s.GuildState().Guild(c.GuildID)
	if err != nil {
		return 0
	}

	return guild.FirstRoleColor(g, message.Member.Roles)
}

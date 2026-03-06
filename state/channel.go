package state

import (
	"errors"
	"slices"
	"sync"

	"github.com/nyttikord/avl"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
)

type ChannelStorage Storage[uint64, channel.Channel]

type Channel struct {
	State
	mu              sync.RWMutex
	storage         ChannelStorage
	privateChannels *avl.KeyAVL[uint64, channel.Channel]
	params          *Params
}

var (
	ErrChannelGuildNotCached = errors.New("channel's guild not cached")
	ErrChannelNotCached      = errors.New("message or thread's channel not cached")
	ErrMemberNotCached       = errors.New("member not cached")
	// ErrMessageIncompletePermissions is returned when the message requested for permissions does not contain enough
	// data to generate the permissions.
	ErrMessageIncompletePermissions = errors.New("message incomplete, unable to determine permissions")
)

func NewChannel(state State, storage ChannelStorage, params *Params) *Channel {
	return &Channel{
		State:   state,
		storage: storage,
		privateChannels: avl.NewKeySimpleImmutable[uint64, channel.Channel](func(c channel.Channel) channel.Channel {
			cp, err := deepCopy(c)
			if err != nil {
				panic(err)
			}
			return cp
		}),
		params: params,
	}
}

// KeyChannel returns the unique key linked with the given [channel.Channel].
func KeyChannel(c *channel.Channel) uint64 {
	return KeyChannelReverse(c.ID)
}

// KeyChannelReverse returns the key linked with the requested [channel.Channel].
func KeyChannelReverse(channelID uint64) uint64 {
	return channelID
}

// AddChannel adds a [channel.Channel] to the current [Channel] state, or updates it if it already exists.
// Channels may exist either as private channels or inside a [guild.Guild].
func (s *Channel) AddChannel(chann *channel.Channel) error {
	g, err := s.GuildState().GetGuild(chann.GuildID)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.Join(err, ErrChannelGuildNotCached)
		}
		return err
	}

	if c, err := s.GetChannel(chann.ID); err == nil {
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
		s.privateChannels.Insert(KeyChannel(chann), *chann)
	} else {
		if chann.IsThread() {
			fn(g.Threads)
		} else {
			fn(g.Channels)
		}
		err = s.GuildState().AddGuild(g)
		if err != nil {
			return err
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.storage.Write(KeyChannel(chann), *chann)
}

// RemoveChannel removes a [channel.Channel] from current [Channel] state.
func (s *Channel) RemoveChannel(chann *channel.Channel) error {
	_, err := s.GetChannel(chann.ID)
	if err != nil {
		return err
	}

	if chann.Type == types.ChannelDM || chann.Type == types.ChannelGroupDM {
		s.mu.Lock()
		defer s.mu.Unlock()

		key := KeyChannel(chann)
		s.privateChannels.Delete(key)
		return s.storage.Delete(key)
	}

	g, err := s.GuildState().GetGuild(chann.GuildID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.Join(err, ErrChannelGuildNotCached)
		}
		return err
	}

	if chann.IsThread() {
		g.Threads = slices.DeleteFunc(g.Threads, func(c *channel.Channel) bool { return c.ID == chann.ID })
	} else {
		g.Channels = slices.DeleteFunc(g.Channels, func(c *channel.Channel) bool { return c.ID == chann.ID })
	}

	err = s.GuildState().AddGuild(g)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.storage.Delete(KeyChannel(chann))
}

// GetChannel returns the [channel.Channel].
func (s *Channel) GetChannel(channelID uint64) (*channel.Channel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, err := s.storage.Get(KeyChannelReverse(channelID))
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// ListPrivateChannels returns all private [channel.Channel]s.
func (s *Channel) ListPrivateChannels() []channel.Channel {
	return s.privateChannels.Sort()
}

// AddMessage adds a [channel.Message] to the current [Channel] state, or updates it if it exists.
// If the [channel.Channel] cannot be found, the message is discarded.
// Messages are kept in state up to [Params.MaxMessageCount] per [channel.Channel].
func (s *Channel) AddMessage(message *channel.Message) error {
	c, err := s.GetChannel(message.ChannelID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.Join(err, ErrChannelNotCached)
		}
		return err
	}

	if m, err := s.GetMessage(message.ChannelID, message.ID); err == nil {
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

	if len(c.Messages) > s.params.MaxMessageCount {
		c.Messages = c.Messages[len(c.Messages)-s.params.MaxMessageCount:]
	}

	return s.AddChannel(c)
}

// RemoveMessage removes a [channel.Message] from the current [Channel] state.
func (s *Channel) RemoveMessage(message *channel.Message) error {
	return s.RemoveMessageByID(message.ChannelID, message.ID)
}

// RemoveMessageByID removes a [channel.Message] by channelID and messageID from the current [Channel] state.
func (s *Channel) RemoveMessageByID(channelID, messageID uint64) error {
	c, err := s.GetChannel(channelID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.Join(err, ErrChannelNotCached)
		}
		return err
	}

	c.Messages = slices.DeleteFunc(c.Messages, func(m *channel.Message) bool { return m.ID == messageID })

	return s.AddChannel(c)
}

// GetMessage gets a [channel.Message] by [channel.Channel.ID] and [channel.Message.ID].
func (s *Channel) GetMessage(channelID, messageID uint64) (*channel.Message, error) {
	c, err := s.GetChannel(channelID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, errors.Join(err, ErrChannelNotCached)
		}
		return nil, err
	}

	for _, m := range c.Messages {
		if m.ID == messageID {
			return m, nil
		}
	}

	return nil, ErrNotFound
}

// OnThreadListSync syncs [guild.Guild] threads with provided ones.
func (s *Channel) OnThreadListSync(guildID uint64, channelIDs []uint64, threads []*channel.Channel, members []*channel.ThreadMember) error {
	g, err := s.GuildState().GetGuild(guildID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.Join(err, ErrChannelGuildNotCached)
		}
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	ths := make([]*channel.Channel, 0, len(g.Threads))
	messages := make(map[uint64][]*channel.Message, len(g.Threads))
	// converting channelIDs to map to have better perf
	var channels map[uint64]struct{} // stored value is never used, we use it like a set
	if len(channelIDs) > 0 {
		channels = make(map[uint64]struct{}, len(channelIDs))
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
		err = s.AddChannel(c)
		if err != nil {
			return err
		}
	}
	g.Threads = ths
	return s.GuildState().AddGuild(g)
}

// OnThreadMembersUpdate updates [channel.ThreadMember]s list.
func (s *Channel) OnThreadMembersUpdate(id uint64, guildID uint64, count int, addedMembers []channel.AddedThreadMember, removedMembers []string) error {
	thread, err := s.GetChannel(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.Join(err, ErrChannelNotCached)
		}
		return err
	}

	for _, removedMember := range removedMembers {
		thread.Members = slices.DeleteFunc(thread.Members, func(m *channel.ThreadMember) bool { return m.ID == removedMember })
	}

	for _, addedMember := range addedMembers {
		thread.Members = append(thread.Members, addedMember.ThreadMember)
		if addedMember.Member != nil {
			err = s.MemberState().AddMember(addedMember.Member)
			if err != nil {
				return err
			}
		}
		if addedMember.Presence != nil {
			err = s.MemberState().AddPresence(guildID, addedMember.Presence)
			if err != nil {
				return err
			}
		}
	}
	thread.MemberCount = count

	return s.AddChannel(thread)
}

// ThreadMemberUpdate sets or updates member data for the current user.
func (s *Channel) ThreadMemberUpdate(tm *channel.ThreadMember) error {
	thread, err := s.GetChannel(tm.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.Join(err, ErrChannelNotCached)
		}
		return err
	}

	thread.Member = tm
	return nil
}

// GetUserChannelPermissions returns the permission of a [user.User] in a [channel.Channel].
func (s *Channel) GetUserChannelPermissions(userID, channelID uint64) (int64, error) {
	c, err := s.GetChannel(channelID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return 0, errors.Join(err, ErrChannelNotCached)
		}
		return 0, err
	}

	g, err := s.GuildState().GetGuild(c.GuildID)
	if err != nil {
		// checking for state.ErrStateNotFound is useless because it is already checked by Channel
		return 0, err
	}

	member, err := s.MemberState().GetMember(g.ID, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return 0, errors.Join(err, ErrMemberNotCached)
		}
		return 0, err
	}

	return guild.MemberPermissions(g, c, userID, member.Roles), nil
}

// GetMessagePermissions returns the permissions of the author of the [channel.Message] in the [channel.Channel] in
// which it was sent.
func (s *Channel) GetMessagePermissions(message *channel.Message) (int64, error) {
	if message.Author == nil || message.Member == nil {
		return 0, ErrMessageIncompletePermissions
	}

	c, err := s.GetChannel(message.ChannelID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return 0, errors.Join(err, ErrChannelNotCached)
		}
		return 0, err
	}

	g, err := s.GuildState().GetGuild(c.GuildID)
	if err != nil {
		// checking for state.ErrStateNotFound is useless because it is already checked by Channel
		return 0, err
	}

	return guild.MemberPermissions(g, c, message.Author.ID, message.Member.Roles), nil
}

// GetMessageColor returns the color of the author's name as displayed in the client associated with this
// [channel.Message].
// Returns 0 in cases of error, which is the color of @everyone.
func (s *Channel) GetMessageColor(message *channel.Message) int {
	if message.Member == nil || message.Member.Roles == nil {
		return 0
	}

	c, err := s.GetChannel(message.ChannelID)
	if err != nil {
		return 0
	}

	g, err := s.GuildState().GetGuild(c.GuildID)
	if err != nil {
		return 0
	}

	return guild.FirstRoleColor(g, message.Member.Roles)
}

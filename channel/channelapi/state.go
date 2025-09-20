package channelapi

import (
	"errors"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/state"
)

type State struct {
	state.State
	channelMap      map[string]*channel.Channel
	privateChannels []*channel.Channel
}

var ErrGuildNotCached = errors.New("channel's guild not cached")

func NewState(state state.State) *State {
	return &State{
		State:           state,
		channelMap:      make(map[string]*channel.Channel),
		privateChannels: make([]*channel.Channel, 0),
	}
}

// AppendGuildChannel is for internal use only.
func (s *State) AppendGuildChannel(c *channel.Channel) {
	s.channelMap[c.ID] = c
}

// ChannelAdd adds a channel.Channel to the current State, or updates it if it already exists.
// Channels may exist either as PrivateChannels or inside a guild.Guild.
func (s *State) ChannelAdd(channel *channel.Channel) error {
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	// If the channel exists, replace it
	if c, ok := s.channelMap[channel.ID]; ok {
		if channel.Messages == nil {
			channel.Messages = c.Messages
		}
		if channel.PermissionOverwrites == nil {
			channel.PermissionOverwrites = c.PermissionOverwrites
		}
		if channel.ThreadMetadata == nil {
			channel.ThreadMetadata = c.ThreadMetadata
		}

		*c = *channel
		return nil
	}

	if channel.Type == types.ChannelDM || channel.Type == types.ChannelGroupDM {
		s.privateChannels = append(s.privateChannels, channel)
		s.channelMap[channel.ID] = channel
		return nil
	}

	g, err := s.GuildState().Guild(channel.GuildID)
	if err != nil {
		if errors.Is(err, state.ErrStateNotFound) {
			return errors.Join(err, ErrGuildNotCached)
		}
		return err
	}

	if channel.IsThread() {
		g.Threads = append(g.Threads, channel)
	} else {
		g.Channels = append(g.Channels, channel)
	}

	s.channelMap[channel.ID] = channel

	return nil
}

// ChannelRemove removes a channel.Channel from current State.
func (s *State) ChannelRemove(channel *channel.Channel) error {
	_, err := s.Channel(channel.ID)
	if err != nil {
		return err
	}

	if channel.Type == types.ChannelDM || channel.Type == types.ChannelGroupDM {
		s.GetMutex().Lock()
		defer s.GetMutex().Unlock()

		for i, c := range s.privateChannels {
			if c.ID == channel.ID {
				s.privateChannels = append(s.privateChannels[:i], s.privateChannels[i+1:]...)
				break
			}
		}
		delete(s.channelMap, channel.ID)
		return nil
	}

	guild, err := s.GuildState().Guild(channel.GuildID)
	if err != nil {
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	if channel.IsThread() {
		for i, t := range guild.Threads {
			if t.ID == channel.ID {
				guild.Threads = append(guild.Threads[:i], guild.Threads[i+1:]...)
				break
			}
		}
	} else {
		for i, c := range guild.Channels {
			if c.ID == channel.ID {
				guild.Channels = append(guild.Channels[:i], guild.Channels[i+1:]...)
				break
			}
		}
	}

	delete(s.channelMap, channel.ID)

	return nil
}

// Channel returns the channel.Channel.
func (s *State) Channel(channelID string) (*channel.Channel, error) {
	s.GetMutex().RLock()
	defer s.GetMutex().RUnlock()

	if c, ok := s.channelMap[channelID]; ok {
		return c, nil
	}

	return nil, state.ErrStateNotFound
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
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	// If the message exists, merge in the new message contents.
	for _, m := range c.Messages {
		if m.ID == message.ID {
			if message.Content != "" {
				m.Content = message.Content
			}
			if message.EditedTimestamp != nil {
				m.EditedTimestamp = message.EditedTimestamp
			}
			if message.Mentions != nil {
				m.Mentions = message.Mentions
			}
			if message.Embeds != nil {
				m.Embeds = message.Embeds
			}
			if message.Attachments != nil {
				m.Attachments = message.Attachments
			}
			if !message.Timestamp.IsZero() {
				m.Timestamp = message.Timestamp
			}
			if message.Author != nil {
				m.Author = message.Author
			}
			if message.Components != nil {
				m.Components = message.Components
			}

			return nil
		}
	}

	c.Messages = append(c.Messages, message)

	if len(c.Messages) > s.GetMaxMessageCount() {
		c.Messages = c.Messages[len(c.Messages)-s.GetMaxMessageCount():]
	}

	return nil
}

// MessageRemove removes a channel.Message from the current State.
func (s *State) MessageRemove(message *channel.Message) error {
	return s.MessageRemoveByID(message.ChannelID, message.ID)
}

// MessageRemoveByID removes a channel.Message by channelID and messageID from the current State.
func (s *State) MessageRemoveByID(channelID, messageID string) error {
	c, err := s.Channel(channelID)
	if err != nil {
		return err
	}

	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	for i, m := range c.Messages {
		if m.ID == messageID {
			c.Messages = append(c.Messages[:i], c.Messages[i+1:]...)
			return nil
		}
	}

	return state.ErrStateNotFound
}

// Message gets a message by channel and message ID.
func (s *State) Message(channelID, messageID string) (*channel.Message, error) {
	c, err := s.Channel(channelID)
	if err != nil {
		return nil, err
	}

	s.GetMutex().RLock()
	defer s.GetMutex().RUnlock()

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

	// This algorithm filters out archived or
	// threads which are children of channels in channelIDs
	// and then it adds all synced threads to guild threads and cache

	// THIS CODE IS UGLY BUT I DON'T HAVE THE STRENGTH TO FIX IT YET
	index := 0
outer:
	for _, t := range g.Threads {
		if !t.ThreadMetadata.Archived && channelIDs != nil {
			for _, v := range channelIDs {
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
	for _, t := range threads {
		s.channelMap[t.ID] = t
		g.Threads = append(g.Threads, t)
	}

	for _, m := range members {
		if c, ok := s.channelMap[m.ID]; ok {
			c.Member = m
		}
	}

	return nil
}

// ThreadMembersUpdate updates thread members list.
// TODO: use gokord.ThreadMembersUpdate when event will be remade
func (s *State) ThreadMembersUpdate(id string, guildID string, count int, addedMembers []channel.AddedThreadMember, removedMembers []string) error {
	thread, err := s.Channel(id)
	if err != nil {
		return err
	}
	s.GetMutex().Lock()
	defer s.GetMutex().Unlock()

	for idx, member := range thread.Members {
		for _, removedMember := range removedMembers {
			if member.ID == removedMember {
				thread.Members = append(thread.Members[:idx], thread.Members[idx+1:]...)
				break
			}
		}
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

	return nil
}

// ThreadMemberUpdate sets or updates member data for the current user.
func (s *State) ThreadMemberUpdate(tm *channel.ThreadMember) error {
	thread, err := s.Channel(tm.ID)
	if err != nil {
		return err
	}

	thread.Member = tm
	return nil
}

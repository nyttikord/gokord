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

// ChannelAdd adds a channel to the current world state, or updates it if it already exists.
// Channels may exist either as PrivateChannels or inside a guild.
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

	g, err := s.Guild(channel.GuildID)
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

// ChannelRemove removes a channel from current world state.
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

	guild, err := s.Guild(channel.GuildID)
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

// Channel gets a channel by ID, it will look in all guilds and private channels.
func (s *State) Channel(channelID string) (*channel.Channel, error) {
	s.GetMutex().RLock()
	defer s.GetMutex().RUnlock()

	if c, ok := s.channelMap[channelID]; ok {
		return c, nil
	}

	return nil, state.ErrStateNotFound
}

func (s *State) PrivateChannels() []*channel.Channel {
	return s.privateChannels
}

// MessageAdd adds a message to the current world state, or updates it if it exists.
// If the channel cannot be found, the message is discarded.
// Messages are kept in state up to s.MaxMessageCount per channel.
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

// MessageRemove removes a message from the world state.
func (s *State) MessageRemove(message *channel.Message) error {
	return s.messageRemoveByID(message.ChannelID, message.ID)
}

// messageRemoveByID removes a message by channelID and messageID from the world state.
func (s *State) messageRemoveByID(channelID, messageID string) error {
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

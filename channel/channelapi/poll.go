package channelapi

import (
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/user"
)

// PollAnswerVoters returns user.Get who voted for a particular channel.PollAnswer in a channel.Poll on the given
// channel.Message.
func (s Requester) PollAnswerVoters(channelID, messageID string, answerID int, options ...discord.RequestOption) ([]*user.User, error) {
	body, err := s.Request(http.MethodGet, discord.EndpointPollAnswerVoters(channelID, messageID, answerID), nil, options...)
	if err != nil {
		return nil, err
	}

	var r struct {
		Users []*user.User `json:"users"`
	}

	err = s.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return r.Users, nil
}

// PollExpire expires channel.Poll on the given channel.Message.
func (s Requester) PollExpire(channelID, messageID string, options ...discord.RequestOption) (*channel.Message, error) {
	body, err := s.Request(http.MethodPost, discord.EndpointPollExpire(channelID, messageID), nil, options...)
	if err != nil {
		return nil, err
	}

	var m channel.Message
	return &m, s.Unmarshal(body, &m)
}

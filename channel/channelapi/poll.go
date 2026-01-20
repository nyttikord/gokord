package channelapi

import (
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/user"
)

// PollAnswerVoters returns user.User who voted for a particular channel.PollAnswer in a channel.Poll on the given
// channel.Message.
func (s Requester) PollAnswerVoters(channelID, messageID string, answerID int) request.Request[[]*user.User] {
	body, err := s.Request(ctx, http.MethodGet, discord.EndpointPollAnswerVoters(channelID, messageID, answerID), nil, options...)
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
func (s Requester) PollExpire(channelID, messageID string) request.Request[*channel.Message] {
	return request.NewData[*channel.Message](
		s, http.MethodPost, discord.EndpointPollExpire(channelID, messageID),
	).WithBucketID(discord.EndpointPollExpire(channelID, ""))
}

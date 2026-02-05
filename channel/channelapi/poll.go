package channelapi

import (
	"context"
	"net/http"

	. "github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/user"
)

// PollAnswerVoters returns user.User who voted for a particular channel.PollAnswer in a channel.Poll on the given
// Message.
func (r Requester) PollAnswerVoters(channelID, messageID string, answerID int) Request[[]*user.User] {
	return NewCustom[[]*user.User](r, http.MethodGet, discord.EndpointPollAnswerVoters(channelID, messageID, answerID)).
		WithBucketID(discord.EndpointPollAnswerVoters(channelID, messageID, 0)).
		WithPost(func(ctx context.Context, b []byte) ([]*user.User, error) {
			var data struct {
				Users []*user.User `json:"users"`
			}
			return data.Users, r.Unmarshal(b, &data)

		})
}

// PollExpire expires channel.Poll on the given channel.Message.
func (r Requester) PollExpire(channelID, messageID string) Request[*Message] {
	return NewData[*Message](
		r, http.MethodPost, discord.EndpointPollExpire(channelID, messageID),
	).WithBucketID(discord.EndpointPollExpire(channelID, ""))
}

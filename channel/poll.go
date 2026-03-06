package channel

import (
	"context"
	"net/http"
	"time"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/user"
)

// PollMedia contains common data used by question and answers.
type PollMedia struct {
	Text  string           `json:"text,omitempty"`
	Emoji *emoji.Component `json:"emoji,omitempty"` // TODO: rename the type
}

// PollAnswer represents a single answer in a [Poll].
type PollAnswer struct {
	// NOTE: should not be set on creation.
	AnswerID int        `json:"answer_id,omitempty"`
	Media    *PollMedia `json:"poll_media"`
}

// PollAnswerCount stores counted [Poll] votes for a single answer.
type PollAnswerCount struct {
	ID      int  `json:"id"`
	Count   int  `json:"count"`
	MeVoted bool `json:"me_voted"`
}

// PollResults contains voting results on a [Poll].
type PollResults struct {
	Finalized    bool               `json:"is_finalized"`
	AnswerCounts []*PollAnswerCount `json:"answer_counts"`
}

// Poll contains all poll related data.
type Poll struct {
	Question         PollMedia        `json:"question"`
	Answers          []PollAnswer     `json:"answers"`
	AllowMultiselect bool             `json:"allow_multiselect"`
	LayoutType       types.PollLayout `json:"layout_type,omitempty"`

	// NOTE: should be set only on creation, when fetching use Expiry.
	Duration int `json:"duration,omitempty"`

	// NOTE: available only when fetching.

	Results *PollResults `json:"results,omitempty"`
	// NOTE: as Discord documentation notes, this field might be null even when fetching.
	Expiry *time.Time `json:"expiry,omitempty"`
}

// GetPollVoters returns [user.User] who voted for a particular [PollAnswer] in a [Poll] on the given [Message].
func GetPollVoters(channelID, messageID uint64, answerID int) Request[[]*user.User] {
	return NewCustom[[]*user.User](http.MethodGet, discord.EndpointPollAnswerVoters(channelID, messageID, answerID)).
		WithBucketID(discord.EndpointPollAnswerVoters(channelID, messageID, 0)).
		WithPost(func(ctx context.Context, b []byte) ([]*user.User, error) {
			var data struct {
				Users []*user.User `json:"users"`
			}
			return data.Users, Unmarshal(ctx, b, &data)
		})
}

// ExpirePoll on the given [Message].
func ExpirePoll(channelID, messageID uint64) Request[*Message] {
	return NewData[*Message](http.MethodPost, discord.EndpointPollExpire(channelID, messageID)).
		WithBucketID(discord.EndpointPollExpire(channelID, 0))
}

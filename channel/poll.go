package channel

import (
	"time"

	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
)

// PollMedia contains common data used by question and answers.
type PollMedia struct {
	Text  string           `json:"text,omitempty"`
	Emoji *emoji.Component `json:"emoji,omitempty"` // TODO: rename the type
}

// PollAnswer represents a single answer in a poll.
type PollAnswer struct {
	// NOTE: should not be set on creation.
	AnswerID int        `json:"answer_id,omitempty"`
	Media    *PollMedia `json:"poll_media"`
}

// PollAnswerCount stores counted poll votes for a single answer.
type PollAnswerCount struct {
	ID      int  `json:"id"`
	Count   int  `json:"count"`
	MeVoted bool `json:"me_voted"`
}

// PollResults contains voting results on a poll.
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

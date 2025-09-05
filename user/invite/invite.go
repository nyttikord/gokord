package invite

import (
	"time"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
)

// An Invite stores all data related to a specific Discord Guild or Channel invite.
type Invite struct {
	Guild             *guild.Guild             `json:"guild"`
	Channel           *channel.Channel         `json:"channel"`
	Inviter           *user.User               `json:"inviter"`
	Code              string                   `json:"code"`
	CreatedAt         time.Time                `json:"created_at"`
	MaxAge            int                      `json:"max_age"`
	Uses              int                      `json:"uses"`
	MaxUses           int                      `json:"max_uses"`
	Revoked           bool                     `json:"revoked"`
	Temporary         bool                     `json:"temporary"`
	Unique            bool                     `json:"unique"`
	TargetUser        *user.User               `json:"target_user"`
	TargetType        types.InviteTarget       `json:"target_type"`
	TargetApplication *application.Application `json:"target_application"`

	// will only be filled when using InviteWithCounts
	ApproximatePresenceCount int `json:"approximate_presence_count"`
	ApproximateMemberCount   int `json:"approximate_member_count"`

	ExpiresAt *time.Time `json:"expires_at"`
}

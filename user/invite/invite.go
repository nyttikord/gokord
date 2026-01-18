// Package invite contains the Invite... and that's all...
package invite

import (
	"time"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
)

// Invite stores all data related to a specific Discord guild.Guild or channel.Channel invite.
type Invite struct {
	Type      types.Invite     `json:"type"`
	Guild     *guild.Guild     `json:"guild"`
	Channel   *channel.Channel `json:"channel"`
	Inviter   *user.User       `json:"inviter"`
	Code      string           `json:"code"`
	CreatedAt time.Time        `json:"created_at"`
	MaxAge    int              `json:"max_age"`
	Uses      int              `json:"uses"`
	MaxUses   int              `json:"max_uses"`
	Revoked   bool             `json:"revoked"`
	Temporary bool             `json:"temporary"`
	Unique    bool             `json:"unique"`
	// See TargetUser and TargetApplication.
	TargetType types.InviteTarget `json:"target_type"`
	// TargetUser is the user.User streaming displayed for this Invite.
	// Requires TargetType to be types.InviteTargetStream.
	// Set TargetUserID when creating an Invite to use this feature.
	TargetUser *user.User `json:"target_user,omitempty"`
	// TargetApplication is the embedded application.Application to open for this Invite.
	// Requires TargetType to be types.InviteTargetEmbeddedApplication.
	// Set TargetApplicationID when creating an Invite to use this feature.
	TargetApplication *application.Application `json:"target_application,omitempty"`
	// TargetUsersFile is a CSV with a single column of user.User IDs for all the user.User able to accept this Invite.
	// Does not work with a channel.Channel Invite.
	TargetUsersFile []byte `json:"target_users_file,omitempty"`
	// Roles are the guild.Role given when the user.User joins the guild.Guild.
	// Does not work with a channel.Channel Invite.
	Roles []string `json:"role_ids,omitempty"`

	// will only be filled when using InviteWithCounts
	ApproximatePresenceCount int `json:"approximate_presence_count"`
	ApproximateMemberCount   int `json:"approximate_member_count"`

	ExpiresAt *time.Time `json:"expires_at"`
}

// The Discord's documentation does not yet provide complete information.
// Check https://discord.com/developers/docs/resources/invite#get-target-users-job-status for more information.
type TargetUsersJobStatus struct {
	Status         types.TargetUsersJobStatus `json:"status"`
	TotalUsers     uint                       `json:"total_users"`
	ProcessedUsers uint                       `json:"processed_users"`
	CreatedAt      time.Time                  `json:"created_at"`
	CompletedAt    time.Time                  `json:"completed_at"`
	ErrorMessage   string                     `json:"error_message"`
}

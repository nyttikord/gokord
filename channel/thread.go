package channel

import (
	"time"

	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
)

// ThreadStart stores all parameters you can use with MessageThreadStartComplex or ThreadStartComplex
type ThreadStart struct {
	Name                string        `json:"name"`
	AutoArchiveDuration int           `json:"auto_archive_duration,omitempty"`
	Type                types.Channel `json:"type,omitempty"`
	Invitable           bool          `json:"invitable"`
	RateLimitPerUser    int           `json:"rate_limit_per_user,omitempty"`

	// NOTE: forum threads only
	AppliedTags []string `json:"applied_tags,omitempty"`
}

// ThreadMetadata contains a number of thread-specific Channel fields that are not needed by other channel types.
type ThreadMetadata struct {
	// Whether the thread is archived.
	Archived bool `json:"archived"`
	// Duration in minutes to automatically archive the thread after recent activity, can be set to: 60, 1440, 4320, 10080.
	AutoArchiveDuration int `json:"auto_archive_duration"`
	// Timestamp when the thread's archive status was last changed, used for calculating recent activity.
	ArchiveTimestamp time.Time `json:"archive_timestamp"`
	// Whether the thread is locked; when a thread is locked, only users with permission discord.PermissionManageThreads
	// can unarchive it.
	Locked bool `json:"locked"`
	// Whether non-moderators can add other non-moderators to a thread; only available on private threads.
	Invitable bool `json:"invitable"`
}

// ThreadMember is used to indicate whether a user.User has joined a thread or not.
//
// NOTE: ID and UserID are empty (omitted) on the user.Member sent within each thread in the GUILD_CREATE event.
type ThreadMember struct {
	// The id of the thread.
	ID string `json:"id,omitempty"`
	// The id of the user.User.
	UserID string `json:"user_id,omitempty"`
	// The time the current user last joined the thread.
	JoinTimestamp time.Time `json:"join_timestamp"`
	// Any user-thread settings, currently only used for notifications.
	Flags int `json:"flags"`
	// Additional information about the user.User.
	//
	// NOTE: only present if the withMember parameter is set to true when calling Session.ThreadMembers or
	// Session.ThreadMember.
	Member *user.Member `json:"member,omitempty"`
}

// ThreadsList represents a list of threads alongisde with ThreadMember for the current user.User.
type ThreadsList struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

// AddedThreadMember holds information about the user.User who was added to the thread
type AddedThreadMember struct {
	*ThreadMember
	Member   *user.Member     `json:"member"`
	Presence *status.Presence `json:"presence"`
}

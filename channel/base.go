// Package channel contains every data structures linked with channels like... Channel or Message.
// It also has helping functions not using gokord.Session.
//
// Use channelapi.Requester to interact with this.
// You can get this with gokord.Session.
package channel

import (
	"fmt"
	"time"

	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// Flags represent flags of a Channel (including threads).
type Flags int

const (
	// FlagPinned indicates whether the thread is pinned in the forum channel.
	// NOTE: forum threads only.
	FlagPinned Flags = 1 << 1
	// FlagRequireTag indicates whether a tag is required to be specified when creating a thread.
	// NOTE: forum channels only.
	FlagRequireTag Flags = 1 << 4
)

// ForumLayout represents layout of a forum channel (Channel with type types.ChannelGuildForum).
type ForumLayout int

const (
	// ForumLayoutNotSet represents no default layout.
	ForumLayoutNotSet ForumLayout = 0
	// ForumLayoutListView displays forum posts as a list.
	ForumLayoutListView ForumLayout = 1
	// ForumLayoutGalleryView displays forum posts as a collection of tiles.
	ForumLayoutGalleryView ForumLayout = 2
)

// PermissionOverwrite holds permission overwrite data for a Channel
type PermissionOverwrite struct {
	ID    string                    `json:"id"`
	Type  types.PermissionOverwrite `json:"type"`
	Deny  int64                     `json:"deny,string"`
	Allow int64                     `json:"allow,string"`
}

// Channel holds all data related to an individual Discord Channel.
type Channel struct {
	// The ID of the Channel.
	ID string `json:"id"`

	// The ID of the guild.Guild to which the Channel belongs, if it is in a guild.
	// Else, this ID is empty (e.g. DM channels).
	GuildID string `json:"guild_id"`

	// The name of the Channel.
	Name string `json:"name"`

	// The topic of the Channel.
	Topic string `json:"topic"`

	// The type of the Channel.
	Type types.Channel `json:"type"`

	// The ID of the last message sent in the Channel.
	// This is not guaranteed to be an ID of a valid Message.
	LastMessageID string `json:"last_message_id"`

	// The timestamp of the last pinned Message in the Channel.
	// nil if the Channel has no pinned messages.
	LastPinTimestamp *time.Time `json:"last_pin_timestamp"`

	// An approximate count of messages in a thread, stops counting at 50
	MessageCount int `json:"message_count"`
	// An approximate count of users in a thread, stops counting at 50
	MemberCount int `json:"member_count"`

	// Whether the Channel is marked as NSFW.
	NSFW bool `json:"nsfw"`

	// Icon of the group DM Channel.
	Icon string `json:"icon"`

	// The position of the Channel, used for sorting in client.
	Position int `json:"position"`

	// The bitrate of the Channel, if it is a voice Channel (types.ChannelGuildVoice).
	Bitrate int `json:"bitrate"`

	// The recipients of the Channel.
	// This is only populated in DM channels.
	Recipients []*user.User `json:"recipients"`

	// The messages in the Channel.
	// This is only present in state-cached channels, and State.MaxMessageCount must be non-zero.
	Messages []*Message `json:"-"`

	// A list of permission overwrites present for the Channel.
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites"`

	// The user limit of the voice Channel (types.ChannelGuildVoice).
	UserLimit int `json:"user_limit"`

	// The ID of the parent Channel, if the Channel is under a category.
	// For threads - id of the Channel thread was created in.
	ParentID string `json:"parent_id"`

	// Amount of seconds a user.User has to wait before sending another Message or creating another thread (0-21600).
	//
	// Bots, as well as users with the permission discord.PermissionManageMessages or discord.PermissionManageChannels,
	// are unaffected
	RateLimitPerUser int `json:"rate_limit_per_user"`

	// ID of the creator of the group DM or thread
	OwnerID string `json:"owner_id"`

	// ApplicationID of the DM creator Zeroed if guild Channel or not a bot user
	ApplicationID string `json:"application_id"`

	// Thread-specific fields not needed by other channels
	ThreadMetadata *ThreadMetadata `json:"thread_metadata,omitempty"`
	// Thread member object for the current user, if they have joined the thread, only included on certain API endpoints
	Member *ThreadMember `json:"thread_member"`

	// All thread members.
	/// State channels only.
	Members []*ThreadMember `json:"-"`

	// Channel flags.
	Flags Flags `json:"flags"`

	// The set of tags that can be used in a forum channel.
	AvailableTags []ForumTag `json:"available_tags"`

	// The IDs of the set of tags that have been applied to a thread in a forum channel.
	AppliedTags []string `json:"applied_tags"`

	// Emoji to use as the default reaction to a forum post.
	DefaultReactionEmoji ForumDefaultReaction `json:"default_reaction_emoji"`

	// The initial RateLimitPerUser to set on newly created threads in a channel.
	// This field is copied to the thread at creation time and does not live update.
	DefaultThreadRateLimitPerUser int `json:"default_thread_rate_limit_per_user"`

	// The default sort order type used to order posts in forum channels.
	// Defaults to null, which indicates a preferred sort order hasn't been set by a channel admin.
	DefaultSortOrder *types.ForumSortOrder `json:"default_sort_order"`

	// The default forum layout view used to display posts in forum channels.
	// Defaults to ForumLayoutNotSet, which indicates a layout view has not been set by a channel admin.
	DefaultForumLayout ForumLayout `json:"default_forum_layout"`
}

// Mention returns a string which mentions the Channel
func (c *Channel) Mention() string {
	return fmt.Sprintf("<#%s>", c.ID)
}

// IsThread is a helper function to determine if Channel is a thread or not
func (c *Channel) IsThread() bool {
	return c.Type == types.ChannelGuildPublicThread ||
		c.Type == types.ChannelGuildPrivateThread ||
		c.Type == types.ChannelGuildNewsThread
}

// Edit holds Channel field data for a channel edit.
type Edit struct {
	Name                          string                 `json:"name,omitempty"`
	Topic                         string                 `json:"topic,omitempty"`
	NSFW                          *bool                  `json:"nsfw,omitempty"`
	Position                      *int                   `json:"position,omitempty"`
	Bitrate                       int                    `json:"bitrate,omitempty"`
	UserLimit                     int                    `json:"user_limit,omitempty"`
	PermissionOverwrites          []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      string                 `json:"parent_id,omitempty"`
	RateLimitPerUser              *int                   `json:"rate_limit_per_user,omitempty"`
	Flags                         *Flags                 `json:"flags,omitempty"`
	DefaultThreadRateLimitPerUser *int                   `json:"default_thread_rate_limit_per_user,omitempty"`

	// NOTE: threads only
	Archived *bool `json:"archived,omitempty"`
	// NOTE: threads only
	AutoArchiveDuration int `json:"auto_archive_duration,omitempty"`
	// NOTE: threads only
	Locked *bool `json:"locked,omitempty"`
	// NOTE: threads only
	Invitable *bool `json:"invitable,omitempty"`

	// NOTE: forum channels only
	AvailableTags *[]ForumTag `json:"available_tags,omitempty"`
	// NOTE: forum channels only
	DefaultReactionEmoji *ForumDefaultReaction `json:"default_reaction_emoji,omitempty"`
	// NOTE: forum channels only
	DefaultSortOrder *types.ForumSortOrder `json:"default_sort_order,omitempty"` // TODO: null
	// NOTE: forum channels only
	DefaultForumLayout *ForumLayout `json:"default_forum_layout,omitempty"`

	// NOTE: forum threads only
	AppliedTags *[]string `json:"applied_tags,omitempty"`
}

// A Follow holds data returned after following a news Channel
type Follow struct {
	ChannelID string `json:"channel_id"`
	WebhookID string `json:"webhook_id"`
}

// ForumDefaultReaction specifies emoji to use as the default reaction to a forum post.
//
// NOTE: Exactly one of EmojiID and EmojiName must be set.
type ForumDefaultReaction struct {
	// The id of a guild's custom emoji.
	EmojiID string `json:"emoji_id,omitempty"`
	// The Unicode character of the emoji.
	EmojiName string `json:"emoji_name,omitempty"`
}

// ForumTag represents a tag that can be applied to a thread in a forum channel.
type ForumTag struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Moderated bool   `json:"moderated"`
	EmojiID   string `json:"emoji_id,omitempty"`
	EmojiName string `json:"emoji_name,omitempty"`
}

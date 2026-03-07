// Package channel contains every data structures linked with channels like [Channel] or [Message].
package channel

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// Flags represent flags of a [Channel] (including threads).
type Flags int

const (
	// FlagPinned indicates whether the thread is pinned in the forum [Channel].
	// NOTE: forum threads only.
	FlagPinned Flags = 1 << 1
	// FlagRequireTag indicates whether a tag is required to be specified when creating a thread.
	// NOTE: forum channels only.
	FlagRequireTag Flags = 1 << 4
)

// ForumLayout represents layout of a forum channel ([Channel] with [types.ChannelGuildForum]).
type ForumLayout int

const (
	// ForumLayoutNotSet represents no default layout.
	ForumLayoutNotSet ForumLayout = 0
	// ForumLayoutListView displays forum posts as a list.
	ForumLayoutListView ForumLayout = 1
	// ForumLayoutGalleryView displays forum posts as a collection of tiles.
	ForumLayoutGalleryView ForumLayout = 2
)

// PermissionOverwrite holds permission overwrite data for a [Channel].
type PermissionOverwrite struct {
	ID    uint64                    `json:"id,string"`
	Type  types.PermissionOverwrite `json:"type"`
	Deny  int64                     `json:"deny,string"`
	Allow int64                     `json:"allow,string"`
}

// Channel holds all data related to an individual Discord
type Channel struct {
	ID uint64 `json:"id,string"`
	// The ID of the guild.Guild to which the [Channel] belongs, if it is in a [guild.Guild].
	// Else, this ID is empty (e.g. DM channels).
	GuildID uint64        `json:"guild_id,string"`
	Name    string        `json:"name"`
	Topic   string        `json:"topic"`
	Type    types.Channel `json:"type"`
	// The ID of the last message sent in the [Channel].
	// This is not guaranteed to be an ID of a valid [Message].
	LastMessageID uint64 `json:"last_message_id,string"`
	// The timestamp of the last pinned [Message] in the [Channel].
	// It is nil if the [Channel] has no pinned messages.
	LastPinTimestamp *time.Time `json:"last_pin_timestamp"`
	// An approximate count of [Message]s in a thread, stops counting at 50
	MessageCount int `json:"message_count"`
	// An approximate count of [user.User]s in a thread, stops counting at 50
	MemberCount int    `json:"member_count"`
	NSFW        bool   `json:"nsfw"`
	Icon        string `json:"icon"` // Icon of the group DM [Channel].
	// Position of the [Channel], used for sorting in client.
	Position int `json:"position"`
	// Bitrate of the [Channel], if it is a voice [Channel] ([types.ChannelGuildVoice]).
	Bitrate int `json:"bitrate"`
	// Recipients of the [Channel].
	// This is only populated in DM [Channel]s.
	Recipients []*user.User `json:"recipients"`
	// The messages in the [Channel].
	// This is only present in state-cached [Channel]s, and MaxMessageCount must be non-zero.
	Messages []*Message `json:"-"`
	// A list of permission overwrites present for the [Channel].
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites"`
	// UserLimit is the maximum number of [user.User]s in the voice [Channel] ([types.ChannelGuildVoice]).
	UserLimit int `json:"user_limit"`
	// The ID of the parent [Channel], if the [Channel] is under a category.
	// For threads, it is the id of the [Channel] thread was created in.
	ParentID uint64 `json:"parent_id,string"`
	// Amount of seconds a [user.User] has to wait before sending another [Message] or creating another Thread
	// (0-21600).
	//
	// Bots, as well as users with the permission [discord.PermissionManageMessages] or
	// [discord.PermissionManageChannels], are unaffected
	RateLimitPerUser int `json:"rate_limit_per_user"`
	// ID of the creator of the group DM or thread
	OwnerID uint64 `json:"owner_id,string"`
	// ApplicationID of the DM creator Zeroed if [guild.Guild] [Channel] or not a bot [user.User].
	ApplicationID uint64 `json:"application_id,string"`
	// Thread-specific fields not needed by other [Channel]s.
	ThreadMetadata *ThreadMetadata `json:"thread_metadata,omitempty"`
	// ThreadMember for the current [user.User], if they have joined the thread, only included on certain API endpoints
	Member *ThreadMember `json:"thread_member"`
	// All [ThreadMember]s.
	// state.State [Channel]s only.
	Members []*ThreadMember `json:"-"`
	Flags   Flags           `json:"flags"`
	// The set of tags that can be used in a forum [Channel].
	AvailableTags []ForumTag `json:"available_tags"`
	// The IDs of the set of tags that have been applied to a thread in a forum [Channel].
	AppliedTags []string `json:"applied_tags"`
	// DefaultReactionEmoji to a forum post.
	DefaultReactionEmoji ForumDefaultReaction `json:"default_reaction_emoji"`
	// The initial RateLimitPerUser to set on newly created threads in a [Channel].
	// This field is copied to the thread at creation time and does not live update.
	DefaultThreadRateLimitPerUser int `json:"default_thread_rate_limit_per_user"`
	// DefaultSortOrder type used to order posts in forum [Channel]s.
	// Defaults to null, which indicates a preferred sort order hasn't been set by a [Channel] admin.
	DefaultSortOrder *types.ForumSortOrder `json:"default_sort_order"`
	// DefaultForumLayout view used to display posts in forum [Channel]s.
	// Defaults to [ForumLayoutNotSet], which indicates a layout view has not been set by a [Channel] admin.
	DefaultForumLayout ForumLayout `json:"default_forum_layout"`
}

// Mention returns a string which mentions the [Channel]
func (c *Channel) Mention() string {
	return fmt.Sprintf("<#%d>", c.ID)
}

// IsThread returns true if the [Channel] is a thread.
func (c *Channel) IsThread() bool {
	return c.Type == types.ChannelGuildPublicThread ||
		c.Type == types.ChannelGuildPrivateThread ||
		c.Type == types.ChannelGuildNewsThread
}

// EditData holds [Channel] field data for a [Edit].
type EditData struct {
	Name                          string                 `json:"name,omitempty"`
	Topic                         string                 `json:"topic,omitempty"`
	NSFW                          *bool                  `json:"nsfw,omitempty"`
	Position                      *int                   `json:"position,omitempty"`
	Bitrate                       int                    `json:"bitrate,omitempty"`
	UserLimit                     int                    `json:"user_limit,omitempty"`
	PermissionOverwrites          []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      uint64                 `json:"parent_id,omitempty,string"`
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

// A Follow holds data returned after following a news [Channel].
type Follow struct {
	ChannelID uint64 `json:"channel_id,string"`
	WebhookID uint64 `json:"webhook_id,string"`
}

// ForumDefaultReaction specifies [emoji.Emoji] to use as the default reaction to a forum post.
//
// NOTE: Exactly one of EmojiID and EmojiName must be set.
type ForumDefaultReaction struct {
	// The id of a guild's custom [emoji.Emoji].
	EmojiID uint64 `json:"emoji_id,omitempty,string"`
	// The Unicode character of the [emoji.Emoji].
	EmojiName string `json:"emoji_name,omitempty"`
}

// ForumTag represents a tag that can be applied to a thread in a forum [Channel].
type ForumTag struct {
	ID        uint64 `json:"id,omitempty,string"`
	Name      string `json:"name"`
	Moderated bool   `json:"moderated"`
	EmojiID   uint64 `json:"emoji_id,omitempty,string"`
	EmojiName string `json:"emoji_name,omitempty"`
}

// Get returns the [Channel] with the given ID.
func Get(channelID uint64) Request[*Channel] {
	return NewData[*Channel](http.MethodGet, discord.EndpointChannel(channelID))
}

// Edit the given [Channel].
func Edit(channelID uint64, data *EditData) Request[*Channel] {
	return NewData[*Channel](http.MethodPatch, discord.EndpointChannel(channelID)).
		WithData(data)
}

// Delete the given [Channel].
func Delete(channelID uint64) Request[*Channel] {
	return NewData[*Channel](http.MethodDelete, discord.EndpointChannel(channelID))
}

// List returns the list of [Channel] in the [guild.Guild].
func List(guildID uint64) Request[[]*Channel] {
	return NewData[[]*Channel](http.MethodGet, discord.EndpointGuildChannels(guildID))
}

// CreateData is provided to [CreateComplex].
type CreateData struct {
	Name                 string                 `json:"name"`
	Type                 types.Channel          `json:"type"`
	Topic                string                 `json:"topic,omitempty"`
	Bitrate              int                    `json:"bitrate,omitempty"`
	UserLimit            int                    `json:"user_limit,omitempty"`
	RateLimitPerUser     int                    `json:"rate_limit_per_user,omitempty"`
	Position             int                    `json:"position,omitempty"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             uint64                 `json:"parent_id,omitempty,string"`
	NSFW                 bool                   `json:"nsfw,omitempty"`
}

// CreateComplex creates a new [Channel] in the given [guild.Guild].
func CreateComplex(guildID uint64, data CreateData) Request[*Channel] {
	return NewData[*Channel](http.MethodPost, discord.EndpointGuildChannels(guildID)).
		WithData(data)
}

// Create a new [Channel] in the given [guild.Guild].
func Create(guildID uint64, name string, ctype types.Channel) Request[*Channel] {
	return CreateComplex(guildID, CreateData{
		Name: name,
		Type: ctype,
	})
}

// Reorder updates the order of [Channel] in a [guild.Guild].
func Reorder(guildID uint64, channels []*Channel) Empty {
	data := make([]struct {
		ID       uint64 `json:"id,string"`
		Position int    `json:"position"`
	}, len(channels))

	for i, c := range channels {
		data[i].ID = c.ID
		data[i].Position = c.Position
	}

	req := NewSimple(http.MethodPatch, discord.EndpointGuildChannels(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// Typing broadcasts to all members that authenticated [user.User] is typing in the given [Channel].
func Typing(channelID uint64) Empty {
	req := NewSimple(http.MethodPost, discord.EndpointChannelTyping(channelID))
	return WrapAsEmpty(req)
}

// SetPermission creates a [PermissionOverwrite] for the given [Channel].
func SetPermission(channelID, targetID uint64, targetType types.PermissionOverwrite, allow, deny int64) Empty {
	data := struct {
		ID    uint64                    `json:"id,string"`
		Type  types.PermissionOverwrite `json:"type"`
		Allow int64                     `json:"allow,string"`
		Deny  int64                     `json:"deny,string"`
	}{targetID, targetType, allow, deny}

	req := NewSimple(http.MethodPut, discord.EndpointChannelPermission(channelID, targetID)).
		WithData(data).
		WithBucketID(discord.EndpointChannelPermission(channelID, 0))
	return WrapAsEmpty(req)
}

// DeletePermission [PermissionOverwrite] for the given [Channel].
func DeletePermission(channelID, targetID uint64) Empty {
	req := NewSimple(http.MethodPut, discord.EndpointChannelPermission(channelID, targetID)).
		WithBucketID(discord.EndpointChannelPermission(channelID, 0))
	return WrapAsEmpty(req)
}

// FollowNews [Channel] in the given [Channel].
//
// channelID is the [Channel] to follow.
// targetID is where the news [Channel] should post to.
func FollowNews(channelID, targetID uint64) Request[*Follow] {
	data := struct {
		WebhookChannelID uint64 `json:"webhook_channel_id,string"`
	}{targetID}

	return NewData[*Follow](http.MethodPost, discord.EndpointChannelFollow(channelID)).
		WithData(data)
}

// Pin a [Message] within the given [Channel].
func Pin(channelID, messageID uint64) Empty {
	req := NewSimple(http.MethodPut, discord.EndpointChannelMessagePin(channelID, messageID)).
		WithBucketID(discord.EndpointChannelMessagePin(channelID, 0))
	return WrapAsEmpty(req)
}

// Unpin a [Message] within the given [Channel].
func Unpin(channelID, messageID uint64) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointChannelMessagePin(channelID, messageID)).
		WithBucketID(discord.EndpointChannelMessagePin(channelID, 0))
	return WrapAsEmpty(req)
}

// ListPinned returns [MessagesPinned] within the given [Channel].
//
// limit is the max number of messages to return (max 50).
// If provided all messages returned will be before the given time.
func ListPinned(channelID uint64, before *time.Time, limit int) Request[*MessagesPinned] {
	uri := discord.EndpointChannelMessagesPins(channelID)

	v := url.Values{}
	if before != nil {
		v.Set("before", before.Format(time.RFC3339))
	}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	return NewData[*MessagesPinned](http.MethodGet, uri)
}

// CreatePrivate [Channel] ([types.ChannelDM]) with another [user.User].
func CreatePrivate(userID uint64) Request[*Channel] {
	data := struct {
		RecipientID uint64 `json:"recipient_id,string"`
	}{userID}

	return NewData[*Channel](http.MethodPost, discord.EndpointUserChannels(0)).
		WithBucketID(discord.EndpointUserChannels(0)).WithData(data)
}

package channel

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
)

// ThreadStart stores all parameters you can use with [StartThreadMessageComplex] or [StartThreadComplex].
type ThreadStart struct {
	Name                string        `json:"name"`
	AutoArchiveDuration int           `json:"auto_archive_duration,omitempty"`
	Type                types.Channel `json:"type,omitempty"`
	Invitable           bool          `json:"invitable"`
	RateLimitPerUser    int           `json:"rate_limit_per_user,omitempty"`

	// NOTE: forum threads only
	AppliedTags []string `json:"applied_tags,omitempty"`
}

// ThreadMetadata contains a number of thread-specific [Channel] fields that are not needed by other channel types.
type ThreadMetadata struct {
	Archived bool `json:"archived"`
	// Duration in minutes to automatically archive the thread after recent activity, can be set to: 60, 1440, 4320, 10080.
	AutoArchiveDuration int `json:"auto_archive_duration"`
	// Timestamp when the thread's archive status was last changed, used for calculating recent activity.
	ArchiveTimestamp time.Time `json:"archive_timestamp"`
	// Whether the thread is locked; when a thread is locked, only users with permission
	// [discord.PermissionManageThreads] can unarchive it.
	Locked bool `json:"locked"`
	// Whether non-moderators can add other non-moderators to a thread; only available on private threads.
	Invitable bool `json:"invitable"`
}

// ThreadMember is used to indicate whether a [user.User] has joined a thread or not.
//
// NOTE: ID and UserID are empty (omitted) on the [user.Member] sent within each thread in the GUILD_CREATE event.
type ThreadMember struct {
	ID     uint64 `json:"id,omitempty,string"`
	UserID uint64 `json:"user_id,omitempty,string"`
	// The time the current user last joined the thread.
	JoinTimestamp time.Time `json:"join_timestamp"`
	// Any user-thread settings, currently only used for notifications.
	Flags int `json:"flags"`
	// Additional information about the [user.User].
	//
	// NOTE: only present if the withMember parameter is set to true
	Member *user.Member `json:"member,omitempty"`
}

// ThreadsList represents a list of threads alongside with [ThreadMember] for the current [user.User].
type ThreadsList struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

// AddedThreadMember holds information about the [user.User] who was added to the thread.
type AddedThreadMember struct {
	*ThreadMember
	Member   *user.Member     `json:"member"`
	Presence *status.Presence `json:"presence"`
}

// StartThreadMessageComplex creates a new thread from an existing [Message].
func StartThreadMessageComplex(channelID, messageID uint64, data *ThreadStart) Request[*Channel] {
	return NewData[*Channel](http.MethodPost, discord.EndpointChannelMessageThread(channelID, messageID)).
		WithBucketID(discord.EndpointChannelMessageThread(channelID, 0)).WithData(data)
}

// StartThreadMessage creates a new thread from an existing [Message].
//
// archiveDuration is the auto archive duration in minutes.
func StartThreadMessage(channelID, messageID uint64, name string, archiveDuration int) Request[*Channel] {
	return StartThreadMessageComplex(channelID, messageID, &ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	})
}

// StartThreadComplex creates a new thread.
func StartThreadComplex(channelID uint64, data *ThreadStart) Request[*Channel] {
	return NewData[*Channel](http.MethodPost, discord.EndpointChannelThreads(channelID)).
		WithData(data)
}

// StartThread creates a new thread.
//
// archiveDuration is the auto archive duration in minutes.
func StartThread(channelID uint64, name string, typ types.Channel, archiveDuration int) Request[*Channel] {
	return StartThreadComplex(channelID, &ThreadStart{
		Name:                name,
		Type:                typ,
		AutoArchiveDuration: archiveDuration,
	})
}

// StartForumThreadComplex in a [types.ChannelGuildForum].
func StartForumThreadComplex(channelID uint64, threadData *ThreadStart, messageData *MessageSend) Request[*Channel] {
	endpoint := discord.EndpointChannelThreads(channelID)

	for _, embed := range messageData.Embeds {
		if embed.Type == "" {
			embed.Type = types.EmbedRich
		}
	}

	files := messageData.Files

	data := struct {
		*ThreadStart
		Message *MessageSend `json:"message"`
	}{ThreadStart: threadData, Message: messageData}

	if len(files) == 0 {
		return NewData[*Channel](http.MethodPost, endpoint).WithData(data)
	}
	return NewMultipart[*Channel](http.MethodPost, endpoint, data, files)
}

// StartForumThread in a [types.ChannelGuildForum].
//
// archiveDuration is the auto archive duration in minutes.
func StartForumThread(channelID uint64, name string, archiveDuration int, content string) Request[*Channel] {
	return StartForumThreadComplex(channelID, &ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &MessageSend{Content: content})
}

// StartForumThreadEmbed in a [types.ChannelGuildForum].
//
// archiveDuration is the auto archive duration in minutes.
func StartForumThreadEmbed(channelID uint64, name string, archiveDuration int, embed *MessageEmbed) Request[*Channel] {
	return StartForumThreadEmbeds(channelID, name, archiveDuration, []*MessageEmbed{embed})
}

// StartForumThreadEmbeds in a [types.ChannelGuildForum].
//
// archiveDuration is the auto archive duration in minutes.
func StartForumThreadEmbeds(channelID uint64, name string, archiveDuration int, embeds []*MessageEmbed) Request[*Channel] {
	return StartForumThreadComplex(channelID, &ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &MessageSend{Embeds: embeds})
}

// JoinThread adds current [user.User] to a thread.
func JoinThread(id uint64) Empty {
	req := NewSimple(http.MethodPut, discord.EndpointThreadMember(id, 0)).
		WithBucketID(discord.EndpointThreadMember(id, 0))
	return WrapAsEmpty(req)
}

// LeaveThread removes current [user.User] to a thread.
func LeaveThread(id uint64) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointThreadMember(id, 0)).
		WithBucketID(discord.EndpointThreadMember(id, 0))
	return WrapAsEmpty(req)
}

// ThreadMemberAdd adds a [user.User] to a thread.
func ThreadMemberAdd(threadID, memberID uint64) Empty {
	req := NewSimple(http.MethodPut, discord.EndpointThreadMember(threadID, memberID)).
		WithBucketID(discord.EndpointThreadMember(threadID, 0))
	return WrapAsEmpty(req)
}

// ThreadMemberRemove removes a [user.User] from a thread.
func ThreadMemberRemove(threadID, memberID uint64) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointThreadMember(threadID, memberID)).
		WithBucketID(discord.EndpointThreadMember(threadID, 0))
	return WrapAsEmpty(req)
}

// GetThreadMember for the specified [user.User] of the thread.
//
// If withMember is true, it includes a guild member object.
func GetThreadMember(threadID, memberID uint64, withMember bool) Request[*ThreadMember] {
	uri := discord.EndpointThreadMember(threadID, memberID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	return NewData[*ThreadMember](http.MethodGet, uri).
		WithBucketID(discord.EndpointThreadMember(threadID, 0))
}

// ListThreadMembers returns all [ThreadMember] of specified thread.
//
// limit is the max number of thread members to return (1-100). Defaults to 100.
// If afterID is set, every member ID will be after this.
// If withMember is true, it includes a guild member object.
func ListThreadMembers(threadID uint64, limit int, withMember bool, afterID string) Request[[]*ThreadMember] {
	uri := discord.EndpointThreadMembers(threadID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		queryParams.Set("after", afterID)
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	return NewData[[]*ThreadMember](http.MethodGet, uri)
}

// ListThreadsActive returns all active threads in the given [Channel].
func ListThreadsActive(channelID uint64) Request[*ThreadsList] {
	return NewData[*ThreadsList](http.MethodGet, discord.EndpointChannelActiveThreads(channelID))
}

// ListThreadsArchived returns archived threads in the given [Channel].
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func ListThreadsArchived(channelID uint64, before *time.Time, limit int) Request[*ThreadsList] {
	endpoint := discord.EndpointChannelPublicArchivedThreads(channelID)
	v := url.Values{}
	if before != nil {
		v.Set("before", before.Format(time.RFC3339))
	}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if len(v) > 0 {
		endpoint += "?" + v.Encode()
	}

	return NewData[*ThreadsList](http.MethodGet, endpoint)
}

// ListThreadsPrivateArchived returns archived private threads in the given [Channel].
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func ListThreadsPrivateArchived(channelID uint64, before *time.Time, limit int) Request[*ThreadsList] {
	endpoint := discord.EndpointChannelPrivateArchivedThreads(channelID)
	v := url.Values{}
	if before != nil {
		v.Set("before", before.Format(time.RFC3339))
	}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if len(v) > 0 {
		endpoint += "?" + v.Encode()
	}

	return NewData[*ThreadsList](http.MethodGet, endpoint)
}

// ListThreadsPrivateJoinedArchived returns archived joined private threads in the given [Channel].
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func ListThreadsPrivateJoinedArchived(channelID uint64, before *time.Time, limit int) Request[*ThreadsList] {
	endpoint := discord.EndpointChannelJoinedPrivateArchivedThreads(channelID)
	v := url.Values{}
	if before != nil {
		v.Set("before", before.Format(time.RFC3339))
	}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if len(v) > 0 {
		endpoint += "?" + v.Encode()
	}

	return NewData[*ThreadsList](http.MethodGet, endpoint)
}

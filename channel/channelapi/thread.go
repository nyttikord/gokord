package channelapi

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	. "github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
)

// MessageThreadStartComplex creates a new thread from an existing channel.Message.
func (r Requester) MessageThreadStartComplex(channelID, messageID string, data *ThreadStart) Request[*Channel] {
	return NewData[*Channel](
		r, http.MethodPost, discord.EndpointChannelMessageThread(channelID, messageID),
	).WithBucketID(discord.EndpointChannelMessageThread(channelID, "")).WithData(data)
}

// MessageThreadStart creates a new thread from an existing channel.Message.
//
// archiveDuration is the auto archive duration in minutes.
func (r Requester) MessageThreadStart(channelID, messageID string, name string, archiveDuration int) Request[*Channel] {
	return r.MessageThreadStartComplex(channelID, messageID, &ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	})
}

// ThreadStartComplex creates a new thread.
func (r Requester) ThreadStartComplex(channelID string, data *ThreadStart) Request[*Channel] {
	return NewData[*Channel](
		r, http.MethodPost, discord.EndpointChannelThreads(channelID),
	).WithData(data)
}

// ThreadStart creates a new thread.
//
// archiveDuration is the auto archive duration in minutes.
func (r Requester) ThreadStart(channelID, name string, typ types.Channel, archiveDuration int) Request[*Channel] {
	return r.ThreadStartComplex(channelID, &ThreadStart{
		Name:                name,
		Type:                typ,
		AutoArchiveDuration: archiveDuration,
	})
}

// ForumThreadStartComplex starts a new thread (creates a post) in a types.ChannelGuildForum.
func (r Requester) ForumThreadStartComplex(channelID string, threadData *ThreadStart, messageData *MessageSend) Request[*Channel] {
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
		return NewData[*Channel](
			r, http.MethodPost, endpoint,
		).WithData(data)
	}
	return NewMultipart[*Channel](r, http.MethodPost, endpoint, data, files)
}

// ForumThreadStart starts a new thread (post) in a types.ChannelGuildForum.
//
// archiveDuration is the auto archive duration in minutes.
func (r Requester) ForumThreadStart(channelID, name string, archiveDuration int, content string) Request[*Channel] {
	return r.ForumThreadStartComplex(channelID, &ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &MessageSend{Content: content})
}

// ForumThreadStartEmbed starts a new thread (post) in a types.ChannelGuildForum.
//
// archiveDuration is the auto archive duration in minutes.
func (r Requester) ForumThreadStartEmbed(channelID, name string, archiveDuration int, embed *MessageEmbed) Request[*Channel] {
	return r.ForumThreadStartComplex(channelID, &ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &MessageSend{Embeds: []*MessageEmbed{embed}})
}

// ForumThreadStartEmbeds starts a new thread (post) in a types.ChannelGuildForum.
//
// archiveDuration is the auto archive duration in minutes.
func (r Requester) ForumThreadStartEmbeds(channelID, name string, archiveDuration int, embeds []*MessageEmbed) Request[*Channel] {
	return r.ForumThreadStartComplex(channelID, &ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &MessageSend{Embeds: embeds})
}

// ThreadJoin adds current user.User to a thread.
func (r Requester) ThreadJoin(id string) Empty {
	req := NewSimple(
		r, http.MethodPut, discord.EndpointThreadMember(id, "@me"),
	).WithBucketID(discord.EndpointThreadMember(id, ""))
	return WrapAsEmpty(req)
}

// ThreadLeave removes current user.User to a thread.
func (r Requester) ThreadLeave(id string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointThreadMember(id, "@me"),
	).WithBucketID(discord.EndpointThreadMember(id, ""))
	return WrapAsEmpty(req)
}

// ThreadMemberAdd adds a user.Member to a thread.
func (r Requester) ThreadMemberAdd(threadID, memberID string) Empty {
	req := NewSimple(
		r, http.MethodPut, discord.EndpointThreadMember(threadID, memberID),
	).WithBucketID(discord.EndpointThreadMember(threadID, ""))
	return WrapAsEmpty(req)
}

// ThreadMemberRemove removes a user.Member from a thread.
func (r Requester) ThreadMemberRemove(threadID, memberID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointThreadMember(threadID, memberID),
	).WithBucketID(discord.EndpointThreadMember(threadID, ""))
	return WrapAsEmpty(req)
}

// ThreadMember returns ThreadMember for the specified user.Member of the thread.
//
// If withMember is true, it includes a guild member object.
func (r Requester) ThreadMember(threadID, memberID string, withMember bool) Request[*ThreadMember] {
	uri := discord.EndpointThreadMember(threadID, memberID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	return NewData[*ThreadMember](
		r, http.MethodGet, uri,
	).WithBucketID(discord.EndpointThreadMember(threadID, ""))
}

// ThreadMembers returns all user.Member of specified thread.
//
// limit is the max number of thread members to return (1-100). Defaults to 100.
// If afterID is set, every member ID will be after this.
// If withMember is true, it includes a guild member object.
func (r Requester) ThreadMembers(threadID string, limit int, withMember bool, afterID string) Request[[]*ThreadMember] {
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

	return NewData[[]*ThreadMember](r, http.MethodGet, uri)
}

// ThreadsActive returns all active threads in the given channel.Channel.
func (r Requester) ThreadsActive(channelID string) Request[*ThreadsList] {
	return NewData[*ThreadsList](
		r, http.MethodGet, discord.EndpointChannelActiveThreads(channelID),
	)
}

// ThreadsArchived returns archived threads in the given channel.Channel.
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func (r Requester) ThreadsArchived(channelID string, before *time.Time, limit int) Request[*ThreadsList] {
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

	return NewData[*ThreadsList](
		r, http.MethodGet, endpoint,
	)
}

// ThreadsPrivateArchived returns archived private threads in the given channel.Channel.
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func (r Requester) ThreadsPrivateArchived(channelID string, before *time.Time, limit int) Request[*ThreadsList] {
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

	return NewData[*ThreadsList](
		r, http.MethodGet, endpoint,
	)
}

// ThreadsPrivateJoinedArchived returns archived joined private threads in the given channel.Channel.
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func (r Requester) ThreadsPrivateJoinedArchived(channelID string, before *time.Time, limit int) Request[*ThreadsList] {
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

	return NewData[*ThreadsList](
		r, http.MethodGet, endpoint,
	)
}

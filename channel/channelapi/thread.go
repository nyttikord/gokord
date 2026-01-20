package channelapi

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
)

// MessageThreadStartComplex creates a new thread from an existing channel.Message.
func (s Requester) MessageThreadStartComplex(channelID, messageID string, data *channel.ThreadStart) Request[*channel.Channel] {
	return NewData[*channel.Channel](
		s, http.MethodPost, discord.EndpointChannelMessageThread(channelID, messageID),
	).WithBucketID(discord.EndpointChannelMessageThread(channelID, "")).WithData(data)
}

// MessageThreadStart creates a new thread from an existing channel.Message.
//
// archiveDuration is the auto archive duration in minutes.
func (s Requester) MessageThreadStart(channelID, messageID string, name string, archiveDuration int) Request[*channel.Channel] {
	return s.MessageThreadStartComplex(channelID, messageID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	})
}

// ThreadStartComplex creates a new thread.
func (s Requester) ThreadStartComplex(channelID string, data *channel.ThreadStart) Request[*channel.Channel] {
	return NewData[*channel.Channel](
		s, http.MethodPost, discord.EndpointChannelThreads(channelID),
	).WithData(data)
}

// ThreadStart creates a new thread.
//
// archiveDuration is the auto archive duration in minutes.
func (s Requester) ThreadStart(channelID, name string, typ types.Channel, archiveDuration int) Request[*channel.Channel] {
	return s.ThreadStartComplex(channelID, &channel.ThreadStart{
		Name:                name,
		Type:                typ,
		AutoArchiveDuration: archiveDuration,
	})
}

// ForumThreadStartComplex starts a new thread (creates a post) in a types.ChannelGuildForum.
func (s Requester) ForumThreadStartComplex(channelID string, threadData *channel.ThreadStart, messageData *channel.MessageSend) Request[*channel.Channel] {
	endpoint := discord.EndpointChannelThreads(channelID)

	for _, embed := range messageData.Embeds {
		if embed.Type == "" {
			embed.Type = types.EmbedRich
		}
	}

	files := messageData.Files

	data := struct {
		*channel.ThreadStart
		Message *channel.MessageSend `json:"message"`
	}{ThreadStart: threadData, Message: messageData}

	if len(files) == 0 {
		return NewData[*channel.Channel](
			s, http.MethodPost, endpoint,
		).WithData(data)
	}
	contentType, body, encodeErr := channel.MultipartBodyWithJSON(data, files)
	if encodeErr != nil {
		return NewError[*channel.Channel](encodeErr)
	}

	response, err = s.RequestRaw(ctx, http.MethodPost, endpoint, contentType, body, endpoint, 0, options...)
}

// ForumThreadStart starts a new thread (post) in a types.ChannelGuildForum.
//
// archiveDuration is the auto archive duration in minutes.
func (s Requester) ForumThreadStart(channelID, name string, archiveDuration int, content string) Request[*channel.Channel] {
	return s.ForumThreadStartComplex(channelID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &channel.MessageSend{Content: content})
}

// ForumThreadStartEmbed starts a new thread (post) in a types.ChannelGuildForum.
//
// archiveDuration is the auto archive duration in minutes.
func (s Requester) ForumThreadStartEmbed(channelID, name string, archiveDuration int, embed *channel.MessageEmbed) Request[*channel.Channel] {
	return s.ForumThreadStartComplex(channelID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &channel.MessageSend{Embeds: []*channel.MessageEmbed{embed}})
}

// ForumThreadStartEmbeds starts a new thread (post) in a types.ChannelGuildForum.
//
// archiveDuration is the auto archive duration in minutes.
func (s Requester) ForumThreadStartEmbeds(channelID, name string, archiveDuration int, embeds []*channel.MessageEmbed) Request[*channel.Channel] {
	return s.ForumThreadStartComplex(channelID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &channel.MessageSend{Embeds: embeds})
}

// ThreadJoin adds current user.User to a thread.
func (s Requester) ThreadJoin(id string) Empty {
	req := NewSimple(
		s, http.MethodPut, discord.EndpointThreadMember(id, "@me"),
	).WithBucketID(discord.EndpointThreadMember(id, ""))
	return WrapAsEmpty(req)
}

// ThreadLeave removes current user.User to a thread.
func (s Requester) ThreadLeave(id string) Empty {
	req := NewSimple(
		s, http.MethodDelete, discord.EndpointThreadMember(id, "@me"),
	).WithBucketID(discord.EndpointThreadMember(id, ""))
	return WrapAsEmpty(req)
}

// ThreadMemberAdd adds a user.Member to a thread.
func (s Requester) ThreadMemberAdd(threadID, memberID string) Empty {
	req := NewSimple(
		s, http.MethodPut, discord.EndpointThreadMember(threadID, memberID),
	).WithBucketID(discord.EndpointThreadMember(threadID, ""))
	return WrapAsEmpty(req)
}

// ThreadMemberRemove removes a user.Member from a thread.
func (s Requester) ThreadMemberRemove(threadID, memberID string) Empty {
	req := NewSimple(
		s, http.MethodDelete, discord.EndpointThreadMember(threadID, memberID),
	).WithBucketID(discord.EndpointThreadMember(threadID, ""))
	return WrapAsEmpty(req)
}

// ThreadMember returns channel.ThreadMember for the specified user.Member of the thread.
//
// If withMember is true, it includes a guild member object.
func (s Requester) ThreadMember(threadID, memberID string, withMember bool) Request[*channel.ThreadMember] {
	uri := discord.EndpointThreadMember(threadID, memberID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	return NewData[*channel.ThreadMember](
		s, http.MethodGet, uri,
	).WithBucketID(discord.EndpointThreadMember(threadID, ""))
}

// ThreadMembers returns all user.Member of specified thread.
//
// limit is the max number of thread members to return (1-100). Defaults to 100.
// If afterID is set, every member ID will be after this.
// If withMember is true, it includes a guild member object.
func (s Requester) ThreadMembers(threadID string, limit int, withMember bool, afterID string) Request[[]*channel.ThreadMember] {
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

	return NewData[[]*channel.ThreadMember](s, http.MethodGet, uri)
}

// ThreadsActive returns all active threads in the given channel.Channel.
func (s Requester) ThreadsActive(channelID string) Request[*channel.ThreadsList] {
	return NewData[*channel.ThreadsList](
		s, http.MethodGet, discord.EndpointChannelActiveThreads(channelID),
	)
}

// ThreadsArchived returns archived threads in the given channel.Channel.
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func (s Requester) ThreadsArchived(channelID string, before *time.Time, limit int) Request[*channel.ThreadsList] {
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

	return NewData[*channel.ThreadsList](
		s, http.MethodGet, endpoint,
	)
}

// ThreadsPrivateArchived returns archived private threads in the given channel.Channel.
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func (s Requester) ThreadsPrivateArchived(channelID string, before *time.Time, limit int) Request[*channel.ThreadsList] {
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

	return NewData[*channel.ThreadsList](
		s, http.MethodGet, endpoint,
	)
}

// ThreadsPrivateJoinedArchived returns archived joined private threads in the given channel.Channel.
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func (s Requester) ThreadsPrivateJoinedArchived(channelID string, before *time.Time, limit int) Request[*channel.ThreadsList] {
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

	return NewData[*channel.ThreadsList](
		s, http.MethodGet, endpoint,
	)
}

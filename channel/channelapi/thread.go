package channelapi

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
)

// MessageThreadStartComplex creates a new thread from an existing channel.Message.
func (s Requester) MessageThreadStartComplex(ctx context.Context, channelID, messageID string, data *channel.ThreadStart, options ...discord.RequestOption) (*channel.Channel, error) {
	body, err := s.Request(ctx, http.MethodPost, discord.EndpointChannelMessageThread(channelID, messageID), data, options...)
	if err != nil {
		return nil, err
	}

	var c channel.Channel
	return &c, s.Unmarshal(body, &c)
}

// MessageThreadStart creates a new thread from an existing channel.Message.
//
// archiveDuration is the auto archive duration in minutes.
func (s Requester) MessageThreadStart(ctx context.Context, channelID, messageID string, name string, archiveDuration int, options ...discord.RequestOption) (*channel.Channel, error) {
	return s.MessageThreadStartComplex(ctx, channelID, messageID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, options...)
}

// ThreadStartComplex creates a new thread.
func (s Requester) ThreadStartComplex(ctx context.Context, channelID string, data *channel.ThreadStart, options ...discord.RequestOption) (*channel.Channel, error) {
	body, err := s.Request(ctx, http.MethodPost, discord.EndpointChannelThreads(channelID), data, options...)
	if err != nil {
		return nil, err
	}

	var c channel.Channel
	return &c, s.Unmarshal(body, &c)
}

// ThreadStart creates a new thread.
//
// archiveDuration is the auto archive duration in minutes.
func (s Requester) ThreadStart(ctx context.Context, channelID, name string, typ types.Channel, archiveDuration int, options ...discord.RequestOption) (*channel.Channel, error) {
	return s.ThreadStartComplex(ctx, channelID, &channel.ThreadStart{
		Name:                name,
		Type:                typ,
		AutoArchiveDuration: archiveDuration,
	}, options...)
}

// ForumThreadStartComplex starts a new thread (creates a post) in a types.ChannelGuildForum.
func (s Requester) ForumThreadStartComplex(ctx context.Context, channelID string, threadData *channel.ThreadStart, messageData *channel.MessageSend, options ...discord.RequestOption) (*channel.Channel, error) {
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

	var err error
	var response []byte
	if len(files) > 0 {
		contentType, body, encodeErr := channel.MultipartBodyWithJSON(data, files)
		if encodeErr != nil {
			return nil, encodeErr
		}

		response, err = s.RequestRaw(ctx, http.MethodPost, endpoint, contentType, body, endpoint, 0, options...)
	} else {
		response, err = s.Request(ctx, http.MethodPost, endpoint, data, options...)
	}
	if err != nil {
		return nil, err
	}

	var c channel.Channel
	return &c, s.Unmarshal(response, &c)
}

// ForumThreadStart starts a new thread (post) in a types.ChannelGuildForum.
//
// archiveDuration is the auto archive duration in minutes.
func (s Requester) ForumThreadStart(ctx context.Context, channelID, name string, archiveDuration int, content string, options ...discord.RequestOption) (th *channel.Channel, err error) {
	return s.ForumThreadStartComplex(ctx, channelID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &channel.MessageSend{Content: content}, options...)
}

// ForumThreadStartEmbed starts a new thread (post) in a types.ChannelGuildForum.
//
// archiveDuration is the auto archive duration in minutes.
func (s Requester) ForumThreadStartEmbed(ctx context.Context, channelID, name string, archiveDuration int, embed *channel.MessageEmbed, options ...discord.RequestOption) (*channel.Channel, error) {
	return s.ForumThreadStartComplex(ctx, channelID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &channel.MessageSend{Embeds: []*channel.MessageEmbed{embed}}, options...)
}

// ForumThreadStartEmbeds starts a new thread (post) in a types.ChannelGuildForum.
//
// archiveDuration is the auto archive duration in minutes.
func (s Requester) ForumThreadStartEmbeds(ctx context.Context, channelID, name string, archiveDuration int, embeds []*channel.MessageEmbed, options ...discord.RequestOption) (*channel.Channel, error) {
	return s.ForumThreadStartComplex(ctx, channelID, &channel.ThreadStart{
		Name:                name,
		AutoArchiveDuration: archiveDuration,
	}, &channel.MessageSend{Embeds: embeds}, options...)
}

// ThreadJoin adds current user.User to a thread.
func (s Requester) ThreadJoin(ctx context.Context, id string, options ...discord.RequestOption) error {
	_, err := s.Request(ctx, http.MethodPut, discord.EndpointThreadMember(id, "@me"), nil, options...)
	return err
}

// ThreadLeave removes current user.User to a thread.
func (s Requester) ThreadLeave(ctx context.Context, id string, options ...discord.RequestOption) error {
	_, err := s.Request(ctx, http.MethodDelete, discord.EndpointThreadMember(id, "@me"), nil, options...)
	return err
}

// ThreadMemberAdd adds a user.Member to a thread.
func (s Requester) ThreadMemberAdd(ctx context.Context, threadID, memberID string, options ...discord.RequestOption) error {
	_, err := s.Request(ctx, http.MethodPut, discord.EndpointThreadMember(threadID, memberID), nil, options...)
	return err
}

// ThreadMemberRemove removes a user.Member from a thread.
func (s Requester) ThreadMemberRemove(ctx context.Context, threadID, memberID string, options ...discord.RequestOption) error {
	_, err := s.Request(ctx, http.MethodDelete, discord.EndpointThreadMember(threadID, memberID), nil, options...)
	return err
}

// ThreadMember returns channel.ThreadMember for the specified user.Member of the thread.
//
// If withMember is true, it includes a guild member object.
func (s Requester) ThreadMember(ctx context.Context, threadID, memberID string, withMember bool, options ...discord.RequestOption) (*channel.ThreadMember, error) {
	uri := discord.EndpointThreadMember(threadID, memberID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	body, err := s.Request(ctx, http.MethodGet, uri, nil, options...)

	if err != nil {
		return nil, err
	}

	var m channel.ThreadMember
	return &m, s.Unmarshal(body, &m)
}

// ThreadMembers returns all user.Member of specified thread.
//
// limit is the max number of thread members to return (1-100). Defaults to 100.
// If afterID is set, every member ID will be after this.
// If withMember is true, it includes a guild member object.
func (s Requester) ThreadMembers(ctx context.Context, threadID string, limit int, withMember bool, afterID string, options ...discord.RequestOption) ([]*channel.ThreadMember, error) {
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

	body, err := s.Request(ctx, http.MethodGet, uri, nil, options...)

	if err != nil {
		return nil, err
	}

	var ms []*channel.ThreadMember
	return ms, s.Unmarshal(body, &ms)
}

// ThreadsActive returns all active threads in the given channel.Channel.
func (s Requester) ThreadsActive(ctx context.Context, channelID string, options ...discord.RequestOption) (*channel.ThreadsList, error) {
	body, err := s.Request(ctx, http.MethodGet, discord.EndpointChannelActiveThreads(channelID), nil, options...)
	if err != nil {
		return nil, err
	}

	var tl channel.ThreadsList
	return &tl, s.Unmarshal(body, &tl)
}

// ThreadsArchived returns archived threads in the given channel.Channel.
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func (s Requester) ThreadsArchived(ctx context.Context, channelID string, before *time.Time, limit int, options ...discord.RequestOption) (*channel.ThreadsList, error) {
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

	body, err := s.Request(ctx, http.MethodGet, endpoint, nil, options...)
	if err != nil {
		return nil, err
	}

	var tl channel.ThreadsList
	return &tl, s.Unmarshal(body, &tl)
}

// ThreadsPrivateArchived returns archived private threads in the given channel.Channel.
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func (s Requester) ThreadsPrivateArchived(ctx context.Context, channelID string, before *time.Time, limit int, options ...discord.RequestOption) (*channel.ThreadsList, error) {
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

	body, err := s.Request(ctx, http.MethodGet, endpoint, nil, options...)
	if err != nil {
		return nil, err
	}

	var tl channel.ThreadsList
	return &tl, s.Unmarshal(body, &tl)
}

// ThreadsPrivateJoinedArchived returns archived joined private threads in the given channel.Channel.
//
// If specified returns only threads before the timestamp
// limit is the optional maximum amount of threads to return.
func (s Requester) ThreadsPrivateJoinedArchived(ctx context.Context, channelID string, before *time.Time, limit int, options ...discord.RequestOption) (*channel.ThreadsList, error) {
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

	body, err := s.Request(ctx, http.MethodGet, endpoint, nil, options...)
	if err != nil {
		return nil, err
	}

	var tl channel.ThreadsList
	return &tl, s.Unmarshal(body, &tl)
}

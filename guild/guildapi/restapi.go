package guildapi

import (
	"bytes"
	"errors"
	"image"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user/invite"
)

var (
	ErrPruneDaysBounds = errors.New("the number of days should be more than or equal to 1")
	ErrGuildNoIcon     = errors.New("guild does not have an icon set")
	ErrGuildNoSplash   = errors.New("guild does not have a splash set")
)

// Requester handles everything inside the guild package.
type Requester struct {
	discord.Requester
	State *State
}

// UserGuilds returns an array of guild.UserGuild structures for all guilds.
//
// limit is the number of guilds that can be returned (max 200).
// If beforeID is set, it will return all guilds before this ID.
// If afterID is set, it will return all guilds after this ID.
// Set withCounts to true if you want to include approximate member and presence counts.
func (r Requester) UserGuilds(limit int, beforeID, afterID string, withCounts bool, options ...discord.RequestOption) ([]*guild.UserGuild, error) {
	v := url.Values{}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if withCounts {
		v.Set("with_counts", "true")
	}

	uri := discord.EndpointUserGuilds("@me")

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	body, err := r.RequestWithBucketID(http.MethodGet, uri, nil, discord.EndpointUserGuilds(""), options...)
	if err != nil {
		return nil, err
	}

	var ug []*guild.UserGuild
	return ug, r.Unmarshal(body, &ug)
}

// Invites returns the list of invite.Invite for the given guild.Guild.
func (r Requester) Invites(guildID string, options ...discord.RequestOption) ([]*invite.Invite, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuildInvites(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var st []*invite.Invite
	return st, r.Unmarshal(body, &st)
}

// Icon returns an image.Image of a guild.Guild icon.
func (r Requester) Icon(guildID string, options ...discord.RequestOption) (image.Image, error) {
	g, err := r.Guild(guildID, options...)
	if err != nil {
		return nil, err
	}

	if g.Icon == "" {
		return nil, ErrGuildNoIcon
	}

	body, err := r.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildIcon(guildID, g.Icon),
		nil,
		discord.EndpointGuildIcon(guildID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	return img, err
}

// Splash returns an image.Image of a guild.Guild splash image.
func (r Requester) Splash(guildID string, options ...discord.RequestOption) (image.Image, error) {
	g, err := r.Guild(guildID, options...)
	if err != nil {
		return nil, err
	}

	if g.Splash == "" {
		return nil, ErrGuildNoSplash
	}

	body, err := r.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildSplash(guildID, g.Splash),
		nil,
		discord.EndpointGuildSplash(guildID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	return img, err
}

// Embed returns the guild.Embed for a guild.Guild.
func (r Requester) Embed(guildID string, options ...discord.RequestOption) (*guild.Embed, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuildEmbed(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var em guild.Embed
	return &em, r.Unmarshal(body, &em)
}

// EmbedEdit edits the guild.Embed of a guild.Guild.
func (r Requester) EmbedEdit(guildID string, data *guild.Embed, options ...discord.RequestOption) error {
	_, err := r.Request(http.MethodPatch, discord.EndpointGuildEmbed(guildID), data, options...)
	return err
}

// AuditLog returns the guild.AuditLog for a guild.Guild.
//
// If provided, all returned guild.AuditLog will be filtered for the given userID.
// If provided all guild.AuditLog entries returned will be before the given beforeID.
// If provided the guild.AuditLog will be filtered for the given actionType.
// limit is the number of messages that can be returned (default 50, min 1, max 100).
func (r Requester) AuditLog(guildID, userID, beforeID string, actionType, limit int, options ...discord.RequestOption) (*guild.AuditLog, error) {
	uri := discord.EndpointGuildAuditLogs(guildID)

	v := url.Values{}
	if userID != "" {
		v.Set("user_id", userID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if actionType > 0 {
		v.Set("action_type", strconv.Itoa(actionType))
	}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	body, err := r.Request(http.MethodGet, uri, nil, options...)
	if err != nil {
		return nil, err
	}

	var al guild.AuditLog
	return &al, r.Unmarshal(body, &al)
}

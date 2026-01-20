package guildapi

import (
	"context"
	"errors"
	"image"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user/invite"
)

var (
	ErrPruneDaysBounds = errors.New("the number of days should be more than or equal to 1")
	ErrGuildNoIcon     = errors.New("guild does not have an icon set")
	ErrGuildNoSplash   = errors.New("guild does not have a splash set")
)

// Requester handles everything inside the guild package.
type Requester struct {
	REST
	Websocket
	State *State
}

// UserGuilds returns an array of guild.UserGuild structures for all guilds.
//
// limit is the number of guilds that can be returned (max 200).
// If beforeID is set, it will return all guilds before this ID.
// If afterID is set, it will return all guilds after this ID.
// Set withCounts to true if you want to include approximate member and presence counts.
func (r Requester) UserGuilds(limit int, beforeID, afterID string, withCounts bool) Request[[]*UserGuild] {
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

	return NewSimpleData[[]*UserGuild](
		r, http.MethodGet, uri,
	).WithBucketID(discord.EndpointUserGuilds(""))
}

// Invites returns the list of invite.Invite for the given guild.Guild.
func (r Requester) Invites(guildID string) Request[[]*invite.Invite] {
	return NewSimpleData[[]*invite.Invite](r, http.MethodGet, discord.EndpointGuildInvites(guildID))
}

// Icon returns an image.Image of a guild.Guild icon.
func (r Requester) Icon(guildID string) Request[image.Image] {
	return NewImage(r, http.MethodGet, "").
		WithBucketID(discord.EndpointGuildIcon(guildID, "")).
		WithPre(func(ctx context.Context, do *Do) error {
			g, err := r.Guild(guildID).Do(ctx)
			if err != nil {
				return err
			}
			if g.Icon == "" {
				return ErrGuildNoIcon
			}
			do.Endpoint = discord.EndpointGuildIcon(guildID, g.Icon)
			return nil
		})
}

// Splash returns an image.Image of a guild.Guild splash image.
func (r Requester) Splash(guildID string) Request[image.Image] {
	return NewImage(r, http.MethodGet, "").
		WithBucketID(discord.EndpointGuildSplash(guildID, "")).
		WithPre(func(ctx context.Context, do *Do) error {
			g, err := r.Guild(guildID).Do(ctx)
			if err != nil {
				return err
			}
			if g.Splash == "" {
				return ErrGuildNoSplash
			}
			do.Endpoint = discord.EndpointGuildSplash(guildID, g.Splash)
			return nil
		})
}

// Embed returns the guild.Embed for a guild.Guild.
func (r Requester) Embed(guildID string) Request[*Embed] {
	return NewSimpleData[*Embed](r, http.MethodGet, discord.EndpointGuildEmbed(guildID))
}

// EmbedEdit edits the guild.Embed of a guild.Guild.
func (r Requester) EmbedEdit(guildID string, data *Embed) Empty {
	req := NewSimple(r, http.MethodPatch, discord.EndpointGuildEmbed(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

// AuditLog returns the guild.AuditLog for a guild.Guild.
//
// If provided, all returned AuditLog will be filtered for the given userID.
// If provided all AuditLog entries returned will be before the given beforeID.
// If provided the AuditLog will be filtered for the given actionType.
// limit is the number of messages that can be returned (default 50, min 1, max 100).
func (r Requester) AuditLog(guildID, userID, beforeID string, actionType, limit int) Request[*AuditLog] {
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

	return NewSimpleData[*AuditLog](r, http.MethodGet, uri)
}

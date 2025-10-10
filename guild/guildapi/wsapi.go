package guildapi

import (
	"context"

	"github.com/nyttikord/gokord/discord"
)

type requestGuildMembersData struct {
	GuildID   string    `json:"guild_id"`
	Query     *string   `json:"query,omitempty"`
	UserIDs   *[]string `json:"user_ids,omitempty"`
	Limit     int       `json:"limit"`
	Nonce     string    `json:"nonce,omitempty"`
	Presences bool      `json:"presences"`
}

type requestGuildMembersOp struct {
	Op   discord.GatewayOpCode   `json:"op"`
	Data requestGuildMembersData `json:"d"`
}

// GatewayMembers requests user.Member from the gateway.
// It responds with event.GuildMembersChunk.
//
// query is a string that username starts with, leave empty to return every user.Member.
// limit is the maximum number of items to return, or 0 to request every user.Member matched.
// nonce to identify the event.GuildMembersChunk response.
// presences indicates whether to request presences of user.Member.
func (r Requester) GatewayMembers(ctx context.Context, guildID, query string, limit int, nonce string, presences bool) error {
	data := requestGuildMembersData{
		GuildID:   guildID,
		Query:     &query,
		Limit:     limit,
		Nonce:     nonce,
		Presences: presences,
	}
	return r.gatewayRequestMembers(ctx, data)
}

// GatewayMembersList requests user.Member from the gateway.
// It responds with event.GuildMembersChunk.
//
// userIDs are the user.Member's IDs to fetch.
// limit is the maximum number of items to return, or 0 to request every user.Member matched.
// nonce to identify the event.GuildMembersChunk response.
// presences indicates whether to request presences of user.Member.
func (r Requester) GatewayMembersList(ctx context.Context, guildID string, userIDs []string, limit int, nonce string, presences bool) error {
	data := requestGuildMembersData{
		GuildID:   guildID,
		UserIDs:   &userIDs,
		Limit:     limit,
		Nonce:     nonce,
		Presences: presences,
	}
	return r.gatewayRequestMembers(ctx, data)
}

func (r Requester) gatewayRequestMembers(ctx context.Context, data requestGuildMembersData) error {
	r.Logger().Debug("requesting guild members via gateway")

	return r.GatewayWriteStruct(ctx, requestGuildMembersOp{discord.GatewayOpCodeRequestGuildMembers, data})
}

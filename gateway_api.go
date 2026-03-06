package gokord

import (
	"context"
	"log/slog"
	"time"

	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user/status"
)

type wsAPI struct {
	*Session
	logger *slog.Logger
}

type updateStatusOp struct {
	Op   discord.GatewayOpCode `json:"op"`
	Data bot.UpdateStatusData  `json:"d"`
}

func newUpdateStatusData(idle bool, activityType types.Activity, name, url string) *bot.UpdateStatusData {
	usd := &bot.UpdateStatusData{
		Status: "online",
	}

	if idle {
		now := int(time.Now().Unix())
		usd.IdleSince = &now
	}

	if name != "" {
		usd.Activities = []*status.Activity{{
			Name: name,
			Type: activityType,
			URL:  url,
		}}
	}

	return usd
}

func (s *wsAPI) UpdateGameStatus(ctx context.Context, idle bool, name string) error {
	return s.UpdateStatusComplex(ctx, *newUpdateStatusData(idle, types.ActivityGame, name, ""))
}

func (s *wsAPI) UpdateWatchStatus(ctx context.Context, idle bool, name string) error {
	return s.UpdateStatusComplex(ctx, *newUpdateStatusData(idle, types.ActivityWatching, name, ""))
}

func (s *wsAPI) UpdateStreamingStatus(ctx context.Context, idle bool, name string, url string) error {
	gameType := types.ActivityGame
	if url != "" {
		gameType = types.ActivityStreaming
	}
	return s.UpdateStatusComplex(ctx, *newUpdateStatusData(idle, gameType, name, url))
}

func (s *wsAPI) UpdateListeningStatus(ctx context.Context, name string) error {
	return s.UpdateStatusComplex(ctx, *newUpdateStatusData(false, types.ActivityListening, name, ""))
}

func (s *wsAPI) UpdateCustomStatus(ctx context.Context, state string) error {
	data := bot.UpdateStatusData{
		Status: "online",
	}

	if state != "" {
		// Discord requires a non-empty activity name, therefore we provide "Custom Status" as a placeholder.
		data.Activities = []*status.Activity{{
			Name:  "Custom Status",
			Type:  types.ActivityCustom,
			State: state,
		}}
	}

	return s.UpdateStatusComplex(ctx, data)
}

func (s *wsAPI) UpdateStatusComplex(ctx context.Context, usd bot.UpdateStatusData) (err error) {
	if len(usd.Activities) == 0 {
		usd.Activities = make([]*status.Activity, 0)
	}
	return s.GatewayWriteStruct(ctx, updateStatusOp{discord.GatewayOpCodePresenceUpdate, usd})
}

type requestGuildMembersData struct {
	GuildID   uint64   `json:"guild_id,string"`
	Query     *string  `json:"query,omitempty"`
	UserIDs   []uint64 `json:"user_ids,omitempty,string"`
	Limit     int      `json:"limit"`
	Nonce     string   `json:"nonce,omitempty"`
	Presences bool     `json:"presences"`
}

type requestGuildMembersOp struct {
	Op   discord.GatewayOpCode   `json:"op"`
	Data requestGuildMembersData `json:"d"`
}

func (s *wsAPI) GatewayMembers(ctx context.Context, guildID uint64, query string, limit int, nonce string, presences bool) error {
	data := requestGuildMembersData{
		GuildID:   guildID,
		Query:     &query,
		Limit:     limit,
		Nonce:     nonce,
		Presences: presences,
	}
	return s.gatewayRequestMembers(ctx, data)
}

func (s *wsAPI) GatewayMembersList(ctx context.Context, guildID uint64, userIDs []uint64, limit int, nonce string, presences bool) error {
	data := requestGuildMembersData{
		GuildID:   guildID,
		UserIDs:   userIDs,
		Limit:     limit,
		Nonce:     nonce,
		Presences: presences,
	}
	return s.gatewayRequestMembers(ctx, data)
}

func (s *wsAPI) gatewayRequestMembers(ctx context.Context, data requestGuildMembersData) error {
	s.logger.Debug("requesting guild members via gateway")

	return s.GatewayWriteStruct(ctx, requestGuildMembersOp{discord.GatewayOpCodeRequestGuildMembers, data})
}

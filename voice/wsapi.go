package voice

import (
	"context"
	"sync"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/state"
	"github.com/nyttikord/gokord/user"
)

// Requester handles everything inside the voice package.
type Requester struct {
	discord.Requester

	sync.RWMutex
	Connections map[string]*Connection
}

type channelJoinData struct {
	GuildID   *string `json:"guild_id"`
	ChannelID *string `json:"channel_id"`
	SelfMute  bool    `json:"self_mute"`
	SelfDeaf  bool    `json:"self_deaf"`
}

type channelJoinOp struct {
	Op   discord.GatewayOpCode `json:"op"`
	Data channelJoinData       `json:"d"`
}

// ChannelJoin joins the requester user to a voice channel.Channel.
//
// mute indicates whether you will be set to muted upon joining.
// deaf indicates whether you will be set to deafened upon joining.
func (r *Requester) ChannelJoin(ctx context.Context, guildID, channelID string, mute, deaf bool) (*Connection, error) {
	r.RLock()
	v, _ := r.Connections[guildID]
	r.RUnlock()

	if v == nil {
		v = &Connection{Logger: r.Requester.Logger().With("module", "voice")}
		r.Lock()
		r.Connections[guildID] = v
		r.Unlock()
	}

	v.Lock()
	v.GuildID = guildID
	v.ChannelID = channelID
	v.deaf = deaf
	v.mute = mute
	v.requester = r
	v.Unlock()

	err := r.ChannelJoinManual(ctx, guildID, channelID, mute, deaf)
	if err != nil {
		return nil, err
	}

	// TODO: doesn't exactly work perfect yet...
	err = v.waitUntilConnected()
	if err != nil {
		r.Logger().Error("waiting for voice to connect", "error", err)
		v.Close()
		return nil, err
	}

	return v, nil
}

// ChannelJoinManual initiates a voice requester to a voice channel.Channel, but does not complete it.
//
// This should only be used when event.VoiceServerUpdate will be intercepted and used elsewhere.
//
// mute indicates whether you will be set to muted upon joining.
// deaf indicates whether you will be set to deafened upon joining.
func (r *Requester) ChannelJoinManual(ctx context.Context, guildID, channelID string, mute, deaf bool) error {
	var cID *string
	if channelID == "" {
		cID = nil
	} else {
		cID = &channelID
	}

	// Send the request to Discord that we want to join the voice channel
	data := channelJoinOp{discord.GatewayOpCodeVoiceStateUpdate, channelJoinData{&guildID, cID, mute, deaf}}
	return r.GatewayWriteStruct(ctx, data)
}

// UpdateState updates the user.VoiceState (received during event.VoiceStateUpdate).
func (r *Requester) UpdateState(u *user.VoiceState, ss state.Bot) {
	// We only care about events that are about us.
	if ss.User().ID != u.UserID {
		return
	}

	if u.ChannelID == "" {
		return
	}

	r.RLock()
	v, exists := r.Connections[u.GuildID]
	r.RUnlock()
	if !exists {
		return
	}

	v.Lock()
	v.UserID = u.UserID
	v.sessionID = u.SessionID
	v.ChannelID = u.ChannelID
	v.Unlock()
}

// UpdateServer handles the event.VoiceServerUpdate.
//
// This is also fired if the guild.Guild's voice region changes while connected to a voice channel.Channel.
// In that case, need to re-establish connection to the new region endpoint.
func (r *Requester) UpdateServer(ctx context.Context, token string, guildID string, endpoint string) {
	r.Logger().Debug("voice server update")

	r.RLock()
	v, exists := r.Connections[guildID]
	r.RUnlock()

	if !exists {
		return
	}

	// If currently connected to voice ws/udp, then disconnect.
	// Has no effect if not connected.
	v.Close()

	if endpoint == "" {
		r.Logger().Warn("endpoint is not defined, voice server was not reallocated?")
		return
	}

	v.Lock()
	v.token = token
	v.endpoint = endpoint
	v.GuildID = guildID
	v.Unlock()

	err := v.open(ctx)
	if err != nil {
		r.Logger().Error("opening voice connection", "error", err)
	}
}

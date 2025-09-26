package voice

import (
	"sync"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/state"
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
	Op   int             `json:"op"`
	Data channelJoinData `json:"d"`
}

// ChannelJoin joins the requester user to a voice channel.Channel.
//
// mute indicates whether you will be set to muted upon joining.
// deaf indicates whether you will be set to deafened upon joining.
func (r *Requester) ChannelJoin(guildID, channelID string, mute, deaf bool) (*Connection, error) {
	r.RLock()
	voice, _ := r.Connections[guildID]
	r.RUnlock()

	if voice == nil {
		voice = &Connection{Logger: r.Requester}
		r.Lock()
		r.Connections[guildID] = voice
		r.Unlock()
	}

	voice.Lock()
	voice.GuildID = guildID
	voice.ChannelID = channelID
	voice.deaf = deaf
	voice.mute = mute
	voice.requester = r
	voice.Unlock()

	err := r.ChannelJoinManual(guildID, channelID, mute, deaf)
	if err != nil {
		return nil, err
	}

	// TODO: doesn't exactly work perfect yet...
	err = voice.waitUntilConnected()
	if err != nil {
		r.LogError(err, "waiting for voice to connect")
		voice.Close()
		return nil, err
	}

	return voice, nil
}

// ChannelJoinManual initiates a voice requester to a voice channel.Channel, but does not complete it.
//
// This should only be used when event.VoiceServerUpdate will be intercepted and used elsewhere.
//
// mute indicates whether you will be set to muted upon joining.
// deaf indicates whether you will be set to deafened upon joining.
func (r *Requester) ChannelJoinManual(guildID, channelID string, mute, deaf bool) error {
	var cID *string
	if channelID == "" {
		cID = nil
	} else {
		cID = &channelID
	}

	// Send the request to Discord that we want to join the voice channel
	data := channelJoinOp{4, channelJoinData{&guildID, cID, mute, deaf}}
	return r.GatewayWriteStruct(data)
}

// OnVoiceStateUpdate handles event.VoiceStateUpdate.
func (r *Requester) OnVoiceStateUpdate(st *event.VoiceStateUpdate, ss state.Bot) {
	// If we don't have a connection for the channel, don't bother
	if st.ChannelID == "" {
		return
	}

	// Check if we have a voice connection to update
	r.RLock()
	v, exists := r.Connections[st.GuildID]
	r.RUnlock()
	if !exists {
		return
	}

	// We only care about events that are about us.
	if ss.User().ID != st.UserID {
		return
	}

	// Store the SessionID for later use.
	v.Lock()
	v.UserID = st.UserID
	v.sessionID = st.SessionID
	v.ChannelID = st.ChannelID
	v.Unlock()
}

// OnVoiceServerUpdate handles the event.VoiceServerUpdate.
//
// This is also fired if the guild's voice region changes while connected to a voice channel.
// In that case, need to re-establish connection to the new region endpoint.
func (r *Requester) OnVoiceServerUpdate(st *event.VoiceServerUpdate) {
	r.LogDebug("voice server update")

	r.RLock()
	v, exists := r.Connections[st.GuildID]
	r.RUnlock()

	// If no VoiceConnection exists, just skip this
	if !exists {
		return
	}

	// If currently connected to v ws/udp, then disconnect.
	// Has no effect if not connected.
	v.Close()

	// Store values for later use
	v.Lock()
	v.token = st.Token
	v.endpoint = st.Endpoint
	v.GuildID = st.GuildID
	v.Unlock()

	// Open a connection to the v server
	err := v.open()
	if err != nil {
		r.LogError(err, "opening v connection")
	}
}

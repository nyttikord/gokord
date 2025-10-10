// Package voice contains everything related to Discord voice Connection.
//
// This package is incomplete.
package voice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/nyttikord/gokord/discord"
)

// Connection holds all the data and functions related to a Discord voice connection.
type Connection struct {
	sync.RWMutex
	Logger *slog.Logger

	Ready        bool // If Ready, Connection is ready to send/receive audio.
	UserID       string
	GuildID      string
	ChannelID    string
	deaf         bool
	mute         bool
	speaking     bool
	reconnecting bool // If reconnecting, voice connection is trying to Reconnect.

	OpusSend chan []byte  // OpusSend is the chan for sending opus audio
	OpusRecv chan *Packet // OpusRecv is the chan for receiving opus audio

	wsConn    *websocket.Conn
	wsMutex   sync.Mutex
	udpConn   *net.UDPConn
	requester *Requester

	sessionID string
	token     string
	endpoint  string

	// close is used to send a close signal to goroutines.
	close chan struct{}

	// connected is used to allow blocking until connected.
	connected chan bool

	// Used to pass the sessionID from UpdateState
	// sessionRecv chan string UNUSED ATM

	op4 op4
	op2 op2

	voiceSpeakingUpdateHandlers []SpeakingUpdateHandler
}

// SpeakingUpdateHandler type provides a function definition for the SpeakingUpdate event.
type SpeakingUpdateHandler func(vc *Connection, vs *SpeakingUpdate)

// Speaking sends a speaking notification to Discord over the voice websocket.
// This must be sent as true prior to sending audio and should be set to false once finished sending audio.
//
// b is true when you speak and false when you don't.
func (v *Connection) Speaking(ctx context.Context, bl bool) error {
	v.Logger.Debug("called", "speaks", bl)

	type speakingData struct {
		Speaking bool `json:"speaking"`
		Delay    int  `json:"delay"`
	}

	type speakingOp struct {
		Op   discord.VoiceOpCode `json:"op"` // Always 5
		Data speakingData        `json:"d"`
	}

	if v.wsConn == nil {
		return fmt.Errorf("no VoiceConnection websocket")
	}

	data := speakingOp{discord.VoiceOpCodeSessionSpeaking, speakingData{bl, 0}}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	v.wsMutex.Lock()
	err = v.wsConn.Write(ctx, websocket.MessageText, b)
	v.wsMutex.Unlock()

	v.Lock()
	defer v.Unlock()
	if err != nil {
		v.speaking = false
		v.Logger.Error("writing json", "error", err)
		return err
	}

	v.speaking = bl

	return err
}

// Disconnect from this voice channel.Channel and closes the websocket and udp connections to Discord.
func (v *Connection) Disconnect(ctx context.Context) error {
	// Send a OP4 with a nil channel to disconnect
	v.Lock()
	defer v.Unlock()
	if v.sessionID != "" {
		data := channelJoinOp{discord.GatewayOpCodeVoiceStateUpdate, channelJoinData{&v.GuildID, nil, true, true}}
		err := v.requester.GatewayWriteStruct(ctx, data)
		if err != nil {
			return err
		}
		v.sessionID = ""
	}

	// Close websocket and udp connections
	v.Unlock()
	v.Close()
	v.Lock()

	v.Logger.Info("Deleting VoiceConnection", "guild", v.GuildID)

	v.requester.Lock()
	delete(v.requester.Connections, v.GuildID)
	v.requester.Unlock()
	return nil
}

// Close closes the voice ws and udp connections.
// Use Disconnect to have a better disconnection process.
func (v *Connection) Close() {
	v.Lock()
	defer v.Unlock()

	v.Ready = false
	v.speaking = false

	if v.close != nil {
		v.Logger.Debug("closing voice goroutines")
		v.close <- struct{}{}
		close(v.close)
		v.close = nil
	}

	if v.udpConn != nil {
		v.Logger.Debug("closing udp")
		err := v.udpConn.Close()
		if err != nil {
			v.Logger.Error("closing udp connection", "error", err)
		}
		v.udpConn = nil
	}

	if v.wsConn != nil {
		v.Logger.Debug("sending close frame")

		// To cleanly close a connection, a client should send a close frame and wait for the server to close the
		// connection.
		v.wsMutex.Lock()
		err := v.wsConn.Close(websocket.StatusNormalClosure, "")
		v.wsMutex.Unlock()
		if err != nil {
			v.Logger.Error("closing websocket", "error", err)
			v.Logger.Debug("closing websocket")
			err = v.wsConn.CloseNow()
			if err != nil {
				v.Logger.Error("closing websocket", "error", err)
			}
		}

		v.wsConn = nil
	}
}

// AddHandler adds a Handler for SpeakingUpdate events.
func (v *Connection) AddHandler(h SpeakingUpdateHandler) {
	v.Lock()
	defer v.Unlock()

	v.voiceSpeakingUpdateHandlers = append(v.voiceSpeakingUpdateHandlers, h)
}

// SpeakingUpdate is a struct for a SpeakingUpdate event.
type SpeakingUpdate struct {
	UserID   string `json:"user_id"`
	SSRC     int    `json:"ssrc"`
	Speaking bool   `json:"speaking"`
}

// op4 stores the data for the voice operation 4 websocket event which provides us with the NaCl SecretBox encryption key.
type op4 struct {
	SecretKey [32]byte `json:"secret_key"`
	Mode      string   `json:"mode"`
}

// op2 stores the data for the voice operation 2 websocket event which is sort of like the voice READY packet.
type op2 struct {
	SSRC              uint32        `json:"ssrc"`
	Port              int           `json:"port"`
	Modes             []string      `json:"modes"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	IP                string        `json:"ip"`
}

// waitUntilConnected waits until the Connection become ready, if it does not become ready it returns an error.
func (v *Connection) waitUntilConnected() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	chann := make(chan struct{}, 1)

	go func() {
		ready := false
		for !ready {
			v.RLock()
			ready = v.Ready
			v.RUnlock()
			time.Sleep(1 * time.Second)
		}
		chann <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-chann:
		return nil
	}
}

// Open opens a voice connection.
// This should be called after Requester.ChannelJoin is used and the data VOICE websocket events are captured.
func (v *Connection) open(ctx context.Context) error {
	v.Lock()
	defer v.Unlock()

	// Don't open a websocket if one is already open
	if v.wsConn != nil {
		v.Logger.Warn("refusing to overwrite non-nil websocket")
		return nil
	}

	// TODO temp? loop to wait for the SessionID
	ctx2, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	chann := make(chan struct{}, 1)

	go func() {
		sID := v.sessionID
		for len(sID) == 0 {
			// Release the lock, so sessionID can be populated upon receiving an event.VoiceStateUpdate.
			v.Unlock()
			time.Sleep(50 * time.Millisecond)
			v.Lock()
			sID = v.sessionID
		}
		chann <- struct{}{}
	}()

	select {
	case <-ctx2.Done():
		return ctx.Err()
	case <-chann:
	}

	vg := "wss://" + strings.TrimSuffix(v.endpoint, ":80")
	v.Logger.Debug("connecting to voice endpoint", "endpoint", vg)
	var err error
	v.wsConn, _, err = v.requester.GatewayDial(context.Background(), vg, nil)
	if err != nil {
		v.Logger.Error("connecting to voice endpoint", "error", err, "endpoint", vg)
		v.Logger.Debug("voice", "struct", v)
		return err
	}

	type handshakeData struct {
		ServerID  string `json:"server_id"`
		UserID    string `json:"user_id"`
		SessionID string `json:"session_id"`
		Token     string `json:"token"`
	}
	type handshakeOp struct {
		Op   discord.VoiceOpCode `json:"op"` // Always 0
		Data handshakeData       `json:"d"`
	}
	data := handshakeOp{discord.VoiceOpCodeIdentify, handshakeData{v.GuildID, v.UserID, v.sessionID, v.token}}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	v.wsMutex.Lock()
	err = v.wsConn.Write(ctx, websocket.MessageText, b)
	v.wsMutex.Unlock()
	if err != nil {
		v.Logger.Error("sending init packet", "error", err)
		return err
	}

	v.close = make(chan struct{})
	go v.wsListen(ctx, v.wsConn, v.close)

	// add loop/check for Ready bool here?
	// then return false if not ready?
	// but then wsListen will also err.
	return nil
}

// Reconnect will close down a voice connection then immediately try to Reconnect to that requester.
//
// NOTE: This func is messy and a WIP while I find what works.
// It will be cleaned up once a proven stable option is flushed out.
// aka: this is ugly shit code, please don't judge too harshly.
func (v *Connection) Reconnect(ctx context.Context) {
	v.Logger.Debug("called")

	v.Lock()
	if v.reconnecting {
		v.Logger.Warn("already reconnecting to channel exiting", "channel", v.ChannelID)
		v.Unlock()
		return
	}
	v.reconnecting = true
	v.Unlock()

	defer func() {
		v.Lock()
		v.reconnecting = false
		v.Unlock()
	}()

	// Close any currently open connections
	v.Close()

	wait := time.Duration(1)
	for {
		time.Sleep(wait * time.Second)
		wait *= min(wait*2, 600)

		v.Logger.Info("trying to reconnect to channel", "channel", v.ChannelID)

		_, err := v.requester.ChannelJoin(ctx, v.GuildID, v.ChannelID, v.mute, v.deaf)
		if err == nil {
			v.Logger.Info("successfully reconnected to channel", "channel", v.ChannelID)
			return
		}

		v.Logger.Error("reconnecting to channel", "error", err, "channel", v.ChannelID)

		// if the Reconnect above didn't work lets just send a disconnect packet to reset things.
		// Send a OP4 with a nil channel to disconnect.
		data := channelJoinOp{discord.GatewayOpCodeVoiceStateUpdate, channelJoinData{&v.GuildID, nil, true, true}}
		err = v.requester.GatewayWriteStruct(ctx, data)
		if err != nil {
			v.Logger.Error("sending disconnect packet", "error", err)
		}
	}
}

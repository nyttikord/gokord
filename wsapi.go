package gokord

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
)

var (
	// ErrWSAlreadyOpen is thrown when you attempt to open a websocket that already is open.
	ErrWSAlreadyOpen = errors.New("web socket already opened")
	// ErrWSNotFound is thrown when you attempt to use a websocket that doesn't exist
	ErrWSNotFound = errors.New("no websocket connection exists")
	// ErrWSShardBounds is thrown when you try to use a shard ID that is more than the total shard count
	ErrWSShardBounds = errors.New("ShardID must be less than ShardCount")
)

type resumePacket struct {
	Op   int `json:"op"`
	Data struct {
		Token     string `json:"token"`
		SessionID string `json:"session_id"`
		Sequence  int64  `json:"seq"`
	} `json:"d"`
}

func (s *Session) GatewayWriteStruct(v any) error {
	s.RLock()
	defer s.RUnlock()
	if s.wsConn == nil {
		return ErrWSNotFound
	}

	s.wsMutex.Lock()
	err := s.wsConn.WriteJSON(v)
	s.wsMutex.Unlock()
	return err
}

// Open creates a websocket connection to Discord.
// https://discord.com/developers/docs/topics/gateway#connecting
func (s *Session) Open() error {
	s.Lock()
	defer s.Unlock()

	// If the websock is already open, bail out here.
	if s.wsConn != nil {
		return ErrWSAlreadyOpen
	}

	sequence := s.sequence.Load()

	var gateway string
	var err error
	// Get the gateway to use for the Websocket connection
	if sequence != 0 && s.sessionID != "" && s.resumeGatewayURL != "" {
		s.LogDebug("using resume gateway %s", s.resumeGatewayURL)
		gateway = s.resumeGatewayURL
	} else {
		if s.gateway == "" {
			s.gateway, err = s.Gateway()
			if err != nil {
				return err
			}
		}

		gateway = s.gateway
	}

	// Add the version and encoding to the URL
	gateway += "?v=" + discord.APIVersion + "&encoding=json"

	// Connect to the Gateway
	s.LogInfo("connecting to gateway %s", gateway)
	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	s.wsConn, _, err = s.Dialer.Dial(gateway, header)
	if err != nil {
		s.LogError(err, "connecting to gateway %s", s.gateway)
		s.gateway = "" // clear cached gateway
		s.wsConn = nil // Just to be safe.
		return err
	}

	s.wsConn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	defer func() {
		if err != nil {
			s.ForceClose(false)
		}
	}()

	// The first response from Discord should be an Op 10 (Hello) Packet.
	// When processed by onEvent the heartbeat goroutine will be started.
	mt, m, err := s.wsConn.ReadMessage()
	if err != nil {
		return err
	}
	e, err := s.onEvent(mt, m)
	if err != nil {
		return err
	}
	if e.Operation != 10 {
		err = fmt.Errorf("expecting Op 10, got Op %d instead", e.Operation)
		return err
	}
	s.LogDebug("Op 10 Hello Packet received from Discord")
	s.LastHeartbeatAck = time.Now().UTC()
	var h helloOp
	if err = json.Unmarshal(e.RawData, &h); err != nil {
		err = fmt.Errorf("error unmarshalling helloOp, %s", err)
		return err
	}

	// Now we send either an Op 2 Identity if this is a brand new connection or Op 6 Resume if we are resuming an
	// existing connection.
	if s.sessionID == "" && sequence == 0 {
		// Send Op 2 Identity Packet
		err = s.identify()
		if err != nil {
			err = fmt.Errorf("error sending identify packet to gateway, %s, %s", s.gateway, err)
			return err
		}
	} else {
		// Send Op 6 Resume Packet
		p := resumePacket{}
		p.Op = 6
		p.Data.Token = s.Identify.Token
		p.Data.SessionID = s.sessionID
		p.Data.Sequence = sequence

		s.LogInfo("sending resume packet to gateway")
		s.wsMutex.Lock()
		err = s.wsConn.WriteJSON(p)
		s.wsMutex.Unlock()
		if err != nil {
			err = fmt.Errorf("error sending gateway resume packet, %s, %s", s.gateway, err)
			return err
		}
	}

	// Now Discord should send us a READY or RESUMED packet.
	mt, m, err = s.wsConn.ReadMessage()
	if err != nil {
		return err
	}
	e, err = s.onEvent(mt, m)
	if err != nil {
		return err
	}
	if e.Type != `READY` && e.Type != `RESUMED` {
		// This is not fatal, but it does not follow their API documentation.
		s.LogWarn("Expected READY/RESUMED, instead got:\n%#v\n", e)
	}

	s.LogDebug("We are now connected to Discord, emitting connect event")
	s.EventManager().EmitEvent(s, event.ConnectType, &event.Connect{})

	// A VoiceConnections map is a hard requirement for Voice.
	// XXX: can this be moved to when opening a voice connection?
	if s.VoiceConnections == nil {
		s.LogDebug("creating new VoiceConnections map")
		s.VoiceConnections = make(map[string]*VoiceConnection)
	}

	// Create listening chan outside of listen, as it needs to happen inside the mutex lock and needs to exist before
	// calling heartbeat and listen goroutines.
	s.listening = make(chan any)

	// Start sending heartbeats and reading messages from Discord.
	go s.heartbeat(s.wsConn, s.listening, h.HeartbeatInterval)
	go s.listen(s.wsConn, s.listening)

	return nil
}

// listen polls the websocket connection for events, it will stop when the listening channel is closed, or an error
// occurs.
func (s *Session) listen(wsConn *websocket.Conn, listening <-chan any) {
	for {
		messageType, message, err := wsConn.ReadMessage()

		if err != nil {
			// Detect if we have been closed manually.
			// If a Close() has already happened, the websocket we are listening on will be different to the current
			// session.
			s.RLock()
			sameConnection := s.wsConn == wsConn
			s.RUnlock()

			if sameConnection {
				s.LogError(err, "reading from gateway %s websocket", s.gateway)
				// There has been an error reading, close the websocket so that OnDisconnect event is emitted.
				err = s.Close()
				if err != nil {
					s.LogError(err, "error closing session connection, force closing")
					s.ForceClose(false)
				}

				s.LogInfo("calling reconnect() now")
				s.reconnect()
			}
			return
		}

		select {
		case <-listening:
			return
		default:
			_, err = s.onEvent(messageType, message)
			if err != nil {
				s.LogError(err, "handling event")
			}
		}
	}
}

type heartbeatOp struct {
	Op   int   `json:"op"`
	Data int64 `json:"d"`
}

type helloOp struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

// FailedHeartbeatAcks is the Number of heartbeat intervals to wait until forcing a connection restart.
const FailedHeartbeatAcks = 5 * time.Millisecond

// HeartbeatLatency returns the latency between heartbeat acknowledgement and heartbeat send.
func (s *Session) HeartbeatLatency() time.Duration {

	return s.LastHeartbeatAck.Sub(s.LastHeartbeatSent)

}

// heartbeat sends regular heartbeats to Discord so it knows the client is still connected.
// If you do not send these heartbeats Discord will disconnect the websocket connection after a few seconds.
func (s *Session) heartbeat(wsConn *websocket.Conn, listening <-chan any, heartbeatIntervalMsec time.Duration) {
	if listening == nil || wsConn == nil {
		return
	}

	var err error
	ticker := time.NewTicker(heartbeatIntervalMsec * time.Millisecond)
	defer ticker.Stop()

	for {
		s.RLock()
		last := s.LastHeartbeatAck
		s.RUnlock()
		sequence := s.sequence.Load()
		s.LogDebug("sending gateway websocket heartbeat seq %d", sequence)
		s.wsMutex.Lock()
		s.LastHeartbeatSent = time.Now().UTC()
		err = wsConn.WriteJSON(heartbeatOp{1, sequence})
		s.wsMutex.Unlock()
		if err != nil || time.Now().UTC().Sub(last) > (heartbeatIntervalMsec*FailedHeartbeatAcks) {
			if err != nil {
				s.LogError(err, "sending heartbeat to gateway %s", s.gateway)
			} else {
				s.LogWarn("haven't gotten a heartbeat ACK in %v, triggering a reconnection", time.Now().UTC().Sub(last))
			}
			err = s.Close()
			if err != nil {
				s.LogError(err, "error closing session connection, force closing")
				s.ForceClose(false)
			}
			s.reconnect()
			return
		}
		s.Lock()
		s.DataReady = true
		s.Unlock()

		select {
		case <-ticker.C:
			// continue loop and send heartbeat
		case <-listening:
			return
		}
	}
}

// onEvent is the "event handler" for all messages received on the
// Discord Gateway API websocket connection.
//
// If you use the AddHandler() function to register a handler for a
// specific event this function will pass the event along to that handler.
//
// If you use the AddHandler() function to register a handler for the
// "OnEvent" event then all events will be passed to that handler.
func (s *Session) onEvent(messageType int, message []byte) (*event.Event, error) {
	var err error
	var reader io.Reader
	reader = bytes.NewBuffer(message)

	// If this is a compressed message, uncompress it.
	if messageType == websocket.BinaryMessage {

		z, err2 := zlib.NewReader(reader)
		if err2 != nil {
			return nil, err2
		}

		defer func() {
			err3 := z.Close()
			if err3 != nil {
				s.LogError(err3, "closing zlib")
			}
		}()

		reader = z
	}

	// Decode the event into an Event struct.
	var e *event.Event
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&e); err != nil {
		return e, err
	}

	//s.LogDebug("Op: %d, Seq: %d, Type: %s, Data: %s\n\n", e.Operation, e.Sequence, e.Type, string(e.RawData))

	// Ping request.
	// Must respond with a heartbeat packet within 5 seconds
	if e.Operation == 1 {
		s.LogDebug("sending heartbeat in response to Op1")
		s.wsMutex.Lock()
		err = s.wsConn.WriteJSON(heartbeatOp{1, s.sequence.Load()})
		s.wsMutex.Unlock()
		if err != nil {
			return e, err
		}

		return e, nil
	}

	// Reconnect
	// Must immediately disconnect from gateway and reconnect to new gateway.
	if e.Operation == 7 {
		s.LogInfo("Closing and reconnecting in response to Op7")
		err = s.CloseWithCode(websocket.CloseServiceRestart)
		if err != nil {
			s.LogError(err, "error closing session connection, force closing")
			s.ForceClose(false)
		}
		s.reconnect()
		return e, nil
	}

	// Invalid Session
	// Must respond with an Identify packet.
	if e.Operation == 9 {
		s.LogWarn("Invalid session received, reconnecting")
		err = s.CloseWithCode(websocket.CloseServiceRestart)
		if err != nil {
			s.LogError(err, "error closing session connection, force closing")
			s.ForceClose(false)
		}

		var resumable bool
		if err = json.Unmarshal(e.RawData, &resumable); err != nil {
			return e, err
		}

		if !resumable {
			s.LogInfo("Gateway session is not resumable, discarding its information")
			s.resumeGatewayURL = ""
			s.sessionID = ""
			s.sequence.Store(0)
		}

		s.reconnect()
		return e, nil
	}

	if e.Operation == 10 {
		// Op10 is handled by Open()
		return e, nil
	}

	if e.Operation == 11 {
		s.Lock()
		s.LastHeartbeatAck = time.Now().UTC()
		s.Unlock()
		s.LogDebug("got heartbeat ACK")
		return e, nil
	}

	// Do not try to Dispatch a non-Dispatch Message
	if e.Operation != 0 {
		// But we probably should be doing something with them.
		// TEMP
		s.LogWarn("unknown Op: %d, Seq: %d, Type: %s, Data: %s, message: %s", e.Operation, e.Sequence, e.Type, string(e.RawData), string(message))
		return e, nil
	}

	// Store the message sequence
	s.sequence.Store(e.Sequence)

	// Map event to registered event handlers and pass it along to any registered handlers.
	if eh, ok := event.GetInterfaceProvider(e.Type); ok {
		e.Struct = eh.New()

		// Attempt to unmarshal our event.
		if err = json.Unmarshal(e.RawData, e.Struct); err != nil {
			s.LogWarn("failed to unmarshal %s event, data: %s", e.Type, e.RawData)
			// READY events are always emitted
			if e.Type != event.ReadyType {
				return nil, err
			}
		}

		s.EventManager().EmitEvent(s, e.Type, e.Struct)
	} else {
		s.LogWarn("unknown event: Op: %d, Seq: %d, Type: %s, Data: %s", e.Operation, e.Sequence, e.Type, string(e.RawData))
		s.EventManager().EmitEvent(s, event.EventType, e)
	}

	return e, nil
}

type voiceChannelJoinData struct {
	GuildID   *string `json:"guild_id"`
	ChannelID *string `json:"channel_id"`
	SelfMute  bool    `json:"self_mute"`
	SelfDeaf  bool    `json:"self_deaf"`
}

type voiceChannelJoinOp struct {
	Op   int                  `json:"op"`
	Data voiceChannelJoinData `json:"d"`
}

// ChannelVoiceJoin joins the session user to a voice channel.Channel.
//
// mute indicates whether you will be set to muted upon joining.
// deaf indicates whether you will be set to deafened upon joining.
func (s *Session) ChannelVoiceJoin(guildID, channelID string, mute, deaf bool) (*VoiceConnection, error) {
	s.RLock()
	voice, _ := s.VoiceConnections[guildID]
	s.RUnlock()

	if voice == nil {
		voice = &VoiceConnection{stdLogger: stdLogger{Level: s.GetLevel()}}
		s.Lock()
		s.VoiceConnections[guildID] = voice
		s.Unlock()
	}

	voice.Lock()
	voice.GuildID = guildID
	voice.ChannelID = channelID
	voice.deaf = deaf
	voice.mute = mute
	voice.session = s
	voice.Unlock()

	err := s.ChannelVoiceJoinManual(guildID, channelID, mute, deaf)
	if err != nil {
		return nil, err
	}

	// TODO: doesn't exactly work perfect yet...
	err = voice.waitUntilConnected()
	if err != nil {
		s.LogError(err, "waiting for voice to connect")
		voice.Close()
		return nil, err
	}

	return voice, nil
}

// ChannelVoiceJoinManual initiates a voice session to a voice channel.Channel, but does not complete it.
//
// This should only be used when the VoiceServerUpdate will be intercepted and used elsewhere.
//
// mute indicates whether you will be set to muted upon joining.
// deaf indicates whether you will be set to deafened upon joining.
func (s *Session) ChannelVoiceJoinManual(guildID, channelID string, mute, deaf bool) error {
	var cID *string
	if channelID == "" {
		cID = nil
	} else {
		cID = &channelID
	}

	// Send the request to Discord that we want to join the voice channel
	data := voiceChannelJoinOp{4, voiceChannelJoinData{&guildID, cID, mute, deaf}}
	s.wsMutex.Lock()
	err := s.wsConn.WriteJSON(data)
	s.wsMutex.Unlock()
	return err
}

// onVoiceStateUpdate handles event.VoiceStateUpdate.
func (s *Session) onVoiceStateUpdate(st *event.VoiceStateUpdate) {
	// If we don't have a connection for the channel, don't bother
	if st.ChannelID == "" {
		return
	}

	// Check if we have a voice connection to update
	s.RLock()
	voice, exists := s.VoiceConnections[st.GuildID]
	s.RUnlock()
	if !exists {
		return
	}

	// We only care about events that are about us.
	if s.SessionState().User().ID != st.UserID {
		return
	}

	// Store the SessionID for later use.
	voice.Lock()
	voice.UserID = st.UserID
	voice.sessionID = st.SessionID
	voice.ChannelID = st.ChannelID
	voice.Unlock()
}

// onVoiceServerUpdate handles the event.VoiceServerUpdate.
//
// This is also fired if the guild's voice region changes while connected to a voice channel.
// In that case, need to re-establish connection to the new region endpoint.
func (s *Session) onVoiceServerUpdate(st *event.VoiceServerUpdate) {
	s.LogDebug("voice server update")

	s.RLock()
	voice, exists := s.VoiceConnections[st.GuildID]
	s.RUnlock()

	// If no VoiceConnection exists, just skip this
	if !exists {
		return
	}

	// If currently connected to voice ws/udp, then disconnect.
	// Has no effect if not connected.
	voice.Close()

	// Store values for later use
	voice.Lock()
	voice.token = st.Token
	voice.endpoint = st.Endpoint
	voice.GuildID = st.GuildID
	voice.Unlock()

	// Open a connection to the voice server
	err := voice.open()
	if err != nil {
		s.LogError(err, "opening voice connection")
	}
}

type identifyOp struct {
	Op   int      `json:"op"`
	Data Identify `json:"d"`
}

// identify sends the identify packet to the gateway
func (s *Session) identify() error {
	if s.Identify.Shard[0] >= s.Identify.Shard[1] {
		return ErrWSShardBounds
	}

	// Send Identify packet to Discord
	op := identifyOp{2, s.Identify}
	s.wsMutex.Lock()
	err := s.wsConn.WriteJSON(op)
	s.wsMutex.Unlock()

	return err
}

func (s *Session) reconnect() {
	if !s.ShouldReconnectOnError {
		return
	}
	s.LogInfo("trying to reconnect to gateway")

	wait := time.Duration(1)
	err := s.Open()

	for err != nil {
		// Certain race conditions can call reconnect() twice.
		// If this happens, we just break out of the reconnect loop
		// TODO: fix this
		if errors.Is(err, ErrWSAlreadyOpen) {
			s.LogDebug("Websocket already exists, no need to reconnect")
			return
		}

		s.LogError(err, "reconnecting to gateway")

		time.Sleep(wait * time.Second)
		wait *= 2
		if wait > 600 {
			wait = 600
		}

		s.LogInfo("trying to reconnect to gateway")

		err = s.Open()
	}
	s.LogInfo("successfully reconnected to gateway")

	// I'm not sure if this is actually needed.
	// If the gw reconnect works properly, voice should stay alive
	// However, there seems to be cases where something "weird" happens.
	// So we're doing this for now just to improve stability in those edge cases.
	if !s.ShouldReconnectVoiceOnSessionError {
		return
	}
	s.RLock()
	defer s.RUnlock()
	for _, v := range s.VoiceConnections {

		s.LogInfo("reconnecting voice connection to guild %s", v.GuildID)
		go v.reconnect()

		// This is here just to prevent violently spamming the
		// voice reconnects
		time.Sleep(1 * time.Second)
	}
}

// Close closes a websocket and stops all listening/heartbeat goroutines.
func (s *Session) Close() error {
	return s.CloseWithCode(websocket.CloseNormalClosure)
}

// CloseWithCode closes a websocket using the provided closeCode and stops all listening/heartbeat goroutines.
// TODO: Add support for Voice WS/UDP connections
func (s *Session) CloseWithCode(closeCode int) error {
	s.LogInfo("closing with code %d", closeCode)
	s.Lock()
	defer s.Unlock()

	if s.wsConn == nil {
		return ErrWSNotFound
	}

	s.DataReady = false

	if s.listening != nil {
		s.LogDebug("closing listening channel")
		close(s.listening)
		s.listening = nil
	}

	for _, v := range s.VoiceConnections {
		err := v.Disconnect()
		if err != nil {
			s.LogError(err, "disconnecting voice from channel %s", v.ChannelID)
		}
	}
	// TODO: Close all active Voice Connections force stop any reconnecting voice channels

	// To cleanly close a connection, a client should send a close frame and wait for the server to close the
	// connection.
	s.LogDebug("sending close frame")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		s.wsMutex.Lock()
		err := s.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(closeCode, ""))
		s.wsMutex.Unlock()
		errChan <- err
		// TODO: waiting for Discord to close the websocket
		// I have searched a way to wait for the wsConn to be closed, but I have found nothing on it.
		// I don't know how to do this.
		// I don't know if this needed.
		// Currently, this work without issues, so it's fine I guess?
	}()

	// we do not handle it because throwing an error while sending a close message is a big error,
	// and we prevent continuing the close
	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	// required
	s.Unlock()
	s.ForceClose(true)
	s.Lock()

	return nil
}

// ForceClose the connection.
// Use Close or CloseWithCode before to have a better closing process.
func (s *Session) ForceClose(emitDisconnect bool) {
	s.Lock()
	defer s.Unlock()
	s.LogInfo("closing gateway websocket")
	err := s.wsConn.Close()
	if err != nil {
		// we handle it here because the websocket is actually closed
		s.LogError(err, "closing websocket")
	}
	s.wsConn = nil

	if emitDisconnect {
		// required
		s.Unlock()
		s.EventManager().EmitEvent(s, event.DisconnectType, &event.Disconnect{})
		s.Lock()
	}
}

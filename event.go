package gokord

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/guild"
)

// setGuildIds will set the GuildID on all the members of a guild.Guild.
// This is done as event data does not have it set.
func setGuildIds(g *guild.Guild) {
	for _, c := range g.Channels {
		c.GuildID = g.ID
	}

	for _, m := range g.Members {
		m.GuildID = g.ID
	}

	for _, vs := range g.VoiceStates {
		vs.GuildID = g.ID
	}
}

// onInterface handles all internal events and routes them to the appropriate internal handler.
func (s *Session) onInterface(i any) {
	switch t := i.(type) {
	case *event.Ready:
		for _, g := range t.Guilds {
			setGuildIds(g)
		}
		s.onReady(t)
	case *event.GuildCreate:
		setGuildIds(t.Guild)
	case *event.GuildUpdate:
		setGuildIds(t.Guild)
	case *event.VoiceServerUpdate:
		go s.voiceAPI.UpdateServer(t.Token, t.GuildID, t.Endpoint)
	case *event.VoiceStateUpdate:
		go s.voiceAPI.UpdateState(t.VoiceState, s.sessionState)
	}
	err := s.sessionState.onInterface(s, i)
	if err != nil {
		s.logger.Debug("dispatching internal event", "error", err)
	}
}

// onReady handles the ready event.
func (s *Session) onReady(r *event.Ready) {

	// Store the SessionID within the Session struct.
	s.sessionID = r.SessionID

	// Store the ResumeGatewayURL within the Session struct.
	s.resumeGatewayURL = r.ResumeGatewayURL
}

// onGatewayEvent is the "event handler" for all messages received on the Discord Gateway API websocket connection.
func (s *Session) onGatewayEvent(messageType int, message []byte) (*discord.Event, error) {
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
				s.logger.Error("closing zlib", "error", err)
			}
		}()

		reader = z
	}

	var e *discord.Event
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&e); err != nil {
		return e, err
	}

	// handle special opcode
	switch e.Operation {
	case discord.GatewayOpCodeHello: // processed by Open()
		return e, nil
	case discord.GatewayOpCodeHeartbeat: // must respond with a heartbeat packet within 5 seconds
		s.logger.Debug("sending heartbeat in response to Op1")
		return e, s.heartbeat(s.ws, s.sequence.Load())
	case discord.GatewayOpCodeReconnect: // must immediately disconnect from gateway and reconnect to new gateway
		s.logger.Info("closing and reconnecting in response to Op7")
		err = s.CloseWithCode(websocket.CloseServiceRestart)
		if err != nil {
			s.logger.Error("closing session connection, force closing", "error", err)
			s.ForceClose()
		}
		s.reconnect()
		return e, nil
	case discord.GatewayOpCodeInvalidSession: // must respond with an Identify packet
		s.logger.Warn("invalid session received, reconnecting")
		err = s.CloseWithCode(websocket.CloseServiceRestart)
		if err != nil {
			s.logger.Error("closing session connection, force closing", "error", err)
			s.ForceClose()
		}

		var resumable bool
		if err = json.Unmarshal(e.RawData, &resumable); err != nil {
			return e, err
		}

		if !resumable {
			s.logger.Info("gateway session is not resumable, discarding its information")
			s.resumeGatewayURL = ""
			s.sessionID = ""
			s.sequence.Store(0)
		}

		s.reconnect()
		return e, nil
	case discord.GatewayOpCodeHeartbeatAck:
		s.Lock()
		s.LastHeartbeatAck = time.Now().UTC()
		s.Unlock()
		s.logger.Debug("got heartbeat ACK")
		return e, nil
	}

	// Do not try to Dispatch a non-Dispatch Message
	if e.Operation != discord.GatewayOpCodeDispatch {
		// But we probably should be doing something with them.
		// TEMP
		s.logger.Warn(
			"unknown opcode",
			"op", e.Operation,
			"seq", e.Sequence,
			"type", e.Type,
			"raw", string(e.RawData),
			"message", string(message),
		)
		return e, nil
	}

	s.sequence.Store(e.Sequence)

	if eh, ok := event.GetInterfaceProvider(e.Type); ok {
		e.Struct = eh.New()

		if err = json.Unmarshal(e.RawData, e.Struct); err != nil {
			s.logger.Warn("failed to unmarshal event", "type", e.Type, "raw", e.RawData)
			// READY events are always emitted
			if e.Type != event.ReadyType {
				return nil, err
			}
		}

		s.eventManager.EmitEvent(s, e.Type, e.Struct)
	} else {
		s.logger.Warn(
			"unknown event",
			"op", e.Operation,
			"seq", e.Sequence,
			"type", e.Type,
			"raw", e.RawData,
			"message", string(message),
		)
		s.eventManager.EmitEvent(s, event.EventType, e)
	}

	return e, nil
}

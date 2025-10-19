package gokord

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/coder/websocket"
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
func (s *Session) onInterface(ctx context.Context, i any) {
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
		go s.voiceAPI.UpdateServer(ctx, t.Token, t.GuildID, t.Endpoint)
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
	//s.logger.Debug("bot ready", "session_id", s.sessionID, "resume_url", r.ResumeGatewayURL)
}

// getGatewayEvent returns the discord.Event associated with the message given.
func getGatewayEvent(messageType websocket.MessageType, message []byte) (*discord.Event, error) {
	var err error
	var reader io.Reader
	reader = bytes.NewBuffer(message)

	if messageType == websocket.MessageBinary {
		// If this is a compressed message, uncompress it.
		z, err := zlib.NewReader(reader)
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(z)

		err = z.Close()
		if err != nil {
			return nil, err
		}
		reader = bytes.NewBuffer(b)
	}

	var e *discord.Event
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&e); err != nil {
		return e, err
	}
	return e, nil
}

// onGatewayEvent is the "event handler" for all messages received on the Discord Gateway API websocket connection.
func (s *Session) onGatewayEvent(ctx context.Context, e *discord.Event) error {
	// handle special opcode
	switch e.Operation {
	case discord.GatewayOpCodeHeartbeat: // must respond with a heartbeat packet within 5 seconds
		s.logger.Debug("sending heartbeat in response to Op1")
		return s.heartbeat(ctx)
	case discord.GatewayOpCodeReconnect: // must immediately disconnect from gateway and reconnect to new gateway
		s.logger.Info("closing and reconnecting in response to Op7")
		err := s.ForceClose() // was already closed
		if err != nil {
			// if we can't close, we must crash the app
			panic(err)
		}
		s.forceReconnect(ctx)
		return nil
	case discord.GatewayOpCodeInvalidSession:
		s.logger.Warn("invalid session received, reconnecting")
		err := s.CloseWithCode(ctx, websocket.StatusServiceRestart)
		if err != nil {
			// if we can't close, we must crash the app
			panic(err)
		}

		var resumable bool
		if err = json.Unmarshal(e.RawData, &resumable); err != nil {
			return err
		}

		if resumable {
			s.forceReconnect(ctx)
			return nil
		}

		s.logger.Info("gateway session is not resumable, discarding its information")
		s.resumeGatewayURL = ""
		s.sessionID = ""
		s.sequence.Store(0)
		if err = s.Open(ctx); err != nil {
			panic(err)
		}
		return nil
	case discord.GatewayOpCodeHeartbeatAck:
		s.Lock()
		s.LastHeartbeatAck = time.Now().UTC()
		s.Unlock()
		s.logger.Debug("got heartbeat ACK", "ping", s.HeartbeatLatency())
		return nil
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
		)
		return nil
	}

	s.sequence.Store(e.Sequence)

	var typ string
	var d any
	if eh, ok := event.GetInterfaceProvider(e.Type); ok {
		e.Struct = eh.New()

		if err := json.Unmarshal(e.RawData, e.Struct); err != nil {
			s.logger.Warn("failed to unmarshal event", "type", e.Type, "raw", e.RawData)
			// READY events are always emitted
			if e.Type != event.ReadyType {
				return err
			}
		}
		typ = e.Type
		d = e.Struct
	} else {
		s.logger.Warn(
			"unknown event",
			"op", e.Operation,
			"seq", e.Sequence,
			"type", e.Type,
			"raw", e.RawData,
		)
		typ = event.EventType
		d = e
	}
	s.eventManager.EmitEvent(ctx, s, typ, d)
	return nil
}

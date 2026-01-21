package gokord

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/coder/websocket"
	"github.com/nyttikord/gokord/bot"
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
	}
	err := s.sessionState.onInterface(s, i)
	if err != nil {
		s.logger.Error("dispatching internal event", "error", err)
	}
}

// onReady handles the ready event.
func (s *Session) onReady(r *event.Ready) {
	// Store the SessionID within the Session struct.
	s.sessionID = r.SessionID

	// Store the ResumeGatewayURL within the Session struct.
	s.resumeGatewayURL = r.ResumeGatewayURL
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

type eventHandlingResult struct {
	restart        bool
	force          bool
	openNewSession bool
}

// onGatewayEvent is the "event handler" for all messages received on the Discord Gateway API websocket connection.
//
// Return nil if everything is good, false if the ws must be restarted and true if it must be force restarted.
func (s *Session) onGatewayEvent(ctx context.Context, e *discord.Event) (*eventHandlingResult, error) {
	// handle special opcode
	switch e.Operation {
	case discord.GatewayOpCodeHeartbeat: // must respond with a heartbeat packet within 5 seconds
		s.logger.Debug("sending heartbeat in response to Op1")
		return nil, s.heartbeat(ctx)
	case discord.GatewayOpCodeReconnect: // must immediately disconnect from gateway and reconnect to new gateway
		s.logger.Info("reconnecting in response to Op7")
		return &eventHandlingResult{restart: true}, nil
	case discord.GatewayOpCodeInvalidSession:
		s.logger.Warn("invalid session received, reconnecting")
		var resumable bool
		if err := json.Unmarshal(e.RawData, &resumable); err != nil {
			return nil, err
		}

		if resumable {
			return &eventHandlingResult{restart: true}, nil
		}

		err := s.CloseWithCode(ctx, websocket.StatusServiceRestart)
		if err != nil {
			// if we can't close, we must crash the app
			panic(err)
		}

		s.logger.Info("gateway session is not resumable, discarding its information")
		s.resumeGatewayURL = ""
		s.sessionID = ""
		s.sequence.Store(0)
		return &eventHandlingResult{openNewSession: true}, nil
	case discord.GatewayOpCodeHeartbeatAck:
		s.lastHeartbeatAck.Store(time.Now().UnixMilli())
		s.logger.Debug("got heartbeat ACK", "ping", s.HeartbeatLatency())
		return nil, nil
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
		return nil, nil
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
				return nil, err
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
			"raw", string(e.RawData),
		)
		typ = event.EventType
		d = e
	}
	if e.Type != event.IntegrationCreateType {
		ctx = bot.SetLogger(ctx, bot.Logger(ctx).With("event", e.Type))
	} else {
		ctx = s.interactionManager.Context(ctx)
	}
	s.eventManager.EmitEvent(ctx, s, typ, d)
	return nil, nil
}

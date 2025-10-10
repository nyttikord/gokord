package gokord

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coder/websocket"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
)

var (
	ErrShouldNotReconnect   = errors.New("session should not reconnect")
	ErrSendingResumePacket  = errors.New("cannot send resume packet")
	ErrHandlingMissedEvents = errors.New("cannot handle missed events")
	ErrInvalidSession       = errors.New("invalid session")
)

type resumePacket struct {
	Op   discord.GatewayOpCode `json:"op"`
	Data struct {
		Token     string `json:"token"`
		SessionID string `json:"session_id"`
		Sequence  int64  `json:"seq"`
	} `json:"d"`
}

func (s *Session) reconnect(ctx context.Context) error {
	if !s.ShouldReconnectOnError {
		return ErrShouldNotReconnect
	}
	s.logger.Info("trying to reconnect to gateway")

	err := s.setupGateway(ctx, s.resumeGatewayURL)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	var p resumePacket
	p.Op = discord.GatewayOpCodeResume
	p.Data.Token = s.Identify.Token
	p.Data.SessionID = s.sessionID
	p.Data.Sequence = s.sequence.Load()

	s.logger.Info("sending resume packet to gateway")
	err = s.GatewayWriteStruct(ctx, p)
	if err != nil {
		return errors.Join(err, ErrSendingResumePacket)
	}
	defer func() {
		if err != nil {
			s.logger.Warn("force closing after error")
			err = s.ForceClose()
			if err != nil {
				// if we can't close, we must crash the app
				panic(err)
			}
		}
	}()
	// handle missed event
	e := new(discord.Event)
	e.Type = ""
	for e.Type != event.ResumedType {
		mt, m, err := s.ws.Read(ctx)
		if err != nil {
			return errors.Join(err, ErrHandlingMissedEvents)
		}
		e, err = getGatewayEvent(mt, m)
		if err != nil {
			return errors.Join(err, ErrHandlingMissedEvents)
		}
		switch e.Operation {
		case discord.GatewayOpCodeHello:
			err = s.handleHello(e)
		case discord.GatewayOpCodeInvalidSession:
			return ErrInvalidSession
		default:
			s.Unlock() // required
			err = s.onGatewayEvent(ctx, e)
			s.Lock()
		}
		if err != nil {
			return errors.Join(err, ErrHandlingMissedEvents)
		}
	}
	s.logger.Info("successfully reconnected to gateway")

	s.finishConnection(ctx)

	// I'm not sure if this is actually needed.
	// If the gw reconnect works properly, voice should stay alive.
	// However, there seems to be cases where something "weird" happens.
	// So we're doing this for now just to improve stability in those edge cases.
	if !s.ShouldReconnectVoiceOnSessionError {
		return nil
	}
	for _, v := range s.voiceAPI.Connections {
		s.logger.Info("reconnecting voice connection to guild", "guild", v.GuildID)
		go v.Reconnect(ctx)

		// This is here just to prevent violently spamming the voice reconnects.
		time.Sleep(1 * time.Second)
	}
	return nil
}

// forceReconnect the session.
// If the reconnection fails, it opens a new session.
// If it cannot create a new session, it panics.
func (s *Session) forceReconnect(ctx context.Context) {
	err := s.reconnect(ctx)
	if err == nil {
		return
	}
	// if the reconnects fail, we close the websocket
	err = s.ForceClose()
	if err != nil {
		// if we can't close, we must crash the app
		panic(err)
	}
	s.logger.Error("reconnecting", "error", err, "gateway", s.gateway)
	s.Logger().Warn("opening a new session")
	err = s.Open(ctx)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("failed to force reconnect"))
		// panic because we can't reconnect
		panic(err)
	}
}

// Close closes a websocket and stops all listening/heartbeat goroutines.
// If it returns an error, the session is not closed.
func (s *Session) Close(ctx context.Context) error {
	return s.CloseWithCode(ctx, websocket.StatusNormalClosure)
}

// CloseWithCode closes a websocket using the provided closeCode and stops all listening/heartbeat goroutines.
// If it returns an error, the session is not closed.
// TODO: Add support for Voice WS/UDP connections
func (s *Session) CloseWithCode(ctx context.Context, closeCode websocket.StatusCode) error {
	s.Lock()
	defer s.Unlock()
	s.logger.Info("closing", "code", closeCode)

	if s.ws == nil {
		return ErrWSNotFound
	}

	if s.listening != nil {
		s.logger.Debug("closing goroutines")
		s.listening.Store(false)
	}

	for _, v := range s.voiceAPI.Connections {
		err := v.Disconnect(ctx)
		if err != nil {
			s.logger.Error("disconnecting voice from channel", "error", err, "channel", v.ChannelID)
		}
	}
	// TODO: stop any reconnecting voice channels

	s.logger.Info("closing websocket")
	// is a clean stop
	s.wsMutex.Lock()
	err := s.ws.Close(closeCode, "")
	s.wsMutex.Unlock()
	if err != nil {
		s.logger.Error("closing websocket", "error", err, "gateway", s.gateway)
		s.ws = nil
		s.Unlock()
		err = s.ForceClose()
		s.Lock()
		return err
	}

	// required
	s.Unlock()
	s.ws = nil
	s.eventManager.EmitEvent(ctx, s, event.DisconnectType, &event.Disconnect{})
	s.Lock()

	return nil
}

// ForceClose the connection.
// Use Close or CloseWithCode before to have a better closing process.
//
// It doesn't send an event.Disconnect, unlike Close or CloseWithCode.
func (s *Session) ForceClose() error {
	s.Lock()
	defer s.Unlock()
	s.logger.Warn("force closing websocket")
	err := s.ws.CloseNow()
	if err != nil {
		// we handle it here because the websocket is actually closed
		return err
	}
	s.ws = nil
	return nil
}

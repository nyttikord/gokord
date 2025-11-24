package gokord

import (
	"context"
	"errors"
	"fmt"
	"net"
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

func (s *Session) reconnect(ctx context.Context, forceClose bool) error {
	if s.restarting.Load() {
		return nil
	}
	if !s.Options.ShouldReconnectOnError {
		return ErrShouldNotReconnect
	}

	s.restarting.Store(true)
	defer s.restarting.Store(false)

	var err error
	if !forceClose {
		err = s.CloseWithCode(ctx, websocket.StatusServiceRestart)
	}
	if forceClose || err != nil {
		if !forceClose && !errors.Is(err, net.ErrClosed) {
			s.logger.Warn("error while closing", "error", err)
		}
		if err = s.ForceClose(); err != nil {
			return err
		}
	}

	s.Lock()
	defer s.Unlock()

	err = s.setupGateway(ctx, s.resumeGatewayURL)
	if err != nil {
		return err
	}

	var p resumePacket
	p.Op = discord.GatewayOpCodeResume
	p.Data.Token = s.Identify.Token
	p.Data.SessionID = s.sessionID
	p.Data.Sequence = s.sequence.Load()

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

	s.setupListen(ctx)

	// handle missed event
	e := new(discord.Event)
	e.Type = ""
	for e.Type != event.ResumedType {
		res := <-s.wsRead
		var err error
		e, err = res.getEvent()
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
	if !s.Options.ShouldReconnectVoiceOnSessionError {
		return nil
	}
	for _, v := range s.voiceAPI.Connections {
		s.logger.Debug("reconnecting voice connection to guild", "guild", v.GuildID)
		go v.Reconnect(ctx)

		// This is here just to prevent violently spamming the voice reconnects.
		time.Sleep(1 * time.Second)
	}
	return nil
}

// forceReconnect the session.
// If the reconnection fails, it opens a new session.
// If it cannot create a new session, it panics.
func (s *Session) forceReconnect(ctx context.Context, forceClose bool) {
	if s.restarting.Load() {
		return
	}
	err := s.reconnect(ctx, forceClose)
	if err == nil {
		return
	}
	if errors.Is(err, ErrShouldNotReconnect) {
		panic(err)
	}
	// if the reconnects fail, we close the websocket
	s.logger.Warn("reconnecting to gateway", "error", err)
	s.logger.Warn("force closing websocket")
	err = s.ForceClose()
	if err != nil && !errors.Is(err, net.ErrClosed) {
		// if we can't close, we must crash the app
		panic(err)
	}
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

	s.waitListen.Close(ctx)
	s.cancelWSRead()

	for _, v := range s.voiceAPI.Connections {
		err := v.Disconnect(ctx)
		if err != nil {
			s.logger.Error("disconnecting voice from channel", "error", err, "channel", v.ChannelID)
		}
	}
	// TODO: stop any reconnecting voice channels

	s.logger.Debug("closing websocket")
	// is a clean stop
	s.wsMutex.Lock()
	err := s.ws.Close(closeCode, "")
	s.wsMutex.Unlock()
	if err != nil {
		s.logger.Warn("closing websocket", "error", err, "gateway", s.gateway)
		s.Unlock()
		err = s.ForceClose()
		s.Lock()
	}
	// required
	s.ws = nil
	if err != nil {
		return err
	}
	if err := s.waitListen.Wait(ctx); err != nil {
		return err
	}

	s.Unlock()
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
	s.logger.Info("force closing")
	var err error
	s.waitListen.Close(context.Background())
	s.cancelWSRead()
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
		}
	}()
	err = s.ws.CloseNow()
	// avoid returning an error is the websocket is closed, because this method must close the websocket and if this is
	// already closed, there is no error
	if err != nil && !errors.Is(err, net.ErrClosed) {
		return err
	}
	if err = s.waitListen.Wait(context.Background()); err != nil {
		return err
	}
	s.ws = nil
	return nil
}

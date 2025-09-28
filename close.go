package gokord

import (
	"context"
	"errors"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nyttikord/gokord/event"
)

func (s *Session) reconnect() {
	if !s.ShouldReconnectOnError {
		return
	}
	s.logger.Info("trying to reconnect to gateway")

	wait := time.Duration(1)
	err := s.Open()

	for err != nil {
		// Certain race conditions can call reconnect() twice.
		// If this happens, we just break out of the reconnect loop
		// TODO: fix this
		if errors.Is(err, ErrWSAlreadyOpen) {
			s.logger.Debug("Websocket already exists, no need to reconnect")
			return
		}

		s.logger.Error("reconnecting to gateway", "error", err)

		time.Sleep(wait * time.Second)
		wait *= min(wait*2, 600)

		s.logger.Info("trying to reconnect to gateway")

		err = s.Open()
	}
	s.logger.Info("successfully reconnected to gateway")

	// I'm not sure if this is actually needed.
	// If the gw reconnect works properly, voice should stay alive.
	// However, there seems to be cases where something "weird" happens.
	// So we're doing this for now just to improve stability in those edge cases.
	if !s.ShouldReconnectVoiceOnSessionError {
		return
	}
	s.RLock()
	defer s.RUnlock()
	for _, v := range s.voiceAPI.Connections {
		s.logger.Info("reconnecting voice connection to guild", "guild", v.GuildID)
		go v.Reconnect()

		// This is here just to prevent violently spamming the voice reconnects.
		time.Sleep(1 * time.Second)
	}
}

// Close closes a websocket and stops all listening/heartbeat goroutines.
// If it returns an error, the session is not closed.
func (s *Session) Close() error {
	return s.CloseWithCode(websocket.CloseNormalClosure)
}

// CloseWithCode closes a websocket using the provided closeCode and stops all listening/heartbeat goroutines.
// If it returns an error, the session is not closed.
// TODO: Add support for Voice WS/UDP connections
func (s *Session) CloseWithCode(closeCode int) error {
	s.logger.Info("closing", "code", closeCode)
	s.Lock()
	defer s.Unlock()

	if s.ws == nil {
		return ErrWSNotFound
	}

	s.DataReady = false

	if s.listening != nil {
		s.logger.Debug("closing listening channel")
		close(s.listening)
		s.listening = nil
	}

	for _, v := range s.voiceAPI.Connections {
		err := v.Disconnect()
		if err != nil {
			s.logger.Error("disconnecting voice from channel", "error", err, "channel", v.ChannelID)
		}
	}
	// TODO: Close all active Voice Connections force stop any reconnecting voice channels

	// To cleanly close a connection, a client should send a close frame and wait for the server to close the
	// connection.
	s.logger.Debug("sending close frame")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		s.wsMutex.Lock()
		err := s.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(closeCode, ""))
		s.wsMutex.Unlock()
		errChan <- err
		// TODO: waiting for Discord to close the websocket
		// I have searched a way to wait for the ws to be closed, but I have found nothing on it.
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
	s.ForceClose()
	s.eventManager.EmitEvent(s, event.DisconnectType, &event.Disconnect{})
	s.Lock()

	return nil
}

// ForceClose the connection.
// Use Close or CloseWithCode before to have a better closing process.
//
// It doesn't send an event.Disconnect, unlike Close or CloseWithCode.
func (s *Session) ForceClose() {
	s.Lock()
	defer s.Unlock()
	s.logger.Info("closing gateway websocket")
	err := s.ws.Close()
	if err != nil {
		// we handle it here because the websocket is actually closed
		s.logger.Error("closing websocket", "error", err)
	}
	s.ws = nil
}

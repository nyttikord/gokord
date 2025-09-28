package gokord

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
)

type resumePacket struct {
	Op   discord.GatewayOpCode `json:"op"`
	Data struct {
		Token     string `json:"token"`
		SessionID string `json:"session_id"`
		Sequence  int64  `json:"seq"`
	} `json:"d"`
}

// Open creates a websocket connection to Discord.
// https://discord.com/developers/docs/topics/gateway#connecting
func (s *Session) Open() error {
	s.Lock()
	defer s.Unlock()

	// If the websock is already open, bail out here.
	if s.ws != nil {
		return ErrWSAlreadyOpen
	}

	sequence := s.sequence.Load()

	var gateway string
	var err error
	// Get the gateway to use for the Websocket connection
	if sequence != 0 && s.sessionID != "" && s.resumeGatewayURL != "" {
		s.logger.Debug("using resume gateway", "gateway", s.resumeGatewayURL)
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
	s.logger.Info("connecting to gateway", "gateway", gateway)
	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	s.ws, _, err = s.Dialer.Dial(gateway, header)
	if err != nil {
		s.logger.Error("connecting to gateway", "error", err, "gateway", s.gateway)
		s.gateway = "" // clear cached gateway
		s.ws = nil     // Just to be safe.
		return err
	}

	s.ws.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	defer func() {
		if err != nil {
			s.ForceClose()
		}
	}()

	// The first response from Discord should be an Op 10 (Hello) Packet.
	// When processed by onGatewayEvent the heartbeat goroutine will be started.
	mt, m, err := s.ws.ReadMessage()
	if err != nil {
		return err
	}
	s.Unlock()
	e, err := s.onGatewayEvent(mt, m)
	s.Lock()
	if err != nil {
		return err
	}
	if e.Operation != discord.GatewayOpCodeHello {
		return fmt.Errorf("expecting Op 10, got Op %d instead", e.Operation)
	}
	s.logger.Debug("Op 10 Hello Packet received from Discord")
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
		var p resumePacket
		p.Op = discord.GatewayOpCodeResume
		p.Data.Token = s.Identify.Token
		p.Data.SessionID = s.sessionID
		p.Data.Sequence = sequence

		s.logger.Info("sending resume packet to gateway")
		s.wsMutex.Lock()
		err = s.ws.WriteJSON(p)
		s.wsMutex.Unlock()
		if err != nil {
			err = fmt.Errorf("error sending gateway resume packet, %s, %s", s.gateway, err)
			return err
		}
	}

	// Now Discord should send us a READY or RESUMED packet.
	mt, m, err = s.ws.ReadMessage()
	if err != nil {
		return err
	}
	s.Unlock()
	e, err = s.onGatewayEvent(mt, m)
	s.Lock()
	if err != nil {
		return err
	}
	if e.Type != `READY` && e.Type != `RESUMED` {
		// This is not fatal, but it does not follow their API documentation.
		s.logger.Warn("expected READY/RESUMED", "got", e)
	}

	s.logger.Debug("We are now connected to Discord, emitting connect event")
	s.eventManager.EmitEvent(s, event.ConnectType, &event.Connect{})

	// Create listening chan outside of listen, as it needs to happen inside the mutex lock and needs to exist before
	// calling heartbeat and listen goroutines.
	s.listening = make(chan any)

	// Start sending heartbeats and reading messages from Discord.
	go s.heartbeats(s.ws, s.listening, h.HeartbeatInterval)
	go s.listen(s.ws, s.listening)

	return nil
}

// listen polls the websocket connection for events, it will stop when the listening channel is closed, or an error
// occurs.
func (s *Session) listen(wsConn *websocket.Conn, listening <-chan any) {
	messageType, message, err := wsConn.ReadMessage()
	for err == nil {
		select {
		case <-listening:
			return
		default:
			_, err = s.onGatewayEvent(messageType, message)
			if err != nil {
				s.logger.Error("handling event", "error", err)
			}
		}

		messageType, message, err = wsConn.ReadMessage()
	}

	// Detect if we have been closed manually.
	// If a Close() has already happened, the websocket we are listening on will be different to the current
	// session.
	s.RLock()
	sameConnection := s.ws == wsConn
	s.RUnlock()

	// everything is fine
	if !sameConnection {
		return
	}
	s.logger.Error("reading from websocket", "error", err, "gateway", s.gateway)
	err = s.Close()
	if err != nil {
		s.logger.Error("closing session connection, force closing", "error", err)
		s.ForceClose()
	}

	s.logger.Info("calling reconnect() now")
	s.reconnect()
}

type heartbeatOp struct {
	Op   discord.GatewayOpCode `json:"op"`
	Data int64                 `json:"d"`
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
func (s *Session) heartbeats(ws *websocket.Conn, listening <-chan any, heartbeatIntervalMsec time.Duration) {
	if listening == nil || ws == nil {
		return
	}

	var err error
	ticker := time.NewTicker(heartbeatIntervalMsec * time.Millisecond)
	defer ticker.Stop()

	last := time.Now().UTC()

	for err == nil && time.Now().UTC().Sub(last) <= (heartbeatIntervalMsec*FailedHeartbeatAcks) {
		s.RLock()
		last = s.LastHeartbeatAck
		s.RUnlock()

		sequence := s.sequence.Load()
		err = s.heartbeat(ws, sequence)
		s.sequence.Add(1)

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

	if err != nil {
		s.logger.Error("sending heartbeat", "error", err, "gateway", s.gateway)
	} else {
		s.logger.Warn(
			"haven't gotten a heartbeat ACK, triggering a reconnection",
			"time since last ACK", time.Now().UTC().Sub(last),
		)
	}
	err = s.Close()
	if err != nil {
		s.logger.Error("closing session connection, force closing", "error", err)
		s.ForceClose()
	}
	s.reconnect()
}

func (s *Session) heartbeat(ws *websocket.Conn, sequence int64) error {
	s.logger.Debug("sending gateway websocket heartbeat", "sequence", sequence)
	s.wsMutex.Lock()
	s.LastHeartbeatSent = time.Now().UTC()
	err := ws.WriteJSON(heartbeatOp{discord.GatewayOpCodeHeartbeat, sequence})
	s.wsMutex.Unlock()
	return err
}

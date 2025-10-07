package gokord

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync/atomic"
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

	// init new sequence
	s.sequence = &atomic.Int64{}
	s.sequence.Store(0)

	var gateway string
	var err error
	if s.gateway == "" {
		s.gateway, err = s.Gateway()
		if err != nil {
			return err
		}
	}

	gateway = s.gateway

	// Add the version and encoding to the URL
	gateway += "?v=" + discord.APIVersion + "&encoding=json"
	return s.connect(gateway)
}

func (s *Session) setupGateway(gateway string) error {
	// Connect to the Gateway
	s.logger.Info("connecting to gateway", "gateway", gateway)
	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	var err error
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

	return nil
}

func (s *Session) connect(gateway string) error {
	var err error
	if err = s.setupGateway(gateway); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			s.ForceClose()
		}
	}()

	// The first response from Discord should be an Op 10 (Hello) Packet.
	mt, m, err := s.ws.ReadMessage()
	if err != nil {
		return err
	}
	e, err := getGatewayEvent(mt, m)
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
	h.HeartbeatInterval *= time.Millisecond

	// Send Op 2 Identity Packet
	err = s.identify()
	if err != nil {
		err = fmt.Errorf("error sending identify packet to gateway, %s, %s", s.gateway, err)
		return err
	}

	// Now Discord should send us a READY or RESUMED packet.
	mt, m, err = s.ws.ReadMessage()
	if err != nil {
		return err
	}
	e, err = getGatewayEvent(mt, m)
	if err != nil {
		return err
	}
	if e.Type != `READY` {
		return fmt.Errorf("expected READY, got %v", e)
	}

	s.logger.Debug("We are now connected to Discord, emitting connect event")
	s.eventManager.EmitEvent(s, event.ConnectType, &event.Connect{})

	// Create listening chan outside of listen, as it needs to happen inside the mutex lock and needs to exist before
	// calling heartbeat and listen goroutines.
	s.listening = make(chan any)

	// Start sending heartbeats and reading messages from Discord.
	go func() {
		time.Sleep(time.Duration(rand.Float32() * float32(h.HeartbeatInterval)))
		s.heartbeats(s.ws, s.listening, h.HeartbeatInterval)
	}()
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
			var e *discord.Event
			e, err = getGatewayEvent(messageType, message)
			if err != nil {
				s.logger.Error("handling event", "error", err, "when", "getting event")
			} else {
				err = s.onGatewayEvent(e)
				if err != nil {
					s.logger.Error("handling event", "error", err, "when", "handling event")
				}
			}
		}
		if err == nil {
			messageType, message, err = wsConn.ReadMessage()
		}
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
	s.forceReconnect()
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
func (s *Session) heartbeats(ws *websocket.Conn, listening <-chan any, heartbeatInterval time.Duration) {
	if listening == nil || ws == nil {
		return
	}

	var err error
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	last := time.Now().UTC()

	for err == nil && time.Now().UTC().Sub(last) <= (heartbeatInterval/time.Millisecond*FailedHeartbeatAcks) {
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
			s.logger.Debug("exiting heartbeats")
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
	s.forceReconnect()
}

func (s *Session) heartbeat(ws *websocket.Conn, sequence int64) error {
	s.logger.Debug("sending gateway websocket heartbeat", "sequence", sequence)
	s.wsMutex.Lock()
	s.LastHeartbeatSent = time.Now().UTC()
	err := ws.WriteJSON(heartbeatOp{discord.GatewayOpCodeHeartbeat, sequence})
	s.wsMutex.Unlock()
	return err
}

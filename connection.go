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

	var err error
	if s.gateway == "" {
		s.gateway, err = s.Gateway()
		if err != nil {
			return err
		}
	}
	return s.connect()
}

func (s *Session) setupGateway(gateway string) error {
	s.Lock()
	defer s.Unlock()
	// Add the version and encoding to the URL
	gateway += "?v=" + discord.APIVersion + "&encoding=json"

	// Connect to the Gateway
	s.logger.Info("connecting to gateway", "gateway", gateway)
	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	var err error
	s.ws, _, err = s.Dialer.Dial(gateway, header)
	if err != nil {
		s.gateway = "" // clear cached gateway
		s.ws = nil     // Just to be safe.
		return err
	}

	s.ws.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	return nil
}

// connect must be called when Session's mutex is locked.
func (s *Session) connect() error {
	s.Unlock() // required
	err := s.setupGateway(s.gateway)
	s.Lock()
	if err != nil {
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

	// Now Discord should send us a READY packet.
	mt, m, err = s.ws.ReadMessage()
	if err != nil {
		return err
	}
	e, err = getGatewayEvent(mt, m)
	if err != nil {
		return err
	}
	if e.Type != event.ReadyType {
		return fmt.Errorf("expected %s, got %v", event.ReadyType, e)
	}

	s.logger.Debug("We are now connected to Discord, emitting connect event")
	s.eventManager.EmitEvent(s, event.ConnectType, &event.Connect{})

	// Create listening chan outside of listen, as it needs to happen inside the mutex lock and needs to exist before
	// calling heartbeat and listen goroutines.
	s.listening = make(chan any, 1)

	// Start sending heartbeats and reading messages from Discord.
	go func() {
		time.Sleep(time.Duration(rand.Float32() * float32(h.HeartbeatInterval)))
		s.heartbeats(s.listening, h.HeartbeatInterval)
	}()
	go s.listen(s.ws, s.listening)

	return nil
}

// listen polls the websocket connection for events, it will stop when the listening channel is closed, or when an error
// occurs.
func (s *Session) listen(ws *websocket.Conn, listening <-chan any) {
	messageType, message, err := ws.ReadMessage()
	for err == nil {
		select {
		case <-listening:
			s.logger.Debug("exiting listen websocket event")
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
			messageType, message, err = ws.ReadMessage()
		}
	}

	// Detect if we have been closed manually.
	// If a Close() has already happened, the websocket we are listening on will be different to the current
	// session.
	// TODO: clean this
	s.RLock()
	sameConnection := s.ws == ws
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
const FailedHeartbeatAcks = 5

// HeartbeatLatency returns the latency between heartbeat acknowledgement and heartbeat send.
func (s *Session) HeartbeatLatency() time.Duration {
	return s.LastHeartbeatAck.Sub(s.LastHeartbeatSent)
}

// heartbeat sends regular heartbeats to Discord so it knows the client is still connected.
// If you do not send these heartbeats Discord will disconnect the websocket connection after a few seconds.
func (s *Session) heartbeats(listening <-chan any, heartbeatInterval time.Duration) {
	s.logger.Debug("starting heartbeats")
	var err error
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	last := time.Now().UTC()

	for err == nil && time.Now().UTC().Sub(last) <= (heartbeatInterval*FailedHeartbeatAcks) {
		s.RLock()
		last = s.LastHeartbeatAck
		s.RUnlock()

		err = s.heartbeat()

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

func (s *Session) heartbeat() error {
	s.Lock()
	defer s.Unlock()
	seq := s.sequence.Load()
	s.LastHeartbeatSent = time.Now().UTC()
	s.logger.Debug("sending websocket heartbeat", "sequence", seq)
	return s.GatewayWriteStruct(heartbeatOp{discord.GatewayOpCodeHeartbeat, seq})
}

package gokord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
)

var (
	ErrReadingReadyPacket = errors.New("cannot read READY packet")
	ErrIdentifying        = errors.New("cannot identify")
)

// Open creates a websocket connection to Discord.
// https://discord.com/developers/docs/topics/gateway#connecting
func (s *Session) Open(ctx context.Context) error {
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

	err = s.setupGateway(ctx, s.gateway)
	if err != nil {
		return err
	}

	err = s.connect(ctx)
	if err != nil {
		return err
	}

	s.finishConnection(ctx)

	return nil
}

func (s *Session) setupGateway(ctx context.Context, gateway string) error {
	// Add the version and encoding to the URL
	gateway = strings.TrimSuffix(gateway, "/")
	gateway += "/?v=" + discord.APIVersion + "&encoding=json"

	// Connect to the Gateway
	s.logger.Debug("connecting to gateway", "gateway", gateway)
	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	var err error
	s.ws, _, err = websocket.Dial(ctx, gateway, &websocket.DialOptions{HTTPHeader: header})
	if err != nil {
		s.gateway = "" // clear cached gateway
		s.ws = nil     // Just to be safe.
		return err
	}

	return nil
}

type heartbeatOp struct {
	Op   discord.GatewayOpCode `json:"op"`
	Data int64                 `json:"d"`
}

func (s *Session) handleHello(e *discord.Event) error {
	s.logger.Debug("Op 10 (Hello) received")
	s.LastHeartbeatAck = time.Now().UTC()
	var h helloOp
	if err := json.Unmarshal(e.RawData, &h); err != nil {
		return errors.Join(err, fmt.Errorf("cannot unmarshal HelloOp"))
	}
	s.heartbeatInterval = h.HeartbeatInterval * time.Millisecond
	return nil
}

func (s *Session) setupListen(ctx context.Context) {
	if s.wsRead != nil {
		s.logger.Info("listen already running")
		return
	}
	ctx2, cancel := context.WithCancel(ctx)
	s.cancelWSRead = cancel

	wsRead := make(chan readResult)
	s.wsRead = wsRead
	go func() {
		s.logger.Info("listening started")
		err := s.listen(ctx2, wsRead)
		select {
		case <-ctx2.Done():
			s.wsRead = nil
			return
		default:
			if errors.Is(err, io.EOF) {
				s.logger.Warn("EOF received, restarting")
				s.wsRead = nil
				s.forceReconnect(ctx)
				return
			}
			s.logger.Error("listening websocket")
			panic(err)
		}
	}()
}

// connect must be called when Session's mutex is locked.
func (s *Session) connect(ctx context.Context) error {
	var err error
	defer func() {
		if err != nil {
			if err = s.ForceClose(); err != nil {
				// if we can't close, we must crash the app
				panic(err)
			}
		}
	}()

	s.setupListen(ctx)

	// if we can't connect in 10s, returns an error
	ctx2, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	errc := make(chan error, 1)
	go func() {
		// The first response from Discord should be an Op 10 (Hello) Packet.
		res := <-s.wsRead
		e, err := res.getEvent()
		if err != nil {
			errc <- err
		}
		if e.Operation != discord.GatewayOpCodeHello {
			errc <- fmt.Errorf("expecting Op 10, got Op %d instead", e.Operation)
		}
		if err = s.handleHello(e); err != nil {
			errc <- err
		}

		// Send Op 2 Identity Packet
		err = s.identify(ctx)
		if err != nil {
			errc <- errors.Join(err, ErrIdentifying)
		}

		// Now Discord should send us a READY packet.
		res = <-s.wsRead
		e, err = res.getEvent()
		if err != nil {
			errc <- errors.Join(err, ErrReadingReadyPacket)
		}
		if e.Type != event.ReadyType {
			s.logger.Error("invalid READY packet", "type got", e.Type)
			errc <- ErrReadingReadyPacket
		}
		s.Unlock() // required to dispatch ready
		err = s.onGatewayEvent(ctx, e)
		s.Lock()
		errc <- err
	}()

	select {
	case <-ctx2.Done():
		err = ctx2.Err()
	case err = <-errc:
	}
	return err
}

// TODO: rename this method
func (s *Session) finishConnection(ctx context.Context) {
	s.logger.Debug("connected to Discord, emitting connect event")
	s.eventManager.EmitEvent(ctx, s, event.ConnectType, &event.Connect{})

	var ctx2 context.Context
	ctx2, s.waitListen.cancel = context.WithCancel(ctx)

	restarting := false
	restart := func() {
		if restarting {
			return
		}
		restarting = true
		/*s.logger.Info("closing websocket")
		err := s.ForceClose() // force closing because the websocket is always unusable in this state according to our tests
		if err != nil {
			panic(err)
		}*/
		s.logger.Info("reconnecting")
		s.forceReconnect(ctx)
	}

	// Start sending heartbeats and reading messages from Discord.
	s.waitListen.Add(func(free func()) {
		s.Logger().Info("new heartbeats")
		last, err := s.heartbeats(ctx2)
		free()
		s.Logger().Info("heartbeats ended")
		select {
		case <-ctx2.Done():
			s.logger.Debug("exiting heartbeats")
			return
		default:
			s.logger.Warn("sending heartbeats", "error", err, "time since last ACK", time.Now().UTC().Sub(last))
			restart()
		}
	})
	s.waitListen.Add(func(free func()) {
		s.Logger().Info("dispatching events started")
		var err error
		for err == nil {
			select {
			case res := <-s.wsRead:
				err = res.dispatch(s, ctx2)
			case <-ctx2.Done():
				free()
				s.logger.Debug("exiting dispatching events")
				return
			}
		}
		free()
		s.logger.Warn("reading from websocket", "error", err, "gateway", s.gateway)
		restart()
	})
}

// listen polls the websocket connection for data, it will stop when an error occurs.
func (s *Session) listen(ctx context.Context, c chan<- readResult) error {
	var messageType websocket.MessageType
	var message []byte
	var err error
	for err == nil {
		messageType, message, err = s.ws.Read(ctx)
		if err == nil {
			c <- readResult{messageType, message}
		}
	}
	return err
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
func (s *Session) heartbeats(ctx context.Context) (time.Time, error) {
	select {
	case <-time.After(time.Duration(rand.Float32() * float32(s.heartbeatInterval))):
	case <-ctx.Done():
		return time.Now().UTC(), nil
	}
	s.logger.Debug("starting heartbeats")
	var err error
	ticker := time.NewTicker(s.heartbeatInterval)
	defer ticker.Stop()

	last := time.Now().UTC()
	// first heartbeat
	err = s.heartbeat(ctx)

	for err == nil && time.Now().UTC().Sub(last) <= s.heartbeatInterval*FailedHeartbeatAcks {
		select {
		case <-ticker.C:
			s.RLock()
			last = s.LastHeartbeatAck
			s.RUnlock()

			err = s.heartbeat(ctx)
		case <-ctx.Done():
			return last, nil
		}
	}
	if err == nil {
		err = errors.New("haven't gotten a heartbeat ACK in time")
	}
	return last, err
}

func (s *Session) heartbeat(ctx context.Context) error {
	seq := s.sequence.Load()
	s.LastHeartbeatSent = time.Now().UTC()
	s.logger.Debug("sending websocket heartbeat", "sequence", seq)
	return s.GatewayWriteStruct(ctx, heartbeatOp{discord.GatewayOpCodeHeartbeat, seq})
}

type readResult struct {
	MessageType websocket.MessageType
	Message     []byte
}

func (r *readResult) getEvent() (*discord.Event, error) {
	return getGatewayEvent(r.MessageType, r.Message)
}

// dispatch the event received
func (r *readResult) dispatch(s *Session, ctx context.Context) error {
	e, err := r.getEvent()
	if err == nil {
		err = s.onGatewayEvent(ctx, e)
	}
	return err
}

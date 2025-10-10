package voice

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/nyttikord/gokord/discord"
)

// wsListen listens on the voice websocket for messages and passes them to the voice event handler.
// This is automatically called by the open.
func (v *Connection) wsListen(ctx context.Context, wsConn *websocket.Conn, close <-chan struct{}) {
	for {
		_, message, err := v.wsConn.Read(ctx)
		if err != nil {
			// 4014 indicates a manual disconnection by someone in the guild;
			// we shouldn't Reconnect.
			var errClose websocket.CloseError
			if errors.As(err, &errClose) && errClose.Code == 4014 {
				v.Logger.Info("received 4014 manual disconnection")

				// Abandon the voice WS connection
				v.Lock()
				v.wsConn = nil
				v.Unlock()

				// Wait for VOICE_SERVER_UPDATE.
				// When the user moves the bot to another voice channel,
				// VOICE_SERVER_UPDATE is received after the code 4014.
				for i := 0; i < 5; i++ { // TODO: temp, wait for VoiceServerUpdate.
					time.Sleep(1 * time.Second)

					v.RLock()
					reconnected := v.wsConn != nil
					v.RUnlock()
					if !reconnected {
						continue
					}
					v.Logger.Info("successfully reconnected after 4014 manual disconnection")
					return
				}

				// When VOICE_SERVER_UPDATE is not received, disconnect as usual.
				v.Logger.Info("disconnect due to 4014 manual disconnection")

				v.requester.Lock()
				delete(v.requester.Connections, v.GuildID)
				v.requester.Unlock()

				v.Close()

				return
			}

			// Detect if we have been closed manually.
			// If a Close() has already happened, the websocket we are listening to on will be different to the current
			// requester.
			v.RLock()
			sameConnection := v.wsConn == wsConn
			v.RUnlock()
			if sameConnection {
				v.Logger.Error("voice websocket closed unexpectantly", "error", err, "endpoint", v.endpoint)
				go v.Reconnect(ctx)
			}
			return
		}

		select {
		case <-close:
			return
		default:
			go v.onEvent(ctx, message)
		}
	}
}

// onEvent handles any voice websocket events.
// This is only called by wsListen.
func (v *Connection) onEvent(ctx context.Context, message []byte) {
	v.Logger.Debug("received", "raw", message)

	var e struct {
		discord.Event
		Operation discord.VoiceOpCode `json:"op"`
	}
	if err := json.Unmarshal(message, &e); err != nil {
		v.Logger.Error("unmarshall event", "error", err)
		return
	}

	switch e.Operation {
	case discord.VoiceOpCodeReady: // READY
		if err := json.Unmarshal(e.RawData, &v.op2); err != nil {
			v.Logger.Error("OP2 unmarshall", "error", err, "raw", e.RawData)
			return
		}

		// Start the voice websocket heartbeat to keep the connection alive
		go v.wsHeartbeat(ctx, v.wsConn, v.close, v.op2.HeartbeatInterval)
		// TODO monitor a chan/bool to verify this was successful

		// Start the UDP connection
		err := v.udpOpen(ctx)
		if err != nil {
			v.Logger.Error("opening udp connection", "error", err)
			return
		}

		// Start the opusSender.
		// TODO: Should we allow 48000/960 values to be user defined?
		if v.OpusSend == nil {
			v.OpusSend = make(chan []byte, 2)
		}
		go v.opusSender(ctx, v.udpConn, v.close, v.OpusSend, 48000, 960)

		// Start the opusReceiver
		if !v.deaf {
			if v.OpusRecv == nil {
				v.OpusRecv = make(chan *Packet, 2)
			}

			go v.opusReceiver(ctx, v.udpConn, v.close, v.OpusRecv)
		}
		return

	case discord.VoiceOpCodeHeartbeat: // HEARTBEAT response
		// add code to use this to track latency?
		return

	case discord.VoiceOpCodeSessionDescription: // udp encryption secret key
		v.Lock()
		defer v.Unlock()

		v.op4 = op4{}
		if err := json.Unmarshal(e.RawData, &v.op4); err != nil {
			v.Logger.Error("OP4 unmarshall", "error", err, "raw", e.RawData)
			return
		}
		return

	case discord.VoiceOpCodeSessionSpeaking:
		if len(v.voiceSpeakingUpdateHandlers) == 0 {
			return
		}

		voiceSpeakingUpdate := &SpeakingUpdate{}
		if err := json.Unmarshal(e.RawData, voiceSpeakingUpdate); err != nil {
			v.Logger.Error("OP5 unmarshall, %s", "error", err, "raw", e.RawData)
			return
		}

		for _, h := range v.voiceSpeakingUpdateHandlers {
			h(v, voiceSpeakingUpdate)
		}

	default:
		v.Logger.Warn("unknown voice operation", "op", e.Operation, "raw", e.RawData)
	}

	return
}

type heartbeatOp struct {
	Op   discord.VoiceOpCode `json:"op"` // Always 3
	Data int                 `json:"d"`
}

// NOTE: When a guild.Guild voice server changes how do we shut this down properly, so a new connection can be setup
// without fuss?
//
// wsHeartbeat sends regular heartbeats to voice Discord so it knows the client is still connected.
// If you do not send these heartbeats Discord will disconnect the websocket connection after a few seconds.
func (v *Connection) wsHeartbeat(ctx context.Context, wsConn *websocket.Conn, close <-chan struct{}, i time.Duration) {
	if close == nil || wsConn == nil {
		return
	}
	ticker := time.NewTicker(i * time.Millisecond)
	defer ticker.Stop()
	for {
		v.Logger.Debug("sending heartbeat packet")
		b, err := json.Marshal(heartbeatOp{discord.VoiceOpCodeHeartbeat, int(time.Now().Unix())})
		if err != nil {
			v.Logger.Error("marshall heartbeat", "error", err)
			return
		}
		v.wsMutex.Lock()
		err = wsConn.Write(ctx, websocket.MessageText, b)
		v.wsMutex.Unlock()
		if err != nil {
			v.Logger.Error("sending heartbeat to voice endpoint", "error", err, "endpoint", v.endpoint)
			return
		}

		select {
		case <-ticker.C:
			// continue loop and send heartbeat
		case <-close:
			return
		}
	}
}

type udpData struct {
	Address string `json:"address"` // Public IP of machine running this code
	Port    uint16 `json:"port"`    // UDP Port of machine running this code
	Mode    string `json:"mode"`    // always "xsalsa20_poly1305"
}

type udpD struct {
	Protocol string  `json:"protocol"` // Always "udp" ?
	Data     udpData `json:"data"`
}

type udpOp struct {
	Op   discord.VoiceOpCode `json:"op"` // Always 1
	Data udpD                `json:"d"`
}

// udpOpen opens a UDP connection to the voice server and completes the initial required handshake.
// This connection is left open in the requester and can be used to send or receive audio.
// This should only be called from onEvent OP2.
func (v *Connection) udpOpen(ctx context.Context) error {
	v.Lock()
	defer v.Unlock()

	if v.wsConn == nil {
		return fmt.Errorf("nil voice websocket")
	}

	if v.udpConn != nil {
		return fmt.Errorf("udp connection already open")
	}

	if v.close == nil {
		return fmt.Errorf("nil close channel")
	}

	if v.endpoint == "" {
		return fmt.Errorf("empty endpoint")
	}

	host := v.op2.IP + ":" + strconv.Itoa(v.op2.Port)
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		v.Logger.Error("resolving udp host", "error", err, "host", host)
		return err
	}

	v.Logger.Debug("connecting to udp addr", "addr", addr.String())
	v.udpConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		v.Logger.Error("connecting to udp addr", "error", err, "addr", addr.String())
		return err
	}

	// Create a 74 byte array to store the packet data
	sb := make([]byte, 74)
	binary.BigEndian.PutUint16(sb, 1)              // Packet type (0x1 is request, 0x2 is response)
	binary.BigEndian.PutUint16(sb[2:], 70)         // Packet length (excluding type and length fields)
	binary.BigEndian.PutUint32(sb[4:], v.op2.SSRC) // The SSRC code from the Op 2 VoiceConnection event

	// And send that data over the UDP connection to Discord.
	_, err = v.udpConn.Write(sb)
	if err != nil {
		v.Logger.Error("udp write", "error", err, "addr", addr.String())
		return err
	}

	// Create a 74 byte array and listen for the initial handshake response from Discord.
	// Once we get it parse the IP and PORT information out of the response.
	// This should be our public IP and PORT as Discord saw us.
	rb := make([]byte, 74)
	rlen, _, err := v.udpConn.ReadFromUDP(rb)
	if err != nil {
		v.Logger.Error("udp read", "error", err, "addr", addr.String())
		return err
	}

	if rlen < 74 {
		// is this warn very useful?
		v.Logger.Warn("received udp packet too small")
		return fmt.Errorf("received udp packet too small")
	}

	// Loop over position 8 through 71 to grab the IP address.
	var ip string
	for i := 8; i < len(rb)-2; i++ {
		if rb[i] == 0 {
			break
		}
		ip += string(rb[i])
	}

	// Grab port from position 72 and 73
	port := binary.BigEndian.Uint16(rb[len(rb)-2:])

	// Take the data from above and send it back to Discord to finalize the UDP connection handshake.
	data := udpOp{discord.VoiceOpCodeSelectProtocol, udpD{"udp", udpData{ip, port, "xsalsa20_poly1305"}}}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	v.wsMutex.Lock()
	err = v.wsConn.Write(ctx, websocket.MessageText, b)
	v.wsMutex.Unlock()
	if err != nil {
		v.Logger.Error("udp write", "error", err, "data", data)
		return err
	}

	// start udpKeepAlive
	go v.udpKeepAlive(v.udpConn, v.close, 5*time.Second)
	// TODO: find a way to check that it fired off okay

	return nil
}

// udpKeepAlive sends an udp packet to keep the udp connection open.
// This is a proof of concept.
func (v *Connection) udpKeepAlive(udpConn *net.UDPConn, close <-chan struct{}, i time.Duration) {
	if udpConn == nil || close == nil {
		return
	}

	var err error
	var sequence uint64

	packet := make([]byte, 8)

	ticker := time.NewTicker(i)
	defer ticker.Stop()
	for {
		binary.LittleEndian.PutUint64(packet, sequence)
		sequence++

		_, err = udpConn.Write(packet)
		if err != nil {
			v.Logger.Error("write", "error", err)
			return
		}

		select {
		case <-ticker.C:
			// continue loop and send keepalive
		case <-close:
			return
		}
	}
}

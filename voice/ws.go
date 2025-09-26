package voice

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nyttikord/gokord/discord"
)

// wsListen listens on the voice websocket for messages and passes them to the voice event handler.
// This is automatically called by the open.
func (v *Connection) wsListen(wsConn *websocket.Conn, close <-chan struct{}) {
	for {
		_, message, err := v.wsConn.ReadMessage()
		if err != nil {
			// 4014 indicates a manual disconnection by someone in the guild;
			// we shouldn't Reconnect.
			if websocket.IsCloseError(err, 4014) {
				v.LogInfo("received 4014 manual disconnection")

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
					v.LogInfo("successfully reconnected after 4014 manual disconnection")
					return
				}

				// When VOICE_SERVER_UPDATE is not received, disconnect as usual.
				v.LogInfo("disconnect due to 4014 manual disconnection")

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
				v.LogError(err, "voice endpoint %s websocket closed unexpectantly", v.endpoint)
				go v.Reconnect()
			}
			return
		}

		select {
		case <-close:
			return
		default:
			go v.onEvent(message)
		}
	}
}

// onEvent handles any voice websocket events.
// This is only called by wsListen.
func (v *Connection) onEvent(message []byte) {
	v.LogDebug("received: %s", string(message))

	var e struct {
		discord.Event
		Operation discord.VoiceOpCode `json:"op"`
	}
	if err := json.Unmarshal(message, &e); err != nil {
		v.LogError(err, "unmarshall event")
		return
	}

	switch e.Operation {
	case discord.VoiceOpCodeReady: // READY
		if err := json.Unmarshal(e.RawData, &v.op2); err != nil {
			v.LogError(err, "OP2 unmarshall, %s", string(e.RawData))
			return
		}

		// Start the voice websocket heartbeat to keep the connection alive
		go v.wsHeartbeat(v.wsConn, v.close, v.op2.HeartbeatInterval)
		// TODO monitor a chan/bool to verify this was successful

		// Start the UDP connection
		err := v.udpOpen()
		if err != nil {
			v.LogError(err, "opening udp connection")
			return
		}

		// Start the opusSender.
		// TODO: Should we allow 48000/960 values to be user defined?
		if v.OpusSend == nil {
			v.OpusSend = make(chan []byte, 2)
		}
		go v.opusSender(v.udpConn, v.close, v.OpusSend, 48000, 960)

		// Start the opusReceiver
		if !v.deaf {
			if v.OpusRecv == nil {
				v.OpusRecv = make(chan *Packet, 2)
			}

			go v.opusReceiver(v.udpConn, v.close, v.OpusRecv)
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
			v.LogError(err, "OP4 unmarshall, %s", string(e.RawData))
			return
		}
		return

	case discord.VoiceOpCodeSessionSpeaking:
		if len(v.voiceSpeakingUpdateHandlers) == 0 {
			return
		}

		voiceSpeakingUpdate := &SpeakingUpdate{}
		if err := json.Unmarshal(e.RawData, voiceSpeakingUpdate); err != nil {
			v.LogError(err, "OP5 unmarshall, %s", string(e.RawData))
			return
		}

		for _, h := range v.voiceSpeakingUpdateHandlers {
			h(v, voiceSpeakingUpdate)
		}

	default:
		v.LogWarn("unknown voice operation, %d, %s", e.Operation, string(e.RawData))
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
func (v *Connection) wsHeartbeat(wsConn *websocket.Conn, close <-chan struct{}, i time.Duration) {
	if close == nil || wsConn == nil {
		return
	}

	var err error
	ticker := time.NewTicker(i * time.Millisecond)
	defer ticker.Stop()
	for {
		v.LogDebug("sending heartbeat packet")
		v.wsMutex.Lock()
		err = wsConn.WriteJSON(heartbeatOp{discord.VoiceOpCodeHeartbeat, int(time.Now().Unix())})
		v.wsMutex.Unlock()
		if err != nil {
			v.LogError(err, "sending heartbeat to voice endpoint %s", v.endpoint)
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
func (v *Connection) udpOpen() (err error) {
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
		v.LogError(err, "resolving udp host %s", host)
		return
	}

	v.LogDebug("connecting to udp addr %s", addr.String())
	v.udpConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		v.LogError(err, "connecting to udp addr %s", addr.String())
		return
	}

	// Create a 74 byte array to store the packet data
	sb := make([]byte, 74)
	binary.BigEndian.PutUint16(sb, 1)              // Packet type (0x1 is request, 0x2 is response)
	binary.BigEndian.PutUint16(sb[2:], 70)         // Packet length (excluding type and length fields)
	binary.BigEndian.PutUint32(sb[4:], v.op2.SSRC) // The SSRC code from the Op 2 VoiceConnection event

	// And send that data over the UDP connection to Discord.
	_, err = v.udpConn.Write(sb)
	if err != nil {
		v.LogError(err, "udp write to %s", addr.String())
		return
	}

	// Create a 74 byte array and listen for the initial handshake response from Discord.
	// Once we get it parse the IP and PORT information out of the response.
	// This should be our public IP and PORT as Discord saw us.
	rb := make([]byte, 74)
	rlen, _, err := v.udpConn.ReadFromUDP(rb)
	if err != nil {
		v.LogError(err, "udp read, %s", addr.String())
		return
	}

	if rlen < 74 {
		v.LogWarn("received udp packet too small")
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

	v.wsMutex.Lock()
	err = v.wsConn.WriteJSON(data)
	v.wsMutex.Unlock()
	if err != nil {
		v.LogError(err, "udp write error, %#v", data)
		return
	}

	// start udpKeepAlive
	go v.udpKeepAlive(v.udpConn, v.close, 5*time.Second)
	// TODO: find a way to check that it fired off okay

	return
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
			v.LogError(err, "write")
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

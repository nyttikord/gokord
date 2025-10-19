package voice

import (
	"context"
	"encoding/binary"
	"net"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
)

// opusSender will listen on the given channel and send any pre-encoded opus audio to Discord.
// This is a proof of concept.
func (v *Connection) opusSender(ctx context.Context, udpConn *net.UDPConn, close <-chan struct{}, opus <-chan []byte, rate, size int) {
	if udpConn == nil || close == nil {
		return
	}

	// VoiceConnection is now ready to receive audio packets
	// TODO: this needs reviewed as I think there must be a better way.
	v.Lock()
	v.Ready = true
	v.Unlock()
	defer func() {
		v.Lock()
		v.Ready = false
		v.Unlock()
	}()

	var sequence uint16
	var timestamp uint32
	var recvbuf []byte
	var ok bool
	udpHeader := make([]byte, 12)
	var nonce [24]byte

	// build the parts that don't change in the udpHeader
	udpHeader[0] = 0x80
	udpHeader[1] = 0x78
	binary.BigEndian.PutUint32(udpHeader[8:], v.op2.SSRC)

	// start a send loop that loops until buf chan is closed
	ticker := time.NewTicker(time.Millisecond * time.Duration(size/(rate/1000)))
	defer ticker.Stop()
	for {
		// Get data from chan.
		// If chan is closed, return.
		select {
		case <-close:
			return
		case recvbuf, ok = <-opus:
			if !ok {
				return
			}
			// else, continue loop
		}

		v.RLock()
		speaking := v.speaking
		v.RUnlock()
		if !speaking {
			err := v.Speaking(ctx, true)
			if err != nil {
				v.Logger.Error("sending speaking packet", "error", err)
			}
		}

		// Add sequence and timestamp to udpPacket
		binary.BigEndian.PutUint16(udpHeader[2:], sequence)
		binary.BigEndian.PutUint32(udpHeader[4:], timestamp)

		// encrypt the opus data
		copy(nonce[:], udpHeader)
		v.RLock()
		sendBuf := secretbox.Seal(udpHeader, recvbuf, &nonce, &v.op4.SecretKey)
		v.RUnlock()

		// block here until we're exactly at the right time :)
		// Then send rtp audio packet to Discord over UDP
		select {
		case <-close:
			return
		case <-ticker.C:
			// continue
		}
		_, err := udpConn.Write(sendBuf)

		if err != nil {
			v.Logger.Error("udp write", "error", err)
			v.Logger.Debug("voice", "struct", v)
			return
		}

		if (sequence) == 0xFFFF {
			sequence = 0
		} else {
			sequence++
		}

		if (timestamp + uint32(size)) >= 0xFFFFFFFF {
			timestamp = 0
		} else {
			timestamp += uint32(size)
		}
	}
}

// A Packet contains the headers and content of a received voice packet.
type Packet struct {
	SSRC      uint32
	Sequence  uint16
	Timestamp uint32
	Type      []byte
	Opus      []byte
	PCM       []int16
}

// opusReceiver listens on the UDP socket for incoming packets and sends them across the given channel.
//
// NOTE: This function may change names later.
func (v *Connection) opusReceiver(ctx context.Context, udpConn *net.UDPConn, close <-chan struct{}, c chan *Packet) {
	if udpConn == nil || close == nil {
		return
	}

	recvbuf := make([]byte, 1024)
	var nonce [24]byte

	for {
		rlen, err := udpConn.Read(recvbuf)
		if err != nil {
			// Detect if we have been closed manually. If a Close() has already
			// happened, the udp connection we are listening on will be different
			// to the current requester.
			v.RLock()
			sameConnection := v.udpConn == udpConn
			v.RUnlock()
			if sameConnection {

				v.Logger.Error("udp read", "error", err, "endpoint", v.endpoint)
				v.Logger.Debug("voice", "struct", v)

				go v.Reconnect(ctx)
			}
			return
		}

		select {
		case <-close:
			return
		default:
			// continue loop
		}

		// For now, skip anything except audio.
		if rlen < 12 || (recvbuf[0] != 0x80 && recvbuf[0] != 0x90) {
			continue
		}

		// build a audio packet struct
		p := Packet{}
		p.Type = recvbuf[0:2]
		p.Sequence = binary.BigEndian.Uint16(recvbuf[2:4])
		p.Timestamp = binary.BigEndian.Uint32(recvbuf[4:8])
		p.SSRC = binary.BigEndian.Uint32(recvbuf[8:12])
		// decrypt opus data
		copy(nonce[:], recvbuf[0:12])

		if opus, ok := secretbox.Open(nil, recvbuf[12:rlen], &nonce, &v.op4.SecretKey); ok {
			p.Opus = opus
		} else {
			continue
		}

		// extension bit set, and not a RTCP packet
		if ((recvbuf[0] & 0x10) == 0x10) && ((recvbuf[1] & 0x80) == 0) {
			// get extended header length
			extLen := binary.BigEndian.Uint16(p.Opus[2:4])
			// 4 bytes (ext header header) + 4*extLen (ext header data)
			shift := int(4 + 4*extLen)
			if len(p.Opus) > shift {
				p.Opus = p.Opus[shift:]
			}
		}

		if c != nil {
			select {
			case c <- &p:
			case <-close:
				return
			}
		}
	}
}

package discord

import "encoding/json"

// Event provides a basic initial struct for all websocket events.
type Event struct {
	Operation GatewayOpCode   `json:"op"`
	Sequence  int64           `json:"s"`
	Type      string          `json:"t"`
	RawData   json.RawMessage `json:"d"`
	// Struct contains one of the other types in this file.
	Struct any `json:"-"`
}

type GatewayOpCode uint

const (
	GatewayOpCodeDispatch                GatewayOpCode = 0
	GatewayOpCodeHeartbeat               GatewayOpCode = 1
	GatewayOpCodeIdentify                GatewayOpCode = 2
	GatewayOpCodePresenceUpdate          GatewayOpCode = 3
	GatewayOpCodeVoiceStateUpdate        GatewayOpCode = 4
	GatewayOpCodeResume                  GatewayOpCode = 6
	GatewayOpCodeReconnect               GatewayOpCode = 7
	GatewayOpCodeRequestGuildMembers     GatewayOpCode = 8
	GatewayOpCodeInvalidSession          GatewayOpCode = 9
	GatewayOpCodeHello                   GatewayOpCode = 10
	GatewayOpCodeHeartbeatAck            GatewayOpCode = 11
	GatewayOpCodeRequestSoundboardSounds GatewayOpCode = 31
)

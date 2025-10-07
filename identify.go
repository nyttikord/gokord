package gokord

import "github.com/nyttikord/gokord/discord"

// Identify is sent during initial handshake with the discord gateway.
// https://discord.com/developers/docs/topics/gateway#identify
type Identify struct {
	Token          string              `json:"token"`
	Properties     IdentifyProperties  `json:"properties"`
	Compress       bool                `json:"compress"`
	LargeThreshold int                 `json:"large_threshold"`
	Shard          *[2]int             `json:"shard,omitempty"`
	Presence       GatewayStatusUpdate `json:"presence"`
	Intents        discord.Intent      `json:"intents"`
}

// IdentifyProperties contains the "properties" portion of an Identify packet.
// https://discord.com/developers/docs/topics/gateway#identify-identify-connection-properties
type IdentifyProperties struct {
	OS              string `json:"$os"`
	Browser         string `json:"$browser"`
	Device          string `json:"$device"`
	Referer         string `json:"$referer"`
	ReferringDomain string `json:"$referring_domain"`
}

type identifyOp struct {
	Op   discord.GatewayOpCode `json:"op"`
	Data Identify              `json:"d"`
}

// identify sends the identify packet to the gateway
func (s *Session) identify() error {
	if s.Identify.Shard[0] >= s.Identify.Shard[1] {
		return ErrWSShardBounds
	}

	// Send Identify packet to Discord
	return s.GatewayWriteStruct(identifyOp{discord.GatewayOpCodeIdentify, s.Identify})
}

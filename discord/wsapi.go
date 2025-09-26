package discord

import "encoding/json"

// Event provides a basic initial struct for all websocket events.
type Event struct {
	Operation GatewayOpCode   `json:"op"`
	Sequence  int64           `json:"s"`
	Type      string          `json:"t"`
	RawData   json.RawMessage `json:"d"`
	// Struct contains one of the other types in this file.
	Struct interface{} `json:"-"`
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

type VoiceOpCode uint

const (
	VoiceOpCodeIdentify                        VoiceOpCode = 0
	VoiceOpCodeSelectProtocol                  VoiceOpCode = 1
	VoiceOpCodeReady                           VoiceOpCode = 2
	VoiceOpCodeHeartbeat                       VoiceOpCode = 3
	VoiceOpCodeSessionDescription              VoiceOpCode = 4
	VoiceOpCodeSessionSpeaking                 VoiceOpCode = 5
	VoiceOpCodeHeartbeatAck                    VoiceOpCode = 6
	VoiceOpCodeHeartbeatResume                 VoiceOpCode = 7
	VoiceOpCodeHello                           VoiceOpCode = 8
	VoiceOpCodeResumed                         VoiceOpCode = 9
	VoiceOpCodeClientsConnect                  VoiceOpCode = 11
	VoiceOpCodeClientDisconnect                VoiceOpCode = 13
	VoiceOpCodeDavePrepareTransition           VoiceOpCode = 21
	VoiceOpCodeDaveExecuteTransition           VoiceOpCode = 22
	VoiceOpCodeDaveTransitionReady             VoiceOpCode = 23
	VoiceOpCodeDavePrepareEpoch                VoiceOpCode = 24
	VoiceOpCodeDaveMlsExternalSender           VoiceOpCode = 25
	VoiceOpCodeDaveMlsKeyPackage               VoiceOpCode = 26
	VoiceOpCodeDaveMlsProposals                VoiceOpCode = 27
	VoiceOpCodeDaveMlsCommitWelcome            VoiceOpCode = 28
	VoiceOpCodeDaveMlsAnnounceCommitTransition VoiceOpCode = 29
	VoiceOpCodeDaveMlsWelcome                  VoiceOpCode = 30
	VoiceOpCodeDaveMlsInvalidCommitWelcome     VoiceOpCode = 31
)

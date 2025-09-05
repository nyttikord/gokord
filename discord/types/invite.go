package types

// Target indicates the type of target of an invite
// https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
type Target uint8

// Invite target types
const (
	TargetStream              Target = 1
	TargetEmbeddedApplication Target = 2
)

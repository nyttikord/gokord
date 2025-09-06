package types

// InviteTarget indicates the type of invite.Invite.
// https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
type InviteTarget uint8

const (
	InviteTargetStream              InviteTarget = 1
	InviteTargetEmbeddedApplication InviteTarget = 2
)

//TODO: InviteType is missing

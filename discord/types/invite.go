package types

// InviteTarget indicates the target type of invite.Invite.
// https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
type InviteTarget uint8

const (
	InviteTargetStream              InviteTarget = 1
	InviteTargetEmbeddedApplication InviteTarget = 2
)

// Invite is the type of invite.Invite.
type Invite int

const (
	InviteGuild   Invite = 0
	InviteGroupDM Invite = 1
	InviteFriend  Invite = 2
)

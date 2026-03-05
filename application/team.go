package application

import "github.com/nyttikord/gokord/user"

// The MembershipState represents whether the user is in the team or has been invited into it.
type MembershipState int

const (
	MembershipStateInvited  MembershipState = 1
	MembershipStateAccepted MembershipState = 2
)

// A TeamMember struct stores values for a single [Team] member, extending the normal [user.User] data.
//
// The [TeamMember.User] field is partial.
type TeamMember struct {
	User            *user.User      `json:"user"`
	TeamID          uint64          `json:"team_id,string"`
	MembershipState MembershipState `json:"membership_state"`
	Permissions     []string        `json:"permissions"`
}

// A Team struct stores the members of a Discord Developer Team as well as some metadata about it.
type Team struct {
	ID          uint64        `json:"id,string"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	OwnerID     uint64        `json:"owner_user_id,string"`
	Members     []*TeamMember `json:"members"`
}

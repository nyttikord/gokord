package application

import "github.com/nyttikord/gokord/user"

// The MembershipState represents whether the user is in the team or has been invited into it.
type MembershipState int

const (
	MembershipStateInvited  MembershipState = 1
	MembershipStateAccepted MembershipState = 2
)

// A TeamMember struct stores values for a single Team Member, extending the normal User data - note that the user field
// is partial.
type TeamMember struct {
	User            *user.User      `json:"user"`
	TeamID          string          `json:"team_id"`
	MembershipState MembershipState `json:"membership_state"`
	Permissions     []string        `json:"permissions"`
}

// A Team struct stores the members of a Discord Developer Team as well as some metadata about it.
type Team struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	OwnerID     string        `json:"owner_user_id"`
	Members     []*TeamMember `json:"members"`
}

package guild

import (
	"fmt"
	"github.com/nyttikord/gokord/discord"
)

// RoleFlags represent the flags of a Role.
// https://discord.com/developers/docs/topics/permissions#role-object-role-flags
type RoleFlags int

// Block containing known RoleFlags values.
const (
	// RoleFlagInPrompt indicates whether the Role is selectable by members in an onboarding prompt.
	RoleFlagInPrompt RoleFlags = 1 << 0
)

// A Role stores information about Discord guild member roles.
type Role struct {
	// The ID of the role.
	ID string `json:"id"`

	// The name of the role.
	Name string `json:"name"`

	// Whether this role is managed by an integration, and
	// thus cannot be manually added to, or taken from, members.
	Managed bool `json:"managed"`

	// Whether this role is mentionable.
	Mentionable bool `json:"mentionable"`

	// Whether this role is hoisted (shows up separately in member list).
	Hoist bool `json:"hoist"`

	// The hex color of this role.
	//
	// Deprecated: use Role.Colors
	Color int `json:"color"`

	// The role's colors
	Colors RoleColors `json:"colors"`

	// The position of this role in the guild's role hierarchy.
	Position int `json:"position"`

	// The permissions of the role on the guild (doesn't include channel overrides).
	// This is a combination of bit masks; the presence of a certain permission can
	// be checked by performing a bitwise AND between this int and the permission.
	Permissions int64 `json:"permissions,string"`

	// The hash of the role icon. Use Role.IconURL to retrieve the icon's URL.
	Icon string `json:"icon"`

	// The emoji assigned to this role.
	UnicodeEmoji string `json:"unicode_emoji"`

	// The flags of the role, which describe its extra features.
	// This is a combination of bit masks; the presence of a certain flag can
	// be checked by performing a bitwise AND between this int and the flag.
	Flags RoleFlags `json:"flags"`
}

// Mention returns a string which mentions the role
func (r *Role) Mention() string {
	return fmt.Sprintf("<@&%s>", r.ID)
}

// IconURL returns the URL of the role's icon.
//
//	size:    The size of the desired role icon as a power of two
//	         Image size can be any power of two between 16 and 4096.
func (r *Role) IconURL(size string) string {
	if r.Icon == "" {
		return ""
	}

	URL := discord.EndpointRoleIcon(r.ID, r.Icon)

	if size != "" {
		return URL + "?size=" + size
	}
	return URL
}

// RoleParams represents the parameters needed to create or update a Role
type RoleParams struct {
	// The role's name
	Name string `json:"name,omitempty"`
	// The color the role should have (as a decimal, not hex)
	Color *int `json:"color,omitempty"`
	// Whether to display the role's users separately
	Hoist *bool `json:"hoist,omitempty"`
	// The overall permissions number of the role
	Permissions *int64 `json:"permissions,omitempty,string"`
	// Whether this role is mentionable
	Mentionable *bool `json:"mentionable,omitempty"`
	// The role's unicode emoji.
	// NOTE: can only be set if the guild has the ROLE_ICONS feature.
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	// The role's icon image encoded in base64.
	// NOTE: can only be set if the guild has the ROLE_ICONS feature.
	Icon *string `json:"icon,omitempty"`
}

// Roles are a collection of Role
type Roles []*Role

func (r Roles) Len() int {
	return len(r)
}

func (r Roles) Less(i, j int) bool {
	return r[i].Position > r[j].Position
}

func (r Roles) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// RoleColors stores colors of the Role
type RoleColors struct {
	// Primary color for the role
	PrimaryColor int `json:"primary_color"`
	// Secondary color for the role, this will make the role a gradient between the other provided colors
	SecondaryColor *int `json:"secondary_color"`
	// Tertiary color for the role, this will turn the gradient into a holographic style
	TertiaryColor *int `json:"tertiary_color"`
}

// A GuildedRole stores data for guild roles.
type GuildedRole struct {
	Role    *Role  `json:"role"`
	GuildID string `json:"guild_id"`
}

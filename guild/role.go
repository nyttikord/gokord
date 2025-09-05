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
	// RoleFlagInPrompt indicates whether the Role is selectable by members in an OnboardingPrompt.
	RoleFlagInPrompt RoleFlags = 1 << 0
)

// A Role stores information about Discord Guild user.Member roles.
type Role struct {
	// The ID of the Role.
	ID string `json:"id"`

	// The Name of the Role.
	Name string `json:"name"`

	// Whether this Role is managed by a types.Integration, and thus cannot be manually added to, or taken from, members.
	Managed bool `json:"managed"`

	// Whether this Role is Mentionable.
	Mentionable bool `json:"mentionable"`

	// Whether this Role is hoisted (shows up separately in member list).
	Hoist bool `json:"hoist"`

	// The hex Color of this Role.
	//
	// Deprecated: use Role.Colors
	Color int `json:"color"`

	// The Role's Colors
	Colors RoleColors `json:"colors"`

	// The Position of this Role in the Guild's role hierarchy.
	Position int `json:"position"`

	// The Permissions of the role on the Guild (doesn't include channel overrides).
	// This is a combination of bit masks;
	// the presence of a certain permission can be checked by performing a bitwise AND between this int and the permission.
	Permissions int64 `json:"permissions,string"`

	// The hash of the Role Icon. Use Role.IconURL to retrieve the icon's URL.
	Icon string `json:"icon"`

	// The UnicodeEmoji assigned to this Role.
	UnicodeEmoji string `json:"unicode_emoji"`

	// The Flags of the Role, which describe its extra features.
	// This is a combination of bit masks;
	// the presence of a certain flag can be checked by performing a bitwise AND between this int and the flag.
	Flags RoleFlags `json:"flags"`
}

// Mention returns a string which mentions the Role.
func (r *Role) Mention() string {
	return fmt.Sprintf("<@&%s>", r.ID)
}

// IconURL returns the URL of the Role's icon.
//
// size is the size of the desired icon image as a power of two.
// It can be any power of two between 16 and 4096.
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
	// The Role's Name
	Name string `json:"name,omitempty"`
	// The Color the Role should have (as a decimal, not hex)
	Color *int `json:"color,omitempty"`
	// Whether to display the Role's users separately
	Hoist *bool `json:"hoist,omitempty"`
	// The overall Permissions number of the Role
	Permissions *int64 `json:"permissions,omitempty,string"`
	// Whether this Role is Mentionable
	Mentionable *bool `json:"mentionable,omitempty"`
	// The Role's UnicodeEmoji.
	//
	// Note: can only be set if the guild has the FeatureRoleIcons feature.
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	// The Role's Icon image encoded in base64.
	//
	// Note: can only be set if the guild has the FeatureRoleIcons feature.
	Icon *string `json:"icon,omitempty"`
}

// Roles are a collection of Role.
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

// RoleColors stores colors of the Role.
type RoleColors struct {
	// PrimaryColor for the Role.
	PrimaryColor int `json:"primary_color"`
	// SecondaryColor for the Role, this will make the role a gradient between the other provided colors.
	SecondaryColor *int `json:"secondary_color"`
	// TertiaryColor for the Role, this will turn the gradient into a holographic style.
	TertiaryColor *int `json:"tertiary_color"`
}

// A GuildedRole stores data for Guild Role.
type GuildedRole struct {
	Role    *Role  `json:"role"`
	GuildID string `json:"guild_id"`
}

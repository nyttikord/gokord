package guild

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
)

// RoleFlags represent the flags of a [Role].
// https://discord.com/developers/docs/topics/permissions#role-object-role-flags
type RoleFlags int

// Block containing known RoleFlags values.
const (
	// RoleFlagInPrompt indicates whether the Role is selectable by members in an OnboardingPrompt.
	RoleFlagInPrompt RoleFlags = 1 << 0
)

// Role stores information about Discord [Guild] [user.Member] role.
type Role struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
	// Whether this [Role] is managed by a [user.Integration], and thus cannot be manually added to, or taken from,
	// [user.Member].
	Managed bool `json:"managed"`
	// Whether this [Role] is Mentionable.
	Mentionable bool `json:"mentionable"`
	// Whether this [Role] is hoisted (shows up separately in [user.Member] list).
	Hoist bool `json:"hoist"`
	// The hex Color of this [Role].
	//
	// DEPRECATED: use Role.Colors.
	// Will be removed after the 1.0.0.
	Color int `json:"color"`
	// The Role's Colors.
	Colors RoleColors `json:"colors"`
	// Position of this [Role] in the [Guild]'s role hierarchy.
	Position int `json:"position"`
	// Permissions of the role on the [Guild] (doesn't include [channel.Channel] overrides).
	Permissions int64 `json:"permissions,string"`
	// The hash of the [Role] Icon. Use [Role.IconURL] to retrieve the icon's URL.
	Icon string `json:"icon"`
	// UnicodeEmoji assigned to this [Role].
	UnicodeEmoji string `json:"unicode_emoji"`
	// Flags of the [Role], which describe its extra features.
	Flags RoleFlags `json:"flags"`
}

// Mention returns a string which mentions the Role.
func (r *Role) Mention() string {
	return fmt.Sprintf("<@&%d>", r.ID)
}

// IconURL returns the URL of the [Role.Icon].
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

// RoleParams represents the parameters needed to create or update a [Role].
type RoleParams struct {
	Name string `json:"name,omitempty"`
	// The Color the Role should have (as a decimal, not hex).
	Color *int `json:"color,omitempty"`
	// Whether to display the Role's users separately.
	Hoist *bool `json:"hoist,omitempty"`
	// The overall Permissions number of the Role.
	Permissions *int64 `json:"permissions,omitempty,string"`
	// Whether this Role is Mentionable.
	Mentionable *bool `json:"mentionable,omitempty"`
	// The Role's UnicodeEmoji.
	//
	// NOTE: can only be set if the guild has the FeatureRoleIcons feature.
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	// The Role's Icon image encoded in base64.
	//
	// NOTE: can only be set if the guild has the FeatureRoleIcons feature.
	Icon *string `json:"icon,omitempty"`
}

// Roles are a collection of [Role].
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

// RoleColors stores colors of the [Role].
type RoleColors struct {
	// PrimaryColor for the [Role].
	PrimaryColor int `json:"primary_color"`
	// SecondaryColor for the [Role], this will make the role a gradient between the other provided colors.
	SecondaryColor *int `json:"secondary_color"`
	// TertiaryColor for the [Role], this will turn the gradient into a holographic style.
	TertiaryColor *int `json:"tertiary_color"`
}

// A GuildedRole stores data for [Guild] [Role].
type GuildedRole struct {
	Role    *Role  `json:"role"`
	GuildID uint64 `json:"guild_id,string"`
}

var (
	ErrVerificationLevelBounds = errors.New("VerificationLevel out of bounds, should be between 0 and 3")
	ErrInvalidColorValue       = errors.New("invalid color value: cannot be larger than 0xFFFFFF")
)

// AddMemberRole adds the specified [Role] to a given [user.Member].
func AddMemberRole(guildID, userID, roleID uint64) Empty {
	req := NewSimple(http.MethodPut, discord.EndpointGuildMemberRole(guildID, userID, roleID)).
		WithBucketID(discord.EndpointGuildMembers(guildID))
	return WrapAsEmpty(req)
}

// RemoveMemberRole removes the specified [Role] to a given [user.Member].
func RemoveMemberRole(guildID, userID, roleID uint64) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointGuildMemberRole(guildID, userID, roleID)).
		WithBucketID(discord.EndpointGuildMembers(guildID))
	return WrapAsEmpty(req)
}

// CreateRole in the given [Guild].
func CreateRole(guildID uint64, data *RoleParams) Request[*Role] {
	return NewData[*Role](http.MethodPost, discord.EndpointGuildRoles(guildID)).
		WithData(data)
}

// ListRoles returns all [Role] for the given [Guild].
func ListRoles(guildID uint64) Request[[]*Role] {
	return NewData[[]*Role](http.MethodPost, discord.EndpointGuildRoles(guildID))
}

// EditRole and returns updated data.
func EditRole(guildID, roleID uint64, data *RoleParams) Request[*Role] {
	// Prevent sending a color int that is too big.
	if data.Color != nil && *data.Color > 0xFFFFFF {
		return NewError[*Role](ErrInvalidColorValue)
	}

	return NewData[*Role](http.MethodPatch, discord.EndpointGuildRole(guildID, roleID)).
		WithBucketID(discord.EndpointGuildRoles(guildID)).WithData(data)
}

// ReorderRole with the given data.
func ReorderRole(guildID uint64, roles []*Role) Empty {
	req := NewSimple(http.MethodPatch, discord.EndpointGuildRoles(guildID)).
		WithData(roles)
	return WrapAsEmpty(req)
}

// DeleteRole with the given data.
func DeleteRole(guildID, roleID uint64) Empty {
	req := NewSimple(http.MethodPatch, discord.EndpointGuildRole(guildID, roleID)).
		WithBucketID(discord.EndpointGuildRoles(guildID))
	return WrapAsEmpty(req)
}

// CountsRoleMember returns a map of [Role.ID] to the number of [user.Member] with this role.
// It doesn't include the @everyone [Role].
func CountsRoleMember(guildID uint64) Request[map[string]uint] {
	return NewData[map[string]uint](http.MethodPost, discord.EndpointGuildRoleMemberCounts(guildID)).
		WithBucketID(discord.EndpointGuildRoles(guildID))
}

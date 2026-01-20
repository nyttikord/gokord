package guildapi

import (
	"fmt"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/guild"
)

// RoleCreate creates a new guild.Role.
func (r Requester) RoleCreate(guildID string, data *RoleParams) Request[*Role] {
	return NewSimpleData[*Role](
		r, http.MethodPost, discord.EndpointGuildRoles(guildID),
	).WithData(data)
}

// Roles returns all guild.Role for the given guild.Guild.
func (r Requester) Roles(guildID string) Request[[]*Role] {
	return NewSimpleData[[]*Role](
		r, http.MethodPost, discord.EndpointGuildRoles(guildID),
	)
}

// RoleEdit updates an existing guild.Role and returns updated data.
func (r Requester) RoleEdit(guildID, roleID string, data *RoleParams) Request[*Role] {
	// Prevent sending a color int that is too big.
	if data.Color != nil && *data.Color > 0xFFFFFF {
		return nil, fmt.Errorf("color value cannot be larger than 0xFFFFFF")
	}

	return NewSimpleData[*Role](
		r, http.MethodPatch, discord.EndpointGuildRole(guildID, roleID),
	).WithBucketID(discord.EndpointGuildRoles(guildID)).WithData(data)
}

// RoleReorder reoders guild.Role.
func (r Requester) RoleReorder(guildID string, roles []*Role) Empty {
	req := NewSimple(
		r, http.MethodPatch, discord.EndpointGuildRoles(guildID),
	).WithData(roles)
	return WrapAsEmpty(req)
}

// RoleDelete deletes a guild.Role.
func (r Requester) RoleDelete(guildID, roleID string) Empty {
	req := NewSimple(
		r, http.MethodPatch, discord.EndpointGuildRole(guildID, roleID),
	).WithBucketID(discord.EndpointGuildRoles(guildID))
	return WrapAsEmpty(req)
}

// RoleMemberCounts returns a map of guild.Role ID to the number of user.Member with the role.
// It doesn't include the @everyone Role.
func (r Requester) RoleMemberCounts(guildID string) Request[map[string]uint] {
	return NewSimpleData[map[string]uint](
		r, http.MethodPost, discord.EndpointGuildRoleMemberCounts(guildID),
	).WithBucketID(discord.EndpointGuildRoles(guildID))
}

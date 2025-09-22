package guild

import (
	"sort"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
)

// MemberPermissions calculates the permissions for a user.Member.
// https://support.discord.com/hc/en-us/articles/206141927-How-is-the-permission-hierarchy-structured-
func MemberPermissions(guild *Guild, channel *channel.Channel, userID string, roles []string) int64 {
	if userID == guild.OwnerID {
		return discord.PermissionAll
	}

	var perms int64
	for _, role := range guild.Roles {
		if role.ID == guild.ID {
			perms |= role.Permissions
			break
		}
	}

	for _, role := range guild.Roles {
		for _, roleID := range roles {
			if role.ID == roleID {
				perms |= role.Permissions
				break
			}
		}
	}

	if perms&discord.PermissionAdministrator == discord.PermissionAdministrator {
		perms |= discord.PermissionAll
	}

	// Apply @everyone overrides from the channel.
	for _, overwrite := range channel.PermissionOverwrites {
		if guild.ID == overwrite.ID {
			perms &= ^overwrite.Deny
			perms |= overwrite.Allow
			break
		}
	}

	var denies, allows int64
	// Member overwrites can override role overrides, so do two passes
	for _, overwrite := range channel.PermissionOverwrites {
		for _, roleID := range roles {
			if overwrite.Type == types.PermissionOverwriteRole && roleID == overwrite.ID {
				denies |= overwrite.Deny
				allows |= overwrite.Allow
				break
			}
		}
	}

	perms &= ^denies
	perms |= allows

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.Type == types.PermissionOverwriteMember && overwrite.ID == userID {
			perms &= ^overwrite.Deny
			perms |= overwrite.Allow
			break
		}
	}

	if perms&discord.PermissionAdministrator == discord.PermissionAdministrator {
		perms |= discord.PermissionAllChannel
	}

	return perms
}

func FirstRoleColor(g *Guild, memberRoles []string) int {
	roles := Roles(g.Roles)
	sort.Sort(roles)

	for _, role := range roles {
		for _, roleID := range memberRoles {
			if role.ID == roleID {
				if role.Colors.PrimaryColor != 0 {
					return role.Colors.PrimaryColor
				}
			}
		}
	}

	for _, role := range roles {
		if role.ID == g.ID {
			return role.Colors.PrimaryColor
		}
	}

	return 0
}

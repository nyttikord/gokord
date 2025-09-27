package guild

import (
	"sort"
	"sync"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
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

func Copy(g Guild) Guild {
	var wg sync.WaitGroup

	// deep copy of everything
	// copy() builtin does not copy pointers
	go func() {
		wg.Add(1)
		defer wg.Done()
		roles := make([]*Role, len(g.Roles))
		for i, role := range g.Roles {
			roles[i] = &*role
		}
		g.Roles = roles
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		emojis := make([]*emoji.Emoji, len(g.Emojis))
		for i, e := range g.Emojis {
			emojis[i] = &*e
		}
		g.Emojis = emojis
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		members := make([]*user.Member, len(g.Members))
		for i, m := range g.Members {
			members[i] = &*m
		}
		g.Members = members
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		presences := make([]*status.Presence, len(g.Presences))
		for i, p := range g.Presences {
			presences[i] = &*p
		}
		g.Presences = presences
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		channels := make([]*channel.Channel, len(g.Channels))
		for i, c := range g.Channels {
			channels[i] = &*c
		}
		g.Channels = channels
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		threads := make([]*channel.Channel, len(g.Threads))
		for i, c := range g.Threads {
			threads[i] = &*c
		}
		g.Threads = threads
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		stages := make([]*channel.StageInstance, len(g.StageInstances))
		for i, c := range g.StageInstances {
			stages[i] = &*c
		}
		g.StageInstances = stages
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		voices := make([]*user.VoiceState, len(g.VoiceStates))
		for i, vs := range g.VoiceStates {
			voices[i] = &*vs
		}
		g.VoiceStates = voices
	}()

	wg.Wait()

	return g
}

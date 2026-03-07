package discord

import "fmt"

func getUser(id uint64) string {
	if id == 0 {
		return "@me"
	}
	return fmt.Sprintf("%d", id)
}

func EndpointUser(uID uint64) string {
	return fmt.Sprintf("%s/%s", EndpointUsers, getUser(uID))
}
func EndpointUserGuilds(uID uint64) string {
	return EndpointUser(uID) + "/guilds"
}
func EndpointUserGuild(uID, gID uint64) string {
	return fmt.Sprintf("%s/guilds/%d", EndpointUser(uID), gID)
}
func EndpointUserGuildMember(uID, gID uint64) string {
	return EndpointUserGuild(uID, gID) + "/member"
}
func EndpointUserChannels(uID uint64) string {
	return EndpointUser(uID) + "/channels"
}
func EndpointUserApplicationRoleConnection(aID uint64) string {
	return fmt.Sprintf("%s/@me/applications/%d/role-connection", EndpointUsers, aID)
}
func EndpointUserConnections(uID uint64) string { return EndpointUser(uID) + "/connections" }

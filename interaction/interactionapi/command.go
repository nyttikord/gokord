package interactionapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/interaction"
)

// CommandCreate creates an interaction.Command and returns it.
//
// Specifies guildID if you want to create guild.Guild Command instead of a global one.
func (r Requester) CommandCreate(appID string, guildID string, cmd *Command) Request[*Command] {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	return NewSimpleData[*Command](r, http.MethodPost, endpoint).WithData(cmd)
}

// CommandEdit edits interaction.Command and returns new command data.
//
// Specifies guildID to edit a guild.Guild Command.
func (r Requester) CommandEdit(appID, guildID, cmdID string, cmd *Command) Request[*Command] {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	return NewSimpleData[*Command](r, http.MethodPatch, endpoint).WithData(cmd)
}

// CommandBulkOverwrite creates interaction.Command overwriting all existing interaction.Command.
func (r Requester) CommandBulkOverwrite(appID string, guildID string, cmds []*Command) Request[[]*Command] {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	return NewSimpleData[[]*Command](r, http.MethodPut, endpoint).WithData(cmds)
}

// CommandDelete deletes interaction.Command.
//
// Specifies guildID to delete a guild.Guild interaction.Command.
func (r Requester) CommandDelete(appID, guildID, cmdID string) EmptyRequest {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}
	bucket := discord.EndpointApplicationGlobalCommand(appID, "")
	if guildID != "" {
		bucket = discord.EndpointApplicationGuildCommand(appID, guildID, "")
	}

	req := NewSimple(r, http.MethodDelete, endpoint).WithBucketID(bucket)
	return WrapAsEmpty(req)
}

// Command retrieves an interaction.Command.
//
// Specifies guildID to retrieve a guild.Guild Command.
func (r Requester) Command(appID, guildID, cmdID string) Request[*Command] {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	return NewSimpleData[*Command](r, http.MethodGet, endpoint)
}

// Commands retrieves all interaction.Command.
//
// Specifies guildID to retrieve all guild.Guild interaction.Command from the specified guild.Guild.
func (r Requester) Commands(appID, guildID string) Request[[]*Command] {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	return NewSimpleData[[]*Command](r, http.MethodGet, endpoint+"?with_localizations=true")
}

// GuildCommandsPermissions returns permissions for interaction.Command in a guild.Guild.
func (r Requester) GuildCommandsPermissions(appID, guildID string) Request[[]*GuildCommandPermissions] {
	endpoint := discord.EndpointApplicationCommandsGuildPermissions(appID, guildID)

	return NewSimpleData[[]*GuildCommandPermissions](r, http.MethodGet, endpoint)
}

// CommandPermissions returns all permissions of an interaction.Command.
//
// guildID is the guild.Guild containing the interaction.Command.
// (I don't know if this is required.)
func (r Requester) CommandPermissions(appID, guildID, cmdID string) Request[*GuildCommandPermissions] {
	endpoint := discord.EndpointApplicationCommandPermissions(appID, guildID, cmdID)

	return NewSimpleData[*GuildCommandPermissions](r, http.MethodGet, endpoint)
}

// CommandPermissionsEdit edits the permissions of an interaction.Command.
//
// guildID is the guild.Guild containing the interaction.Command.
// (I don't know if this is required.)
//
// NOTE: Requires OAuth2 token with applications.commands.permissions.update scope.
func (r Requester) CommandPermissionsEdit(appID, guildID, cmdID string, permissions *CommandPermissionsList) EmptyRequest {
	req := NewSimple(
		r, http.MethodPut, discord.EndpointApplicationCommandPermissions(appID, guildID, cmdID),
	).WithBucketID(discord.EndpointApplicationCommandPermissions(appID, guildID, "")).WithData(permissions)
	return WrapAsEmpty(req)
}

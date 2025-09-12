package interactionapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/interaction"
)

// CommandCreate creates an interaction.Command and returns it.
//
// Specifies guildID if you want to create guild.Guild interaction.Command instead of a global one.
func (s Requester) CommandCreate(appID string, guildID string, cmd *interaction.Command, options ...discord.RequestOption) (*interaction.Command, error) {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	body, err := s.Request(http.MethodPost, endpoint, *cmd, options...)
	if err != nil {
		return nil, err
	}

	var c interaction.Command
	return &c, s.Unmarshal(body, &c)
}

// CommandEdit edits interaction.Command and returns new command data.
//
// Specifies guildID to edit a guild.Guild interaction.Command.
func (s Requester) CommandEdit(appID, guildID, cmdID string, cmd *interaction.Command, options ...discord.RequestOption) (*interaction.Command, error) {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	body, err := s.Request(http.MethodPatch, endpoint, *cmd, options...)
	if err != nil {
		return nil, err
	}

	var c interaction.Command
	return &c, s.Unmarshal(body, &c)
}

// CommandBulkOverwrite creates interaction.Command overwriting existing interaction.Command.
func (s Requester) CommandBulkOverwrite(appID string, guildID string, commands []*interaction.Command, options ...discord.RequestOption) (createdCommands []*interaction.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	body, err := s.Request(http.MethodPut, endpoint, commands, options...)
	if err != nil {
		return nil, err
	}

	var c []*interaction.Command
	return c, s.Unmarshal(body, &c)
}

// CommandDelete deletes interaction.Command.
//
// Specifies guildID to delete a guild.Guild interaction.Command.
func (s Requester) CommandDelete(appID, guildID, cmdID string, options ...discord.RequestOption) error {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	_, err := s.Request(http.MethodDelete, endpoint, nil, options...)
	return err
}

// Command retrieves an interaction.Command.
//
// Specifies guildID to retrieve a guild.Guild interaction.Command.
func (s Requester) Command(appID, guildID, cmdID string, options ...discord.RequestOption) (cmd *interaction.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	body, err := s.Request(http.MethodGet, endpoint, nil, options...)
	if err != nil {
		return nil, err
	}

	var c interaction.Command
	return &c, s.Unmarshal(body, &c)
}

// Commands retrieves all interaction.Command.
//
// Specifies guildID to retrieve all guild.Guild interaction.Command from the specified guild.Guild.
func (s Requester) Commands(appID, guildID string, options ...discord.RequestOption) ([]*interaction.Command, error) {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	body, err := s.Request(http.MethodGet, endpoint+"?with_localizations=true", nil, options...)
	if err != nil {
		return nil, err
	}

	var c []*interaction.Command
	return c, s.Unmarshal(body, &c)
}

// GuildCommandsPermissions returns permissions for interaction.Command in a guild.Guild.
func (s Requester) GuildCommandsPermissions(appID, guildID string, options ...discord.RequestOption) ([]*interaction.GuildCommandPermissions, error) {
	endpoint := discord.EndpointApplicationCommandsGuildPermissions(appID, guildID)

	body, err := s.Request(http.MethodGet, endpoint, nil, options...)
	if err != nil {
		return nil, err
	}

	var p []*interaction.GuildCommandPermissions
	return p, s.Unmarshal(body, &p)
}

// CommandPermissions returns all permissions of an interaction.Command.
//
// guildID is the guild.Guild containing the interaction.Command. (I don't know if this is required.)
func (s Requester) CommandPermissions(appID, guildID, cmdID string, options ...discord.RequestOption) (*interaction.GuildCommandPermissions, error) {
	endpoint := discord.EndpointApplicationCommandPermissions(appID, guildID, cmdID)

	body, err := s.Request(http.MethodGet, endpoint, nil, options...)
	if err != nil {
		return nil, err
	}

	var p *interaction.GuildCommandPermissions
	return p, s.Unmarshal(body, &p)
}

// CommandPermissionsEdit edits the permissions of an interaction.Command.
//
// guildID is the guild.Guild containing the interaction.Command. (I don't know if this is required.)
//
// Note: Requires OAuth2 token with applications.commands.permissions.update scope.
func (s Requester) CommandPermissionsEdit(appID, guildID, cmdID string, permissions *interaction.CommandPermissionsList, options ...discord.RequestOption) error {
	_, err := s.Request(
		http.MethodPut,
		discord.EndpointApplicationCommandPermissions(appID, guildID, cmdID),
		permissions,
		options...,
	)
	return err
}

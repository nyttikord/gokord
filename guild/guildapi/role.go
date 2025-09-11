package guildapi

import (
	"fmt"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/guild"
)

// RoleCreate creates a new guild.Role.
func (r Requester) GuildRoleCreate(guildID string, data *guild.RoleParams, options ...discord.RequestOption) (*guild.Role, error) {
	body, err := r.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointGuildRoles(guildID),
		data,
		discord.EndpointGuildRoles(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st guild.Role
	return &st, r.Unmarshal(body, &st)
}

// Roles returns all guild.Role for a given guild.Guild.
func (r Requester) Roles(guildID string, options ...discord.RequestOption) ([]*guild.Role, error) {
	body, err := r.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointGuildRoles(guildID),
		nil,
		discord.EndpointGuildRoles(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st []*guild.Role
	return st, r.Unmarshal(body, st)
}

// RoleEdit updates an existing guild.Role and returns updated data.
func (r Requester) RoleEdit(guildID, roleID string, data *guild.RoleParams, options ...discord.RequestOption) (*guild.Role, error) {
	// Prevent sending a color int that is too big.
	if data.Color != nil && *data.Color > 0xFFFFFF {
		return nil, fmt.Errorf("color value cannot be larger than 0xFFFFFF")
	}

	body, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildRole(guildID, roleID),
		data,
		discord.EndpointGuildRole(guildID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st guild.Role
	return &st, r.Unmarshal(body, &st)
}

// RoleReorder reoders guild.Role.
func (r Requester) RoleReorder(guildID string, roles []*guild.Role, options ...discord.RequestOption) ([]*guild.Role, error) {
	body, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildRoles(guildID),
		roles,
		discord.EndpointGuildRoles(guildID),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var st []*guild.Role
	return st, r.Unmarshal(body, st)
}

// RoleDelete deletes a guild.Role.
func (r Requester) GuildRoleDelete(guildID, roleID string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildRole(guildID, roleID),
		nil,
		discord.EndpointGuildRole(guildID, ""),
		options...,
	)
	return err
}

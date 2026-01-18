package guildapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/guild"
)

// RoleCreate creates a new guild.Role.
func (r Requester) RoleCreate(ctx context.Context, guildID string, data *guild.RoleParams, options ...discord.RequestOption) (*guild.Role, error) {
	body, err := r.Request(ctx, http.MethodPost, discord.EndpointGuildRoles(guildID), data, options...)
	if err != nil {
		return nil, err
	}

	var st guild.Role
	return &st, r.Unmarshal(body, &st)
}

// Roles returns all guild.Role for a given guild.Guild.
func (r Requester) Roles(ctx context.Context, guildID string, options ...discord.RequestOption) ([]*guild.Role, error) {
	body, err := r.Request(ctx, http.MethodGet, discord.EndpointGuildRoles(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var st []*guild.Role
	return st, r.Unmarshal(body, &st)
}

// RoleEdit updates an existing guild.Role and returns updated data.
func (r Requester) RoleEdit(ctx context.Context, guildID, roleID string, data *guild.RoleParams, options ...discord.RequestOption) (*guild.Role, error) {
	// Prevent sending a color int that is too big.
	if data.Color != nil && *data.Color > 0xFFFFFF {
		return nil, fmt.Errorf("color value cannot be larger than 0xFFFFFF")
	}

	body, err := r.RequestWithBucketID(
		ctx,
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
func (r Requester) RoleReorder(ctx context.Context, guildID string, roles []*guild.Role, options ...discord.RequestOption) ([]*guild.Role, error) {
	body, err := r.Request(ctx, http.MethodPatch, discord.EndpointGuildRoles(guildID), roles, options...)
	if err != nil {
		return nil, err
	}

	var st []*guild.Role
	return st, r.Unmarshal(body, &st)
}

// RoleDelete deletes a guild.Role.
func (r Requester) RoleDelete(ctx context.Context, guildID, roleID string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		ctx,
		http.MethodDelete,
		discord.EndpointGuildRole(guildID, roleID),
		nil,
		discord.EndpointGuildRole(guildID, ""),
		options...,
	)
	return err
}

// RoleMemberCounts returns a map of guild.Role ID to the number of user.Member with the role.
// It doesn't include the @everyone guild.Role.
func (r Requester) RoleMemberCounts(ctx context.Context, guildID string, options ...discord.RequestOption) (map[string]uint, error) {
	body, err := r.Request(ctx, http.MethodGet, discord.EndpointGuildRoleMemberCounts(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var res map[string]uint
	return res, r.Unmarshal(body, &res)
}

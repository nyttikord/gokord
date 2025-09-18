// Package application handles everything related with Discord's Application and Team.
//
// Use applicationapi.Requester to interact with this.
// You can get this with gokord.Session.
package application

import (
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// InstallParams represents application's installation parameters for default in-app oauth2 authorization link.
type InstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions int64    `json:"permissions,string"`
}

// IntegrationTypeConfig represents application's configuration for a particular integration type.
type IntegrationTypeConfig struct {
	OAuth2InstallParams *InstallParams `json:"oauth2_install_params,omitempty"`
}

// Application stores values for a Discord application
type Application struct {
	ID                     string                                              `json:"id,omitempty"`
	Name                   string                                              `json:"name"`
	Icon                   string                                              `json:"icon,omitempty"`
	Description            string                                              `json:"description,omitempty"`
	RPCOrigins             []string                                            `json:"rpc_origins,omitempty"`
	BotPublic              bool                                                `json:"bot_public,omitempty"`
	BotRequireCodeGrant    bool                                                `json:"bot_require_code_grant,omitempty"`
	TermsOfServiceURL      string                                              `json:"terms_of_service_url"`
	PrivacyProxyURL        string                                              `json:"privacy_policy_url"`
	Owner                  *user.User                                          `json:"owner"`
	Summary                string                                              `json:"summary"`
	VerifyKey              string                                              `json:"verify_key"`
	Team                   *Team                                               `json:"team"`
	GuildID                string                                              `json:"guild_id"`
	PrimarySKUID           string                                              `json:"primary_sku_id"`
	Slug                   string                                              `json:"slug"`
	CoverImage             string                                              `json:"cover_image"`
	Flags                  int                                                 `json:"flags,omitempty"`
	IntegrationTypesConfig map[types.IntegrationInstall]*IntegrationTypeConfig `json:"integration_types,omitempty"`
}

// RoleConnectionMetadata stores application role connection metadata.
type RoleConnectionMetadata struct {
	Type                     types.RoleConnectionMetadata `json:"type"`
	Key                      string                       `json:"key"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[discord.Locale]string    `json:"name_localizations"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[discord.Locale]string    `json:"description_localizations"`
}

// RoleConnection represents the role connection that an application has attached to a user.
type RoleConnection struct {
	PlatformName     string            `json:"platform_name"`
	PlatformUsername string            `json:"platform_username"`
	Metadata         map[string]string `json:"metadata"`
}

// Asset struct stores values for an asset of an Application
type Asset struct {
	Type int    `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

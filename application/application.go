// Package application handles everything related with Discord's Application and Team.
//
// Use applicationapi.Requester to interact with this.
// You can get this with gokord.Session.
package application

import (
	"context"
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/user"
)

// InstallParams represents [Application]'s installation parameters for default in-app oauth2 authorization link.
type InstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions int64    `json:"permissions,string"`
}

// IntegrationTypeConfig represents [Application]'s configuration for a particular integration type.
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

// RoleConnectionMetadata stores [Application] role connection metadata.
type RoleConnectionMetadata struct {
	Type                     types.RoleConnectionMetadata `json:"type"`
	Key                      string                       `json:"key"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[discord.Locale]string    `json:"name_localizations"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[discord.Locale]string    `json:"description_localizations"`
}

// RoleConnection represents the role connection that an [Application] has attached to a user.
type RoleConnection struct {
	PlatformName     string            `json:"platform_name"`
	PlatformUsername string            `json:"platform_username"`
	Metadata         map[string]string `json:"metadata"`
}

// Asset struct stores values for an asset of an [Application].
type Asset struct {
	Type int    `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Get returns an [Application].
func Get(appID string) Request[*Application] {
	return NewData[*Application](http.MethodGet, discord.EndpointOAuth2Application(appID)).
		WithBucketID(discord.EndpointOAuth2Application(""))
}

// List returns all [Application]s for the authenticated user.
func List() Request[[]*Application] {
	return NewData[[]*Application](http.MethodGet, discord.EndpointOAuth2Applications)
}

// Create a new [Application].
//
// uris are the redirect URIs (not required).
func Create(ap *Application) Request[*Application] {
	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{ap.Name, ap.Description}

	return NewData[*Application](http.MethodPost, discord.EndpointOAuth2Applications).
		WithData(data)
}

// Edit an existing [Application].
func Edit(appID string, ap *Application) Request[*Application] {
	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{ap.Name, ap.Description}

	return NewData[*Application](http.MethodPut, discord.EndpointOAuth2Application(appID)).
		WithData(data).WithBucketID(discord.EndpointOAuth2Application(""))
}

// Delete an existing [Application].
func Delete(appID string) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointOAuth2Application(appID)).
		WithBucketID(discord.EndpointOAuth2Application(""))
	return WrapAsEmpty(req)
}

// ListAssets returns all [Asset]s.
func ListAssets(appID string) Request[[]*Asset] {
	return NewData[[]*Asset](http.MethodGet, discord.EndpointOAuth2Application(appID)).
		WithBucketID(discord.EndpointOAuth2Application(""))
}

// CreateBot creates an [Application] Bot Account.
//
// NOTE: func name may change, if I can think up something better.
func CreateBot(appID string) Request[*user.User] {
	return NewData[*user.User](http.MethodPost, discord.EndpointOAuth2ApplicationsBot(appID)).
		WithBucketID(discord.EndpointOAuth2ApplicationsBot(""))
}

// ListEmojis returns all [emoji.Emoji] for the given [Application].
func ListEmojis(appID string) Request[[]*emoji.Emoji] {
	return NewCustom[[]*emoji.Emoji](http.MethodGet, discord.EndpointApplicationEmojis(appID)).
		WithPost(func(ctx context.Context, b []byte) ([]*emoji.Emoji, error) {
			var data struct {
				Items []*emoji.Emoji `json:"items"`
			}
			return data.Items, Unmarshal(ctx, b, &data)
		})
}

// GetEmoji for the given [Application].
func Emoji(appID, emojiID string) Request[*emoji.Emoji] {
	return NewData[*emoji.Emoji](http.MethodGet, discord.EndpointApplicationEmoji(appID, emojiID)).
		WithBucketID(discord.EndpointApplicationEmojis(appID))
}

// CreateEmoji for the given [Application].
func CreateEmoji(appID string, data *emoji.Params) Request[*emoji.Emoji] {
	return NewData[*emoji.Emoji](http.MethodPost, discord.EndpointApplicationEmojis(appID)).
		WithData(data)
}

// EditEmoji for the given [Application].
func EditEmoji(appID string, emojiID string, data *emoji.Params) Request[*emoji.Emoji] {
	return NewData[*emoji.Emoji](http.MethodPatch, discord.EndpointApplicationEmoji(appID, emojiID)).
		WithData(data).WithBucketID(discord.EndpointApplicationEmojis(appID))
}

// DeleteEmoji for the given [Application].
func DeleteEmoji(appID, emojiID string) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointApplicationEmoji(appID, emojiID)).
		WithBucketID(discord.EndpointApplicationEmojis(appID))
	return WrapAsEmpty(req)
}

// GetRoleConnectionMetadata for the given [Application].
func GetRoleConnectionMetadata(appID string) Request[[]*RoleConnectionMetadata] {
	return NewData[[]*RoleConnectionMetadata](http.MethodGet, discord.EndpointApplicationRoleConnectionMetadata(appID))
}

// EditRoleConnectionMetadata for the given [Application].
func EditRoleConnectionMetadata(appID string, metadata []*RoleConnectionMetadata) Request[[]*RoleConnectionMetadata] {
	return NewData[[]*RoleConnectionMetadata](http.MethodPut, discord.EndpointApplicationRoleConnectionMetadata(appID)).
		WithData(metadata)
}

// GetRoleConnection for the given [Application].
func GetRoleConnection(appID string) Request[*RoleConnection] {
	return NewData[*RoleConnection](http.MethodGet, discord.EndpointUserApplicationRoleConnection(appID))

}

// EditRoleConnection for the specified [Application].
func EditRoleConnection(appID string, rconn *RoleConnection) Request[*RoleConnection] {
	return NewData[*RoleConnection](http.MethodPut, discord.EndpointUserApplicationRoleConnection(appID)).
		WithData(rconn)
}

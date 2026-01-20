// Package applicationapi contains everything to interact with everything located in the application package.
package applicationapi

import (
	"net/http"

	. "github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/user"
)

// Requester handles everything inside the application package.
type Requester struct {
	REST
}

// Application returns an application.Application.
func (r Requester) Application(appID string) Request[*Application] {
	return NewData[*Application](
		r, http.MethodGet, discord.EndpointOAuth2Application(appID),
	).WithBucketID(discord.EndpointOAuth2Application(""))
}

// Applications returns all application.Application for the authenticated user.Application.
func (r Requester) Applications() Request[[]*Application] {
	return NewData[[]*Application](
		r.REST, http.MethodGet, discord.EndpointOAuth2Applications,
	)
}

// ApplicationCreate creates a new application.Application.
//
// uris are the redirect URIs (not required).
func (r Requester) ApplicationCreate(ap *Application) Request[*Application] {
	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{ap.Name, ap.Description}

	return NewData[*Application](
		r, http.MethodPost, discord.EndpointOAuth2Applications,
	).WithData(data)
}

// ApplicationUpdate updates an existing application.Application.
func (r Requester) ApplicationUpdate(appID string, ap *Application) Request[*Application] {
	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{ap.Name, ap.Description}

	return NewData[*Application](
		r, http.MethodPut, discord.EndpointOAuth2Application(appID),
	).WithData(data).WithBucketID(discord.EndpointOAuth2Application(""))
}

// ApplicationDelete deletes an existing application.Application.
func (r Requester) ApplicationDelete(appID string) Empty {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointOAuth2Application(appID),
	).WithBucketID(discord.EndpointOAuth2Application(""))
	return WrapAsEmpty(req)
}

// Assets returns application.Asset.
func (r Requester) Assets(appID string) Request[[]*Asset] {
	return NewData[[]*Asset](
		r, http.MethodGet, discord.EndpointOAuth2Application(appID),
	).WithBucketID(discord.EndpointOAuth2Application(""))
}

// BotCreate creates an application.Application Bot Account.
//
// NOTE: func name may change, if I can think up something better.
func (r Requester) BotCreate(appID string) Request[*user.User] {
	return NewData[*user.User](
		r, http.MethodPost, discord.EndpointOAuth2ApplicationsBot(appID),
	).WithBucketID(discord.EndpointOAuth2ApplicationsBot(""))
}

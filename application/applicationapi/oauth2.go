// Package applicationapi contains everything to interact with everything located in the application package.
package applicationapi

import (
	"net/http"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/user"
)

// Requester handles everything inside the application package.
type Requester struct {
	discord.Requester
}

// Application returns an application.Application.
func (r Requester) Application(appID string, options ...discord.RequestOption) (*application.Application, error) {
	body, err := r.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointOAuth2Application(appID),
		nil,
		discord.EndpointOAuth2Application(""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var app application.Application
	return &app, r.Unmarshal(body, &app)
}

// Applications returns all application.Application for the authenticated user.Application.
func (r Requester) Applications(options ...discord.RequestOption) ([]*application.Application, error) {
	body, err := r.Request("GET", discord.EndpointOAuth2Applications, nil, options...)
	if err != nil {
		return nil, err
	}

	var app []*application.Application
	return app, r.Unmarshal(body, &app)
}

// ApplicationCreate creates a new application.Application.
//
// uris are the redirect URIs (not required).
func (r Requester) ApplicationCreate(ap *application.Application, options ...discord.RequestOption) (*application.Application, error) {
	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{ap.Name, ap.Description}

	body, err := r.Request(http.MethodPost, discord.EndpointOAuth2Applications, data, options...)
	if err != nil {
		return nil, err
	}

	var app application.Application
	return &app, r.Unmarshal(body, &app)
}

// ApplicationUpdate updates an existing application.Application.
func (r Requester) ApplicationUpdate(appID string, ap *application.Application, options ...discord.RequestOption) (*application.Application, error) {
	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{ap.Name, ap.Description}

	body, err := r.RequestWithBucketID(
		http.MethodPut,
		discord.EndpointOAuth2Application(appID),
		data,
		discord.EndpointOAuth2Application(""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var app application.Application
	return &app, r.Unmarshal(body, &app)
}

// ApplicationDelete deletes an existing application.Application.
func (r Requester) ApplicationDelete(appID string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointOAuth2Application(appID),
		nil,
		discord.EndpointOAuth2Application(""),
		options...,
	)
	return err
}

// Assets returns application.Asset.
func (r Requester) Assets(appID string, options ...discord.RequestOption) ([]*application.Asset, error) {
	body, err := r.RequestWithBucketID(
		http.MethodGet,
		discord.EndpointOAuth2ApplicationAssets(appID),
		nil,
		discord.EndpointOAuth2ApplicationAssets(""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var a []*application.Asset
	return a, r.Unmarshal(body, &a)
}

// BotCreate creates an application.Application Bot Account.
//
// NOTE: func name may change, if I can think up something better.
func (r Requester) BotCreate(appID string, options ...discord.RequestOption) (*user.User, error) {
	body, err := r.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointOAuth2ApplicationsBot(appID),
		nil,
		discord.EndpointOAuth2ApplicationsBot(""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var u user.User
	return &u, r.Unmarshal(body, &u)
}

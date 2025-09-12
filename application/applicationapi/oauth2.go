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

// Get returns an application.Application.
func (s Requester) Get(appID string, options ...discord.RequestOption) (*application.Application, error) {
	body, err := s.RequestWithBucketID(
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
	return &app, s.Unmarshal(body, &app)
}

// GetAll returns all application.Application for the authenticated user.Get.
func (s Requester) GetAll(options ...discord.RequestOption) ([]*application.Application, error) {
	body, err := s.Request("GET", discord.EndpointOAuth2Applications, nil, options...)
	if err != nil {
		return nil, err
	}

	var app []*application.Application
	return app, s.Unmarshal(body, &app)
}

// Create creates a new application.Application.
//
// uris are the redirect URIs (not required).
func (s Requester) Create(ap *application.Application, options ...discord.RequestOption) (*application.Application, error) {
	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{ap.Name, ap.Description}

	body, err := s.Request(http.MethodPost, discord.EndpointOAuth2Applications, data, options...)
	if err != nil {
		return nil, err
	}

	var app application.Application
	return &app, s.Unmarshal(body, &app)
}

// Update updates an existing application.Application.
func (s Requester) Update(appID string, ap *application.Application, options ...discord.RequestOption) (*application.Application, error) {
	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{ap.Name, ap.Description}

	body, err := s.RequestWithBucketID(
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
	return &app, s.Unmarshal(body, &app)
}

// Delete deletes an existing application.Application.
func (s Requester) Delete(appID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointOAuth2Application(appID),
		nil,
		discord.EndpointOAuth2Application(""),
		options...,
	)
	return err
}

// Assets returns application.Asset.
func (s Requester) Assets(appID string, options ...discord.RequestOption) ([]*application.Asset, error) {
	body, err := s.RequestWithBucketID(
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
	return a, s.Unmarshal(body, &a)
}

// BotCreate creates an application.Application Bot Account.
//
// Note: func name may change, if I can think up something better.
func (s Requester) BotCreate(appID string, options ...discord.RequestOption) (*user.User, error) {
	body, err := s.RequestWithBucketID(
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
	return &u, s.Unmarshal(body, &u)
}

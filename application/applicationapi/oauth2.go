package applicationapi

import (
	"net/http"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/user"
)

type Requester struct {
	discord.Requester
}

// Application returns an application.Application.
func (s Requester) Application(appID string, options ...discord.RequestOption) (*application.Application, error) {
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

// Applications returns all application.Application for the authenticated user.User.
func (s Requester) Applications(options ...discord.RequestOption) ([]*application.Application, error) {
	body, err := s.Request("GET", discord.EndpointOAuth2Applications, nil, options...)
	if err != nil {
		return nil, err
	}

	var app []*application.Application
	return app, s.Unmarshal(body, &app)
}

// ApplicationCreate creates a new Application.
//
// uris are the redirect URIs (not required).
func (s Requester) ApplicationCreate(ap *application.Application, options ...discord.RequestOption) (*application.Application, error) {
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

// ApplicationUpdate updates an existing application.Application.
func (s Requester) ApplicationUpdate(appID string, ap *application.Application, options ...discord.RequestOption) (*application.Application, error) {
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

// ApplicationDelete deletes an existing application.Application.
func (s Requester) ApplicationDelete(appID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointOAuth2Application(appID),
		nil,
		discord.EndpointOAuth2Application(""),
		options...,
	)
	return err
}

// ApplicationAssets returns an application.Asset.
func (s Requester) ApplicationAssets(appID string, options ...discord.RequestOption) ([]*application.Asset, error) {
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

// ApplicationBotCreate creates an application.Application Bot Account.
//
// Note: func name may change, if I can think up something better.
func (s Requester) ApplicationBotCreate(appID string, options ...discord.RequestOption) (*user.User, error) {
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

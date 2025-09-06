package gokord

import (
	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/user"
)

// Application returns an Application structure of a specific Application
//
//	appID : The ID of an Application
func (s *Session) Application(appID string) (st *application.Application, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointOAuth2Application(appID), nil, discord.EndpointOAuth2Application(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// Applications returns all applications for the authenticated user
func (s *Session) Applications() (st []*application.Application, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointOAuth2Applications, nil, discord.EndpointOAuth2Applications)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ApplicationCreate creates a new Application
//
//	name : Name of Application / Bot
//	uris : Redirect URIs (Not required)
func (s *Session) ApplicationCreate(ap *application.Application) (st *application.Application, err error) {

	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{ap.Name, ap.Description}

	body, err := s.RequestWithBucketID("POST", discord.EndpointOAuth2Applications, data, discord.EndpointOAuth2Applications)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ApplicationUpdate updates an existing Application
//
//	var : desc
func (s *Session) ApplicationUpdate(appID string, ap *application.Application) (st *application.Application, err error) {

	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{ap.Name, ap.Description}

	body, err := s.RequestWithBucketID("PUT", discord.EndpointOAuth2Application(appID), data, discord.EndpointOAuth2Application(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ApplicationDelete deletes an existing Application
//
//	appID : The ID of an Application
func (s *Session) ApplicationDelete(appID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", discord.EndpointOAuth2Application(appID), nil, discord.EndpointOAuth2Application(""))
	if err != nil {
		return
	}

	return
}

// Asset struct stores values for an asset of an application
type Asset struct {
	Type int    `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ApplicationAssets returns an application's assets
func (s *Session) ApplicationAssets(appID string) (ass []*Asset, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointOAuth2ApplicationAssets(appID), nil, discord.EndpointOAuth2ApplicationAssets(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &ass)
	return
}

// ------------------------------------------------------------------------------------------------
// Code specific to Discord OAuth2 Application Bots
// ------------------------------------------------------------------------------------------------

// ApplicationBotCreate creates an Application Bot Account
//
//	appID : The ID of an Application
//
// NOTE: func name may change, if I can think up something better.
func (s *Session) ApplicationBotCreate(appID string) (st *user.User, err error) {

	body, err := s.RequestWithBucketID("POST", discord.EndpointOAuth2ApplicationsBot(appID), nil, discord.EndpointOAuth2ApplicationsBot(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

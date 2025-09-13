package gokord

import (
	_ "image/jpeg" // For JPEG decoding
	_ "image/png"  // For PNG decoding
	"net/http"
	"strings"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/discord"
)

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Voice
// ------------------------------------------------------------------------------------------------

func (s *Session) VoiceRegions(options ...discord.RequestOption) ([]*discord.VoiceRegion, error) {
	body, err := s.Request(http.MethodGet, discord.EndpointVoiceRegions, nil, options...)
	if err != nil {
		return nil, err
	}

	var vc []*discord.VoiceRegion
	return vc, unmarshal(body, &vc)
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Websockets
// ------------------------------------------------------------------------------------------------

// Gateway returns the websocket Gateway address
func (s *Session) Gateway(options ...discord.RequestOption) (string, error) {
	response, err := s.Request(http.MethodGet, discord.EndpointGateway, nil, options...)
	if err != nil {
		return "", err
	}

	temp := struct {
		URL string `json:"url"`
	}{}

	err = unmarshal(response, &temp)
	if err != nil {
		return "", err
	}

	gateway := temp.URL

	// Ensure the gateway always has a trailing slash.
	// MacOS will fail to connect if we add query params without a trailing slash on the base domain.
	if !strings.HasSuffix(gateway, "/") {
		gateway += "/"
	}

	return gateway, nil
}

// GatewayBot returns the websocket Gateway address and the recommended number of shards
func (s *Session) GatewayBot(options ...discord.RequestOption) (*GatewayBotResponse, error) {
	response, err := s.Request(http.MethodGet, discord.EndpointGatewayBot, nil, options...)
	if err != nil {
		return nil, err
	}

	var resp GatewayBotResponse
	err = unmarshal(response, &resp)
	if err != nil {
		return nil, err
	}

	// Ensure the gateway always has a trailing slash.
	// MacOS will fail to connect if we add query params without a trailing slash on the base domain.
	if !strings.HasSuffix(resp.URL, "/") {
		resp.URL += "/"
	}

	return &resp, nil
}

// I don't know what this does, so I let it here.

// ApplicationRoleConnectionMetadata returns application.RoleConnectionMetadata.
func (s *Session) ApplicationRoleConnectionMetadata(appID string, options ...discord.RequestOption) ([]*application.RoleConnectionMetadata, error) {
	body, err := s.Request(http.MethodGet, discord.EndpointApplicationRoleConnectionMetadata(appID), nil, options...)
	if err != nil {
		return nil, err
	}

	var m []*application.RoleConnectionMetadata
	return m, unmarshal(body, &m)
}

// ApplicationRoleConnectionMetadataUpdate updates and returns application.RoleConnectionMetadata.
func (s *Session) ApplicationRoleConnectionMetadataUpdate(appID string, metadata []*application.RoleConnectionMetadata, options ...discord.RequestOption) ([]*application.RoleConnectionMetadata, error) {
	body, err := s.Request(http.MethodPut, discord.EndpointApplicationRoleConnectionMetadata(appID), metadata, options...)
	if err != nil {
		return nil, err
	}

	var m []*application.RoleConnectionMetadata
	return m, unmarshal(body, &m)
}

// UserApplicationRoleConnection returns application.RoleConnection to the specified application.Application.
func (s *Session) UserApplicationRoleConnection(appID string, options ...discord.RequestOption) (*application.RoleConnection, error) {
	body, err := s.Request(http.MethodGet, discord.EndpointUserApplicationRoleConnection(appID), nil, options...)
	if err != nil {
		return nil, err
	}

	var c application.RoleConnection
	return &c, unmarshal(body, &c)
}

// UserApplicationRoleConnectionUpdate updates and returns application.RoleConnection to the specified application.Application.
func (s *Session) UserApplicationRoleConnectionUpdate(appID string, rconn *application.RoleConnection, options ...discord.RequestOption) (*application.RoleConnection, error) {
	body, err := s.Request(http.MethodPut, discord.EndpointUserApplicationRoleConnection(appID), rconn, options...)
	if err != nil {
		return nil, err
	}

	var c application.RoleConnection
	return &c, unmarshal(body, &c)
}

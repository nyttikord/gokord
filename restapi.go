package gokord

import (
	_ "image/jpeg" // For JPEG decoding
	_ "image/png"  // For PNG decoding
	"net/http"
	"strings"

	"github.com/nyttikord/gokord/discord"
)

// Gateway returns the websocket Gateway address
func (s *Session) Gateway() (string, error) {
	response, err := s.REST.Request(ctx, http.MethodGet, discord.EndpointGateway, nil, options...)
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
func (s *Session) GatewayBot() (*GatewayBotResponse, error) {
	response, err := s.REST.Request(ctx, http.MethodGet, discord.EndpointGatewayBot, nil, options...)
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

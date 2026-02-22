package gokord

import (
	"context"
	"net/http"
	"strings"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
)

// Gateway returns the websocket Gateway address.
func (s *Session) Gateway() Request[string] {
	return NewCustom[string](s.rest, http.MethodGet, discord.EndpointGateway).
		WithPost(func(ctx context.Context, b []byte) (string, error) {
			var data struct {
				URL string `json:"url"`
			}
			err := unmarshal(b, &data)
			if err != nil {
				return "", err
			}

			gateway := data.URL

			// Ensure the gateway always has a trailing slash.
			// MacOS will fail to connect if we add query params without a trailing slash on the base domain.
			if !strings.HasSuffix(gateway, "/") {
				gateway += "/"
			}
			return gateway, nil
		})
}

// GatewayBot returns the websocket Gateway address and the recommended number of shards.
func (s *Session) GatewayBot() Request[*GatewayBotResponse] {
	return NewCustom[*GatewayBotResponse](s.rest, http.MethodGet, discord.EndpointGatewayBot).
		WithPost(func(ctx context.Context, b []byte) (*GatewayBotResponse, error) {
			var data GatewayBotResponse
			err := unmarshal(b, &data)
			if err != nil {
				return nil, err
			}

			// Ensure the gateway always has a trailing slash.
			// MacOS will fail to connect if we add query params without a trailing slash on the base domain.
			if !strings.HasSuffix(data.URL, "/") {
				data.URL += "/"
			}
			return &data, nil
		})
}

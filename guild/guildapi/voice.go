package guildapi

import (
	"context"
	"net/http"

	"github.com/nyttikord/gokord/discord"
)

func (r *Requester) VoiceRegions(ctx context.Context, options ...discord.RequestOption) ([]*discord.VoiceRegion, error) {
	body, err := r.Request(ctx, http.MethodGet, discord.EndpointVoiceRegions, nil, options...)
	if err != nil {
		return nil, err
	}

	var vc []*discord.VoiceRegion
	return vc, r.Unmarshal(body, &vc)
}

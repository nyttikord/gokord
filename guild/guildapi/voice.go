package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
)

func (r *Requester) VoiceRegions() Request[[]*discord.VoiceRegion] {
	return NewSimpleData[[]*discord.VoiceRegion](r, http.MethodGet, discord.EndpointVoiceRegions)
}

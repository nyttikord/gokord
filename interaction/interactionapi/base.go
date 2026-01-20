// Package interactionapi contains everything to interact with everything located in the interaction package.
package interactionapi

import (
	"context"
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/channel/channelapi"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/interaction"
)

// Requester handles everything inside the interaction package.
type Requester struct {
	request.REST
	ChannelAPI func() *channelapi.Requester
}

// Respond creates the response to an interaction.Interaction.
func (r Requester) Respond(ctx context.Context, i *Interaction, resp *Response) request.Empty {
	endpoint := discord.EndpointInteractionResponse(i.ID, i.Token)

	if resp.Data == nil || len(resp.Data.Files) == 0 {
		req := request.NewSimple(r, http.MethodPost, endpoint).WithData(resp)
		return request.WrapAsEmpty(req)
	}

	contentType, body, err := channel.MultipartBodyWithJSON(resp, resp.Data.Files)
	if err != nil {
		return request.WrapErrorAsEmpty(err)
	}

	_, err = r.RequestRaw(ctx, http.MethodPost, endpoint, contentType, body, endpoint, 0, options...)
	return err

}

// Response gets the response to an interaction.Interaction.
func (r Requester) Response(i *Interaction) request.Request[*channel.Message] {
	return r.ChannelAPI().WebhookMessage(i.AppID, i.Token, "@original")
}

// ResponseEdit edits the response to an interaction.Interaction.
func (r Requester) ResponseEdit(i *Interaction, newresp *channel.WebhookEdit) request.Request[*channel.Message] {
	return r.ChannelAPI().WebhookMessageEdit(i.AppID, i.Token, "@original", newresp)
}

// ResponseDelete deletes the response to an interaction.Interaction.
func (r Requester) ResponseDelete(i *Interaction) request.Empty {
	req := request.NewSimple(r, http.MethodDelete, discord.EndpointInteractionResponseActions(i.AppID, i.Token))
	return request.WrapAsEmpty(req)
}

// FollowupMessageCreate creates the followup message for an interaction.Interaction.
//
// wait if the function waits for server confirmation of message send and ensures that the return struct is populated
// (it is nil otherwise).
func (r Requester) FollowupMessageCreate(i *Interaction, wait bool, data *channel.WebhookParams) request.Request[*channel.Message] {
	return r.ChannelAPI().WebhookExecute(i.AppID, i.Token, wait, data)
}

// FollowupMessageEdit edits a followup message of an interaction.Interaction.
func (r Requester) FollowupMessageEdit(i *Interaction, messageID string, data *channel.WebhookEdit) request.Request[*channel.Message] {
	return r.ChannelAPI().WebhookMessageEdit(i.AppID, i.Token, messageID, data)
}

// FollowupMessageDelete deletes a followup message of an interaction.Interaction.
func (r Requester) FollowupMessageDelete(i *Interaction, messageID string) request.Empty {
	return r.ChannelAPI().WebhookMessageDelete(i.AppID, i.Token, messageID)
}

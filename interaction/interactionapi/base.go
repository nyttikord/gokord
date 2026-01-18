// Package interactionapi contains everything to interact with everything located in the interaction package.
package interactionapi

import (
	"context"
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/channel/channelapi"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/interaction"
)

// Requester handles everything inside the interaction package.
type Requester struct {
	discord.RESTRequester
	ChannelAPI func() *channelapi.Requester
}

// Respond creates the response to an interaction.Interaction.
func (s Requester) Respond(ctx context.Context, i *interaction.Interaction, resp *interaction.Response, options ...discord.RequestOption) error {
	endpoint := discord.EndpointInteractionResponse(i.ID, i.Token)

	if resp.Data != nil && len(resp.Data.Files) > 0 {
		contentType, body, err := channel.MultipartBodyWithJSON(resp, resp.Data.Files)
		if err != nil {
			return err
		}

		_, err = s.RequestRaw(ctx, http.MethodPost, endpoint, contentType, body, endpoint, 0, options...)
		return err
	}

	_, err := s.Request(ctx, http.MethodPost, endpoint, *resp, options...)
	return err
}

// Response gets the response to an interaction.Interaction.
func (s Requester) Response(ctx context.Context, i *interaction.Interaction, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelAPI().WebhookMessage(ctx, i.AppID, i.Token, "@original", options...)
}

// ResponseEdit edits the response to an interaction.Interaction.
func (s Requester) ResponseEdit(ctx context.Context, i *interaction.Interaction, newresp *channel.WebhookEdit, options ...discord.RequestOption) (*channel.Message, error) {
	return s.
		ChannelAPI().
		WebhookMessageEdit(ctx, i.AppID, i.Token, "@original", newresp, options...)
}

// ResponseDelete deletes the response to an interaction.Interaction.
func (s Requester) ResponseDelete(ctx context.Context, i *interaction.Interaction, options ...discord.RequestOption) error {
	_, err := s.Request(
		ctx,
		http.MethodDelete,
		discord.EndpointInteractionResponseActions(i.AppID, i.Token),
		nil,
		options...,
	)
	return err
}

// FollowupMessageCreate creates the followup message for an interaction.Interaction.
//
// wait if the function waits for server confirmation of message send and ensures that the return struct is populated
// (it is nil otherwise).
func (s Requester) FollowupMessageCreate(ctx context.Context, i *interaction.Interaction, wait bool, data *channel.WebhookParams, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelAPI().WebhookExecute(ctx, i.AppID, i.Token, wait, data, options...)
}

// FollowupMessageEdit edits a followup message of an interaction.Interaction.
func (s Requester) FollowupMessageEdit(ctx context.Context, i *interaction.Interaction, messageID string, data *channel.WebhookEdit, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelAPI().WebhookMessageEdit(ctx, i.AppID, i.Token, messageID, data, options...)
}

// FollowupMessageDelete deletes a followup message of an interaction.Interaction.
func (s Requester) FollowupMessageDelete(ctx context.Context, i *interaction.Interaction, messageID string, options ...discord.RequestOption) error {
	return s.ChannelAPI().WebhookMessageDelete(ctx, i.AppID, i.Token, messageID, options...)
}

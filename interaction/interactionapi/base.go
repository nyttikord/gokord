// Package interactionapi contains everything to interact with everything located in the interaction package.
package interactionapi

import (
	"net/http"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/channel/channelapi"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/interaction"
)

// API adds methods to discord.Requester to be used in the interaction package.
type API interface {
	discord.Requester
	ChannelAPI() channelapi.Requester
}

// Requester handles everything inside the interaction package.
type Requester struct {
	API
}

// Respond creates the response to an interaction.Interaction.
func (s Requester) Respond(interaction *interaction.Interaction, resp *interaction.Response, options ...discord.RequestOption) error {
	endpoint := discord.EndpointInteractionResponse(interaction.ID, interaction.Token)

	if resp.Data != nil && len(resp.Data.Files) > 0 {
		contentType, body, err := channel.MultipartBodyWithJSON(resp, resp.Data.Files)
		if err != nil {
			return err
		}

		_, err = s.RequestRaw(http.MethodPost, endpoint, contentType, body, endpoint, 0, options...)
		return err
	}

	_, err := s.Request(http.MethodPost, endpoint, *resp, options...)
	return err
}

// Response gets the response to an interaction.Interaction.
func (s Requester) Response(interaction *interaction.Interaction, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelAPI().WebhookMessage(interaction.AppID, interaction.Token, "@original", options...)
}

// ResponseEdit edits the response to an interaction.Interaction.
func (s Requester) ResponseEdit(interaction *interaction.Interaction, newresp *channel.WebhookEdit, options ...discord.RequestOption) (*channel.Message, error) {
	return s.
		ChannelAPI().
		WebhookMessageEdit(interaction.AppID, interaction.Token, "@original", newresp, options...)
}

// ResponseDelete deletes the response to an interaction.Interaction.
func (s Requester) ResponseDelete(interaction *interaction.Interaction, options ...discord.RequestOption) error {
	_, err := s.Request(
		http.MethodDelete,
		discord.EndpointInteractionResponseActions(interaction.AppID, interaction.Token),
		nil,
		options...,
	)
	return err
}

// FollowupMessageCreate creates the followup message for an interaction.Interaction.
//
// wait if the function waits for server confirmation of message send and ensures that the return struct is populated
// (it is nil otherwise)
func (s Requester) FollowupMessageCreate(interaction *interaction.Interaction, wait bool, data *channel.WebhookParams, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelAPI().WebhookExecute(interaction.AppID, interaction.Token, wait, data, options...)
}

// FollowupMessageEdit edits a followup message of an interaction.Interaction.
func (s Requester) FollowupMessageEdit(interaction *interaction.Interaction, messageID string, data *channel.WebhookEdit, options ...discord.RequestOption) (*channel.Message, error) {
	return s.ChannelAPI().WebhookMessageEdit(interaction.AppID, interaction.Token, messageID, data, options...)
}

// FollowupMessageDelete deletes a followup message of an interaction.Interaction.
func (s Requester) FollowupMessageDelete(interaction *interaction.Interaction, messageID string, options ...discord.RequestOption) error {
	return s.ChannelAPI().WebhookMessageDelete(interaction.AppID, interaction.Token, messageID, options...)
}

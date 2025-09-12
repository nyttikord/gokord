package gokord

import (
	_ "image/jpeg" // For JPEG decoding
	_ "image/png"  // For PNG decoding
	"net/http"
	"strings"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/user"
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

// ApplicationRoleConnectionMetadata returns application role connection metadata.
// appID : ID of the application
func (s *Session) ApplicationRoleConnectionMetadata(appID string) (st []*application.RoleConnectionMetadata, err error) {
	endpoint := discord.EndpointApplicationRoleConnectionMetadata(appID)
	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ApplicationRoleConnectionMetadataUpdate updates and returns application role connection metadata.
// appID    : ID of the application
// metadata : New metadata
func (s *Session) ApplicationRoleConnectionMetadataUpdate(appID string, metadata []*application.RoleConnectionMetadata) (st []*application.RoleConnectionMetadata, err error) {
	endpoint := discord.EndpointApplicationRoleConnectionMetadata(appID)
	var body []byte
	body, err = s.RequestWithBucketID("PUT", endpoint, metadata, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserApplicationRoleConnection returns user role connection to the specified application.
// appID : ID of the application
func (s *Session) UserApplicationRoleConnection(appID string) (st *application.RoleConnection, err error) {
	endpoint := discord.EndpointUserApplicationRoleConnection(appID)
	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return

}

// UserApplicationRoleConnectionUpdate updates and returns user role connection to the specified application.
// appID      : ID of the application
// connection : New ApplicationRoleConnection data
func (s *Session) UserApplicationRoleConnectionUpdate(appID string, rconn *application.RoleConnection) (st *application.RoleConnection, err error) {
	endpoint := discord.EndpointUserApplicationRoleConnection(appID)
	var body []byte
	body, err = s.RequestWithBucketID("PUT", endpoint, rconn, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ----------------------------------------------------------------------
// Functions specific to polls
// ----------------------------------------------------------------------

// PollAnswerVoters returns users who voted for a particular answer in a poll on the specified message.
// channelID : ID of the channel.
// messageID : ID of the message.
// answerID  : ID of the answer.
func (s *Session) PollAnswerVoters(channelID, messageID string, answerID int) (voters []*user.User, err error) {
	endpoint := discord.EndpointPollAnswerVoters(channelID, messageID, answerID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	var r struct {
		Users []*user.User `json:"users"`
	}

	err = unmarshal(body, &r)
	if err != nil {
		return
	}

	voters = r.Users
	return
}

// PollExpire expires poll on the specified message.
// channelID : ID of the channel.
// messageID : ID of the message.
func (s *Session) PollExpire(channelID, messageID string) (msg *channel.Message, err error) {
	endpoint := discord.EndpointPollExpire(channelID, messageID)

	var body []byte
	body, err = s.RequestWithBucketID("POST", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &msg)
	return
}

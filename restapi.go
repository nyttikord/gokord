package gokord

import (
	_ "image/jpeg" // For JPEG decoding
	_ "image/png"  // For PNG decoding
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/premium"
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

// ----------------------------------------------------------------------
// Functions specific to auto moderation
// ----------------------------------------------------------------------

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

// ----------------------------------------------------------------------
// Functions specific to monetization
// ----------------------------------------------------------------------

// SKUs returns all SKUs for a given application.
// appID : The ID of the application.
func (s *Session) SKUs(appID string) (skus []*premium.SKU, err error) {
	endpoint := discord.EndpointApplicationSKUs(appID)

	body, err := s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &skus)
	return
}

// Entitlements returns all Entitlements for a given app, active and expired.
// appID			: The ID of the application.
// filterOptions	: Optional filter options; otherwise set it to nil.
func (s *Session) Entitlements(appID string, filterOptions *premium.EntitlementFilterOptions, options ...discord.RequestOption) (entitlements []*premium.Entitlement, err error) {
	endpoint := discord.EndpointEntitlements(appID)

	queryParams := url.Values{}
	if filterOptions != nil {
		if filterOptions.UserID != "" {
			queryParams.Set("user_id", filterOptions.UserID)
		}
		if filterOptions.SkuIDs != nil && len(filterOptions.SkuIDs) > 0 {
			queryParams.Set("sku_ids", strings.Join(filterOptions.SkuIDs, ","))
		}
		if filterOptions.Before != nil {
			queryParams.Set("before", filterOptions.Before.Format(time.RFC3339))
		}
		if filterOptions.After != nil {
			queryParams.Set("after", filterOptions.After.Format(time.RFC3339))
		}
		if filterOptions.Limit > 0 {
			queryParams.Set("limit", strconv.Itoa(filterOptions.Limit))
		}
		if filterOptions.GuildID != "" {
			queryParams.Set("guild_id", filterOptions.GuildID)
		}
		if filterOptions.ExcludeEnded {
			queryParams.Set("exclude_ended", "true")
		}
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?"+queryParams.Encode(), nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &entitlements)
	return
}

// EntitlementConsume marks a given One-Time Purchase for the user as consumed.
func (s *Session) EntitlementConsume(appID, entitlementID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID("POST", discord.EndpointEntitlementConsume(appID, entitlementID), nil, discord.EndpointEntitlementConsume(appID, ""), options...)
	return
}

// EntitlementTestCreate creates a test entitlement to a given SKU for a given guild or user.
// Discord will act as though that user or guild has entitlement to your premium offering.
func (s *Session) EntitlementTestCreate(appID string, data *premium.EntitlementTest, options ...discord.RequestOption) (err error) {
	endpoint := discord.EndpointEntitlements(appID)

	_, err = s.RequestWithBucketID("POST", endpoint, data, endpoint, options...)
	return
}

// EntitlementTestDelete deletes a currently-active test entitlement. Discord will act as though
// that user or guild no longer has entitlement to your premium offering.
func (s *Session) EntitlementTestDelete(appID, entitlementID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID("DELETE", discord.EndpointEntitlement(appID, entitlementID), nil, discord.EndpointEntitlement(appID, ""), options...)
	return
}

// Subscriptions returns all subscriptions containing the SKU.
// skuID : The ID of the SKU.
// userID : User ID for which to return subscriptions. Required except for OAuth queries.
// before : Optional timestamp to retrieve subscriptions before this time.
// after : Optional timestamp to retrieve subscriptions after this time.
// limit : Optional maximum number of subscriptions to return (1-100, default 50).
func (s *Session) Subscriptions(skuID string, userID string, before, after *time.Time, limit int, options ...discord.RequestOption) (subscriptions []*premium.Subscription, err error) {
	endpoint := discord.EndpointSubscriptions(skuID)

	queryParams := url.Values{}
	if before != nil {
		queryParams.Set("before", before.Format(time.RFC3339))
	}
	if after != nil {
		queryParams.Set("after", after.Format(time.RFC3339))
	}
	if userID != "" {
		queryParams.Set("user_id", userID)
	}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?"+queryParams.Encode(), nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &subscriptions)
	return
}

// Subscription returns a subscription by its SKU and subscription ID.
// skuID : The ID of the SKU.
// subscriptionID : The ID of the subscription.
// userID : User ID for which to return the subscription. Required except for OAuth queries.
func (s *Session) Subscription(skuID, subscriptionID, userID string, options ...discord.RequestOption) (subscription *premium.Subscription, err error) {
	endpoint := discord.EndpointSubscription(skuID, subscriptionID)

	queryParams := url.Values{}
	if userID != "" {
		// Unlike stated in the documentation, the user_id parameter is required here.
		queryParams.Set("user_id", userID)
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?"+queryParams.Encode(), nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &subscription)
	return
}

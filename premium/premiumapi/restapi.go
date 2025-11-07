// Package premiumapi contains everything to interact with everything located in the premium package.
package premiumapi

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/premium"
)

// Requester handles everything inside the premium package.
type Requester struct {
	discord.RESTRequester
}

// SKUs returns all premium.SKU for a given application.Application.
func (s *Requester) SKUs(appID string, options ...discord.RequestOption) ([]*premium.SKU, error) {
	body, err := s.Request(http.MethodGet, discord.EndpointApplicationSKUs(appID), nil, options...)
	if err != nil {
		return nil, err
	}

	var skus []*premium.SKU
	return skus, s.Unmarshal(body, &skus)
}

// Entitlements returns all premium.Entitlement for a given application.Application, active and expired.
//
// filterOptions is the optional filter options; otherwise set it to nil.
func (s *Requester) Entitlements(appID string, filterOptions *premium.EntitlementFilterOptions, options ...discord.RequestOption) ([]*premium.Entitlement, error) {
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

	body, err := s.Request(http.MethodGet, endpoint+"?"+queryParams.Encode(), nil, options...)
	if err != nil {
		return nil, err
	}

	var e []*premium.Entitlement
	return e, s.Unmarshal(body, &e)
}

// EntitlementConsume marks a given One-Time Purchase for the user.User as consumed.
func (s *Requester) EntitlementConsume(appID, entitlementID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodPost,
		discord.EndpointEntitlementConsume(appID, entitlementID),
		nil,
		discord.EndpointEntitlementConsume(appID, ""),
		options...,
	)
	return err
}

// EntitlementTestCreate creates a test premium.Entitlement to a given premium.SKU for a given guild.Guild or user.User.
//
// Discord will act as though that user or guild has premium.Entitlement to your premium offering.
func (s *Requester) EntitlementTestCreate(appID string, data *premium.EntitlementTest, options ...discord.RequestOption) error {
	_, err := s.Request(http.MethodPost, discord.EndpointEntitlements(appID), data, options...)
	return err
}

// EntitlementTestDelete deletes a currently-active test premium.Entitlement.
//
// Discord will act as though that user.User or guild.Guild no longer has premium.Entitlement to your premium offering.
func (s *Requester) EntitlementTestDelete(appID, entitlementID string, options ...discord.RequestOption) error {
	_, err := s.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointEntitlement(appID, entitlementID),
		nil,
		discord.EndpointEntitlement(appID, ""),
		options...,
	)
	return err
}

// Subscriptions returns all premium.Subscription containing the premium.SKU.
//
// before is an optional timestamp to retrieve premium.Subscription before this time.
// after is an optional timestamp to retrieve premium.Subscription after this time.
// limit is an optional maximum number of premium.Subscription to return (1-100, default 50).
func (s *Requester) Subscriptions(skuID string, userID string, before, after *time.Time, limit int, options ...discord.RequestOption) ([]*premium.Subscription, error) {
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

	body, err := s.Request("GET", endpoint+"?"+queryParams.Encode(), nil, options...)
	if err != nil {
		return nil, err
	}

	var sub []*premium.Subscription
	return sub, s.Unmarshal(body, &sub)
}

// Subscription returns a premium.Subscription by its premium.SKU and premium.Subscription ID.
//
// userID for which to return the premium.Subscription. Required except for OAuth queries.
func (s *Requester) Subscription(skuID, subscriptionID, userID string, options ...discord.RequestOption) (*premium.Subscription, error) {
	endpoint := discord.EndpointSubscription(skuID, subscriptionID)

	queryParams := url.Values{}
	if userID != "" {
		// Unlike stated in the documentation, the user_id parameter is required here.
		queryParams.Set("user_id", userID) //TODO: check if this is true
	}

	body, err := s.Request(http.MethodGet, endpoint+"?"+queryParams.Encode(), nil, options...)
	if err != nil {
		return nil, err
	}

	var sub premium.Subscription
	return &sub, s.Unmarshal(body, &sub)
}

// Package premiumapi contains everything to interact with everything located in the premium package.
package premiumapi

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/premium"
)

// Requester handles everything inside the premium package.
type Requester struct {
	RESTRequester
}

// SKUs returns all premium.SKU for a given application.Application.
func (r *Requester) SKUs(appID string) Request[[]*SKU] {
	return NewSimpleData[[]*SKU](r, http.MethodGet, discord.EndpointApplicationSKUs(appID))
}

// Entitlements returns all premium.Entitlement for a given application.Application, active and expired.
//
// filterOptions is the optional filter options; otherwise set it to nil.
func (r *Requester) Entitlements(appID string, filterOptions *EntitlementFilterOptions) Request[[]*Entitlement] {
	endpoint := discord.EndpointEntitlements(appID)

	queryParams := url.Values{}
	if filterOptions != nil {
		if filterOptions.UserID != "" {
			queryParams.Set("user_id", filterOptions.UserID)
		}
		if len(filterOptions.SkuIDs) > 0 {
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
		endpoint += "?" + queryParams.Encode()
	}

	return NewSimpleData[[]*Entitlement](r, http.MethodGet, endpoint)
}

// EntitlementConsume marks a given One-Time Purchase for the user.User as consumed.
func (r *Requester) EntitlementConsume(appID, entitlementID string) EmptyRequest {
	req := NewSimple(
		r, http.MethodPost, discord.EndpointEntitlementConsume(appID, entitlementID),
	).WithBucketID(discord.EndpointEntitlementConsume(appID, ""))
	return WrapAsEmpty(req)
}

// EntitlementTestCreate creates a test premium.Entitlement to a given premium.SKU for a given guild.Guild or user.User.
//
// Discord will act as though that user or guild has premium.Entitlement to your premium offering.
func (r *Requester) EntitlementTestCreate(appID string, data *EntitlementTest) EmptyRequest {
	req := NewSimple(r, http.MethodPost, discord.EndpointEntitlements(appID)).WithData(data)
	return WrapAsEmpty(req)
}

// EntitlementTestDelete deletes a currently-active test premium.Entitlement.
//
// Discord will act as though that user.User or guild.Guild no longer has Entitlement to your premium offering.
func (r *Requester) EntitlementTestDelete(appID, entitlementID string) EmptyRequest {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointEntitlement(appID, entitlementID),
	).WithBucketID(discord.EndpointEntitlement(appID, ""))
	return WrapAsEmpty(req)
}

// Subscriptions returns all premium.Subscription containing the premium.SKU.
//
// before is an optional timestamp to retrieve Subscription before this time.
// after is an optional timestamp to retrieve Subscription after this time.
// limit is an optional maximum number of Subscription to return (1-100, default 50).
func (r *Requester) Subscriptions(skuID string, userID string, before, after *time.Time, limit int) Request[[]*Subscription] {
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

	return NewSimpleData[[]*Subscription](r, http.MethodGet, endpoint+"?"+queryParams.Encode())
}

// Subscription returns a premium.Subscription by its premium.SKU and premium.Subscription ID.
//
// userID for which to return the premium.Subscription.
// Required except for OAuth queries.
func (r *Requester) Subscription(skuID, subscriptionID, userID string) Request[*Subscription] {
	endpoint := discord.EndpointSubscription(skuID, subscriptionID)

	queryParams := url.Values{}
	if userID != "" {
		// Unlike stated in the documentation, the user_id parameter is required here.
		queryParams.Set("user_id", userID) //TODO: check if this is true
		endpoint += "?" + queryParams.Encode()
	}

	return NewSimpleData[*Subscription](r, http.MethodGet, endpoint)
}

// Package premium contains everything related with premium things in your application.Application.
//
// Use premiumapi.Requester to interact with this.
// You can get this with gokord.Session.
package premium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/internal/structs"
)

// SKUFlags is a bitfield of flags used to differentiate user and server subscriptions (see SKUFlag* consts)
// https://discord.com/developers/docs/monetization/skus#sku-object-sku-flags
type SKUFlags int

const (
	// SKUFlagAvailable indicates that the [SKU] is available for purchase.
	SKUFlagAvailable SKUFlags = 1 << 2
	// SKUFlagGuildSubscription indicates that the [SKU] is a guild subscription.
	SKUFlagGuildSubscription SKUFlags = 1 << 7
	// SKUFlagUserSubscription indicates that the [SKU] is a user subscription.
	SKUFlagUserSubscription SKUFlags = 1 << 8
)

// SKU (stock-keeping units) represent premium offerings.
type SKU struct {
	ID   uint64    `json:"id,string"`
	Type types.SKU `json:"type"`
	// The id of the parent [application.Application].
	ApplicationID uint64 `json:"application_id,string"`
	// Customer-facing name of the SKU.
	Name string `json:"name"`
	// System-generated URL Slug based on the [SKU]'s name.
	Slug string `json:"slug"`
	// SKUFlags combined as a bitfield.
	Flags SKUFlags `json:"flags"`
}

// Subscription represents a user making recurring payments for at least one [SKU] over an ongoing period.
// https://discord.com/developers/docs/resources/subscription#subscription-object
type Subscription struct {
	ID uint64 `json:"id,string"`
	// ID of the [user.User] who is subscribed.
	UserID uint64 `json:"user_id,string"`
	// List of [SKU]s subscribed to.
	SKUIDs []uint64 `json:"-"`
	// List of [Entitlement] granted for this [Subscription].
	EntitlementIDs []uint64 `json:"-"`
	// List of [SKU]s that this [user.User] will be subscribed to at renewal.
	RenewalSKUIDs []uint64 `json:"-"`
	// Start of the current [Subscription] period.
	CurrentPeriodStart time.Time `json:"current_period_start"`
	// End of the current [Subscription] period.
	CurrentPeriodEnd time.Time `json:"current_period_end"`
	// Current Status of the [Subscription].
	Status SubscriptionStatus `json:"status"`
	// Time when the [Subscription] was canceled.
	// Only present if the [Subscription] has been canceled.
	CanceledAt *time.Time `json:"canceled_at,omitempty"`
	// ISO3166-1 alpha-2 Country code of the payment source used to purchase the [Subscription].
	// Missing unless queried with a private OAuth scope.
	Country string `json:"country,omitempty"`
}

func (s *Subscription) MarshalJSON() ([]byte, error) {
	type t Subscription
	v := struct {
		t
		SKUIDs         []string `json:"sku_ids"`
		EntitlementIDs []string `json:"entitlement_ids"`
		RenewalSKUIDs  []string `json:"renewal_sku_ids,omitempty"`
	}{
		t(*s),
		structs.UintsToSnowflakes(s.SKUIDs),
		structs.UintsToSnowflakes(s.EntitlementIDs),
		structs.UintsToSnowflakes(s.RenewalSKUIDs),
	}
	return json.Marshal(v)
}

func (s *Subscription) UnmarshalJSON(data []byte) error {
	type t Subscription
	var v struct {
		t
		SKUIDs         []string `json:"sku_ids"`
		EntitlementIDs []string `json:"entitlement_ids"`
		RenewalSKUIDs  []string `json:"renewal_sku_ids,omitempty"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*s = Subscription(v.t)
	s.SKUIDs = structs.SnowflakesToUints(v.SKUIDs)
	s.EntitlementIDs = structs.SnowflakesToUints(v.EntitlementIDs)
	s.RenewalSKUIDs = structs.SnowflakesToUints(v.RenewalSKUIDs)
	return nil
}

// SubscriptionStatus is the current status of a [Subscription].
// https://discord.com/developers/docs/resources/subscription#subscription-statuses
type SubscriptionStatus int

const (
	SubscriptionStatusActive   = 0
	SubscriptionStatusEnding   = 1
	SubscriptionStatusInactive = 2
)

// Entitlement represents that a [user.User] or [guild.Guild] has access to a premium offering in your
// [application.Application].
type Entitlement struct {
	ID uint64 `json:"id,string"`
	// The ID of the [SKU].
	SKUID uint64 `json:"sku_id,string"`
	// The ID of the parent [application.Application].
	ApplicationID uint64 `json:"application_id,string"`
	// The id of the [user.User] that is granted access to the [Entitlement]'s [SKU].
	// Only available for [user.User] [Subscription]s.
	UserID uint64 `json:"user_id,omitempty,string"`
	// The Type of the [Entitlement].
	Type types.Entitlement `json:"type"`
	// If fhe [Entitlement] was deleted.
	Deleted bool `json:"deleted"`
	// The start date at which the [Entitlement] is valid.
	// Not present when using test [Entitlement].
	StartsAt *time.Time `json:"starts_at,omitempty"`
	// The date at which the [Entitlement] is no longer valid.
	// Not present when using test [Entitlement] or when receiving an ENTITLEMENT_CREATE event.
	EndsAt *time.Time `json:"ends_at,omitempty"`
	// The id of the [guild.Guild] that is granted access to the [Entitlement]'s [SKU].
	// Only available for [guild.Guild] [Subscription].
	GuildID uint64 `json:"guild_id,omitempty,string"`
	// Whether the [Entitlement] has been consumed.
	// Only available for consumable items.
	Consumed *bool `json:"consumed,omitempty"`
	// The SubscriptionID of the [Entitlement].
	// Not present when using test [Entitlement].
	SubscriptionID uint64 `json:"subscription_id,omitempty,string"`
}

// EntitlementTest is used to test granting an [Entitlement] to a [user.User] or [guild.Guild].
type EntitlementTest struct {
	// The ID of the [SKU] to grant the [Entitlement] to.
	SKUID uint64 `json:"sku_id,string"`
	// The ID of the [guild.Guild] or [user.User] to grant the [Entitlement] to.
	OwnerID uint64 `json:"owner_id,string"`
	// OwnerType is the type of which the [Entitlement] should be created.
	OwnerType types.EntitlementOwner `json:"owner_type"`
}

// EntitlementFilterOptions are the options for filtering [Entitlement].
type EntitlementFilterOptions struct {
	// Optional [user.User] ID to look up for.
	UserID uint64
	// Optional array of [SKU.ID]s to check for.
	SkuIDs []uint64
	// Optional timestamp to retrieve [Entitlement] before this time.
	Before *time.Time
	// Optional timestamp to retrieve [Entitlement] after this time.
	After *time.Time
	// Optional maximum number of [Entitlement] to return (1-100, default 100).
	Limit int
	// Optional [guild.Guild] ID to look up for.
	GuildID uint64
	// Optional whether ended [Entitlement] should be omitted.
	ExcludeEnded bool
}

// ListSKUs returns all [SKU] for a given [application.Application].
func ListSKUs(appID uint64) Request[[]*SKU] {
	return NewData[[]*SKU](http.MethodGet, discord.EndpointApplicationSKUs(appID))
}

// ListEntitlements returns all [Entitlement] for a given [application.Application], active and expired.
//
// filterOptions is the optional filter options; otherwise set it to nil.
func ListEntitlements(appID uint64, filterOptions *EntitlementFilterOptions) Request[[]*Entitlement] {
	endpoint := discord.EndpointEntitlements(appID)

	queryParams := url.Values{}
	if filterOptions != nil {
		if filterOptions.UserID != 0 {
			queryParams.Set("user_id", fmt.Sprintf("%d", filterOptions.UserID))
		}
		if len(filterOptions.SkuIDs) > 0 {
			queryParams.Set("sku_ids", strings.Join(structs.UintsToSnowflakes(filterOptions.SkuIDs), ","))
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
		if filterOptions.GuildID != 0 {
			queryParams.Set("guild_id", fmt.Sprintf("%d", filterOptions.GuildID))
		}
		if filterOptions.ExcludeEnded {
			queryParams.Set("exclude_ended", "true")
		}
		endpoint += "?" + queryParams.Encode()
	}

	return NewData[[]*Entitlement](http.MethodGet, endpoint)
}

// ConsumeEntitlement marks a given One-Time Purchase for the [user.User] as consumed.
func ConsumeEntitlement(appID, entitlementID uint64) Empty {
	req := NewSimple(http.MethodPost, discord.EndpointEntitlementConsume(appID, entitlementID)).
		WithBucketID(discord.EndpointEntitlementConsume(appID, 0))
	return WrapAsEmpty(req)
}

// CreateEntitlementTest to a given [SKU] for a given [guild.Guild] or [user.User].
//
// Discord will act as though that user or guild has [Entitlement] to your premium offering.
func CreateEntitlementTest(appID uint64, data *EntitlementTest) Empty {
	req := NewSimple(http.MethodPost, discord.EndpointEntitlements(appID)).WithData(data)
	return WrapAsEmpty(req)
}

// DeleteEntitlementTest deletes a currently-active test [Entitlement].
//
// Discord will act as though that [user.User] or [guild.Guild] no longer has [Entitlement] to your premium offering.
func DeleteEntitlementTest(appID, entitlementID uint64) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointEntitlement(appID, entitlementID)).
		WithBucketID(discord.EndpointEntitlement(appID, 0))
	return WrapAsEmpty(req)
}

// ListSubscriptions returns all [Subscription] containing the [SKU].
//
// before is an optional timestamp to retrieve Subscription before this time.
// after is an optional timestamp to retrieve Subscription after this time.
// limit is an optional maximum number of Subscription to return (1-100, default 50).
func ListSubscriptions(skuID, userID uint64, before, after *time.Time, limit int) Request[[]*Subscription] {
	endpoint := discord.EndpointSubscriptions(skuID)

	queryParams := url.Values{}
	if before != nil {
		queryParams.Set("before", before.Format(time.RFC3339))
	}
	if after != nil {
		queryParams.Set("after", after.Format(time.RFC3339))
	}
	if userID != 0 {
		queryParams.Set("user_id", fmt.Sprintf("%d", userID))
	}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}

	return NewData[[]*Subscription](http.MethodGet, endpoint+"?"+queryParams.Encode())
}

// GetSubscription by its [SKU.ID] and [Subscription.ID].
//
// userID for which to return the [Subscription].
// Required except for OAuth queries.
func GetSubscription(skuID, subscriptionID, userID uint64) Request[*Subscription] {
	endpoint := discord.EndpointSubscription(skuID, subscriptionID)

	queryParams := url.Values{}
	if userID != 0 {
		// Unlike stated in the documentation, the user_id parameter is required here.
		queryParams.Set("user_id", fmt.Sprintf("%d", userID)) //TODO: check if this is true
		endpoint += "?" + queryParams.Encode()
	}

	return NewData[*Subscription](http.MethodGet, endpoint)
}

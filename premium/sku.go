// Package premium contains everything related with premium things in your application.Application.
package premium

import (
	"time"

	"github.com/nyttikord/gokord/discord/types"
)

// SKUFlags is a bitfield of flags used to differentiate user and server subscriptions (see SKUFlag* consts)
// https://discord.com/developers/docs/monetization/skus#sku-object-sku-flags
type SKUFlags int

const (
	// SKUFlagAvailable indicates that the SKU is available for purchase.
	SKUFlagAvailable SKUFlags = 1 << 2
	// SKUFlagGuildSubscription indicates that the SKU is a guild subscription.
	SKUFlagGuildSubscription SKUFlags = 1 << 7
	// SKUFlagUserSubscription indicates that the SKU is a user subscription.
	SKUFlagUserSubscription SKUFlags = 1 << 8
)

// SKU (stock-keeping units) represent premium offerings
type SKU struct {
	// The ID of the SKU.
	ID string `json:"id"`

	// The Type of the SKU.
	Type types.SKU `json:"type"`

	// The id of the parent application.Application.
	ApplicationID string `json:"application_id"`

	// Customer-facing name of the SKU.
	Name string `json:"name"`

	// System-generated URL Slug based on the SKU's name.
	Slug string `json:"slug"`

	// SKUFlags combined as a bitfield.
	// The presence of a certain flag can be checked by performing a bitwise AND operation between this int and the flag.
	Flags SKUFlags `json:"flags"`
}

// Subscription represents a user making recurring payments for at least one SKU over an ongoing period.
// https://discord.com/developers/docs/resources/subscription#subscription-object
type Subscription struct {
	// ID of the Subscription.
	ID string `json:"id"`

	// id of the user.User who is subscribed.
	UserID string `json:"user_id"`

	// List of SKU subscribed to.
	SKUIDs []string `json:"sku_ids"`

	// List of Entitlement granted for this Subscription.
	EntitlementIDs []string `json:"entitlement_ids"`

	// List of SKU that this user.User will be subscribed to at renewal.
	RenewalSKUIDs []string `json:"renewal_sku_ids,omitempty"`

	// Start of the current Subscription period.
	CurrentPeriodStart time.Time `json:"current_period_start"`

	// End of the current Subscription period.
	CurrentPeriodEnd time.Time `json:"current_period_end"`

	// Current Status of the Subscription.
	Status SubscriptionStatus `json:"status"`

	// Time when the Subscription was canceled.
	// Only present if the subscription has been canceled.
	CanceledAt *time.Time `json:"canceled_at,omitempty"`

	// ISO3166-1 alpha-2 Country code of the payment source used to purchase the Subscription.
	// Missing unless queried with a private OAuth scope.
	Country string `json:"country,omitempty"`
}

// SubscriptionStatus is the current status of a Subscription.
// https://discord.com/developers/docs/resources/subscription#subscription-statuses
type SubscriptionStatus int

const (
	SubscriptionStatusActive   = 0
	SubscriptionStatusEnding   = 1
	SubscriptionStatusInactive = 2
)

// Entitlement represents that a user or guild has access to a premium offering in your application.Application.
type Entitlement struct {
	// The ID of the Entitlement.
	ID string `json:"id"`

	// The ID of the SKU.
	SKUID string `json:"sku_id"`

	// The id of the parent application.Application.
	ApplicationID string `json:"application_id"`

	// The id of the user.User that is granted access to the Entitlement's SKU.
	// Only available for user subscriptions.
	UserID string `json:"user_id,omitempty"`

	// The Type of the Entitlement.
	Type types.Entitlement `json:"type"`

	// If fhe Entitlement was deleted.
	Deleted bool `json:"deleted"`

	// The start date at which the Entitlement is valid.
	// Not present when using test Entitlement.
	StartsAt *time.Time `json:"starts_at,omitempty"`

	// The date at which the Entitlement is no longer valid.
	// Not present when using test Entitlement or when receiving an ENTITLEMENT_CREATE event.
	EndsAt *time.Time `json:"ends_at,omitempty"`

	// The id of the guild.Guild that is granted access to the Entitlement's SKU.
	// Only available for guild Subscription.
	GuildID string `json:"guild_id,omitempty"`

	// Whether the Entitlement has been consumed.
	// Only available for consumable items.
	Consumed *bool `json:"consumed,omitempty"`

	// The SubscriptionID of the Entitlement.
	// Not present when using test Entitlement.
	SubscriptionID string `json:"subscription_id,omitempty"`
}

// EntitlementTest is used to test granting an Entitlement to a user.User or guild.Guild.
type EntitlementTest struct {
	// The ID of the SKU to grant the Entitlement to.
	SKUID string `json:"sku_id"`

	// The ID of the guild.Guild or user.User to grant the Entitlement to.
	OwnerID string `json:"owner_id"`

	// OwnerType is the type of which the Entitlement should be created.
	OwnerType types.EntitlementOwner `json:"owner_type"`
}

// EntitlementFilterOptions are the options for filtering Entitlement.
type EntitlementFilterOptions struct {
	// Optional user.User ID to look up for.
	UserID string

	// Optional array of SKU IDs to check for.
	SkuIDs []string

	// Optional timestamp to retrieve Entitlement before this time.
	Before *time.Time

	// Optional timestamp to retrieve Entitlement after this time.
	After *time.Time

	// Optional maximum number of Entitlement to return (1-100, default 100).
	Limit int

	// Optional guild.Guild ID to look up for.
	GuildID string

	// Optional whether ended Entitlement should be omitted.
	ExcludeEnded bool
}

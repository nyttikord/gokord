package premium

import "time"

// SKUType is the type of SKU (see SKUType* consts)
// https://discord.com/developers/docs/monetization/skus
type SKUType int

// Valid SKUType values
const (
	SKUTypeDurable      SKUType = 2
	SKUTypeConsumable   SKUType = 3
	SKUTypeSubscription SKUType = 5
	// SKUTypeSubscriptionGroup is a system-generated group for each subscription SKU.
	SKUTypeSubscriptionGroup SKUType = 6
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
	// The ID of the SKU
	ID string `json:"id"`

	// The Type of the SKU
	Type SKUType `json:"type"`

	// The ID of the parent application
	ApplicationID string `json:"application_id"`

	// Customer-facing name of the SKU.
	Name string `json:"name"`

	// System-generated URL slug based on the SKU's name.
	Slug string `json:"slug"`

	// SKUFlags combined as a bitfield. The presence of a certain flag can be checked
	// by performing a bitwise AND operation between this int and the flag.
	Flags SKUFlags `json:"flags"`
}

// Subscription represents a user making recurring payments for at least one SKU over an ongoing period.
// https://discord.com/developers/docs/resources/subscription#subscription-object
type Subscription struct {
	// ID of the subscription
	ID string `json:"id"`

	// ID of the user who is subscribed
	UserID string `json:"user_id"`

	// List of SKUs subscribed to
	SKUIDs []string `json:"sku_ids"`

	// List of entitlements granted for this subscription
	EntitlementIDs []string `json:"entitlement_ids"`

	// List of SKUs that this user will be subscribed to at renewal
	RenewalSKUIDs []string `json:"renewal_sku_ids,omitempty"`

	// Start of the current subscription period
	CurrentPeriodStart time.Time `json:"current_period_start"`

	// End of the current subscription period
	CurrentPeriodEnd time.Time `json:"current_period_end"`

	// Current status of the subscription
	Status SubscriptionStatus `json:"status"`

	// When the subscription was canceled. Only present if the subscription has been canceled.
	CanceledAt *time.Time `json:"canceled_at,omitempty"`

	// ISO3166-1 alpha-2 country code of the payment source used to purchase the subscription. Missing unless queried with a private OAuth scope.
	Country string `json:"country,omitempty"`
}

// SubscriptionStatus is the current status of a Subscription Object
// https://discord.com/developers/docs/resources/subscription#subscription-statuses
type SubscriptionStatus int

// Valid SubscriptionStatus values
const (
	SubscriptionStatusActive   = 0
	SubscriptionStatusEnding   = 1
	SubscriptionStatusInactive = 2
)

// EntitlementType is the type of entitlement (see EntitlementType* consts)
// https://discord.com/developers/docs/monetization/entitlements#entitlement-object-entitlement-types
type EntitlementType int

// Valid EntitlementType values
const (
	EntitlementTypePurchase                = 1
	EntitlementTypePremiumSubscription     = 2
	EntitlementTypeDeveloperGift           = 3
	EntitlementTypeTestModePurchase        = 4
	EntitlementTypeFreePurchase            = 5
	EntitlementTypeUserGift                = 6
	EntitlementTypePremiumPurchase         = 7
	EntitlementTypeApplicationSubscription = 8
)

// Entitlement represents that a user or guild has access to a premium offering
// in your application.
type Entitlement struct {
	// The ID of the entitlement
	ID string `json:"id"`

	// The ID of the SKU
	SKUID string `json:"sku_id"`

	// The ID of the parent application
	ApplicationID string `json:"application_id"`

	// The ID of the user that is granted access to the entitlement's sku
	// Only available for user subscriptions.
	UserID string `json:"user_id,omitempty"`

	// The type of the entitlement
	Type EntitlementType `json:"type"`

	// The entitlement was deleted
	Deleted bool `json:"deleted"`

	// The start date at which the entitlement is valid.
	// Not present when using test entitlements.
	StartsAt *time.Time `json:"starts_at,omitempty"`

	// The date at which the entitlement is no longer valid.
	// Not present when using test entitlements or when receiving an ENTITLEMENT_CREATE event.
	EndsAt *time.Time `json:"ends_at,omitempty"`

	// The ID of the guild that is granted access to the entitlement's sku.
	// Only available for guild subscriptions.
	GuildID string `json:"guild_id,omitempty"`

	// Whether or not the entitlement has been consumed.
	// Only available for consumable items.
	Consumed *bool `json:"consumed,omitempty"`

	// The SubscriptionID of the entitlement.
	// Not present when using test entitlements.
	SubscriptionID string `json:"subscription_id,omitempty"`
}

// EntitlementOwnerType is the type of entitlement (see EntitlementOwnerType* consts)
type EntitlementOwnerType int

// Valid EntitlementOwnerType values
const (
	EntitlementOwnerTypeGuildSubscription EntitlementOwnerType = 1
	EntitlementOwnerTypeUserSubscription  EntitlementOwnerType = 2
)

// EntitlementTest is used to test granting an entitlement to a user or guild
type EntitlementTest struct {
	// The ID of the SKU to grant the entitlement to
	SKUID string `json:"sku_id"`

	// The ID of the guild or user to grant the entitlement to
	OwnerID string `json:"owner_id"`

	// OwnerType is the type of which the entitlement should be created
	OwnerType EntitlementOwnerType `json:"owner_type"`
}

// EntitlementFilterOptions are the options for filtering Entitlements
type EntitlementFilterOptions struct {
	// Optional user ID to look up for.
	UserID string

	// Optional array of SKU IDs to check for.
	SkuIDs []string

	// Optional timestamp to retrieve Entitlements before this time.
	Before *time.Time

	// Optional timestamp to retrieve Entitlements after this time.
	After *time.Time

	// Optional maximum number of entitlements to return (1-100, default 100).
	Limit int

	// Optional guild ID to look up for.
	GuildID string

	// Optional whether or not ended entitlements should be omitted.
	ExcludeEnded bool
}

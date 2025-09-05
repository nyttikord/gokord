package types

// SKU is the type of SKU (see SKU* consts)
// https://discord.com/developers/docs/monetization/skus
type SKU int

// Valid SKU values
const (
	SKUDurable      SKU = 2
	SKUConsumable   SKU = 3
	SKUSubscription SKU = 5
	// SKUSubscriptionGroup is a system-generated group for each subscription SKU.
	SKUSubscriptionGroup SKU = 6
)

// Entitlement is the type of entitlement (see Entitlement* consts)
// https://discord.com/developers/docs/monetization/entitlements#entitlement-object-entitlement-types
type Entitlement int

// Valid Entitlement values
const (
	EntitlementPurchase                = 1
	EntitlementPremiumSubscription     = 2
	EntitlementDeveloperGift           = 3
	EntitlementTestModePurchase        = 4
	EntitlementFreePurchase            = 5
	EntitlementUserGift                = 6
	EntitlementPremiumPurchase         = 7
	EntitlementApplicationSubscription = 8
)

// EntitlementOwner is the type of entitlement (see EntitlementOwner* consts)
type EntitlementOwner int

// Valid EntitlementOwner values
const (
	EntitlementOwnerGuildSubscription EntitlementOwner = 1
	EntitlementOwnerUserSubscription  EntitlementOwner = 2
)

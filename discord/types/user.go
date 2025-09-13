package types

// Premium is the type of premium (nitro) subscription a user.User has.
// https://discord.com/developers/docs/resources/user#user-object-premium-types
type Premium int

const (
	PremiumNone Premium = 0
	// PremiumNitroClassic is the old subscription that cost $5
	PremiumNitroClassic Premium = 1
	// PremiumNitro is the nitro "boost" subscription that cost $10
	PremiumNitro Premium = 2
	// PremiumNitroBasic is the nitro costing $3
	PremiumNitroBasic Premium = 3
)

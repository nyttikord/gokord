package types

// Premium is the type of premium (nitro) subscription a user has (see Premium* consts).
// https://discord.com/developers/docs/resources/user#user-object-premium-types
type Premium int

// Valid Premium values.
const (
	PremiumNone         Premium = 0
	PremiumNitroClassic Premium = 1
	PremiumNitro        Premium = 2
	PremiumNitroBasic   Premium = 3
)

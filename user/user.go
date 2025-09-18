package user

import (
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"

	"strconv"
)

// Flags is the flags for User.
// https://discord.com/developers/docs/resources/user#user-object-user-flags
type Flags int

const (
	FlagDiscordEmployee           Flags = 1 << 0
	FlagDiscordPartner            Flags = 1 << 1
	FlagHypeSquadEvents           Flags = 1 << 2
	FlagBugHunterLevel1           Flags = 1 << 3
	FlagHouseBravery              Flags = 1 << 6
	FlagHouseBrilliance           Flags = 1 << 7
	FlagHouseBalance              Flags = 1 << 8
	FlagEarlySupporter            Flags = 1 << 9
	FlagTeamUser                  Flags = 1 << 10
	FlagSystem                    Flags = 1 << 12
	FlagBugHunterLevel2           Flags = 1 << 14
	FlagVerifiedBot               Flags = 1 << 16
	FlagVerifiedBotDeveloper      Flags = 1 << 17
	FlagDiscordCertifiedModerator Flags = 1 << 18
	FlagBotHTTPInteractions       Flags = 1 << 19
	FlagActiveBotDeveloper        Flags = 1 << 22
)

// User stores all data for an individual Discord user.
type User struct {
	// The ID of the User.
	ID string `json:"id"`

	// The Email of the User.
	// This is only present when the application possesses the email scope for the User.
	Email string `json:"email"`

	// The User's Username.
	Username string `json:"username"`

	// The hash of the User's Avatar.
	// Use User.AvatarURL to retrieve the avatar itself.
	Avatar string `json:"avatar"`

	// The User's chosen language option.
	Locale string `json:"locale"`

	// The discriminator of the User (4 numbers after name).
	Discriminator string `json:"discriminator"`

	// The User's display name, if it is set.
	// For bots, this is the application.Application name.
	GlobalName string `json:"global_name"`

	// The Token of the User.
	// This is only present for the user represented by the current session.
	Token string `json:"token"`

	// Whether the User's email is Verified.
	Verified bool `json:"verified"`

	// Whether the User has multifactor authentication enabled.
	MFAEnabled bool `json:"mfa_enabled"`

	// The hash of the User's Banner image.
	Banner string `json:"banner"`

	// User's Banner color, encoded as an integer representation of hexadecimal color code.
	AccentColor int `json:"accent_color"`

	// Whether the User is a bot.
	Bot bool `json:"bot"`

	// The PublicFlags on a User's account.
	// This is a combination of bit masks; the presence of a certain flag can be checked by performing a bitwise AND
	// between this int and the flag.
	PublicFlags Flags `json:"public_flags"`

	// The type of Nitro subscription on a User's account.
	// Only available when the request is authorized via a Bearer token.
	PremiumType types.Premium `json:"premium_type"`

	// Whether the user is an Official Discord System User (part of the urgent message system).
	System bool `json:"system"`

	// The Flags on a User's account.
	// Only available when the request is authorized via a Bearer token.
	Flags int `json:"flags"`

	// Data for the User's AvatarDecoration.
	AvatarDecorationData *AvatarDecoration `json:"avatar_decoration_data"`

	// Data for the User's Collectibles.
	Collectibles *Collectibles `json:"collectibles"`

	// User's PrimaryGuild (the tag).
	PrimaryGuild *PrimaryGuild `json:"primary_guild"`
}

// String returns a unique identifier of the form username#discriminator or just username, if the discriminator is set
// to "0".
func (u *User) String() string {
	// If the user has been migrated from the legacy username system, their discriminator is "0".
	// See https://support-dev.discord.com/hc/en-us/articles/13667755828631
	if u.Discriminator == "0" {
		return u.Username
	}

	return u.Username + "#" + u.Discriminator
}

// Mention return a string which mentions the User.
func (u *User) Mention() string {
	return "<@" + u.ID + ">"
}

// AvatarURL returns a URL to the User.Avatar.
//
// size is the size of the avatar as a power of two if size is an empty string, no size parameter will be added to the
// URL (between 16 and 4096).
func (u *User) AvatarURL(size string) string {
	return discord.AvatarURL(
		u.Avatar,
		discord.EndpointDefaultUserAvatar(u.DefaultAvatarIndex()),
		discord.EndpointUserAvatar(u.ID, u.Avatar),
		discord.EndpointUserAvatarAnimated(u.ID, u.Avatar),
		size,
	)
}

// BannerURL returns the URL of the User.Banner.
//
// size is the size of the banner as a power of two if size is an empty string, no size parameter will be added to the
// URL (between 16 and 4096).
func (u *User) BannerURL(size string) string {
	return discord.BannerURL(u.Banner, discord.EndpointUserBanner(u.ID, u.Banner), discord.EndpointUserBannerAnimated(u.ID, u.Banner), size)
}

// DefaultAvatarIndex returns the index of the user's default avatar.
func (u *User) DefaultAvatarIndex() int {
	if u.Discriminator == "0" {
		id, _ := strconv.ParseUint(u.ID, 10, 64)
		return int((id >> 22) % 6)
	}

	id, _ := strconv.Atoi(u.Discriminator)
	return id % 5
}

// DisplayName returns the User's global name if they have one, otherwise it returns their User.Username.
func (u *User) DisplayName() string {
	if u.GlobalName != "" {
		return u.GlobalName
	}
	return u.Username
}

type AvatarDecoration struct {
	// Avatar decoration hash.
	Asset string `json:"asset"`
	// ID of the AvatarDecoration's premium.SKU.
	SkuID string `json:"sku_id"`
}

type Collectibles struct {
	Nameplate *Nameplate `json:"nameplate"`
}

type NameplatePalette string

const (
	NameplatePaletteCrimson   = "crimson"
	NameplatePaletteBerry     = "berry"
	NameplatePaletteSky       = "sky"
	NameplatePaletteTeal      = "teal"
	NameplatePaletteForest    = "forest"
	NameplatePaletteBubbleGum = "bubble_gum"
	NameplatePaletteViolet    = "violet"
	NameplatePaletteCobalt    = "cobalt"
	NameplatePaletteClover    = "clover"
	NameplatePaletteLemon     = "lemon"
	NameplatePaletteWhite     = "white"
)

type Nameplate struct {
	// ID of the Nameplate premium.SKU.
	SkuID string `json:"sku_id"`
	// Path to the Nameplate Asset
	Asset string `json:"asset"`
	// Label of this Nameplate.
	// Currently unused.
	Label string `json:"label"`
	// Background color of the Nameplate.
	Palette NameplatePalette `json:"palette"`
}

type PrimaryGuild struct {
	// ID of the User's primary guild.Guild.
	GuildID string `json:"identity_guild_id"`
	// Whether the User is displaying the primary guild.Guild's server tag.
	//
	// This can be null if the system clears the identity, e.g. the server no longer supports tags.
	// This will be false if the user manually removes their tag.
	Enabled *bool `json:"identity_enabled"`
	// Text of the User's server Tag.
	// Limited to 4 characters
	Tag *string `json:"tag"`
	// Server tag Badge hash.
	Badge *string `json:"badge"`
}

func (upg *PrimaryGuild) IsEnabled() bool {
	return upg.Enabled != nil && *upg.Enabled
}

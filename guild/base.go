package guild

import (
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/user"
	"time"
)

// VerificationLevel type definition
type VerificationLevel int

// Constants for VerificationLevel levels from 0 to 4 inclusive
const (
	VerificationLevelNone     VerificationLevel = 0
	VerificationLevelLow      VerificationLevel = 1
	VerificationLevelMedium   VerificationLevel = 2
	VerificationLevelHigh     VerificationLevel = 3
	VerificationLevelVeryHigh VerificationLevel = 4
)

// ExplicitContentFilterLevel type definition
type ExplicitContentFilterLevel int

// Constants for ExplicitContentFilterLevel levels from 0 to 2 inclusive
const (
	ExplicitContentFilterDisabled            ExplicitContentFilterLevel = 0
	ExplicitContentFilterMembersWithoutRoles ExplicitContentFilterLevel = 1
	ExplicitContentFilterAllMembers          ExplicitContentFilterLevel = 2
)

// GuildNSFWLevel type definition
type GuildNSFWLevel int

// Constants for GuildNSFWLevel levels from 0 to 3 inclusive
const (
	NSFWLevelDefault       GuildNSFWLevel = 0
	NSFWLevelExplicit      GuildNSFWLevel = 1
	NSFWLevelSafe          GuildNSFWLevel = 2
	NSFWLevelAgeRestricted GuildNSFWLevel = 3
)

// MfaLevel type definition
type MfaLevel int

// Constants for MfaLevel levels from 0 to 1 inclusive
const (
	MfaLevelNone     MfaLevel = 0
	MfaLevelElevated MfaLevel = 1
)

// PremiumTier type definition
type PremiumTier int

// Constants for PremiumTier levels from 0 to 3 inclusive
const (
	PremiumTierNone PremiumTier = 0
	PremiumTier1    PremiumTier = 1
	PremiumTier2    PremiumTier = 2
	PremiumTier3    PremiumTier = 3
)

// MessageNotifications is the notification level for a guild
// https://discord.com/developers/docs/resources/guild#guild-object-default-message-notification-level
type MessageNotifications int

// Block containing known MessageNotifications values
const (
	MessageNotificationsAllMessages  MessageNotifications = 0
	MessageNotificationsOnlyMentions MessageNotifications = 1
)

// SystemChannelFlag is the type of flags in the system channel (see SystemChannelFlag* consts)
// https://discord.com/developers/docs/resources/guild#guild-object-system-channel-flags
type SystemChannelFlag int

// Block containing known SystemChannelFlag values
const (
	SystemChannelFlagsSuppressJoinNotifications          SystemChannelFlag = 1 << 0
	SystemChannelFlagsSuppressPremium                    SystemChannelFlag = 1 << 1
	SystemChannelFlagsSuppressGuildReminderNotifications SystemChannelFlag = 1 << 2
	SystemChannelFlagsSuppressJoinNotificationReplies    SystemChannelFlag = 1 << 3
)

// Feature indicates the presence of a feature in a guild
type Feature string

// Constants for Feature
const (
	FeatureAnimatedBanner                        Feature = "ANIMATED_BANNER"
	FeatureAnimatedIcon                          Feature = "ANIMATED_ICON"
	FeatureApplicationCommandPermissionV2        Feature = "APPLICATION_COMMAND_PERMISSIONS_V2"
	FeatureAutoModeration                        Feature = "AUTO_MODERATION"
	FeatureBanner                                Feature = "BANNER"
	FeatureCommunity                             Feature = "COMMUNITY"
	FeatureCreatorMonetizableProvisional         Feature = "CREATOR_MONETIZABLE_PROVISIONAL"
	FeatureCreatorStorePage                      Feature = "CREATOR_STORE_PAGE"
	FeatureDeveloperSupportServer                Feature = "DEVELOPER_SUPPORT_SERVER"
	FeatureDiscoverable                          Feature = "DISCOVERABLE"
	FeatureFeaturable                            Feature = "FEATURABLE"
	FeatureInvitesDisabled                       Feature = "INVITES_DISABLED"
	FeatureInviteSplash                          Feature = "INVITE_SPLASH"
	FeatureMemberVerificationGateEnabled         Feature = "MEMBER_VERIFICATION_GATE_ENABLED"
	FeatureMoreSoundboard                        Feature = "MORE_SOUNDBOARD"
	FeatureMoreStickers                          Feature = "MORE_STICKERS"
	FeatureNews                                  Feature = "NEWS"
	FeaturePartnered                             Feature = "PARTNERED"
	FeaturePreviewEnabled                        Feature = "PREVIEW_ENABLED"
	FeatureRaidAlertsDisabled                    Feature = "RAID_ALERTS_DISABLED"
	FeatureRoleIcons                             Feature = "ROLE_ICONS"
	FeatureRoleSubscriptionsAvailableForPurchase Feature = "ROLE_SUBSCRIPTIONS_AVAILABLE_FOR_PURCHASE"
	FeatureRoleSubscriptionsEnabled              Feature = "ROLE_SUBSCRIPTIONS_ENABLED"
	FeatureSoundboard                            Feature = "SOUNDBOARD"
	FeatureTicketedEventsEnabled                 Feature = "TICKETED_EVENTS_ENABLED"
	FeatureVanityURL                             Feature = "VANITY_URL"
	FeatureVerified                              Feature = "VERIFIED"
	FeatureVipRegions                            Feature = "VIP_REGIONS"
	FeatureWelcomeScreenEnabled                  Feature = "WELCOME_SCREEN_ENABLED"
	FeatureEnhancedRoleColors                    Feature = "ENHANCED_ROLE_COLORS"
)

// A Guild holds all data related to a specific Discord Guild.  Guilds are also
// sometimes referred to as Servers in the Discord client.
type Guild struct {
	// The ID of the guild.
	ID string `json:"id"`

	// The name of the guild. (2–100 characters)
	Name string `json:"name"`

	// The hash of the guild's icon. Use Session.GuildIcon
	// to retrieve the icon itself.
	Icon string `json:"icon"`

	// The voice region of the guild.
	Region string `json:"region"`

	// The ID of the AFK voice channel.
	AfkChannelID string `json:"afk_channel_id"`

	// The user ID of the owner of the guild.
	OwnerID string `json:"owner_id"`

	// If we are the owner of the guild
	Owner bool `json:"owner"`

	// The time at which the current user joined the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	JoinedAt time.Time `json:"joined_at"`

	// The hash of the guild's discovery splash.
	DiscoverySplash string `json:"discovery_splash"`

	// The hash of the guild's splash.
	Splash string `json:"splash"`

	// The timeout, in seconds, before a user is considered AFK in voice.
	AfkTimeout int `json:"afk_timeout"`

	// The number of members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	MemberCount int `json:"member_count"`

	// The verification level required for the guild.
	VerificationLevel VerificationLevel `json:"verification_level"`

	// Whether the guild is considered large. This is
	// determined by a member threshold in the identify packet,
	// and is currently hard-coded at 250 members in the library.
	Large bool `json:"large"`

	// The default message notification setting for the guild.
	DefaultMessageNotifications MessageNotifications `json:"default_message_notifications"`

	// A list of roles in the guild.
	Roles []*Role `json:"roles"`

	// A list of the custom emojis present in the guild.
	Emojis []*user.Emoji `json:"emojis"`

	// A list of the custom stickers present in the guild.
	Stickers []*user.Sticker `json:"stickers"`

	// A list of the members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Members []*user.Member `json:"members"`

	// A list of partial presence objects for members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Presences []*user.Presence `json:"presences"`

	// The maximum number of presences for the guild (the default value, currently 25000, is in effect when null is returned)
	MaxPresences int `json:"max_presences"`

	// The maximum number of members for the guild
	MaxMembers int `json:"max_members"`

	// A list of channels in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Channels []*Channel `json:"channels"`

	// A list of all active threads in the guild that current user has permission to view
	// This field is only present in GUILD_CREATE events and websocket
	// update events and thus is only present in state-cached guilds.
	Threads []*Channel `json:"threads"`

	// A list of voice states for the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	VoiceStates []*user.VoiceState `json:"voice_states"`

	// Whether this guild is currently unavailable (most likely due to outage).
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Unavailable bool `json:"unavailable"`

	// The explicit content filter level
	ExplicitContentFilter ExplicitContentFilterLevel `json:"explicit_content_filter"`

	// The NSFW Level of the guild
	NSFWLevel GuildNSFWLevel `json:"nsfw_level"`

	// The list of enabled guild features
	Features []Feature `json:"features"`

	// Required MFA level for the guild
	MfaLevel MfaLevel `json:"mfa_level"`

	// The application id of the guild if bot created.
	ApplicationID string `json:"application_id"`

	// Whether the Server Widget is enabled
	WidgetEnabled bool `json:"widget_enabled"`

	// The Channel ID for the Server Widget
	WidgetChannelID string `json:"widget_channel_id"`

	// The Channel ID to which system messages are sent (eg join and leave messages)
	SystemChannelID string `json:"system_channel_id"`

	// The System channel flags
	SystemChannelFlags SystemChannelFlag `json:"system_channel_flags"`

	// The ID of the rules channel ID, used for rules.
	RulesChannelID string `json:"rules_channel_id"`

	// the vanity url code for the guild
	VanityURLCode string `json:"vanity_url_code"`

	// the description for the guild
	Description string `json:"description"`

	// The hash of the guild's banner
	Banner string `json:"banner"`

	// The premium tier of the guild
	PremiumTier PremiumTier `json:"premium_tier"`

	// The total number of users currently boosting this server
	PremiumSubscriptionCount int `json:"premium_subscription_count"`

	// The preferred locale of a guild with the "PUBLIC" feature; used in server discovery and notices from Discord; defaults to "en-US"
	PreferredLocale string `json:"preferred_locale"`

	// The id of the channel where admins and moderators of guilds with the "PUBLIC" feature receive notices from Discord
	PublicUpdatesChannelID string `json:"public_updates_channel_id"`

	// The maximum amount of users in a video channel
	MaxVideoChannelUsers int `json:"max_video_channel_users"`

	// Approximate number of members in this guild, returned from the GET /guild/<id> endpoint when with_counts is true
	ApproximateMemberCount int `json:"approximate_member_count"`

	// Approximate number of non-offline members in this guild, returned from the GET /guild/<id> endpoint when with_counts is true
	ApproximatePresenceCount int `json:"approximate_presence_count"`

	// Permissions of our user
	Permissions int64 `json:"permissions,string"`

	// Stage instances in the guild
	StageInstances []*StageInstance `json:"stage_instances"`
}

// IconURL returns a URL to the guild's icon.
//
//	size:    The size of the desired icon image as a power of two
//	         Image size can be any power of two between 16 and 4096.
func (g *Guild) IconURL(size string) string {
	return discord.IconURL(g.Icon, discord.EndpointGuildIcon(g.ID, g.Icon), discord.EndpointGuildIconAnimated(g.ID, g.Icon), size)
}

// BannerURL returns a URL to the guild's banner.
//
//	size:    The size of the desired banner image as a power of two
//	         Image size can be any power of two between 16 and 4096.
func (g *Guild) BannerURL(size string) string {
	return discord.BannerURL(g.Banner, discord.EndpointGuildBanner(g.ID, g.Banner), discord.EndpointGuildBannerAnimated(g.ID, g.Banner), size)
}

// A Preview holds data related to a specific public Discord Guild, even if the user is not in the guild.
type Preview struct {
	// The ID of the guild.
	ID string `json:"id"`

	// The name of the guild. (2–100 characters)
	Name string `json:"name"`

	// The hash of the guild's icon. Use Session.GuildIcon
	// to retrieve the icon itself.
	Icon string `json:"icon"`

	// The hash of the guild's splash.
	Splash string `json:"splash"`

	// The hash of the guild's discovery splash.
	DiscoverySplash string `json:"discovery_splash"`

	// A list of the custom emojis present in the guild.
	Emojis []*user.Emoji `json:"emojis"`

	// The list of enabled guild features
	Features []string `json:"features"`

	// Approximate number of members in this guild
	// NOTE: this field is only filled when using GuildWithCounts
	ApproximateMemberCount int `json:"approximate_member_count"`

	// Approximate number of non-offline members in this guild
	// NOTE: this field is only filled when using GuildWithCounts
	ApproximatePresenceCount int `json:"approximate_presence_count"`

	// the description for the guild
	Description string `json:"description"`
}

// IconURL returns a URL to the guild's icon.
//
//	size:    The size of the desired icon image as a power of two
//	         Image size can be any power of two between 16 and 4096.
func (g *Preview) IconURL(size string) string {
	return discord.IconURL(g.Icon, discord.EndpointGuildIcon(g.ID, g.Icon), discord.EndpointGuildIconAnimated(g.ID, g.Icon), size)
}

// A Template represents a replicable template for guild creation
type Template struct {
	// The unique code for the guild template
	Code string `json:"code"`

	// The name of the template
	Name string `json:"name,omitempty"`

	// The description for the template
	Description *string `json:"description,omitempty"`

	// The number of times this template has been used
	UsageCount int `json:"usage_count"`

	// The ID of the user who created the template
	CreatorID string `json:"creator_id"`

	// The user who created the template
	Creator *user.User `json:"creator"`

	// The timestamp of when the template was created
	CreatedAt time.Time `json:"created_at"`

	// The timestamp of when the template was last synced
	UpdatedAt time.Time `json:"updated_at"`

	// The ID of the guild the template was based on
	SourceGuildID string `json:"source_guild_id"`

	// The guild 'snapshot' this template contains
	SerializedSourceGuild *Guild `json:"serialized_source_guild"`

	// Whether the template has unsynced changes
	IsDirty bool `json:"is_dirty"`
}

// TemplateParams stores the data needed to create or update a Template.
type TemplateParams struct {
	// The name of the template (1-100 characters)
	Name string `json:"name,omitempty"`
	// The description of the template (0-120 characters)
	Description string `json:"description,omitempty"`
}

// A UserGuild holds a brief version of a Guild
type UserGuild struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	Owner       bool      `json:"owner"`
	Permissions int64     `json:"permissions,string"`
	Features    []Feature `json:"features"`

	// Approximate number of members in this guild.
	// NOTE: this field is only filled when withCounts is true.
	ApproximateMemberCount int `json:"approximate_member_count"`

	// Approximate number of non-offline members in this guild.
	// NOTE: this field is only filled when withCounts is true.
	ApproximatePresenceCount int `json:"approximate_presence_count"`
}

// A Params stores all the data needed to update discord guild settings
type Params struct {
	Name                        string             `json:"name,omitempty"`
	Region                      string             `json:"region,omitempty"`
	VerificationLevel           *VerificationLevel `json:"verification_level,omitempty"`
	DefaultMessageNotifications int                `json:"default_message_notifications,omitempty"` // TODO: Separate type?
	ExplicitContentFilter       int                `json:"explicit_content_filter,omitempty"`
	AfkChannelID                string             `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int                `json:"afk_timeout,omitempty"`
	Icon                        string             `json:"icon,omitempty"`
	OwnerID                     string             `json:"owner_id,omitempty"`
	Splash                      string             `json:"splash,omitempty"`
	DiscoverySplash             string             `json:"discovery_splash,omitempty"`
	Banner                      string             `json:"banner,omitempty"`
	SystemChannelID             string             `json:"system_channel_id,omitempty"`
	SystemChannelFlags          SystemChannelFlag  `json:"system_channel_flags,omitempty"`
	RulesChannelID              string             `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      string             `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             discord.Locale     `json:"preferred_locale,omitempty"`
	Features                    []Feature          `json:"features,omitempty"`
	Description                 string             `json:"description,omitempty"`
	PremiumProgressBarEnabled   *bool              `json:"premium_progress_bar_enabled,omitempty"`
}

// A Embed stores data for a guild embed.
type Embed struct {
	Enabled   *bool  `json:"enabled,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
}

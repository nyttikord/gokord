// Package guild contains every data structures linked to guilds like... Guild or ScheduledEvent.
// It also has helping functions not using gokord.Session.
//
// Use guildapi.Requester to interact with this.
// You can get this with gokord.Session.
package guild

import (
	"time"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/status"
)

// VerificationLevel of the Guild
type VerificationLevel int

const (
	VerificationLevelNone     VerificationLevel = 0
	VerificationLevelLow      VerificationLevel = 1
	VerificationLevelMedium   VerificationLevel = 2
	VerificationLevelHigh     VerificationLevel = 3
	VerificationLevelVeryHigh VerificationLevel = 4
)

// ExplicitContentFilterLevel of the Guild
type ExplicitContentFilterLevel int

const (
	ExplicitContentFilterDisabled            ExplicitContentFilterLevel = 0
	ExplicitContentFilterMembersWithoutRoles ExplicitContentFilterLevel = 1
	ExplicitContentFilterAllMembers          ExplicitContentFilterLevel = 2
)

// NSFWLevel of the Guild
type NSFWLevel int

const (
	NSFWLevelDefault       NSFWLevel = 0
	NSFWLevelExplicit      NSFWLevel = 1
	NSFWLevelSafe          NSFWLevel = 2
	NSFWLevelAgeRestricted NSFWLevel = 3
)

// MfaLevel of the Guild
type MfaLevel int

const (
	MfaLevelNone     MfaLevel = 0
	MfaLevelElevated MfaLevel = 1
)

// PremiumTier of the Guild.
// Is the level of boosts.
type PremiumTier int

const (
	PremiumTierNone PremiumTier = 0
	PremiumTier1    PremiumTier = 1
	PremiumTier2    PremiumTier = 2
	PremiumTier3    PremiumTier = 3
)

// MessageNotifications is the notification level for a guild.
// https://discord.com/developers/docs/resources/guild#guild-object-default-message-notification-level
type MessageNotifications int

const (
	MessageNotificationsAllMessages  MessageNotifications = 0
	MessageNotificationsOnlyMentions MessageNotifications = 1
)

// SystemChannelFlag is the type of flags in the system channel.
// https://discord.com/developers/docs/resources/guild#guild-object-system-channel-flags
type SystemChannelFlag int

const (
	SystemChannelFlagsSuppressJoinNotifications          SystemChannelFlag = 1 << 0
	SystemChannelFlagsSuppressPremium                    SystemChannelFlag = 1 << 1
	SystemChannelFlagsSuppressGuildReminderNotifications SystemChannelFlag = 1 << 2
	SystemChannelFlagsSuppressJoinNotificationReplies    SystemChannelFlag = 1 << 3
)

// Feature indicates the presence of a feature in a guild.
type Feature string

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

// A Guild holds all data related to a specific Discord Guild.
// Guilds are also sometimes referred to as Servers in the Discord client.
type Guild struct {
	// The ID of the Guild.
	ID string `json:"id"`

	// The Name of the Guild (2–100 characters).
	Name string `json:"name"`

	// The hash of the Guild's Icon.
	// Use Guild.IconURL to retrieve the icon itself.
	Icon string `json:"icon"`

	// The voice Region of the Guild.
	Region string `json:"region"`

	// The ID of the AFK voice channel.Channel.
	AfkChannelID string `json:"afk_channel_id"`

	// The user.User ID of the owner of the Guild.
	OwnerID string `json:"owner_id"`

	// If we are the owner of the Guild
	Owner bool `json:"owner"`

	// The time at which the current user.User joined the Guild.
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	JoinedAt time.Time `json:"joined_at"`

	// The hash of the Guild's DiscoverySplash.
	DiscoverySplash string `json:"discovery_splash"`

	// The hash of the Guild's Splash.
	Splash string `json:"splash"`

	// The timeout, in seconds, before a user.User is considered AFK in voice.
	AfkTimeout int `json:"afk_timeout"`

	// The number of Members in the Guild.
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	MemberCount int `json:"member_count"`

	// The VerificationLevel required for the Guild.
	VerificationLevel VerificationLevel `json:"verification_level"`

	// Whether the Guild is considered large.
	//
	// This is determined by a member threshold in the identify packet, and is currently hard-coded at 250 members in
	// the library.
	Large bool `json:"large"`

	// The DefaultMessageNotifications setting for the Guild.
	DefaultMessageNotifications MessageNotifications `json:"default_message_notifications"`

	// A list of Roles in the Guild.
	Roles []*Role `json:"roles"`

	// A list of the custom Emojis present in the Guild.
	Emojis []*emoji.Emoji `json:"emojis"`

	// A list of the custom Stickers present in the Guild.
	Stickers []*emoji.Sticker `json:"stickers"`

	// A list of the Members in the Guild.
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	Members []*user.Member `json:"members"`

	// A list of partial Presences for members in the Guild.
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	Presences []*status.Presence `json:"presences"`

	// The maximum number of Presences for the Guild (the default value, currently 25000, is in effect when null is returned)
	MaxPresences int `json:"max_presences"`

	// The maximum number of Members for the Guild
	MaxMembers int `json:"max_members"`

	// A list of Channels in the Guild.
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	Channels []*channel.Channel `json:"channels"`

	// A list of all active Threads in the Guild that current user.User has permission to view.
	//
	// This field is only present in GUILD_CREATE events and websocket update events and thus is only present in
	// state-cached guilds.
	Threads []*channel.Channel `json:"threads"`

	// A list of VoiceStates for the Guild.
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	VoiceStates []*user.VoiceState `json:"voice_states"`

	// Whether this Guild is currently unavailable (most likely due to outage).
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	Unavailable bool `json:"unavailable"`

	// The ExplicitContentFilter of the Guild.
	ExplicitContentFilter ExplicitContentFilterLevel `json:"explicit_content_filter"`

	// The NSFWLevel of the Guild.
	NSFWLevel NSFWLevel `json:"nsfw_level"`

	// The list of enabled Guild Features.
	Features []Feature `json:"features"`

	// Required MfaLevel for the Guild.
	MfaLevel MfaLevel `json:"mfa_level"`

	// The application.Application ID of the Guild if bot created.
	ApplicationID string `json:"application_id"`

	// Whether the Server Widget is enabled
	WidgetEnabled bool `json:"widget_enabled"`

	// The channel.Channel ID for the Server Widget
	WidgetChannelID string `json:"widget_channel_id"`

	// The channel.Channel ID to which system messages are sent (e.g., join and leave messages)
	SystemChannelID string `json:"system_channel_id"`

	// The SystemChannelFlags for the Guild.
	SystemChannelFlags SystemChannelFlag `json:"system_channel_flags"`

	// The ID of the rules channel.Channel.
	RulesChannelID string `json:"rules_channel_id"`

	// The VanityURLCode for the Guild.
	VanityURLCode string `json:"vanity_url_code"`

	// The Description for the Guild.
	Description string `json:"description"`

	// The hash of the Guild's Banner.
	// Use Guild.BannerURL to retrieve the banner itself.
	Banner string `json:"banner"`

	// The PremiumTier of the Guild.
	PremiumTier PremiumTier `json:"premium_tier"`

	// The total number of users currently boosting this server.
	PremiumSubscriptionCount int `json:"premium_subscription_count"`

	// The preferred locale of a guild with the "PUBLIC" feature;
	// used in server discovery and notices from Discord;
	// defaults to discord.LocaleEnglishUS.
	PreferredLocale discord.Locale `json:"preferred_locale"`

	// The ID of the channel.Channel where admins and moderators of guilds with the "PUBLIC" feature receive notices
	// from Discord.
	PublicUpdatesChannelID string `json:"public_updates_channel_id"`

	// The maximum amount of users in a video channel.Channel.
	MaxVideoChannelUsers int `json:"max_video_channel_users"`

	// Approximate number of Members in this Guild.
	//
	// NOTE: returned from the GET /guild/<id> endpoint when with_counts is true.
	ApproximateMemberCount int `json:"approximate_member_count"`

	// Approximate number of non-offline Members in this Guild.
	//
	// NOTE: returned from the GET /guild/<id> endpoint when with_counts is true.
	ApproximatePresenceCount int `json:"approximate_presence_count"`

	// Permissions of our user.User.
	Permissions int64 `json:"permissions,string"`

	// StageInstances in the Guild.
	StageInstances []*channel.StageInstance `json:"stage_instances"`
}

// IconURL returns a URL to the Guild.Icon.
//
// size is the size of the desired icon image as a power of two.
// It can be any power of two between 16 and 4096.
func (g *Guild) IconURL(size string) string {
	return discord.IconURL(g.Icon, discord.EndpointGuildIcon(g.ID, g.Icon), discord.EndpointGuildIconAnimated(g.ID, g.Icon), size)
}

// BannerURL returns a URL to the Guild.Banner.
//
// size is the size of the desired icon image as a power of two.
// It can be any power of two between 16 and 4096.
func (g *Guild) BannerURL(size string) string {
	return discord.BannerURL(g.Banner, discord.EndpointGuildBanner(g.ID, g.Banner), discord.EndpointGuildBannerAnimated(g.ID, g.Banner), size)
}

// A Preview holds data related to a specific public Discord Guild, even if the user.User is not in the Guild.
type Preview struct {
	// The ID of the Guild.
	ID string `json:"id"`

	// The Name of the Guild. (2–100 characters)
	Name string `json:"name"`

	// The hash of the Guild's Icon.
	//
	// Use Preview.IconURL to retrieve the icon itself.
	Icon string `json:"icon"`

	// The hash of the Guild's Splash.
	Splash string `json:"splash"`

	// The hash of the Guild's DiscoverySplash.
	DiscoverySplash string `json:"discovery_splash"`

	// A list of the custom emojis present in the Guild.
	Emojis []*emoji.Emoji `json:"emojis"`

	// The list of enabled Guild Features
	Features []string `json:"features"`

	// Approximate number of members in this Guild.
	//
	// NOTE: this field is only filled when using gokord.Session.GetWithCounts.
	ApproximateMemberCount int `json:"approximate_member_count"`

	// Approximate number of non-offline members in this Guild.
	//
	// NOTE: this field is only filled when using GetWithCounts.
	ApproximatePresenceCount int `json:"approximate_presence_count"`

	// The Description for the Guild.
	Description string `json:"description"`
}

// IconURL returns a URL to the Preview.Icon.
//
// size is the size of the desired icon image as a power of two.
// It can be any power of two between 16 and 4096.
func (g *Preview) IconURL(size string) string {
	return discord.IconURL(g.Icon, discord.EndpointGuildIcon(g.ID, g.Icon), discord.EndpointGuildIconAnimated(g.ID, g.Icon), size)
}

// A Template represents a replicable template for Guild creation.
type Template struct {
	// The unique Code for the Guild Template.
	Code string `json:"code"`

	// The Name of the Template.
	Name string `json:"name,omitempty"`

	// The Description for the Template.
	Description *string `json:"description,omitempty"`

	// The number of times this Template has been used.
	UsageCount int `json:"usage_count"`

	// The ID of the user.User who created the Template.
	CreatorID string `json:"creator_id"`

	// The user.User who created the Template.
	Creator *user.User `json:"creator"`

	// The timestamp of when the Template was created.
	CreatedAt time.Time `json:"created_at"`

	// The timestamp of when the Template was last synced.
	UpdatedAt time.Time `json:"updated_at"`

	// The ID of the Guild the Template was based on.
	SourceGuildID string `json:"source_guild_id"`

	// The Guild 'snapshot' this Template contains.
	SerializedSourceGuild *Guild `json:"serialized_source_guild"`

	// Whether the Template has unsynced changes.
	IsDirty bool `json:"is_dirty"`
}

// TemplateParams stores the data needed to create or update a Template.
type TemplateParams struct {
	// The Name of the Template (1-100 characters)
	Name string `json:"name,omitempty"`
	// The Description of the Template (0-120 characters)
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

	// Approximate number of members in this Guild.
	//
	// NOTE: this field is only filled when withCounts is true.
	ApproximateMemberCount int `json:"approximate_member_count"`

	// Approximate number of non-offline members in this Guild.
	//
	// NOTE: this field is only filled when withCounts is true.
	ApproximatePresenceCount int `json:"approximate_presence_count"`
}

// A Params stores all the data needed to update discord Guild settings.
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

// Embed stores data for a Guild embed.
type Embed struct {
	Enabled   *bool  `json:"enabled,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
}

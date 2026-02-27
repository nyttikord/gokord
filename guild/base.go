// Package guild contains every data structures linked to guilds like... Guild or ScheduledEvent.
package guild

import (
	"context"
	"errors"
	"fmt"
	"image"
	"net/http"
	"time"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
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

// A Guild holds all data related to a specific Discord guild.
// Guilds are also sometimes referred to as Servers in the Discord client.
type Guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// The hash of the Guild's Icon.
	// Use [Guild.IconURL] to retrieve the icon itself.
	Icon         string `json:"icon"`
	Region       string `json:"region"`         // The voice region of the [Guild].
	AfkChannelID string `json:"afk_channel_id"` // The ID of the AFK voice [channel.Channel].
	OwnerID      string `json:"owner_id"`       // The [user.User] ID of the owner of the [Guild].
	Owner        bool   `json:"owner"`          // If we are the owner of the [Guild].
	// The time at which the current [user.User] joined the [Guild].
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	JoinedAt        time.Time `json:"joined_at"`
	DiscoverySplash string    `json:"discovery_splash"` // The hash of the [Guild]'s DiscoverySplash.
	Splash          string    `json:"splash"`           // The hash of the [Guild]'s Splash.
	// The timeout, in seconds, before a [user.User] is considered AFK in voice.
	AfkTimeout int `json:"afk_timeout"`
	// The number of [user.Member]s in the [Guild].
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	MemberCount       int               `json:"member_count"`
	VerificationLevel VerificationLevel `json:"verification_level"` // The [VerificationLevel] required for the [Guild].
	// Whether the [Guild] is considered large.
	//
	// This is determined by a member threshold in the identify packet, and is currently hard-coded at 250 members in
	// the library.
	Large bool `json:"large"`
	// DefaultMessageNotifications setting for the [Guild].
	DefaultMessageNotifications MessageNotifications `json:"default_message_notifications"`
	Roles                       []*Role              `json:"roles"` // A list of [Roles] in the [Guild].
	// A list of the custom [emoji.Emoji]s present in the [Guild].
	Emojis   []*emoji.Emoji   `json:"emojis"`
	Stickers []*emoji.Sticker `json:"stickers"` // A list of the custom [emoji.Sticker]s present in the [Guild].
	// A list of the [user.Member]s in the [Guild].
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	Members []*user.Member `json:"members"`
	// A list of partial [status.Presence]s of [user.Member]s in the [Guild].
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	Presences []*status.Presence `json:"presences"`
	// The maximum number of [status.Presence]s for the [Guild] (the default value, currently 25000, is in effect when
	// null is returned)
	MaxPresences int `json:"max_presences"`
	MaxMembers   int `json:"max_members"` // The maximum number of [user.Member]s for the Guild.
	// A list of [channel.Channel]s in the [Guild].
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	Channels []*channel.Channel `json:"channels"`
	// A list of all active Threads in the [Guild] that current [user.User] has permission to view.
	//
	// This field is only present in GUILD_CREATE events and websocket update events and thus is only present in
	// state-cached guilds.
	Threads []*channel.Channel `json:"threads"`
	// A list of [user.VoiceState]s for the [Guild].
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	VoiceStates []*user.VoiceState `json:"voice_states"`
	// Whether this [Guild] is currently unavailable (most likely due to outage).
	//
	// This field is only present in GUILD_CREATE events and websocket update events, and thus is only present in
	// state-cached guilds.
	Unavailable bool `json:"unavailable"`
	// ExplicitContentFilter of the [Guild].
	ExplicitContentFilter ExplicitContentFilterLevel `json:"explicit_content_filter"`
	NSFWLevel             NSFWLevel                  `json:"nsfw_level"` // The NSFWLevel of the Guild.
	// The list of enabled [Guild] [Feature]s.
	Features []Feature `json:"features"`
	MfaLevel MfaLevel  `json:"mfa_level"` // Required [MfaLevel] for the [Guild].
	// The [application.Application] ID of the Guild if bot created.
	ApplicationID   string `json:"application_id"`
	WidgetEnabled   bool   `json:"widget_enabled"`    // Whether the Server Widget is enabled
	WidgetChannelID string `json:"widget_channel_id"` // The [channel.Channel] ID for the Server Widget
	// The [channel.Channel] ID to which system messages are sent (e.g., join and leave messages)
	SystemChannelID    string            `json:"system_channel_id"`
	SystemChannelFlags SystemChannelFlag `json:"system_channel_flags"` // [SystemChannelFlag]s for the [Guild].
	RulesChannelID     string            `json:"rules_channel_id"`     // The ID of the rules [channel.Channel].
	VanityURLCode      string            `json:"vanity_url_code"`      // The VanityURLCode for the [Guild].
	Description        string            `json:"description"`
	// The hash of the [Guild]'s Banner.
	// Use [Guild.BannerURL] to retrieve the banner itself.
	Banner      string      `json:"banner"`
	PremiumTier PremiumTier `json:"premium_tier"` // [PremiumTier] of the [Guild].
	// The total number of users currently boosting this server.
	PremiumSubscriptionCount int `json:"premium_subscription_count"`
	// The preferred [discord.Locale] of a guild with the "PUBLIC" feature;
	// used in server discovery and notices from Discord;
	// defaults to [discord.LocaleEnglishUS].
	PreferredLocale discord.Locale `json:"preferred_locale"`
	// The ID of the [channel.Channel] where admins and moderators of guilds with the "PUBLIC" feature receive notices
	// from Discord.
	PublicUpdatesChannelID string `json:"public_updates_channel_id"`
	// The maximum amount of [user.Member]s in a video [channel.Channel].
	MaxVideoChannelUsers int `json:"max_video_channel_users"`
	// Approximate number of [user.Member]s in this [Guild].
	//
	// NOTE: returned with [GetWithCounts].
	ApproximateMemberCount int `json:"approximate_member_count"`
	// Approximate number of non-offline [user.Member]s in this [Guild].
	//
	// NOTE: returned with [GetWithCounts].
	ApproximatePresenceCount int   `json:"approximate_presence_count"`
	Permissions              int64 `json:"permissions,string"` // Permissions of our [user.User].
	// [channel.StageInstance]s in the [Guild].
	StageInstances []*channel.StageInstance `json:"stage_instances"`
}

// IconURL returns an URL to the [Guild.Icon].
//
// size is the size of the desired icon image as a power of two.
// It can be any power of two between 16 and 4096.
func (g *Guild) IconURL(size string) string {
	return discord.IconURL(g.Icon, discord.EndpointGuildIcon(g.ID, g.Icon), discord.EndpointGuildIconAnimated(g.ID, g.Icon), size)
}

// BannerURL returns an URL to the [Guild.Banner].
//
// size is the size of the desired icon image as a power of two.
// It can be any power of two between 16 and 4096.
func (g *Guild) BannerURL(size string) string {
	return discord.BannerURL(g.Banner, discord.EndpointGuildBanner(g.ID, g.Banner), discord.EndpointGuildBannerAnimated(g.ID, g.Banner), size)
}

// A Params stores all the data needed to update discord Guild settings.
type Params struct {
	Name                        string               `json:"name,omitempty"`
	Region                      string               `json:"region,omitempty"`
	VerificationLevel           *VerificationLevel   `json:"verification_level,omitempty"`
	DefaultMessageNotifications MessageNotifications `json:"default_message_notifications,omitempty"` // TODO: Separate type?
	ExplicitContentFilter       int                  `json:"explicit_content_filter,omitempty"`
	AfkChannelID                string               `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int                  `json:"afk_timeout,omitempty"`
	Icon                        string               `json:"icon,omitempty"`
	OwnerID                     string               `json:"owner_id,omitempty"`
	Splash                      string               `json:"splash,omitempty"`
	DiscoverySplash             string               `json:"discovery_splash,omitempty"`
	Banner                      string               `json:"banner,omitempty"`
	SystemChannelID             string               `json:"system_channel_id,omitempty"`
	SystemChannelFlags          SystemChannelFlag    `json:"system_channel_flags,omitempty"`
	RulesChannelID              string               `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      string               `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             discord.Locale       `json:"preferred_locale,omitempty"`
	Features                    []Feature            `json:"features,omitempty"`
	Description                 string               `json:"description,omitempty"`
	PremiumProgressBarEnabled   *bool                `json:"premium_progress_bar_enabled,omitempty"`
}

// Embed stores data for a [Guild] embed.
type Embed struct {
	Enabled   *bool  `json:"enabled,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
}

var (
	ErrInvalidVoiceRegions = errors.New("invalid voice regions")
	ErrGuildNoIcon         = errors.New("guild does not have an icon set")
	ErrGuildNoSplash       = errors.New("guild does not have a splash set")
)

// Get returns the [Guild] with the given guildID.
func Get(guildID string) Request[*Guild] {
	return NewData[*Guild](http.MethodGet, discord.EndpointGuild(guildID))
}

// GetWithCounts returns the guild.Guild with the given guildID with approximate user.Member and status.Presence counts.
func GetWithCounts(guildID string) Request[*Guild] {
	return NewData[*Guild](http.MethodGet, discord.EndpointGuild(guildID)+"?with_counts=true")
}

// Edit a [Guild] with the given params.
func Edit(guildID string, params *Params) Request[*Guild] {
	return NewData[*Guild](http.MethodPatch, discord.EndpointGuild(guildID)).
		WithData(params).
		WithPre(func(ctx context.Context, do *Do) error {
			if params.Region == "" {
				return nil
			}
			valid := false
			regions, err := ListVoiceRegions().Do(ctx)
			if err != nil {
				return err
			}
			for _, r := range regions {
				if params.Region == r.ID {
					valid = true
				}
			}
			if valid {
				return nil
			}
			var validRegions []string
			for _, r := range regions {
				validRegions = append(validRegions, r.ID)
			}
			return errors.Join(
				ErrInvalidVoiceRegions, fmt.Errorf("%s is not a voice region (%q)", params.Region, validRegions),
			)
		})
}

// Delete a [Guild].
func Delete(guildID string) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointGuild(guildID))
	return WrapAsEmpty(req)
}

// Leave a [Guild].
func Leave(guildID string) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointUserGuild("@e", guildID)).
		WithBucketID(discord.EndpointUserGuild("", guildID))
	return WrapAsEmpty(req)
}

// ListVoiceRegions returns available [discord.VoiceRegion].
func ListVoiceRegions() Request[[]*discord.VoiceRegion] {
	return NewData[[]*discord.VoiceRegion](http.MethodGet, discord.EndpointVoiceRegions)
}

// GetIcon returns an [image.Image] of a [Guild.GetIcon].
func GetIcon(guildID string) Request[image.Image] {
	return NewImage(http.MethodGet, "").
		WithBucketID(discord.EndpointGuildIcon(guildID, "")).
		WithPre(func(ctx context.Context, do *Do) error {
			g, err := Get(guildID).Do(ctx)
			if err != nil {
				return err
			}
			if g.Icon == "" {
				return ErrGuildNoIcon
			}
			do.Endpoint = discord.EndpointGuildIcon(guildID, g.Icon)
			return nil
		})
}

// GetSplash returns an [image.Image] of a [Guild.GetSplash].
func GetSplash(guildID string) Request[image.Image] {
	return NewImage(http.MethodGet, "").
		WithBucketID(discord.EndpointGuildSplash(guildID, "")).
		WithPre(func(ctx context.Context, do *Do) error {
			g, err := Get(guildID).Do(ctx)
			if err != nil {
				return err
			}
			if g.Splash == "" {
				return ErrGuildNoSplash
			}
			do.Endpoint = discord.EndpointGuildSplash(guildID, g.Splash)
			return nil
		})
}

// GetEmbed returns the [Embed] for a [Guild].
func GetEmbed(guildID string) Request[*Embed] {
	return NewData[*Embed](http.MethodGet, discord.EndpointGuildEmbed(guildID))
}

// UpdateEmbed of a [Guild].
func UpdateEmbed(guildID string, data *Embed) Empty {
	req := NewSimple(http.MethodPatch, discord.EndpointGuildEmbed(guildID)).WithData(data)
	return WrapAsEmpty(req)
}

/*
// ListInvites returns the list of [invite.Invite] for the given [Guild].
func ListInvites(guildID string) Request[[]*invite.Invite] {
	return NewData[[]*invite.Invite](http.MethodGet, discord.EndpointGuildInvites(guildID))
}
*/

// ListThreadsActive returns all active threads in the given [guild.Guild].
func ListThreadsActive(guildID string) Request[*channel.ThreadsList] {
	return NewData[*channel.ThreadsList](http.MethodGet, discord.EndpointGuildActiveThreads(guildID))
}

// ListWebhooks returns all [channel.Webhook] for a given [guild.Guild].
func ListWebhooks(guildID string) Request[[]*channel.Webhook] {
	return NewData[[]*channel.Webhook](http.MethodGet, discord.EndpointGuildWebhooks(guildID))
}

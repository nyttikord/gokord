package gokord

import (
	"fmt"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/logger"
	"github.com/nyttikord/gokord/user"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// A Session represents a connection to the Discord API.
type Session struct {
	sync.RWMutex

	// General configurable settings.

	// Authentication token for this session
	// TODO: Remove Below, Deprecated, Use Identify struct
	Token string

	MFA bool

	LogLevel logger.Level

	// Should the session reconnect the websocket on errors.
	ShouldReconnectOnError bool

	// Should voice connections reconnect on a session reconnect.
	ShouldReconnectVoiceOnSessionError bool

	// Should the session retry requests when rate limited.
	ShouldRetryOnRateLimit bool

	// Identify is sent during initial handshake with the discord gateway.
	// https://discord.com/developers/docs/topics/gateway#identify
	Identify Identify

	// TODO: Remove Below, Deprecated, Use Identify struct
	// Should the session request compressed websocket data.
	Compress bool

	// Sharding
	ShardID    int
	ShardCount int

	// Should state tracking be enabled.
	// State tracking is the best way for getting the users
	// active guilds and the members of the guilds.
	StateEnabled bool

	// Whether or not to call event handlers synchronously.
	// e.g. false = launch event handlers in their own goroutines.
	SyncEvents bool

	// Exposed but should not be modified by User.

	// Whether the Data Websocket is ready
	DataReady bool // NOTE: Maye be deprecated soon

	// Max number of REST API retries
	MaxRestRetries int

	// Status stores the current status of the websocket connection
	// this is being tested, may stay, may go away.
	status int32

	// Whether the Voice Websocket is ready
	VoiceReady bool // NOTE: Deprecated.

	// Whether the UDP Connection is ready
	UDPReady bool // NOTE: Deprecated

	// Stores a mapping of guild id's to VoiceConnections
	VoiceConnections map[string]*VoiceConnection

	// Managed state object, updated internally with events when
	// StateEnabled is true.
	State *State

	// The http client used for REST requests
	Client *http.Client

	// The dialer used for WebSocket connection
	Dialer *websocket.Dialer

	// The user agent used for REST APIs
	UserAgent string

	// Stores the last HeartbeatAck that was received (in UTC)
	LastHeartbeatAck time.Time

	// Stores the last Heartbeat sent (in UTC)
	LastHeartbeatSent time.Time

	// used to deal with rate limits
	Ratelimiter *RateLimiter

	// Event handlers
	handlersMu   sync.RWMutex
	handlers     map[string][]*eventHandlerInstance
	onceHandlers map[string][]*eventHandlerInstance

	// The websocket connection.
	wsConn *websocket.Conn

	// When nil, the session is not listening.
	listening chan interface{}

	// sequence tracks the current gateway api websocket sequence number
	sequence *int64

	// stores sessions current Discord Resume Gateway
	resumeGatewayURL string

	// stores sessions current Discord Gateway
	gateway string

	// stores session ID of current Gateway connection
	sessionID string

	// used to make sure gateway websocket writes do not happen concurrently
	wsMutex sync.Mutex
}

// ApplicationIntegrationType dictates where application can be installed and its available interaction contexts.
type ApplicationIntegrationType uint

const (
	// ApplicationIntegrationGuildInstall indicates that app is installable to guilds.
	ApplicationIntegrationGuildInstall ApplicationIntegrationType = 0
	// ApplicationIntegrationUserInstall indicates that app is installable to users.
	ApplicationIntegrationUserInstall ApplicationIntegrationType = 1
)

// ApplicationInstallParams represents application's installation parameters
// for default in-app oauth2 authorization link.
type ApplicationInstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions int64    `json:"permissions,string"`
}

// ApplicationIntegrationTypeConfig represents application's configuration for a particular integration type.
type ApplicationIntegrationTypeConfig struct {
	OAuth2InstallParams *ApplicationInstallParams `json:"oauth2_install_params,omitempty"`
}

// Application stores values for a Discord Application
type Application struct {
	ID                     string                                                           `json:"id,omitempty"`
	Name                   string                                                           `json:"name"`
	Icon                   string                                                           `json:"icon,omitempty"`
	Description            string                                                           `json:"description,omitempty"`
	RPCOrigins             []string                                                         `json:"rpc_origins,omitempty"`
	BotPublic              bool                                                             `json:"bot_public,omitempty"`
	BotRequireCodeGrant    bool                                                             `json:"bot_require_code_grant,omitempty"`
	TermsOfServiceURL      string                                                           `json:"terms_of_service_url"`
	PrivacyProxyURL        string                                                           `json:"privacy_policy_url"`
	Owner                  *user.User                                                       `json:"owner"`
	Summary                string                                                           `json:"summary"`
	VerifyKey              string                                                           `json:"verify_key"`
	Team                   *Team                                                            `json:"team"`
	GuildID                string                                                           `json:"guild_id"`
	PrimarySKUID           string                                                           `json:"primary_sku_id"`
	Slug                   string                                                           `json:"slug"`
	CoverImage             string                                                           `json:"cover_image"`
	Flags                  int                                                              `json:"flags,omitempty"`
	IntegrationTypesConfig map[ApplicationIntegrationType]*ApplicationIntegrationTypeConfig `json:"integration_types,omitempty"`
}

// ApplicationRoleConnectionMetadataType represents the type of application role connection metadata.
type ApplicationRoleConnectionMetadataType int

// Application role connection metadata types.
const (
	ApplicationRoleConnectionMetadataIntegerLessThanOrEqual     ApplicationRoleConnectionMetadataType = 1
	ApplicationRoleConnectionMetadataIntegerGreaterThanOrEqual  ApplicationRoleConnectionMetadataType = 2
	ApplicationRoleConnectionMetadataIntegerEqual               ApplicationRoleConnectionMetadataType = 3
	ApplicationRoleConnectionMetadataIntegerNotEqual            ApplicationRoleConnectionMetadataType = 4
	ApplicationRoleConnectionMetadataDatetimeLessThanOrEqual    ApplicationRoleConnectionMetadataType = 5
	ApplicationRoleConnectionMetadataDatetimeGreaterThanOrEqual ApplicationRoleConnectionMetadataType = 6
	ApplicationRoleConnectionMetadataBooleanEqual               ApplicationRoleConnectionMetadataType = 7
	ApplicationRoleConnectionMetadataBooleanNotEqual            ApplicationRoleConnectionMetadataType = 8
)

// ApplicationRoleConnectionMetadata stores application role connection metadata.
type ApplicationRoleConnectionMetadata struct {
	Type                     ApplicationRoleConnectionMetadataType `json:"type"`
	Key                      string                                `json:"key"`
	Name                     string                                `json:"name"`
	NameLocalizations        map[discord.Locale]string             `json:"name_localizations"`
	Description              string                                `json:"description"`
	DescriptionLocalizations map[discord.Locale]string             `json:"description_localizations"`
}

// ApplicationRoleConnection represents the role connection that an application has attached to a user.
type ApplicationRoleConnection struct {
	PlatformName     string            `json:"platform_name"`
	PlatformUsername string            `json:"platform_username"`
	Metadata         map[string]string `json:"metadata"`
}

// UserConnection is a Connection returned from the UserConnections endpoint
type UserConnection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Revoked      bool           `json:"revoked"`
	Integrations []*Integration `json:"integrations"`
}

// Integration stores integration information
type Integration struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Type              string             `json:"type"`
	Enabled           bool               `json:"enabled"`
	Syncing           bool               `json:"syncing"`
	RoleID            string             `json:"role_id"`
	EnableEmoticons   bool               `json:"enable_emoticons"`
	ExpireBehavior    ExpireBehavior     `json:"expire_behavior"`
	ExpireGracePeriod int                `json:"expire_grace_period"`
	User              *user.User         `json:"user"`
	Account           IntegrationAccount `json:"account"`
	SyncedAt          time.Time          `json:"synced_at"`
}

// ExpireBehavior of Integration
// https://discord.com/developers/docs/resources/guild#integration-object-integration-expire-behaviors
type ExpireBehavior int

// Block of valid ExpireBehaviors
const (
	ExpireBehaviorRemoveRole ExpireBehavior = 0
	ExpireBehaviorKick       ExpireBehavior = 1
)

// IntegrationAccount is integration account information
// sent by the UserConnections endpoint
type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// A VoiceRegion stores data for a specific voice region server.
// https://discord.com/developers/docs/resources/voice#voice-region-object
type VoiceRegion struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Optimal    bool   `json:"optimal"`
	Deprecated bool   `json:"deprecated"`
	Custom     bool   `json:"custom"`
}

// InviteTargetType indicates the type of target of an invite
// https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
type InviteTargetType uint8

// Invite target types
const (
	InviteTargetStream              InviteTargetType = 1
	InviteTargetEmbeddedApplication InviteTargetType = 2
)

// A Invite stores all data related to a specific Discord Guild or Channel invite.
type Invite struct {
	Guild             *Guild           `json:"guild"`
	Channel           *Channel         `json:"channel"`
	Inviter           *user.User       `json:"inviter"`
	Code              string           `json:"code"`
	CreatedAt         time.Time        `json:"created_at"`
	MaxAge            int              `json:"max_age"`
	Uses              int              `json:"uses"`
	MaxUses           int              `json:"max_uses"`
	Revoked           bool             `json:"revoked"`
	Temporary         bool             `json:"temporary"`
	Unique            bool             `json:"unique"`
	TargetUser        *user.User       `json:"target_user"`
	TargetType        InviteTargetType `json:"target_type"`
	TargetApplication *Application     `json:"target_application"`

	// will only be filled when using InviteWithCounts
	ApproximatePresenceCount int `json:"approximate_presence_count"`
	ApproximateMemberCount   int `json:"approximate_member_count"`

	ExpiresAt *time.Time `json:"expires_at"`
}

// A TooManyRequests struct holds information received from Discord
// when receiving a HTTP 429 response.
type TooManyRequests struct {
	Bucket     string        `json:"bucket"`
	Message    string        `json:"message"`
	RetryAfter time.Duration `json:"retry_after"`
}

// UnmarshalJSON helps support translation of a milliseconds-based float
// into a time.Duration on TooManyRequests.
func (t *TooManyRequests) UnmarshalJSON(b []byte) error {
	u := struct {
		Bucket     string  `json:"bucket"`
		Message    string  `json:"message"`
		RetryAfter float64 `json:"retry_after"`
	}{}
	err := Unmarshal(b, &u)
	if err != nil {
		return err
	}

	t.Bucket = u.Bucket
	t.Message = u.Message
	whole, frac := math.Modf(u.RetryAfter)
	t.RetryAfter = time.Duration(whole)*time.Second + time.Duration(frac*1000)*time.Millisecond
	return nil
}

// A ReadState stores data on the read state of channels.
type ReadState struct {
	MentionCount  int    `json:"mention_count"`
	LastMessageID string `json:"last_message_id"`
	ID            string `json:"id"`
}

// An APIErrorMessage is an api error message returned from discord
type APIErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MessageReaction stores the data for a message reaction.
type MessageReaction struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Emoji     Emoji  `json:"emoji"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
}

// GatewayBotResponse stores the data for the gateway/bot response
type GatewayBotResponse struct {
	URL               string             `json:"url"`
	Shards            int                `json:"shards"`
	SessionStartLimit SessionInformation `json:"session_start_limit"`
}

// SessionInformation provides the information for max concurrency sharding
type SessionInformation struct {
	Total          int `json:"total,omitempty"`
	Remaining      int `json:"remaining,omitempty"`
	ResetAfter     int `json:"reset_after,omitempty"`
	MaxConcurrency int `json:"max_concurrency,omitempty"`
}

// GatewayStatusUpdate is sent by the client to indicate a presence or status update
// https://discord.com/developers/docs/topics/gateway#update-status-gateway-status-update-structure
type GatewayStatusUpdate struct {
	Since  int      `json:"since"`
	Game   Activity `json:"game"`
	Status string   `json:"status"`
	AFK    bool     `json:"afk"`
}

// Identify is sent during initial handshake with the discord gateway.
// https://discord.com/developers/docs/topics/gateway#identify
type Identify struct {
	Token          string              `json:"token"`
	Properties     IdentifyProperties  `json:"properties"`
	Compress       bool                `json:"compress"`
	LargeThreshold int                 `json:"large_threshold"`
	Shard          *[2]int             `json:"shard,omitempty"`
	Presence       GatewayStatusUpdate `json:"presence,omitempty"`
	Intents        discord.Intent      `json:"intents"`
}

// IdentifyProperties contains the "properties" portion of an Identify packet
// https://discord.com/developers/docs/topics/gateway#identify-identify-connection-properties
type IdentifyProperties struct {
	OS              string `json:"$os"`
	Browser         string `json:"$browser"`
	Device          string `json:"$device"`
	Referer         string `json:"$referer"`
	ReferringDomain string `json:"$referring_domain"`
}

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

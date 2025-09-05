package gokord

import (
	"encoding/json"
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

// ChannelType is the type of a Channel
type ChannelType int

// Block contains known ChannelType values
const (
	ChannelTypeGuildText          ChannelType = 0
	ChannelTypeDM                 ChannelType = 1
	ChannelTypeGuildVoice         ChannelType = 2
	ChannelTypeGroupDM            ChannelType = 3
	ChannelTypeGuildCategory      ChannelType = 4
	ChannelTypeGuildNews          ChannelType = 5
	ChannelTypeGuildStore         ChannelType = 6
	ChannelTypeGuildNewsThread    ChannelType = 10
	ChannelTypeGuildPublicThread  ChannelType = 11
	ChannelTypeGuildPrivateThread ChannelType = 12
	ChannelTypeGuildStageVoice    ChannelType = 13
	ChannelTypeGuildDirectory     ChannelType = 14
	ChannelTypeGuildForum         ChannelType = 15
	ChannelTypeGuildMedia         ChannelType = 16
)

// ChannelFlags represent flags of a channel/thread.
type ChannelFlags int

// Block containing known ChannelFlags values.
const (
	// ChannelFlagPinned indicates whether the thread is pinned in the forum channel.
	// NOTE: forum threads only.
	ChannelFlagPinned ChannelFlags = 1 << 1
	// ChannelFlagRequireTag indicates whether a tag is required to be specified when creating a thread.
	// NOTE: forum channels only.
	ChannelFlagRequireTag ChannelFlags = 1 << 4
)

// ForumSortOrderType represents sort order of a forum channel.
type ForumSortOrderType int

const (
	// ForumSortOrderLatestActivity sorts posts by activity.
	ForumSortOrderLatestActivity ForumSortOrderType = 0
	// ForumSortOrderCreationDate sorts posts by creation time (from most recent to oldest).
	ForumSortOrderCreationDate ForumSortOrderType = 1
)

// ForumLayout represents layout of a forum channel.
type ForumLayout int

const (
	// ForumLayoutNotSet represents no default layout.
	ForumLayoutNotSet ForumLayout = 0
	// ForumLayoutListView displays forum posts as a list.
	ForumLayoutListView ForumLayout = 1
	// ForumLayoutGalleryView displays forum posts as a collection of tiles.
	ForumLayoutGalleryView ForumLayout = 2
)

// A Channel holds all data related to an individual Discord channel.
type Channel struct {
	// The ID of the channel.
	ID string `json:"id"`

	// The ID of the guild to which the channel belongs, if it is in a guild.
	// Else, this ID is empty (e.g. DM channels).
	GuildID string `json:"guild_id"`

	// The name of the channel.
	Name string `json:"name"`

	// The topic of the channel.
	Topic string `json:"topic"`

	// The type of the channel.
	Type ChannelType `json:"type"`

	// The ID of the last message sent in the channel. This is not
	// guaranteed to be an ID of a valid message.
	LastMessageID string `json:"last_message_id"`

	// The timestamp of the last pinned message in the channel.
	// nil if the channel has no pinned messages.
	LastPinTimestamp *time.Time `json:"last_pin_timestamp"`

	// An approximate count of messages in a thread, stops counting at 50
	MessageCount int `json:"message_count"`
	// An approximate count of users in a thread, stops counting at 50
	MemberCount int `json:"member_count"`

	// Whether the channel is marked as NSFW.
	NSFW bool `json:"nsfw"`

	// Icon of the group DM channel.
	Icon string `json:"icon"`

	// The position of the channel, used for sorting in client.
	Position int `json:"position"`

	// The bitrate of the channel, if it is a voice channel.
	Bitrate int `json:"bitrate"`

	// The recipients of the channel. This is only populated in DM channels.
	Recipients []*user.User `json:"recipients"`

	// The messages in the channel. This is only present in state-cached channels,
	// and State.MaxMessageCount must be non-zero.
	Messages []*Message `json:"-"`

	// A list of permission overwrites present for the channel.
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites"`

	// The user limit of the voice channel.
	UserLimit int `json:"user_limit"`

	// The ID of the parent channel, if the channel is under a category. For threads - id of the channel thread was created in.
	ParentID string `json:"parent_id"`

	// Amount of seconds a user has to wait before sending another message or creating another thread (0-21600)
	// bots, as well as users with the permission manage_messages or manage_channel, are unaffected
	RateLimitPerUser int `json:"rate_limit_per_user"`

	// ID of the creator of the group DM or thread
	OwnerID string `json:"owner_id"`

	// ApplicationID of the DM creator Zeroed if guild channel or not a bot user
	ApplicationID string `json:"application_id"`

	// Thread-specific fields not needed by other channels
	ThreadMetadata *ThreadMetadata `json:"thread_metadata,omitempty"`
	// Thread member object for the current user, if they have joined the thread, only included on certain API endpoints
	Member *ThreadMember `json:"thread_member"`

	// All thread members. State channels only.
	Members []*ThreadMember `json:"-"`

	// Channel flags.
	Flags ChannelFlags `json:"flags"`

	// The set of tags that can be used in a forum channel.
	AvailableTags []ForumTag `json:"available_tags"`

	// The IDs of the set of tags that have been applied to a thread in a forum channel.
	AppliedTags []string `json:"applied_tags"`

	// Emoji to use as the default reaction to a forum post.
	DefaultReactionEmoji ForumDefaultReaction `json:"default_reaction_emoji"`

	// The initial RateLimitPerUser to set on newly created threads in a channel.
	// This field is copied to the thread at creation time and does not live update.
	DefaultThreadRateLimitPerUser int `json:"default_thread_rate_limit_per_user"`

	// The default sort order type used to order posts in forum channels.
	// Defaults to null, which indicates a preferred sort order hasn't been set by a channel admin.
	DefaultSortOrder *ForumSortOrderType `json:"default_sort_order"`

	// The default forum layout view used to display posts in forum channels.
	// Defaults to ForumLayoutNotSet, which indicates a layout view has not been set by a channel admin.
	DefaultForumLayout ForumLayout `json:"default_forum_layout"`
}

// Mention returns a string which mentions the channel
func (c *Channel) Mention() string {
	return fmt.Sprintf("<#%s>", c.ID)
}

// IsThread is a helper function to determine if channel is a thread or not
func (c *Channel) IsThread() bool {
	return c.Type == ChannelTypeGuildPublicThread || c.Type == ChannelTypeGuildPrivateThread || c.Type == ChannelTypeGuildNewsThread
}

// A ChannelEdit holds Channel Field data for a channel edit.
type ChannelEdit struct {
	Name                          string                 `json:"name,omitempty"`
	Topic                         string                 `json:"topic,omitempty"`
	NSFW                          *bool                  `json:"nsfw,omitempty"`
	Position                      *int                   `json:"position,omitempty"`
	Bitrate                       int                    `json:"bitrate,omitempty"`
	UserLimit                     int                    `json:"user_limit,omitempty"`
	PermissionOverwrites          []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      string                 `json:"parent_id,omitempty"`
	RateLimitPerUser              *int                   `json:"rate_limit_per_user,omitempty"`
	Flags                         *ChannelFlags          `json:"flags,omitempty"`
	DefaultThreadRateLimitPerUser *int                   `json:"default_thread_rate_limit_per_user,omitempty"`

	// NOTE: threads only

	Archived            *bool `json:"archived,omitempty"`
	AutoArchiveDuration int   `json:"auto_archive_duration,omitempty"`
	Locked              *bool `json:"locked,omitempty"`
	Invitable           *bool `json:"invitable,omitempty"`

	// NOTE: forum channels only

	AvailableTags        *[]ForumTag           `json:"available_tags,omitempty"`
	DefaultReactionEmoji *ForumDefaultReaction `json:"default_reaction_emoji,omitempty"`
	DefaultSortOrder     *ForumSortOrderType   `json:"default_sort_order,omitempty"` // TODO: null
	DefaultForumLayout   *ForumLayout          `json:"default_forum_layout,omitempty"`

	// NOTE: forum threads only
	AppliedTags *[]string `json:"applied_tags,omitempty"`
}

// A ChannelFollow holds data returned after following a news channel
type ChannelFollow struct {
	ChannelID string `json:"channel_id"`
	WebhookID string `json:"webhook_id"`
}

// PermissionOverwriteType represents the type of resource on which
// a permission overwrite acts.
type PermissionOverwriteType int

// The possible permission overwrite types.
const (
	PermissionOverwriteTypeRole   PermissionOverwriteType = 0
	PermissionOverwriteTypeMember PermissionOverwriteType = 1
)

// A PermissionOverwrite holds permission overwrite data for a Channel
type PermissionOverwrite struct {
	ID    string                  `json:"id"`
	Type  PermissionOverwriteType `json:"type"`
	Deny  int64                   `json:"deny,string"`
	Allow int64                   `json:"allow,string"`
}

// ThreadStart stores all parameters you can use with MessageThreadStartComplex or ThreadStartComplex
type ThreadStart struct {
	Name                string      `json:"name"`
	AutoArchiveDuration int         `json:"auto_archive_duration,omitempty"`
	Type                ChannelType `json:"type,omitempty"`
	Invitable           bool        `json:"invitable"`
	RateLimitPerUser    int         `json:"rate_limit_per_user,omitempty"`

	// NOTE: forum threads only
	AppliedTags []string `json:"applied_tags,omitempty"`
}

// ThreadMetadata contains a number of thread-specific channel fields that are not needed by other channel types.
type ThreadMetadata struct {
	// Whether the thread is archived
	Archived bool `json:"archived"`
	// Duration in minutes to automatically archive the thread after recent activity, can be set to: 60, 1440, 4320, 10080
	AutoArchiveDuration int `json:"auto_archive_duration"`
	// Timestamp when the thread's archive status was last changed, used for calculating recent activity
	ArchiveTimestamp time.Time `json:"archive_timestamp"`
	// Whether the thread is locked; when a thread is locked, only users with MANAGE_THREADS can unarchive it
	Locked bool `json:"locked"`
	// Whether non-moderators can add other non-moderators to a thread; only available on private threads
	Invitable bool `json:"invitable"`
}

// ThreadMember is used to indicate whether a user has joined a thread or not.
// NOTE: ID and UserID are empty (omitted) on the member sent within each thread in the GUILD_CREATE event.
type ThreadMember struct {
	// The id of the thread
	ID string `json:"id,omitempty"`
	// The id of the user
	UserID string `json:"user_id,omitempty"`
	// The time the current user last joined the thread
	JoinTimestamp time.Time `json:"join_timestamp"`
	// Any user-thread settings, currently only used for notifications
	Flags int `json:"flags"`
	// Additional information about the user.
	// NOTE: only present if the withMember parameter is set to true
	// when calling Session.ThreadMembers or Session.ThreadMember.
	Member *Member `json:"member,omitempty"`
}

// ThreadsList represents a list of threads alongisde with thread member objects for the current user.
type ThreadsList struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

// AddedThreadMember holds information about the user who was added to the thread
type AddedThreadMember struct {
	*ThreadMember
	Member   *Member   `json:"member"`
	Presence *Presence `json:"presence"`
}

// ForumDefaultReaction specifies emoji to use as the default reaction to a forum post.
// NOTE: Exactly one of EmojiID and EmojiName must be set.
type ForumDefaultReaction struct {
	// The id of a guild's custom emoji.
	EmojiID string `json:"emoji_id,omitempty"`
	// The unicode character of the emoji.
	EmojiName string `json:"emoji_name,omitempty"`
}

// ForumTag represents a tag that is able to be applied to a thread in a forum channel.
type ForumTag struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Moderated bool   `json:"moderated"`
	EmojiID   string `json:"emoji_id,omitempty"`
	EmojiName string `json:"emoji_name,omitempty"`
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
	Intents        Intent              `json:"intents"`
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

// StageInstance holds information about a live stage.
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-resource
type StageInstance struct {
	// The id of this Stage instance
	ID string `json:"id"`
	// The guild id of the associated Stage channel
	GuildID string `json:"guild_id"`
	// The id of the associated Stage channel
	ChannelID string `json:"channel_id"`
	// The topic of the Stage instance (1-120 characters)
	Topic string `json:"topic"`
	// The privacy level of the Stage instance
	// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
	PrivacyLevel StageInstancePrivacyLevel `json:"privacy_level"`
	// Whether or not Stage Discovery is disabled (deprecated)
	DiscoverableDisabled bool `json:"discoverable_disabled"`
	// The id of the scheduled event for this Stage instance
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
}

// StageInstanceParams represents the parameters needed to create or edit a stage instance
type StageInstanceParams struct {
	// ChannelID represents the id of the Stage channel
	ChannelID string `json:"channel_id,omitempty"`
	// Topic of the Stage instance (1-120 characters)
	Topic string `json:"topic,omitempty"`
	// PrivacyLevel of the Stage instance (default GUILD_ONLY)
	PrivacyLevel StageInstancePrivacyLevel `json:"privacy_level,omitempty"`
	// SendStartNotification will notify @everyone that a Stage instance has started
	SendStartNotification bool `json:"send_start_notification,omitempty"`
}

// StageInstancePrivacyLevel represents the privacy level of a Stage instance
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
type StageInstancePrivacyLevel int

const (
	// StageInstancePrivacyLevelPublic The Stage instance is visible publicly. (deprecated)
	StageInstancePrivacyLevelPublic StageInstancePrivacyLevel = 1
	// StageInstancePrivacyLevelGuildOnly The Stage instance is visible to only guild members.
	StageInstancePrivacyLevelGuildOnly StageInstancePrivacyLevel = 2
)

// PollLayoutType represents the layout of a poll.
type PollLayoutType int

// Valid PollLayoutType values.
const (
	PollLayoutTypeDefault PollLayoutType = 1
)

// PollMedia contains common data used by question and answers.
type PollMedia struct {
	Text  string          `json:"text,omitempty"`
	Emoji *ComponentEmoji `json:"emoji,omitempty"` // TODO: rename the type
}

// PollAnswer represents a single answer in a poll.
type PollAnswer struct {
	// NOTE: should not be set on creation.
	AnswerID int        `json:"answer_id,omitempty"`
	Media    *PollMedia `json:"poll_media"`
}

// PollAnswerCount stores counted poll votes for a single answer.
type PollAnswerCount struct {
	ID      int  `json:"id"`
	Count   int  `json:"count"`
	MeVoted bool `json:"me_voted"`
}

// PollResults contains voting results on a poll.
type PollResults struct {
	Finalized    bool               `json:"is_finalized"`
	AnswerCounts []*PollAnswerCount `json:"answer_counts"`
}

// Poll contains all poll related data.
type Poll struct {
	Question         PollMedia      `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	AllowMultiselect bool           `json:"allow_multiselect"`
	LayoutType       PollLayoutType `json:"layout_type,omitempty"`

	// NOTE: should be set only on creation, when fetching use Expiry.
	Duration int `json:"duration,omitempty"`

	// NOTE: available only when fetching.

	Results *PollResults `json:"results,omitempty"`
	// NOTE: as Discord documentation notes, this field might be null even when fetching.
	Expiry *time.Time `json:"expiry,omitempty"`
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

// Constants for the different bit offsets of text channel permissions
const (
	// Deprecated: PermissionReadMessages has been replaced with PermissionViewChannel for text and voice channels
	PermissionReadMessages = 1 << 10

	// Allows for sending messages in a channel and creating threads in a forum (does not allow sending messages in threads).
	PermissionSendMessages = 1 << 11

	// Allows for sending of /tts messages.
	PermissionSendTTSMessages = 1 << 12

	// Allows for deletion of other users messages.
	PermissionManageMessages = 1 << 13

	// Links sent by users with this permission will be auto-embedded.
	PermissionEmbedLinks = 1 << 14

	// Allows for uploading images and files.
	PermissionAttachFiles = 1 << 15

	// Allows for reading of message history.
	PermissionReadMessageHistory = 1 << 16

	// Allows for using the @everyone tag to notify all users in a channel, and the @here tag to notify all online users in a channel.
	PermissionMentionEveryone = 1 << 17

	// Allows the usage of custom emojis from other servers.
	PermissionUseExternalEmojis = 1 << 18

	// Deprecated: PermissionUseSlashCommands has been replaced by PermissionUseApplicationCommands
	PermissionUseSlashCommands = 1 << 31

	// Allows members to use application commands, including slash commands and context menu commands.
	PermissionUseApplicationCommands = 1 << 31

	// Allows for deleting and archiving threads, and viewing all private threads.
	PermissionManageThreads = 1 << 34

	// Allows for creating public and announcement threads.
	PermissionCreatePublicThreads = 1 << 35

	// Allows for creating private threads.
	PermissionCreatePrivateThreads = 1 << 36

	// Allows the usage of custom stickers from other servers.
	PermissionUseExternalStickers = 1 << 37

	// Allows for sending messages in threads.
	PermissionSendMessagesInThreads = 1 << 38

	// Allows sending voice messages.
	PermissionSendVoiceMessages = 1 << 46

	// Allows sending polls.
	PermissionSendPolls = 1 << 49

	// Allows user-installed apps to send public responses. When disabled, users will still be allowed to use their apps but the responses will be ephemeral. This only applies to apps not also installed to the server.
	PermissionUseExternalApps = 1 << 50
)

// Constants for the different bit offsets of voice permissions
const (
	// Allows for using priority speaker in a voice channel.
	PermissionVoicePrioritySpeaker = 1 << 8

	// Allows the user to go live.
	PermissionVoiceStreamVideo = 1 << 9

	// Allows for joining of a voice channel.
	PermissionVoiceConnect = 1 << 20

	// Allows for speaking in a voice channel.
	PermissionVoiceSpeak = 1 << 21

	// Allows for muting members in a voice channel.
	PermissionVoiceMuteMembers = 1 << 22

	// Allows for deafening of members in a voice channel.
	PermissionVoiceDeafenMembers = 1 << 23

	// Allows for moving of members between voice channels.
	PermissionVoiceMoveMembers = 1 << 24

	// Allows for using voice-activity-detection in a voice channel.
	PermissionVoiceUseVAD = 1 << 25

	// Allows for requesting to speak in stage channels.
	PermissionVoiceRequestToSpeak = 1 << 32

	// Deprecated: PermissionUseActivities has been replaced by PermissionUseEmbeddedActivities.
	PermissionUseActivities = 1 << 39

	// Allows for using Activities (applications with the EMBEDDED flag) in a voice channel.
	PermissionUseEmbeddedActivities = 1 << 39

	// Allows for using soundboard in a voice channel.
	PermissionUseSoundboard = 1 << 42

	// Allows the usage of custom soundboard sounds from other servers.
	PermissionUseExternalSounds = 1 << 45
)

// Constants for general management.
const (
	// Allows for modification of own nickname.
	PermissionChangeNickname = 1 << 26

	// Allows for modification of other users nicknames.
	PermissionManageNicknames = 1 << 27

	// Allows management and editing of roles.
	PermissionManageRoles = 1 << 28

	// Allows management and editing of webhooks.
	PermissionManageWebhooks = 1 << 29

	// Deprecated: PermissionManageEmojis has been replaced by PermissionManageGuildExpressions.
	PermissionManageEmojis = 1 << 30

	// Allows for editing and deleting emojis, stickers, and soundboard sounds created by all users.
	PermissionManageGuildExpressions = 1 << 30

	// Allows for editing and deleting scheduled events created by all users.
	PermissionManageEvents = 1 << 33

	// Allows for viewing role subscription insights.
	PermissionViewCreatorMonetizationAnalytics = 1 << 41

	// Allows for creating emojis, stickers, and soundboard sounds, and editing and deleting those created by the current user.
	PermissionCreateGuildExpressions = 1 << 43

	// Allows for creating scheduled events, and editing and deleting those created by the current user.
	PermissionCreateEvents = 1 << 44
)

// Constants for the different bit offsets of general permissions
const (
	// Allows creation of instant invites.
	PermissionCreateInstantInvite = 1 << 0

	// Allows kicking members.
	PermissionKickMembers = 1 << 1

	// Allows banning members.
	PermissionBanMembers = 1 << 2

	// Allows all permissions and bypasses channel permission overwrites.
	PermissionAdministrator = 1 << 3

	// Allows management and editing of channels.
	PermissionManageChannels = 1 << 4

	// Deprecated: PermissionManageServer has been replaced by PermissionManageGuild.
	PermissionManageServer = 1 << 5

	// Allows management and editing of the guild.
	PermissionManageGuild = 1 << 5

	// Allows for the addition of reactions to messages.
	PermissionAddReactions = 1 << 6

	// Allows for viewing of audit logs.
	PermissionViewAuditLogs = 1 << 7

	// Allows guild members to view a channel, which includes reading messages in text channels and joining voice channels.
	PermissionViewChannel = 1 << 10

	// Allows for viewing guild insights.
	PermissionViewGuildInsights = 1 << 19

	// Allows for timing out users to prevent them from sending or reacting to messages in chat and threads, and from speaking in voice and stage channels.
	PermissionModerateMembers = 1 << 40

	PermissionAllText = PermissionViewChannel |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone
	PermissionAllVoice = PermissionViewChannel |
		PermissionVoiceConnect |
		PermissionVoiceSpeak |
		PermissionVoiceMuteMembers |
		PermissionVoiceDeafenMembers |
		PermissionVoiceMoveMembers |
		PermissionVoiceUseVAD |
		PermissionVoicePrioritySpeaker
	PermissionAllChannel = PermissionAllText |
		PermissionAllVoice |
		PermissionCreateInstantInvite |
		PermissionManageRoles |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewAuditLogs
	PermissionAll = PermissionAllChannel |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageServer |
		PermissionAdministrator |
		PermissionManageWebhooks |
		PermissionManageEmojis
)

// Block contains Discord JSON Error Response codes
const (
	ErrCodeGeneralError = 0

	ErrCodeUnknownAccount                        = 10001
	ErrCodeUnknownApplication                    = 10002
	ErrCodeUnknownChannel                        = 10003
	ErrCodeUnknownGuild                          = 10004
	ErrCodeUnknownIntegration                    = 10005
	ErrCodeUnknownInvite                         = 10006
	ErrCodeUnknownMember                         = 10007
	ErrCodeUnknownMessage                        = 10008
	ErrCodeUnknownOverwrite                      = 10009
	ErrCodeUnknownProvider                       = 10010
	ErrCodeUnknownRole                           = 10011
	ErrCodeUnknownToken                          = 10012
	ErrCodeUnknownUser                           = 10013
	ErrCodeUnknownEmoji                          = 10014
	ErrCodeUnknownWebhook                        = 10015
	ErrCodeUnknownWebhookService                 = 10016
	ErrCodeUnknownSession                        = 10020
	ErrCodeUnknownBan                            = 10026
	ErrCodeUnknownSKU                            = 10027
	ErrCodeUnknownStoreListing                   = 10028
	ErrCodeUnknownEntitlement                    = 10029
	ErrCodeUnknownBuild                          = 10030
	ErrCodeUnknownLobby                          = 10031
	ErrCodeUnknownBranch                         = 10032
	ErrCodeUnknownStoreDirectoryLayout           = 10033
	ErrCodeUnknownRedistributable                = 10036
	ErrCodeUnknownGiftCode                       = 10038
	ErrCodeUnknownStream                         = 10049
	ErrCodeUnknownPremiumServerSubscribeCooldown = 10050
	ErrCodeUnknownGuildTemplate                  = 10057
	ErrCodeUnknownDiscoveryCategory              = 10059
	ErrCodeUnknownSticker                        = 10060
	ErrCodeUnknownInteraction                    = 10062
	ErrCodeUnknownApplicationCommand             = 10063
	ErrCodeUnknownApplicationCommandPermissions  = 10066
	ErrCodeUnknownStageInstance                  = 10067
	ErrCodeUnknownGuildMemberVerificationForm    = 10068
	ErrCodeUnknownGuildWelcomeScreen             = 10069
	ErrCodeUnknownGuildScheduledEvent            = 10070
	ErrCodeUnknownGuildScheduledEventUser        = 10071
	ErrUnknownTag                                = 10087

	ErrCodeBotsCannotUseEndpoint                                            = 20001
	ErrCodeOnlyBotsCanUseEndpoint                                           = 20002
	ErrCodeExplicitContentCannotBeSentToTheDesiredRecipients                = 20009
	ErrCodeYouAreNotAuthorizedToPerformThisActionOnThisApplication          = 20012
	ErrCodeThisActionCannotBePerformedDueToSlowmodeRateLimit                = 20016
	ErrCodeOnlyTheOwnerOfThisAccountCanPerformThisAction                    = 20018
	ErrCodeMessageCannotBeEditedDueToAnnouncementRateLimits                 = 20022
	ErrCodeChannelHasHitWriteRateLimit                                      = 20028
	ErrCodeTheWriteActionYouArePerformingOnTheServerHasHitTheWriteRateLimit = 20029
	ErrCodeStageTopicContainsNotAllowedWordsForPublicStages                 = 20031
	ErrCodeGuildPremiumSubscriptionLevelTooLow                              = 20035

	ErrCodeMaximumGuildsReached                                     = 30001
	ErrCodeMaximumPinsReached                                       = 30003
	ErrCodeMaximumNumberOfRecipientsReached                         = 30004
	ErrCodeMaximumGuildRolesReached                                 = 30005
	ErrCodeMaximumNumberOfWebhooksReached                           = 30007
	ErrCodeMaximumNumberOfEmojisReached                             = 30008
	ErrCodeTooManyReactions                                         = 30010
	ErrCodeMaximumNumberOfGuildChannelsReached                      = 30013
	ErrCodeMaximumNumberOfAttachmentsInAMessageReached              = 30015
	ErrCodeMaximumNumberOfInvitesReached                            = 30016
	ErrCodeMaximumNumberOfAnimatedEmojisReached                     = 30018
	ErrCodeMaximumNumberOfServerMembersReached                      = 30019
	ErrCodeMaximumNumberOfGuildDiscoverySubcategoriesReached        = 30030
	ErrCodeGuildAlreadyHasATemplate                                 = 30031
	ErrCodeMaximumNumberOfThreadParticipantsReached                 = 30033
	ErrCodeMaximumNumberOfBansForNonGuildMembersHaveBeenExceeded    = 30035
	ErrCodeMaximumNumberOfBansFetchesHasBeenReached                 = 30037
	ErrCodeMaximumNumberOfUncompletedGuildScheduledEventsReached    = 30038
	ErrCodeMaximumNumberOfStickersReached                           = 30039
	ErrCodeMaximumNumberOfPruneRequestsHasBeenReached               = 30040
	ErrCodeMaximumNumberOfGuildWidgetSettingsUpdatesHasBeenReached  = 30042
	ErrCodeMaximumNumberOfEditsToMessagesOlderThanOneHourReached    = 30046
	ErrCodeMaximumNumberOfPinnedThreadsInForumChannelHasBeenReached = 30047
	ErrCodeMaximumNumberOfTagsInForumChannelHasBeenReached          = 30048

	ErrCodeUnauthorized                           = 40001
	ErrCodeActionRequiredVerifiedAccount          = 40002
	ErrCodeOpeningDirectMessagesTooFast           = 40003
	ErrCodeSendMessagesHasBeenTemporarilyDisabled = 40004
	ErrCodeRequestEntityTooLarge                  = 40005
	ErrCodeFeatureTemporarilyDisabledServerSide   = 40006
	ErrCodeUserIsBannedFromThisGuild              = 40007
	ErrCodeTargetIsNotConnectedToVoice            = 40032
	ErrCodeMessageAlreadyCrossposted              = 40033
	ErrCodeAnApplicationWithThatNameAlreadyExists = 40041
	ErrCodeInteractionHasAlreadyBeenAcknowledged  = 40060
	ErrCodeTagNamesMustBeUnique                   = 40061

	ErrCodeMissingAccess                                                = 50001
	ErrCodeInvalidAccountType                                           = 50002
	ErrCodeCannotExecuteActionOnDMChannel                               = 50003
	ErrCodeEmbedDisabled                                                = 50004
	ErrCodeGuildWidgetDisabled                                          = 50004
	ErrCodeCannotEditFromAnotherUser                                    = 50005
	ErrCodeCannotSendEmptyMessage                                       = 50006
	ErrCodeCannotSendMessagesToThisUser                                 = 50007
	ErrCodeCannotSendMessagesInVoiceChannel                             = 50008
	ErrCodeChannelVerificationLevelTooHigh                              = 50009
	ErrCodeOAuth2ApplicationDoesNotHaveBot                              = 50010
	ErrCodeOAuth2ApplicationLimitReached                                = 50011
	ErrCodeInvalidOAuthState                                            = 50012
	ErrCodeMissingPermissions                                           = 50013
	ErrCodeInvalidAuthenticationToken                                   = 50014
	ErrCodeTooFewOrTooManyMessagesToDelete                              = 50016
	ErrCodeCanOnlyPinMessageToOriginatingChannel                        = 50019
	ErrCodeInviteCodeWasEitherInvalidOrTaken                            = 50020
	ErrCodeCannotExecuteActionOnSystemMessage                           = 50021
	ErrCodeCannotExecuteActionOnThisChannelType                         = 50024
	ErrCodeInvalidOAuth2AccessTokenProvided                             = 50025
	ErrCodeMissingRequiredOAuth2Scope                                   = 50026
	ErrCodeInvalidWebhookTokenProvided                                  = 50027
	ErrCodeInvalidRole                                                  = 50028
	ErrCodeInvalidRecipients                                            = 50033
	ErrCodeMessageProvidedTooOldForBulkDelete                           = 50034
	ErrCodeInvalidFormBody                                              = 50035
	ErrCodeInviteAcceptedToGuildApplicationsBotNotIn                    = 50036
	ErrCodeInvalidAPIVersionProvided                                    = 50041
	ErrCodeFileUploadedExceedsTheMaximumSize                            = 50045
	ErrCodeInvalidFileUploaded                                          = 50046
	ErrCodeInvalidGuild                                                 = 50055
	ErrCodeInvalidMessageType                                           = 50068
	ErrCodeCannotDeleteAChannelRequiredForCommunityGuilds               = 50074
	ErrCodeInvalidStickerSent                                           = 50081
	ErrCodePerformedOperationOnArchivedThread                           = 50083
	ErrCodeBeforeValueIsEarlierThanThreadCreationDate                   = 50085
	ErrCodeCommunityServerChannelsMustBeTextChannels                    = 50086
	ErrCodeThisServerIsNotAvailableInYourLocation                       = 50095
	ErrCodeThisServerNeedsMonetizationEnabledInOrderToPerformThisAction = 50097
	ErrCodeThisServerNeedsMoreBoostsToPerformThisAction                 = 50101
	ErrCodeTheRequestBodyContainsInvalidJSON                            = 50109

	ErrCodeNoUsersWithDiscordTagExist = 80004

	ErrCodeReactionBlocked = 90001

	ErrCodeAPIResourceIsCurrentlyOverloaded = 130000

	ErrCodeTheStageIsAlreadyOpen = 150006

	ErrCodeCannotReplyWithoutPermissionToReadMessageHistory = 160002
	ErrCodeThreadAlreadyCreatedForThisMessage               = 160004
	ErrCodeThreadIsLocked                                   = 160005
	ErrCodeMaximumNumberOfActiveThreadsReached              = 160006
	ErrCodeMaximumNumberOfActiveAnnouncementThreadsReached  = 160007

	ErrCodeInvalidJSONForUploadedLottieFile                    = 170001
	ErrCodeUploadedLottiesCannotContainRasterizedImages        = 170002
	ErrCodeStickerMaximumFramerateExceeded                     = 170003
	ErrCodeStickerFrameCountExceedsMaximumOfOneThousandFrames  = 170004
	ErrCodeLottieAnimationMaximumDimensionsExceeded            = 170005
	ErrCodeStickerFrameRateOutOfRange                          = 170006
	ErrCodeStickerAnimationDurationExceedsMaximumOfFiveSeconds = 170007

	ErrCodeCannotUpdateAFinishedEvent             = 180000
	ErrCodeFailedToCreateStageNeededForStageEvent = 180002

	ErrCodeCannotEnableOnboardingRequirementsAreNotMet  = 350000
	ErrCodeCannotUpdateOnboardingWhileBelowRequirements = 350001
)

// Intent is the type of a Gateway Intent
// https://discord.com/developers/docs/topics/gateway#gateway-intents
type Intent int

// Constants for the different bit offsets of intents
const (
	IntentGuilds                      Intent = 1 << 0
	IntentGuildMembers                Intent = 1 << 1
	IntentGuildModeration             Intent = 1 << 2
	IntentGuildEmojis                 Intent = 1 << 3
	IntentGuildIntegrations           Intent = 1 << 4
	IntentGuildWebhooks               Intent = 1 << 5
	IntentGuildInvites                Intent = 1 << 6
	IntentGuildVoiceStates            Intent = 1 << 7
	IntentGuildPresences              Intent = 1 << 8
	IntentGuildMessages               Intent = 1 << 9
	IntentGuildMessageReactions       Intent = 1 << 10
	IntentGuildMessageTyping          Intent = 1 << 11
	IntentDirectMessages              Intent = 1 << 12
	IntentDirectMessageReactions      Intent = 1 << 13
	IntentDirectMessageTyping         Intent = 1 << 14
	IntentMessageContent              Intent = 1 << 15
	IntentGuildScheduledEvents        Intent = 1 << 16
	IntentAutoModerationConfiguration Intent = 1 << 20
	IntentAutoModerationExecution     Intent = 1 << 21
	IntentGuildMessagePolls           Intent = 1 << 24
	IntentDirectMessagePolls          Intent = 1 << 25

	// TODO: remove when compatibility is not needed

	IntentGuildBans Intent = IntentGuildModeration

	IntentsGuilds                 Intent = 1 << 0
	IntentsGuildMembers           Intent = 1 << 1
	IntentsGuildBans              Intent = 1 << 2
	IntentsGuildEmojis            Intent = 1 << 3
	IntentsGuildIntegrations      Intent = 1 << 4
	IntentsGuildWebhooks          Intent = 1 << 5
	IntentsGuildInvites           Intent = 1 << 6
	IntentsGuildVoiceStates       Intent = 1 << 7
	IntentsGuildPresences         Intent = 1 << 8
	IntentsGuildMessages          Intent = 1 << 9
	IntentsGuildMessageReactions  Intent = 1 << 10
	IntentsGuildMessageTyping     Intent = 1 << 11
	IntentsDirectMessages         Intent = 1 << 12
	IntentsDirectMessageReactions Intent = 1 << 13
	IntentsDirectMessageTyping    Intent = 1 << 14
	IntentsMessageContent         Intent = 1 << 15
	IntentsGuildScheduledEvents   Intent = 1 << 16

	IntentsAllWithoutPrivileged = IntentGuilds |
		IntentGuildBans |
		IntentGuildEmojis |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping |
		IntentGuildScheduledEvents |
		IntentAutoModerationConfiguration |
		IntentAutoModerationExecution

	IntentsAll = IntentsAllWithoutPrivileged |
		IntentGuildMembers |
		IntentGuildPresences |
		IntentMessageContent

	IntentsNone Intent = 0
)

// MakeIntent used to help convert a gateway intent value for use in the Identify structure;
// this was useful to help support the use of a pointer type when intents were optional.
// This is now a no-op, and is not necessary to use.
func MakeIntent(intents Intent) Intent {
	return intents
}

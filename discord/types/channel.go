package types

// Channel is the type of channel.Channel
type Channel int

// Block contains known Channel type values
const (
	ChannelGuildText          Channel = 0
	ChannelDM                 Channel = 1
	ChannelGuildVoice         Channel = 2
	ChannelGroupDM            Channel = 3
	ChannelGuildCategory      Channel = 4
	ChannelGuildNews          Channel = 5
	ChannelGuildStore         Channel = 6
	ChannelGuildNewsThread    Channel = 10
	ChannelGuildPublicThread  Channel = 11
	ChannelGuildPrivateThread Channel = 12
	ChannelGuildStageVoice    Channel = 13
	ChannelGuildDirectory     Channel = 14
	ChannelGuildForum         Channel = 15
	ChannelGuildMedia         Channel = 16
)

// Webhook is the type of Webhook (see Webhook* consts) in the Webhook struct
// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
type Webhook int

// Valid Webhook values
const (
	WebhookIncoming        Webhook = 1
	WebhookChannelFollower Webhook = 2
)

// AllowedMention describes the types of mentions used
// in the MessageAllowedMentions type.
type AllowedMention string

// The types of mentions used in MessageAllowedMentions.
const (
	AllowedMentionRoles    AllowedMention = "roles"
	AllowedMentionUsers    AllowedMention = "users"
	AllowedMentionEveryone AllowedMention = "everyone"
)

// PollLayout represents the layout of a poll.
type PollLayout int

// Valid PollLayout values.
const (
	PollLayoutDefault PollLayout = 1
)

// Embed is the type of embed
// https://discord.com/developers/docs/resources/channel#embed-object-embed-types
type Embed string

// Block of valid EmbedTypes
const (
	EmbedRich    Embed = "rich"
	EmbedImage   Embed = "image"
	EmbedVideo   Embed = "video"
	EmbedGifv    Embed = "gifv"
	EmbedArticle Embed = "article"
	EmbedLink    Embed = "link"
)

// Message is the type of Message
// https://discord.com/developers/docs/resources/channel#message-object-message-types
type Message int

// Block contains the valid known Message values
const (
	MessageDefault                               Message = 0
	MessageRecipientAdd                          Message = 1
	MessageRecipientRemove                       Message = 2
	MessageCall                                  Message = 3
	MessageChannelNameChange                     Message = 4
	MessageChannelIconChange                     Message = 5
	MessageChannelPinnedMessage                  Message = 6
	MessageGuildMemberJoin                       Message = 7
	MessageUserPremiumGuildSubscription          Message = 8
	MessageUserPremiumGuildSubscriptionTierOne   Message = 9
	MessageUserPremiumGuildSubscriptionTierTwo   Message = 10
	MessageUserPremiumGuildSubscriptionTierThree Message = 11
	MessageChannelFollowAdd                      Message = 12
	MessageGuildDiscoveryDisqualified            Message = 14
	MessageGuildDiscoveryRequalified             Message = 15
	MessageThreadCreated                         Message = 18
	MessageReply                                 Message = 19
	MessageChatInputCommand                      Message = 20
	MessageThreadStarterMessage                  Message = 21
	MessageContextMenuCommand                    Message = 23
)

// MessageActivity is the type of message activity
type MessageActivity int

// Constants for the different types of Message Activity
const (
	MessageActivityJoin        MessageActivity = 1
	MessageActivitySpectate    MessageActivity = 2
	MessageActivityListen      MessageActivity = 3
	MessageActivityJoinRequest MessageActivity = 5
)

// MessageReference is a type of MessageReference
type MessageReference int

// Known valid MessageReference values
// https://discord.com/developers/docs/resources/message#message-reference-types
const (
	MessageReferenceDefault MessageReference = 0
	MessageReferenceForward MessageReference = 1
)

// ForumSortOrder represents sort order of a forum channel.
type ForumSortOrder int

const (
	// ForumSortOrderLatestActivity sorts posts by activity.
	ForumSortOrderLatestActivity ForumSortOrder = 0
	// ForumSortOrderCreationDate sorts posts by creation time (from most recent to oldest).
	ForumSortOrderCreationDate ForumSortOrder = 1
)

// PermissionOverwrite represents the type of resource on which
// a permission overwrite acts.
type PermissionOverwrite int

// The possible permission overwrite types.
const (
	PermissionOverwriteRole   PermissionOverwrite = 0
	PermissionOverwriteMember PermissionOverwrite = 1
)

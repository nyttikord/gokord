package event

import (
	"encoding/json"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/interaction"
	"github.com/nyttikord/gokord/premium"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/invite"
	"github.com/nyttikord/gokord/user/status"
)

// This file contains all the possible structs that can be
// handled by AddHandler/Handler.
// DO NOT ADD ANYTHING BUT EVENT HANDLER STRUCTS TO THIS FILE.
//go:generate go run ../tools/cmd/eventhandlers/main.go

// Connect is the data for a Connect event.
// This is a synthetic event and is not dispatched by Discord.
type Connect struct{}

// Disconnect is the data for a Disconnect event.
// This is a synthetic event and is not dispatched by Discord.
type Disconnect struct{}

// RateLimit is the data for a RateLimit event.
// This is a synthetic event and is not dispatched by Discord.
type RateLimit struct {
	*discord.TooManyRequests
	URL string
}

// A Ready stores all data for the websocket READY event.
type Ready struct {
	Version          int                      `json:"v"`
	SessionID        string                   `json:"session_id"`
	User             *user.User               `json:"user"`
	Shard            *[2]int                  `json:"shard"`
	ResumeGatewayURL string                   `json:"resume_gateway_url"`
	Application      *application.Application `json:"application"`
	Guilds           []*guild.Guild           `json:"guilds"`
	PrivateChannels  []*channel.Channel       `json:"private_channels"`
}

// ChannelCreate is the data for a ChannelCreate event.
type ChannelCreate struct {
	*channel.Channel
}

// ChannelUpdate is the data for a ChannelUpdate event.
type ChannelUpdate struct {
	*channel.Channel
	BeforeUpdate *channel.Channel `json:"-"`
}

// ChannelDelete is the data for a ChannelDelete event.
type ChannelDelete struct {
	*channel.Channel
	BeforeDelete *channel.Channel `json:"-"`
}

// ChannelPinsUpdate stores data for a ChannelPinsUpdate event.
type ChannelPinsUpdate struct {
	LastPinTimestamp string `json:"last_pin_timestamp"`
	ChannelID        string `json:"channel_id"`
	GuildID          string `json:"guild_id,omitempty"`
}

// ThreadCreate is the data for a ThreadCreate event.
type ThreadCreate struct {
	*channel.Channel
	NewlyCreated bool `json:"newly_created"`
}

// ThreadUpdate is the data for a ThreadUpdate event.
type ThreadUpdate struct {
	*channel.Channel
	BeforeUpdate *channel.Channel `json:"-"`
}

// ThreadDelete is the data for a ThreadDelete event.
type ThreadDelete struct {
	*channel.Channel
	BeforeDelete *channel.Channel `json:"-"`
}

// ThreadListSync is the data for a ThreadListSync event.
type ThreadListSync struct {
	// The id of the guild
	GuildID string `json:"guild_id"`
	// The parent channel ids whose threads are being synced.
	// If omitted, then threads were synced for the entire guild.
	// This array may contain channel_ids that have no active threads as well, so you know to clear that data.
	ChannelIDs []string `json:"channel_ids"`
	// All active threads in the given channels that the current user can access
	Threads []*channel.Channel `json:"threads"`
	// All thread member objects from the synced threads for the current user,
	// indicating which threads the current user has been added to
	Members []*channel.ThreadMember `json:"members"`
}

// ThreadMemberUpdate is the data for a ThreadMemberUpdate event.
type ThreadMemberUpdate struct {
	*channel.ThreadMember
	GuildID string `json:"guild_id"`
}

// ThreadMembersUpdate is the data for a ThreadMembersUpdate event.
type ThreadMembersUpdate struct {
	ID             string                      `json:"id"`
	GuildID        string                      `json:"guild_id"`
	MemberCount    int                         `json:"member_count"`
	AddedMembers   []channel.AddedThreadMember `json:"added_members"`
	RemovedMembers []string                    `json:"removed_member_ids"`
}

// GuildCreate is the data for a GuildCreate event.
type GuildCreate struct {
	*guild.Guild
}

// GuildUpdate is the data for a GuildUpdate event.
type GuildUpdate struct {
	*guild.Guild
}

// GuildDelete is the data for a GuildDelete event.
type GuildDelete struct {
	*guild.Guild
	BeforeDelete *guild.Guild `json:"-"`
}

// GuildBanAdd is the data for a GuildBanAdd event.
type GuildBanAdd struct {
	User    *user.User `json:"user"`
	GuildID string     `json:"guild_id"`
}

// GuildBanRemove is the data for a GuildBanRemove event.
type GuildBanRemove struct {
	User    *user.User `json:"user"`
	GuildID string     `json:"guild_id"`
}

// GuildMemberAdd is the data for a GuildMemberAdd event.
type GuildMemberAdd struct {
	*user.Member
}

// GuildMemberUpdate is the data for a GuildMemberUpdate event.
type GuildMemberUpdate struct {
	*user.Member
	BeforeUpdate *user.Member `json:"-"`
}

// GuildMemberRemove is the data for a GuildMemberRemove event.
type GuildMemberRemove struct {
	*user.Member
	BeforeDelete *user.Member `json:"-"`
}

// GuildRoleCreate is the data for a GuildRoleCreate event.
type GuildRoleCreate struct {
	*guild.GuildedRole
}

// GuildRoleUpdate is the data for a GuildRoleUpdate event.
type GuildRoleUpdate struct {
	*guild.GuildedRole
	BeforeUpdate *guild.Role `json:"-"`
}

// A GuildRoleDelete is the data for a GuildRoleDelete event.
type GuildRoleDelete struct {
	RoleID       string      `json:"role_id"`
	GuildID      string      `json:"guild_id"`
	BeforeDelete *guild.Role `json:"-"`
}

// A GuildEmojisUpdate is the data for a guild emoji update event.
type GuildEmojisUpdate struct {
	GuildID string         `json:"guild_id"`
	Emojis  []*emoji.Emoji `json:"emojis"`
}

// A GuildStickersUpdate is the data for a GuildStickersUpdate event.
type GuildStickersUpdate struct {
	GuildID  string           `json:"guild_id"`
	Stickers []*emoji.Sticker `json:"stickers"`
}

// A GuildMembersChunk is the data for a GuildMembersChunk event.
type GuildMembersChunk struct {
	GuildID    string             `json:"guild_id"`
	Members    []*user.Member     `json:"members"`
	ChunkIndex int                `json:"chunk_index"`
	ChunkCount int                `json:"chunk_count"`
	NotFound   []string           `json:"not_found,omitempty"`
	Presences  []*status.Presence `json:"presences,omitempty"`
	Nonce      string             `json:"nonce,omitempty"`
}

// GuildIntegrationsUpdate is the data for a GuildIntegrationsUpdate event.
type GuildIntegrationsUpdate struct {
	GuildID string `json:"guild_id"`
}

// StageInstanceEventCreate is the data for a StageInstanceEventCreate event.
type StageInstanceEventCreate struct {
	*channel.StageInstance
}

// StageInstanceEventUpdate is the data for a StageInstanceEventUpdate event.
type StageInstanceEventUpdate struct {
	*channel.StageInstance
}

// StageInstanceEventDelete is the data for a StageInstanceEventDelete event.
type StageInstanceEventDelete struct {
	*channel.StageInstance
}

// GuildScheduledEventCreate is the data for a GuildScheduledEventCreate event.
type GuildScheduledEventCreate struct {
	*guild.ScheduledEvent
}

// GuildScheduledEventUpdate is the data for a GuildScheduledEventUpdate event.
type GuildScheduledEventUpdate struct {
	*guild.ScheduledEvent
}

// GuildScheduledEventDelete is the data for a GuildScheduledEventDelete event.
type GuildScheduledEventDelete struct {
	*guild.ScheduledEvent
}

// GuildScheduledEventUserAdd is the data for a GuildScheduledEventUserAdd event.
type GuildScheduledEventUserAdd struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

// GuildScheduledEventUserRemove is the data for a GuildScheduledEventUserRemove event.
type GuildScheduledEventUserRemove struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

// IntegrationCreate is the data for a IntegrationCreate event.
type IntegrationCreate struct {
	*user.Integration
	GuildID string `json:"guild_id"`
}

// IntegrationUpdate is the data for a IntegrationUpdate event.
type IntegrationUpdate struct {
	*user.Integration
	GuildID string `json:"guild_id"`
}

// IntegrationDelete is the data for a IntegrationDelete event.
type IntegrationDelete struct {
	ID            string `json:"id"`
	GuildID       string `json:"guild_id"`
	ApplicationID string `json:"application_id,omitempty"`
}

// MessageCreate is the data for a MessageCreate event.
type MessageCreate struct {
	*channel.Message
}

// UnmarshalJSON is a helper function to unmarshal MessageCreate object.
func (m *MessageCreate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

// MessageUpdate is the data for a MessageUpdate event.
type MessageUpdate struct {
	*channel.Message
	// BeforeUpdate will be nil if the Message was not previously cached in the state cache.
	BeforeUpdate *channel.Message `json:"-"`
}

// UnmarshalJSON is a helper function to unmarshal MessageUpdate object.
func (m *MessageUpdate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

// MessageDelete is the data for a MessageDelete event.
type MessageDelete struct {
	*channel.Message
	BeforeDelete *channel.Message `json:"-"`
}

// UnmarshalJSON is a helper function to unmarshal MessageDelete object.
func (m *MessageDelete) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

// MessageReactionAdd is the data for a MessageReactionAdd event.
type MessageReactionAdd struct {
	*channel.MessageReaction
	Member *user.Member `json:"member,omitempty"`
}

// MessageReactionRemove is the data for a MessageReactionRemove event.
type MessageReactionRemove struct {
	*channel.MessageReaction
}

// MessageReactionRemoveAll is the data for a MessageReactionRemoveAll event.
type MessageReactionRemoveAll struct {
	*channel.MessageReaction
}

// PresencesReplace is the data for a PresencesReplace event.
type PresencesReplace []*status.Presence

// PresenceUpdate is the data for a PresenceUpdate event.
type PresenceUpdate struct {
	status.Presence
	GuildID string `json:"guild_id"`
}

// Resumed is the data for a Resumed event.
type Resumed struct {
	Trace []string `json:"_trace"`
}

// TypingStart is the data for a TypingStart event.
type TypingStart struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Timestamp int    `json:"timestamp"`
}

// UserUpdate is the data for a UserUpdate event.
type UserUpdate struct {
	*user.User
}

// VoiceServerUpdate is the data for a VoiceServerUpdate event.
type VoiceServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

// VoiceStateUpdate is the data for a VoiceStateUpdate event.
type VoiceStateUpdate struct {
	*user.VoiceState
	// BeforeUpdate will be nil if the VoiceState was not previously cached in the state cache.
	BeforeUpdate *user.VoiceState `json:"-"`
}

// MessageDeleteBulk is the data for a MessageDeleteBulk event
type MessageDeleteBulk struct {
	Messages  []string `json:"ids"`
	ChannelID string   `json:"channel_id"`
	GuildID   string   `json:"guild_id"`
}

// WebhooksUpdate is the data for a WebhooksUpdate event
type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

// InteractionCreate is the data for a InteractionCreate event
type InteractionCreate struct {
	*interaction.Interaction
}

// UnmarshalJSON is a helper function to unmarshal Interaction object.
func (i *InteractionCreate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &i.Interaction)
}

// InviteCreate is the data for a InviteCreate event
type InviteCreate struct {
	*invite.Invite
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

// InviteDelete is the data for a InviteDelete event
type InviteDelete struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	Code      string `json:"code"`
}

// ApplicationCommandPermissionsUpdate is the data for an ApplicationCommandPermissionsUpdate event
type ApplicationCommandPermissionsUpdate struct {
	*interaction.GuildCommandPermissions
}

// AutoModerationRuleCreate is the data for an AutoModerationRuleCreate event.
type AutoModerationRuleCreate struct {
	*guild.AutoModerationRule
}

// AutoModerationRuleUpdate is the data for an AutoModerationRuleUpdate event.
type AutoModerationRuleUpdate struct {
	*guild.AutoModerationRule
}

// AutoModerationRuleDelete is the data for an AutoModerationRuleDelete event.
type AutoModerationRuleDelete struct {
	*guild.AutoModerationRule
}

// AutoModerationActionExecution is the data for an AutoModerationActionExecution event.
type AutoModerationActionExecution struct {
	GuildID              string                          `json:"guild_id"`
	Action               guild.AutoModerationAction      `json:"action"`
	RuleID               string                          `json:"rule_id"`
	RuleTriggerType      guild.AutoModerationRuleTrigger `json:"rule_trigger_type"`
	UserID               string                          `json:"user_id"`
	ChannelID            string                          `json:"channel_id"`
	MessageID            string                          `json:"message_id"`
	AlertSystemMessageID string                          `json:"alert_system_message_id"`
	Content              string                          `json:"content"`
	MatchedKeyword       string                          `json:"matched_keyword"`
	MatchedContent       string                          `json:"matched_content"`
}

// GuildAuditLogEntryCreate is the data for a GuildAuditLogEntryCreate event.
type GuildAuditLogEntryCreate struct {
	*guild.AuditLogEntry
	GuildID string `json:"guild_id"`
}

// MessagePollVoteAdd is the data for a MessagePollVoteAdd event.
type MessagePollVoteAdd struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
	AnswerID  int    `json:"answer_id"`
}

// MessagePollVoteRemove is the data for a MessagePollVoteRemove event.
type MessagePollVoteRemove struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
	AnswerID  int    `json:"answer_id"`
}

// EntitlementCreate is the data for an EntitlementCreate event.
type EntitlementCreate struct {
	*premium.Entitlement
}

// EntitlementUpdate is the data for an EntitlementUpdate event.
type EntitlementUpdate struct {
	*premium.Entitlement
}

// EntitlementDelete is the data for an EntitlementDelete event.
// NOTE: Entitlements are not deleted when they expire.
type EntitlementDelete struct {
	*premium.Entitlement
}

// SubscriptionCreate is the data for an SubscriptionCreate event.
// https://discord.com/developers/docs/monetization/implementing-app-subscriptions#using-subscription-events-for-the-subscription-lifecycle
type SubscriptionCreate struct {
	*premium.Subscription
}

// SubscriptionUpdate is the data for an SubscriptionUpdate event.
// https://discord.com/developers/docs/monetization/implementing-app-subscriptions#using-subscription-events-for-the-subscription-lifecycle
type SubscriptionUpdate struct {
	*premium.Subscription
}

// SubscriptionDelete is the data for an SubscriptionDelete event.
// https://discord.com/developers/docs/monetization/implementing-app-subscriptions#using-subscription-events-for-the-subscription-lifecycle
type SubscriptionDelete struct {
	*premium.Subscription
}

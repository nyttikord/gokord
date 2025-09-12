package gokord

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg" // For JPEG decoding
	_ "image/png"  // For PNG decoding
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/interactions"
	"github.com/nyttikord/gokord/premium"
	"github.com/nyttikord/gokord/user"
	"github.com/nyttikord/gokord/user/invite"
)

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Invites
// ------------------------------------------------------------------------------------------------

// Invite returns an Invite structure of the given invite
// inviteID : The invite code
func (s *Session) Invite(inviteID string, options ...discord.RequestOption) (st *invite.Invite, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointInvite(inviteID), nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// InviteWithCounts returns an Invite structure of the given invite including approximate member counts
// inviteID : The invite code
func (s *Session) InviteWithCounts(inviteID string, options ...discord.RequestOption) (st *invite.Invite, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointInvite(inviteID)+"?with_counts=true", nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// InviteComplex returns an Invite structure of the given invite including specified fields.
// inviteID                  : The invite code
// guildScheduledEventID     : If specified, includes specified guild scheduled event.
// withCounts                : Whether to include approximate member counts or not
// withExpiration            : Whether to include expiration time or not
func (s *Session) InviteComplex(inviteID, guildScheduledEventID string, withCounts, withExpiration bool, options ...discord.RequestOption) (st *invite.Invite, err error) {
	endpoint := discord.EndpointInvite(inviteID)
	v := url.Values{}
	if guildScheduledEventID != "" {
		v.Set("guild_scheduled_event_id", guildScheduledEventID)
	}
	if withCounts {
		v.Set("with_counts", "true")
	}
	if withExpiration {
		v.Set("with_expiration", "true")
	}

	if len(v) != 0 {
		endpoint += "?" + v.Encode()
	}

	body, err := s.RequestWithBucketID("GET", endpoint, nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// InviteDelete deletes an existing invite
// inviteID   : the code of an invite
func (s *Session) InviteDelete(inviteID string, options ...discord.RequestOption) (st *invite.Invite, err error) {

	body, err := s.RequestWithBucketID("DELETE", discord.EndpointInvite(inviteID), nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// InviteAccept accepts an Invite to a Guild or Channel
// inviteID : The invite code
func (s *Session) InviteAccept(inviteID string, options ...discord.RequestOption) (st *invite.Invite, err error) {

	body, err := s.RequestWithBucketID("POST", discord.EndpointInvite(inviteID), nil, discord.EndpointInvite(""), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Voice
// ------------------------------------------------------------------------------------------------

func (s *Session) VoiceRegions(options ...discord.RequestOption) (st []*discord.VoiceRegion, err error) {

	body, err := s.RequestWithBucketID("GET", discord.EndpointVoiceRegions, nil, discord.EndpointVoiceRegions, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Websockets
// ------------------------------------------------------------------------------------------------

// Gateway returns the websocket Gateway address
func (s *Session) Gateway(options ...discord.RequestOption) (gateway string, err error) {

	response, err := s.RequestWithBucketID("GET", discord.EndpointGateway, nil, discord.EndpointGateway, options...)
	if err != nil {
		return
	}

	temp := struct {
		URL string `json:"url"`
	}{}

	err = unmarshal(response, &temp)
	if err != nil {
		return
	}

	gateway = temp.URL

	// Ensure the gateway always has a trailing slash.
	// MacOS will fail to connect if we add query params without a trailing slash on the base domain.
	if !strings.HasSuffix(gateway, "/") {
		gateway += "/"
	}

	return
}

// GatewayBot returns the websocket Gateway address and the recommended number of shards
func (s *Session) GatewayBot(options ...discord.RequestOption) (st *GatewayBotResponse, err error) {

	response, err := s.RequestWithBucketID("GET", discord.EndpointGatewayBot, nil, discord.EndpointGatewayBot, options...)
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	if err != nil {
		return
	}

	// Ensure the gateway always has a trailing slash.
	// MacOS will fail to connect if we add query params without a trailing slash on the base domain.
	if !strings.HasSuffix(st.URL, "/") {
		st.URL += "/"
	}

	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to application (slash) commands
// ------------------------------------------------------------------------------------------------

// ApplicationCommandCreate creates a global application command and returns it.
// appID       : The application ID.
// guildID     : Guild ID to create guild-specific application command. If empty - creates global application command.
// cmd         : New application command data.
func (s *Session) ApplicationCommandCreate(appID string, guildID string, cmd *interactions.Command, options ...discord.RequestOption) (ccmd *interactions.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	body, err := s.RequestWithBucketID("POST", endpoint, *cmd, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &ccmd)

	return
}

// ApplicationCommandEdit edits application command and returns new command data.
// appID       : The application ID.
// cmdID       : Application command ID to edit.
// guildID     : Guild ID to edit guild-specific application command. If empty - edits global application command.
// cmd         : Updated application command data.
func (s *Session) ApplicationCommandEdit(appID, guildID, cmdID string, cmd *interactions.Command, options ...discord.RequestOption) (updated *interactions.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	body, err := s.RequestWithBucketID("PATCH", endpoint, *cmd, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &updated)

	return
}

// ApplicationCommandBulkOverwrite Creates commands overwriting existing commands. Returns a list of commands.
// appID    : The application ID.
// commands : The commands to create.
func (s *Session) ApplicationCommandBulkOverwrite(appID string, guildID string, commands []*interactions.Command, options ...discord.RequestOption) (createdCommands []*interactions.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	body, err := s.RequestWithBucketID("PUT", endpoint, commands, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &createdCommands)

	return
}

// ApplicationCommandDelete deletes application command by ID.
// appID       : The application ID.
// cmdID       : Application command ID to delete.
// guildID     : Guild ID to delete guild-specific application command. If empty - deletes global application command.
func (s *Session) ApplicationCommandDelete(appID, guildID, cmdID string, options ...discord.RequestOption) error {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	_, err := s.RequestWithBucketID("DELETE", endpoint, nil, endpoint, options...)

	return err
}

// ApplicationCommand retrieves an application command by given ID.
// appID       : The application ID.
// cmdID       : Application command ID.
// guildID     : Guild ID to retrieve guild-specific application command. If empty - retrieves global application command.
func (s *Session) ApplicationCommand(appID, guildID, cmdID string, options ...discord.RequestOption) (cmd *interactions.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommand(appID, cmdID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommand(appID, guildID, cmdID)
	}

	body, err := s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &cmd)

	return
}

// ApplicationCommands retrieves all commands in application.
// appID       : The application ID.
// guildID     : Guild ID to retrieve all guild-specific application commands. If empty - retrieves global application commands.
func (s *Session) ApplicationCommands(appID, guildID string, options ...discord.RequestOption) (cmd []*interactions.Command, err error) {
	endpoint := discord.EndpointApplicationGlobalCommands(appID)
	if guildID != "" {
		endpoint = discord.EndpointApplicationGuildCommands(appID, guildID)
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?with_localizations=true", nil, "GET "+endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &cmd)

	return
}

// GuildApplicationCommandsPermissions returns permissions for application commands in a guild.
// appID       : The application ID
// guildID     : Guild ID to retrieve application commands permissions for.
func (s *Session) GuildApplicationCommandsPermissions(appID, guildID string, options ...discord.RequestOption) (permissions []*interactions.GuildCommandPermissions, err error) {
	endpoint := discord.EndpointApplicationCommandsGuildPermissions(appID, guildID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &permissions)
	return
}

// ApplicationCommandPermissions returns all permissions of an application command
// appID       : The Application ID
// guildID     : The guild ID containing the application command
// cmdID       : The command ID to retrieve the permissions of
func (s *Session) ApplicationCommandPermissions(appID, guildID, cmdID string, options ...discord.RequestOption) (permissions *interactions.GuildCommandPermissions, err error) {
	endpoint := discord.EndpointApplicationCommandPermissions(appID, guildID, cmdID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &permissions)
	return
}

// ApplicationCommandPermissionsEdit edits the permissions of an application command
// appID       : The Application ID
// guildID     : The guild ID containing the application command
// cmdID       : The command ID to edit the permissions of
// permissions : An object containing a list of permissions for the application command
//
// NOTE: Requires OAuth2 token with applications.commands.permissions.update scope
func (s *Session) ApplicationCommandPermissionsEdit(appID, guildID, cmdID string, permissions *interactions.CommandPermissionsList, options ...discord.RequestOption) (err error) {
	endpoint := discord.EndpointApplicationCommandPermissions(appID, guildID, cmdID)

	_, err = s.RequestWithBucketID("PUT", endpoint, permissions, endpoint, options...)
	return
}

// ApplicationCommandPermissionsBatchEdit edits the permissions of a batch of commands
// appID       : The Application ID
// guildID     : The guild ID to batch edit commands of
// permissions : A list of permissions paired with a command ID, guild ID, and application ID per application command
//
// NOTE: This endpoint has been disabled with updates to command permissions (Permissions v2). Please use ApplicationCommandPermissionsEdit instead.
func (s *Session) ApplicationCommandPermissionsBatchEdit(appID, guildID string, permissions []*interactions.GuildCommandPermissions, options ...discord.RequestOption) (err error) {
	endpoint := discord.EndpointApplicationCommandsGuildPermissions(appID, guildID)

	_, err = s.RequestWithBucketID("PUT", endpoint, permissions, endpoint, options...)
	return
}

// InteractionRespond creates the response to an interaction.
// interaction : Interaction instance.
// resp        : Response message data.
func (s *Session) InteractionRespond(interaction *interactions.Interaction, resp *interactions.InteractionResponse, options ...discord.RequestOption) error {
	endpoint := discord.EndpointInteractionResponse(interaction.ID, interaction.Token)

	if resp.Data != nil && len(resp.Data.Files) > 0 {
		contentType, body, err := channel.MultipartBodyWithJSON(resp, resp.Data.Files)
		if err != nil {
			return err
		}

		_, err = s.RequestRaw("POST", endpoint, contentType, body, endpoint, 0, options...)
		return err
	}

	_, err := s.RequestWithBucketID("POST", endpoint, *resp, endpoint, options...)
	return err
}

// InteractionResponse gets the response to an interaction.
// interaction : Interaction instance.
func (s *Session) InteractionResponse(interaction *interactions.Interaction, options ...discord.RequestOption) (*channel.Message, error) {
	return s.WebhookMessage(interaction.AppID, interaction.Token, "@original", options...)
}

// InteractionResponseEdit edits the response to an interaction.
// interaction : Interaction instance.
// newresp     : Updated response message data.
func (s *Session) InteractionResponseEdit(interaction *interactions.Interaction, newresp *channel.WebhookEdit, options ...discord.RequestOption) (*channel.Message, error) {
	return s.WebhookMessageEdit(interaction.AppID, interaction.Token, "@original", newresp, options...)
}

// InteractionResponseDelete deletes the response to an interaction.
// interaction : Interaction instance.
func (s *Session) InteractionResponseDelete(interaction *interactions.Interaction, options ...discord.RequestOption) error {
	endpoint := discord.EndpointInteractionResponseActions(interaction.AppID, interaction.Token)

	_, err := s.RequestWithBucketID("DELETE", endpoint, nil, endpoint, options...)

	return err
}

// FollowupMessageCreate creates the followup message for an interaction.
// interaction : Interaction instance.
// wait        : Waits for server confirmation of message send and ensures that the return struct is populated (it is nil otherwise)
// data        : Data of the message to send.
func (s *Session) FollowupMessageCreate(interaction *interactions.Interaction, wait bool, data *channel.WebhookParams, options ...discord.RequestOption) (*channel.Message, error) {
	return s.WebhookExecute(interaction.AppID, interaction.Token, wait, data, options...)
}

// FollowupMessageEdit edits a followup message of an interaction.
// interaction : Interaction instance.
// messageID   : The followup message ID.
// data        : Data to update the message
func (s *Session) FollowupMessageEdit(interaction *interactions.Interaction, messageID string, data *channel.WebhookEdit, options ...discord.RequestOption) (*channel.Message, error) {
	return s.WebhookMessageEdit(interaction.AppID, interaction.Token, messageID, data, options...)
}

// FollowupMessageDelete deletes a followup message of an interaction.
// interaction : Interaction instance.
// messageID   : The followup message ID.
func (s *Session) FollowupMessageDelete(interaction *interactions.Interaction, messageID string, options ...discord.RequestOption) error {
	return s.WebhookMessageDelete(interaction.AppID, interaction.Token, messageID, options...)
}

// ------------------------------------------------------------------------------------------------
// Functions specific to guilds scheduled events
// ------------------------------------------------------------------------------------------------

// GuildScheduledEvents returns an array of GuildScheduledEvent for a guild
// guildID        : The ID of a Guild
// userCount      : Whether to include the user count in the response
func (s *Session) GuildScheduledEvents(guildID string, userCount bool, options ...discord.RequestOption) (st []*guild.ScheduledEvent, err error) {
	uri := discord.EndpointGuildScheduledEvents(guildID)
	if userCount {
		uri += "?with_user_count=true"
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, discord.EndpointGuildScheduledEvents(guildID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEvent returns a specific GuildScheduledEvent in a guild
// guildID        : The ID of a Guild
// eventID        : The ID of the event
// userCount      : Whether to include the user count in the response
func (s *Session) GuildScheduledEvent(guildID, eventID string, userCount bool, options ...discord.RequestOption) (st *guild.ScheduledEvent, err error) {
	uri := discord.EndpointGuildScheduledEvent(guildID, eventID)
	if userCount {
		uri += "?with_user_count=true"
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, discord.EndpointGuildScheduledEvent(guildID, eventID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEventCreate creates a GuildScheduledEvent for a guild and returns it
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventCreate(guildID string, event *guild.ScheduledEventParams, options ...discord.RequestOption) (st *guild.ScheduledEvent, err error) {
	body, err := s.RequestWithBucketID("POST", discord.EndpointGuildScheduledEvents(guildID), event, discord.EndpointGuildScheduledEvents(guildID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEventEdit updates a specific event for a guild and returns it.
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventEdit(guildID, eventID string, event *guild.ScheduledEventParams, options ...discord.RequestOption) (st *guild.ScheduledEvent, err error) {
	body, err := s.RequestWithBucketID("PATCH", discord.EndpointGuildScheduledEvent(guildID, eventID), event, discord.EndpointGuildScheduledEvent(guildID, eventID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEventDelete deletes a specific GuildScheduledEvent in a guild
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventDelete(guildID, eventID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID("DELETE", discord.EndpointGuildScheduledEvent(guildID, eventID), nil, discord.EndpointGuildScheduledEvent(guildID, eventID), options...)
	return
}

// GuildScheduledEventUsers returns an array of GuildScheduledEventUser for a particular event in a guild
// guildID    : The ID of a Guild
// eventID    : The ID of the event
// limit      : The maximum number of users to return (Max 100)
// withMember : Whether to include the member object in the response
// beforeID   : If is not empty all returned users entries will be before the given ID
// afterID    : If is not empty all returned users entries will be after the given ID
func (s *Session) GuildScheduledEventUsers(guildID, eventID string, limit int, withMember bool, beforeID, afterID string, options ...discord.RequestOption) (st []*guild.ScheduledEventUser, err error) {
	uri := discord.EndpointGuildScheduledEventUsers(guildID, eventID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}
	if beforeID != "" {
		queryParams.Set("before", beforeID)
	}
	if afterID != "" {
		queryParams.Set("after", afterID)
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, discord.EndpointGuildScheduledEventUsers(guildID, eventID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildOnboarding returns onboarding configuration of a guild.
// guildID   : The ID of the guild
func (s *Session) GuildOnboarding(guildID string, options ...discord.RequestOption) (onboarding *guild.Onboarding, err error) {
	endpoint := discord.EndpointGuildOnboarding(guildID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &onboarding)
	return
}

// GuildOnboardingEdit edits onboarding configuration of a guild.
// guildID   : The ID of the guild
// o         : New GuildOnboarding data
func (s *Session) GuildOnboardingEdit(guildID string, o *guild.Onboarding, options ...discord.RequestOption) (onboarding *guild.Onboarding, err error) {
	endpoint := discord.EndpointGuildOnboarding(guildID)

	var body []byte
	body, err = s.RequestWithBucketID("PUT", endpoint, o, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &onboarding)
	return
}

// ----------------------------------------------------------------------
// Functions specific to auto moderation
// ----------------------------------------------------------------------

// AutoModerationRules returns a list of auto moderation rules.
// guildID : ID of the guild
func (s *Session) AutoModerationRules(guildID string, options ...discord.RequestOption) (st []*guild.AutoModerationRule, err error) {
	endpoint := discord.EndpointGuildAutoModerationRules(guildID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// AutoModerationRule returns an auto moderation rule.
// guildID : ID of the guild
// ruleID  : ID of the auto moderation rule
func (s *Session) AutoModerationRule(guildID, ruleID string, options ...discord.RequestOption) (st *guild.AutoModerationRule, err error) {
	endpoint := discord.EndpointGuildAutoModerationRule(guildID, ruleID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// AutoModerationRuleCreate creates an auto moderation rule with the given data and returns it.
// guildID : ID of the guild
// rule    : Rule data
func (s *Session) AutoModerationRuleCreate(guildID string, rule *guild.AutoModerationRule, options ...discord.RequestOption) (st *guild.AutoModerationRule, err error) {
	endpoint := discord.EndpointGuildAutoModerationRules(guildID)

	var body []byte
	body, err = s.RequestWithBucketID("POST", endpoint, rule, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// AutoModerationRuleEdit edits and returns the updated auto moderation rule.
// guildID : ID of the guild
// ruleID  : ID of the auto moderation rule
// rule    : New rule data
func (s *Session) AutoModerationRuleEdit(guildID, ruleID string, rule *guild.AutoModerationRule, options ...discord.RequestOption) (st *guild.AutoModerationRule, err error) {
	endpoint := discord.EndpointGuildAutoModerationRule(guildID, ruleID)

	var body []byte
	body, err = s.RequestWithBucketID("PATCH", endpoint, rule, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// AutoModerationRuleDelete deletes an auto moderation rule.
// guildID : ID of the guild
// ruleID  : ID of the auto moderation rule
func (s *Session) AutoModerationRuleDelete(guildID, ruleID string, options ...discord.RequestOption) (err error) {
	endpoint := discord.EndpointGuildAutoModerationRule(guildID, ruleID)
	_, err = s.RequestWithBucketID("DELETE", endpoint, nil, endpoint, options...)
	return
}

// ApplicationRoleConnectionMetadata returns application role connection metadata.
// appID : ID of the application
func (s *Session) ApplicationRoleConnectionMetadata(appID string) (st []*application.RoleConnectionMetadata, err error) {
	endpoint := discord.EndpointApplicationRoleConnectionMetadata(appID)
	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ApplicationRoleConnectionMetadataUpdate updates and returns application role connection metadata.
// appID    : ID of the application
// metadata : New metadata
func (s *Session) ApplicationRoleConnectionMetadataUpdate(appID string, metadata []*application.RoleConnectionMetadata) (st []*application.RoleConnectionMetadata, err error) {
	endpoint := discord.EndpointApplicationRoleConnectionMetadata(appID)
	var body []byte
	body, err = s.RequestWithBucketID("PUT", endpoint, metadata, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserApplicationRoleConnection returns user role connection to the specified application.
// appID : ID of the application
func (s *Session) UserApplicationRoleConnection(appID string) (st *application.RoleConnection, err error) {
	endpoint := discord.EndpointUserApplicationRoleConnection(appID)
	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return

}

// UserApplicationRoleConnectionUpdate updates and returns user role connection to the specified application.
// appID      : ID of the application
// connection : New ApplicationRoleConnection data
func (s *Session) UserApplicationRoleConnectionUpdate(appID string, rconn *application.RoleConnection) (st *application.RoleConnection, err error) {
	endpoint := discord.EndpointUserApplicationRoleConnection(appID)
	var body []byte
	body, err = s.RequestWithBucketID("PUT", endpoint, rconn, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ----------------------------------------------------------------------
// Functions specific to polls
// ----------------------------------------------------------------------

// PollAnswerVoters returns users who voted for a particular answer in a poll on the specified message.
// channelID : ID of the channel.
// messageID : ID of the message.
// answerID  : ID of the answer.
func (s *Session) PollAnswerVoters(channelID, messageID string, answerID int) (voters []*user.User, err error) {
	endpoint := discord.EndpointPollAnswerVoters(channelID, messageID, answerID)

	var body []byte
	body, err = s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	var r struct {
		Users []*user.User `json:"users"`
	}

	err = unmarshal(body, &r)
	if err != nil {
		return
	}

	voters = r.Users
	return
}

// PollExpire expires poll on the specified message.
// channelID : ID of the channel.
// messageID : ID of the message.
func (s *Session) PollExpire(channelID, messageID string) (msg *channel.Message, err error) {
	endpoint := discord.EndpointPollExpire(channelID, messageID)

	var body []byte
	body, err = s.RequestWithBucketID("POST", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &msg)
	return
}

// ----------------------------------------------------------------------
// Functions specific to monetization
// ----------------------------------------------------------------------

// SKUs returns all SKUs for a given application.
// appID : The ID of the application.
func (s *Session) SKUs(appID string) (skus []*premium.SKU, err error) {
	endpoint := discord.EndpointApplicationSKUs(appID)

	body, err := s.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &skus)
	return
}

// Entitlements returns all Entitlements for a given app, active and expired.
// appID			: The ID of the application.
// filterOptions	: Optional filter options; otherwise set it to nil.
func (s *Session) Entitlements(appID string, filterOptions *premium.EntitlementFilterOptions, options ...discord.RequestOption) (entitlements []*premium.Entitlement, err error) {
	endpoint := discord.EndpointEntitlements(appID)

	queryParams := url.Values{}
	if filterOptions != nil {
		if filterOptions.UserID != "" {
			queryParams.Set("user_id", filterOptions.UserID)
		}
		if filterOptions.SkuIDs != nil && len(filterOptions.SkuIDs) > 0 {
			queryParams.Set("sku_ids", strings.Join(filterOptions.SkuIDs, ","))
		}
		if filterOptions.Before != nil {
			queryParams.Set("before", filterOptions.Before.Format(time.RFC3339))
		}
		if filterOptions.After != nil {
			queryParams.Set("after", filterOptions.After.Format(time.RFC3339))
		}
		if filterOptions.Limit > 0 {
			queryParams.Set("limit", strconv.Itoa(filterOptions.Limit))
		}
		if filterOptions.GuildID != "" {
			queryParams.Set("guild_id", filterOptions.GuildID)
		}
		if filterOptions.ExcludeEnded {
			queryParams.Set("exclude_ended", "true")
		}
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?"+queryParams.Encode(), nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &entitlements)
	return
}

// EntitlementConsume marks a given One-Time Purchase for the user as consumed.
func (s *Session) EntitlementConsume(appID, entitlementID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID("POST", discord.EndpointEntitlementConsume(appID, entitlementID), nil, discord.EndpointEntitlementConsume(appID, ""), options...)
	return
}

// EntitlementTestCreate creates a test entitlement to a given SKU for a given guild or user.
// Discord will act as though that user or guild has entitlement to your premium offering.
func (s *Session) EntitlementTestCreate(appID string, data *premium.EntitlementTest, options ...discord.RequestOption) (err error) {
	endpoint := discord.EndpointEntitlements(appID)

	_, err = s.RequestWithBucketID("POST", endpoint, data, endpoint, options...)
	return
}

// EntitlementTestDelete deletes a currently-active test entitlement. Discord will act as though
// that user or guild no longer has entitlement to your premium offering.
func (s *Session) EntitlementTestDelete(appID, entitlementID string, options ...discord.RequestOption) (err error) {
	_, err = s.RequestWithBucketID("DELETE", discord.EndpointEntitlement(appID, entitlementID), nil, discord.EndpointEntitlement(appID, ""), options...)
	return
}

// Subscriptions returns all subscriptions containing the SKU.
// skuID : The ID of the SKU.
// userID : User ID for which to return subscriptions. Required except for OAuth queries.
// before : Optional timestamp to retrieve subscriptions before this time.
// after : Optional timestamp to retrieve subscriptions after this time.
// limit : Optional maximum number of subscriptions to return (1-100, default 50).
func (s *Session) Subscriptions(skuID string, userID string, before, after *time.Time, limit int, options ...discord.RequestOption) (subscriptions []*premium.Subscription, err error) {
	endpoint := discord.EndpointSubscriptions(skuID)

	queryParams := url.Values{}
	if before != nil {
		queryParams.Set("before", before.Format(time.RFC3339))
	}
	if after != nil {
		queryParams.Set("after", after.Format(time.RFC3339))
	}
	if userID != "" {
		queryParams.Set("user_id", userID)
	}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?"+queryParams.Encode(), nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &subscriptions)
	return
}

// Subscription returns a subscription by its SKU and subscription ID.
// skuID : The ID of the SKU.
// subscriptionID : The ID of the subscription.
// userID : User ID for which to return the subscription. Required except for OAuth queries.
func (s *Session) Subscription(skuID, subscriptionID, userID string, options ...discord.RequestOption) (subscription *premium.Subscription, err error) {
	endpoint := discord.EndpointSubscription(skuID, subscriptionID)

	queryParams := url.Values{}
	if userID != "" {
		// Unlike stated in the documentation, the user_id parameter is required here.
		queryParams.Set("user_id", userID)
	}

	body, err := s.RequestWithBucketID("GET", endpoint+"?"+queryParams.Encode(), nil, endpoint, options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &subscription)
	return
}

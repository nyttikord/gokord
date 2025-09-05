package discord

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

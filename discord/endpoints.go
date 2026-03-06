// Package discord contains every object used to interact with the Discord API.
package discord

import (
	"fmt"
)

// APIVersion is the Discord API version used for the REST and Websocket API.
var APIVersion uint = 10

var (
	EndpointStatus     = "https://status.discord.com/api/v2/"
	EndpointSm         = EndpointStatus + "scheduled-maintenances/"
	EndpointSmActive   = EndpointSm + "active.json"
	EndpointSmUpcoming = EndpointSm + "upcoming.json"

	EndpointDiscord          = "https://discord.com/"
	EndpointAPI              = fmt.Sprintf("%sapi/v%d", EndpointDiscord, APIVersion)
	EndpointGuilds           = EndpointAPI + "/guilds"
	EndpointChannels         = EndpointAPI + "/channels"
	EndpointUsers            = EndpointAPI + "/users"
	EndpointGateway          = EndpointAPI + "/gateway"
	EndpointGatewayBot       = EndpointGateway + "/bot"
	EndpointWebhooks         = EndpointAPI + "/webhooks"
	EndpointStickers         = EndpointAPI + "/stickers"
	EndpointStageInstances   = EndpointAPI + "/stage-instances"
	EndpointSKUs             = EndpointAPI + "/skus"
	EndpointSoundboardSounds = EndpointAPI + "/soundboard-default-sounds"
	EndpointApplications     = EndpointAPI + "/applications"
	EndpointOAuth2           = EndpointAPI + "/oauth2"
	EndpointInvites          = EndpointAPI + "/invites"

	EndpointVoice        = EndpointAPI + "/voice"
	EndpointVoiceRegions = EndpointVoice + "/regions"

	EndpointGroupIcon = func(cID, hash string) string { return EndpointCDNChannelIcons + cID + "/" + hash + ".png" }

	EndpointSticker            = func(sID uint64) string { return fmt.Sprintf("%s/%d", EndpointStickers, sID) }
	EndpointNitroStickersPacks = EndpointAPI + "/sticker-packs"

	EndpointApplicationSKUs = func(aID uint64) string {
		return EndpointApplication(aID) + "/skus"
	}

	EndpointEntitlements = func(aID uint64) string {
		return EndpointApplication(aID) + "/entitlements"
	}
	EndpointEntitlement = func(aID, eID uint64) string {
		return fmt.Sprintf("%s/%d", EndpointEntitlements(aID), eID)
	}
	EndpointEntitlementConsume = func(aID, eID uint64) string {
		return EndpointEntitlement(aID, eID) + "/consume"
	}

	EndpointSubscriptions = func(skuID string) string {
		return EndpointSKUs + "/" + skuID + "/subscriptions"
	}
	EndpointSubscription = func(skuID, subID string) string {
		return EndpointSubscriptions(skuID) + "/" + subID
	}

	EndpointInvite                     = func(iID string) string { return EndpointAPI + "/invites/" + iID }
	EndpointInviteTargetUsers          = func(iID string) string { return EndpointInvite(iID) + "/target-users" }
	EndpointInviteTargetUsersJobStatus = func(iID string) string { return EndpointInviteTargetUsers(iID) + "/job-status" }

	EndpointEmoji         = func(eID string) string { return EndpointCDN + "emojis/" + eID + ".png" }
	EndpointEmojiAnimated = func(eID string) string { return EndpointCDN + "emojis/" + eID + ".gif" }

	EndpointApplication = func(aID uint64) string {
		return fmt.Sprintf("%s/%d", EndpointApplications, +aID)
	}
	EndpointApplicationRoleConnectionMetadata = func(aID uint64) string {
		return EndpointApplication(aID) + "/role-connections/metadata"
	}

	EndpointApplicationEmojis = func(aID uint64) string { return EndpointApplication(aID) + "/emojis" }
	EndpointApplicationEmoji  = func(aID, eID uint64) string {
		return fmt.Sprintf("%s/emojis/%d", EndpointApplication(aID), eID)
	}

	EndpointOAuth2Applications = EndpointOAuth2 + "applications"
	EndpointOAuth2Application  = func(aID uint64) string {
		return fmt.Sprintf("%s/%d", EndpointOAuth2Applications, aID)
	}
	EndpointOAuth2ApplicationsBot = func(aID uint64) string {
		return fmt.Sprintf("%s/%d/bot", EndpointOAuth2Applications, aID)
	}
	EndpointOAuth2ApplicationAssets = func(aID string) string { return EndpointOAuth2Applications + "/" + aID + "/assets" }
)

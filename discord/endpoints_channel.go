package discord

import (
	"fmt"
	"strconv"
)

func EndpointChannel(cID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointChannels, cID)
}

func EndpointChannelThreads(cID uint64) string {
	return EndpointChannel(cID) + "/threads"
}
func EndpointChannelActiveThreads(cID uint64) string {
	return EndpointChannelThreads(cID) + "/active"
}
func EndpointChannelPublicArchivedThreads(cID uint64) string {
	return EndpointChannelThreads(cID) + "/archived/public"
}
func EndpointChannelPrivateArchivedThreads(cID uint64) string {
	return EndpointChannelThreads(cID) + "/archived/private"
}
func EndpointChannelJoinedPrivateArchivedThreads(cID uint64) string {
	return EndpointChannel(cID) + "/users/@me/threads/archived/private"
}
func EndpointThreadMembers(tID uint64) string {
	return EndpointChannel(tID) + "/thread-members"
}
func EndpointThreadMember(tID, mID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointThreadMembers(tID), mID)
}

func EndpointChannelPermissions(cID uint64) string {
	return EndpointChannel(cID) + "/permissions"
}
func EndpointChannelPermission(cID, tID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointChannelPermissions(cID), tID)
}

func EndpointChannelInvites(cID uint64) string {
	return EndpointChannel(cID) + "/invites"
}

func EndpointChannelTyping(cID uint64) string {
	return EndpointChannel(cID) + "/typing"
}
func EndpointChannelMessages(cID uint64) string {
	return EndpointChannel(cID) + "/messages"
}
func EndpointChannelMessage(cID, mID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointChannelMessages(cID), mID)
}
func EndpointChannelMessageThread(cID, mID uint64) string {
	return EndpointChannelMessage(cID, mID) + "/threads"
}
func EndpointChannelMessagesBulkDelete(cID uint64) string {
	return EndpointChannelMessages(cID) + "/bulk-delete"
}
func EndpointChannelMessagesPins(cID uint64) string {
	return EndpointChannelMessages(cID) + "/pins"
}
func EndpointChannelMessagePin(cID, mID uint64) string {
	return fmt.Sprintf("%s/messages/pins/%d", EndpointChannel(cID), mID)
}
func EndpointChannelMessageCrosspost(cID, mID uint64) string {
	return fmt.Sprintf("%s/messages/%d/crosspost", EndpointChannel(cID), mID)
}

func EndpointChannelFollow(cID uint64) string {
	return EndpointChannel(cID) + "/followers"
}

func EndpointChannelSoundboardSoundSend(cID uint64) string {
	return EndpointChannel(cID) + "/send-soundboard-sound"
}

func EndpointChannelWebhooks(cID uint64) string {
	return EndpointChannel(cID) + "/webhooks"
}

func EndpointWebhook(wID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointWebhooks, wID)
}
func EndpointWebhookToken(wID uint64, token string) string {
	return EndpointWebhook(wID) + "/" + token
}
func EndpointWebhookMessage(wID uint64, token, messageID string) string {
	return EndpointWebhookToken(wID, token) + "/messages/" + messageID
}

func EndpointMessageReactionsAll(cID, mID uint64) string {
	return EndpointChannelMessage(cID, mID) + "/reactions"
}
func EndpointMessageReactions(cID, mID, eID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointMessageReactionsAll(cID, mID), eID)
}
func EndpointMessageReaction(cID, mID, eID, uID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointMessageReactions(cID, mID, eID), uID)
}

func EndpointPoll(cID, mID uint64) string {
	return fmt.Sprintf("%s/polls/%d", EndpointChannel(cID), mID)
}
func EndpointPollAnswerVoters(cID, mID uint64, aID int) string {
	return EndpointPoll(cID, mID) + "/answers/" + strconv.Itoa(aID)
}
func EndpointPollExpire(cID, mID uint64) string {
	return EndpointPoll(cID, mID) + "/expire"
}

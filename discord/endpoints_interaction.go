package discord

import "fmt"

func EndpointApplicationGlobalCommands(aID uint64) string {
	return EndpointApplication(aID) + "/commands"
}
func EndpointApplicationGlobalCommand(aID, cID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointApplicationGlobalCommands(aID), cID)
}

func EndpointApplicationGuildCommands(aID, gID uint64) string {
	return fmt.Sprintf("%s/guilds/%d/commands", EndpointApplication(aID), gID)
}
func EndpointApplicationGuildCommand(aID, gID, cID uint64) string {
	return fmt.Sprintf("%s/%d", EndpointApplicationGuildCommands(aID, gID), cID)
}
func EndpointApplicationCommandPermissions(aID, gID, cID uint64) string {
	return EndpointApplicationGuildCommand(aID, gID, cID) + "/permissions"
}
func EndpointApplicationCommandsGuildPermissions(aID, gID uint64) string {
	return EndpointApplicationGuildCommands(aID, gID) + "/permissions"
}

func EndpointInteraction(aID uint64, iToken string) string {
	return fmt.Sprintf("%s/interactions/%d/%s", EndpointAPI, aID, iToken)
}
func EndpointInteractionResponse(iID uint64, iToken string) string {
	return EndpointInteraction(iID, iToken) + "/callback"
}
func EndpointInteractionResponseActions(aID uint64, iToken string) string {
	return EndpointWebhookMessage(aID, iToken, 0)
}

func EndpointFollowupMessage(aID uint64, iToken string) string {
	return EndpointWebhookToken(aID, iToken)
}
func EndpointFollowupMessageActions(aID uint64, iToken string, mID uint64) string {
	return EndpointWebhookMessage(aID, iToken, mID)
}

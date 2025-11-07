package botapi

import (
	"context"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user/status"
)

// Requester handles everything inside the bot package.
type Requester struct {
	discord.RESTRequester
	discord.WSRequester
	//State *State
}

// UpdateStatusData is provided to Requester.UpdateStatusComplex
type UpdateStatusData struct {
	IdleSince  *int               `json:"since"`
	Activities []*status.Activity `json:"activities"`
	AFK        bool               `json:"afk"`
	Status     string             `json:"status"`
}

type updateStatusOp struct {
	Op   discord.GatewayOpCode `json:"op"`
	Data UpdateStatusData      `json:"d"`
}

func newUpdateStatusData(idle int, activityType types.Activity, name, url string) *UpdateStatusData {
	usd := &UpdateStatusData{
		Status: "online",
	}

	if idle > 0 {
		usd.IdleSince = &idle
	}

	if name != "" {
		usd.Activities = []*status.Activity{{
			Name: name,
			Type: activityType,
			URL:  url,
		}}
	}

	return usd
}

// UpdateGameStatus is used to update the user's status.
// If idle>0 then set status to idle.
// If name!="" then set game.
// if otherwise, set status to active, and no activity.
func (r Requester) UpdateGameStatus(ctx context.Context, idle int, name string) (err error) {
	return r.UpdateStatusComplex(ctx, *newUpdateStatusData(idle, types.ActivityGame, name, ""))
}

// UpdateWatchStatus is used to update the user's watch status.
// If idle>0 then set status to idle.
// If name!="" then set movie/stream.
// if otherwise, set status to active, and no activity.
func (r Requester) UpdateWatchStatus(ctx context.Context, idle int, name string) (err error) {
	return r.UpdateStatusComplex(ctx, *newUpdateStatusData(idle, types.ActivityWatching, name, ""))
}

// UpdateStreamingStatus is used to update the user's streaming status.
// If idle>0 then set status to idle.
// If name!="" then set game.
// If name!="" and url!="" then set the status type to streaming with the URL set.
// if otherwise, set status to active, and no game.
func (r Requester) UpdateStreamingStatus(ctx context.Context, idle int, name string, url string) (err error) {
	gameType := types.ActivityGame
	if url != "" {
		gameType = types.ActivityStreaming
	}
	return r.UpdateStatusComplex(ctx, *newUpdateStatusData(idle, gameType, name, url))
}

// UpdateListeningStatus is used to set the user to "Listening to..."
// If name!="" then set to what user is listening to
// Else, set user to active and no activity.
func (r Requester) UpdateListeningStatus(ctx context.Context, name string) (err error) {
	return r.UpdateStatusComplex(ctx, *newUpdateStatusData(0, types.ActivityListening, name, ""))
}

// UpdateCustomStatus is used to update the user's custom status.
// If state!="" then set the custom status.
// Else, set user to active and remove the custom status.
func (r Requester) UpdateCustomStatus(ctx context.Context, state string) (err error) {
	data := UpdateStatusData{
		Status: "online",
	}

	if state != "" {
		// Discord requires a non-empty activity name, therefore we provide "Custom Status" as a placeholder.
		data.Activities = []*status.Activity{{
			Name:  "Custom Status",
			Type:  types.ActivityCustom,
			State: state,
		}}
	}

	return r.UpdateStatusComplex(ctx, data)
}

// UpdateStatusComplex allows for sending the raw status update data untouched by discordgo.
func (r Requester) UpdateStatusComplex(ctx context.Context, usd UpdateStatusData) (err error) {
	if len(usd.Activities) == 0 {
		usd.Activities = make([]*status.Activity, 0)
	}
	return r.GatewayWriteStruct(ctx, updateStatusOp{discord.GatewayOpCodePresenceUpdate, usd})
}

package user

import (
	"encoding/json"
	"github.com/nyttikord/gokord/emoji"
	"time"
)

// ClientStatus stores the online, offline, idle, or dnd status of each device of a Guild member.
type ClientStatus struct {
	Desktop Status `json:"desktop"`
	Mobile  Status `json:"mobile"`
	Web     Status `json:"web"`
}

// Status type definition
type Status string

// Constants for Status with the different current available status
const (
	StatusOnline       Status = "online"
	StatusIdle         Status = "idle"
	StatusDoNotDisturb Status = "dnd"
	StatusInvisible    Status = "invisible"
	StatusOffline      Status = "offline"
)

// A TimeStamps struct contains start and end times used in the rich presence "playing .." Game
type TimeStamps struct {
	EndTimestamp   int64 `json:"end,omitempty"`
	StartTimestamp int64 `json:"start,omitempty"`
}

// UnmarshalJSON unmarshals JSON into TimeStamps struct
func (t *TimeStamps) UnmarshalJSON(b []byte) error {
	temp := struct {
		End   float64 `json:"end,omitempty"`
		Start float64 `json:"start,omitempty"`
	}{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	t.EndTimestamp = int64(temp.End)
	t.StartTimestamp = int64(temp.Start)
	return nil
}

// An Assets struct contains assets and labels used in the rich presence "playing .." Game
type Assets struct {
	LargeImageID string `json:"large_image,omitempty"`
	SmallImageID string `json:"small_image,omitempty"`
	LargeText    string `json:"large_text,omitempty"`
	SmallText    string `json:"small_text,omitempty"`
}

// Activity defines the Activity sent with GatewayStatusUpdate
// https://discord.com/developers/docs/topics/gateway#activity-object
type Activity struct {
	Name          string       `json:"name"`
	Type          ActivityType `json:"type"`
	URL           string       `json:"url,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	ApplicationID string       `json:"application_id,omitempty"`
	State         string       `json:"state,omitempty"`
	Details       string       `json:"details,omitempty"`
	Timestamps    TimeStamps   `json:"timestamps,omitempty"`
	Emoji         emoji.Emoji  `json:"emoji,omitempty"`
	Party         Party        `json:"party,omitempty"`
	Assets        Assets       `json:"assets,omitempty"`
	Secrets       Secrets      `json:"secrets,omitempty"`
	Instance      bool         `json:"instance,omitempty"`
	Flags         int          `json:"flags,omitempty"`
}

// UnmarshalJSON is a custom unmarshaljson to make CreatedAt a time.Time instead of an int
func (activity *Activity) UnmarshalJSON(b []byte) error {
	temp := struct {
		Name          string       `json:"name"`
		Type          ActivityType `json:"type"`
		URL           string       `json:"url,omitempty"`
		CreatedAt     int64        `json:"created_at"`
		ApplicationID json.Number  `json:"application_id,omitempty"`
		State         string       `json:"state,omitempty"`
		Details       string       `json:"details,omitempty"`
		Timestamps    TimeStamps   `json:"timestamps,omitempty"`
		Emoji         emoji.Emoji  `json:"emoji,omitempty"`
		Party         Party        `json:"party,omitempty"`
		Assets        Assets       `json:"assets,omitempty"`
		Secrets       Secrets      `json:"secrets,omitempty"`
		Instance      bool         `json:"instance,omitempty"`
		Flags         int          `json:"flags,omitempty"`
	}{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	activity.ApplicationID = temp.ApplicationID.String()
	activity.CreatedAt = time.Unix(0, temp.CreatedAt*1000000)
	activity.Assets = temp.Assets
	activity.Details = temp.Details
	activity.Emoji = temp.Emoji
	activity.Flags = temp.Flags
	activity.Instance = temp.Instance
	activity.Name = temp.Name
	activity.Party = temp.Party
	activity.Secrets = temp.Secrets
	activity.State = temp.State
	activity.Timestamps = temp.Timestamps
	activity.Type = temp.Type
	activity.URL = temp.URL
	return nil
}

// Party defines the Party field in the Activity struct
// https://discord.com/developers/docs/topics/gateway#activity-object
type Party struct {
	ID   string `json:"id,omitempty"`
	Size []int  `json:"size,omitempty"`
}

// Secrets defines the Secrets field for the Activity struct
// https://discord.com/developers/docs/topics/gateway#activity-object
type Secrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

// ActivityType is the type of Activity (see ActivityType* consts) in the Activity struct
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-types
type ActivityType int

// Valid ActivityType values
const (
	ActivityTypeGame      ActivityType = 0
	ActivityTypeStreaming ActivityType = 1
	ActivityTypeListening ActivityType = 2
	ActivityTypeWatching  ActivityType = 3
	ActivityTypeCustom    ActivityType = 4
	ActivityTypeCompeting ActivityType = 5
)

// A Presence stores the online, offline, or idle and game status of Guild members.
type Presence struct {
	User         *User        `json:"user"`
	Status       Status       `json:"status"`
	Activities   []*Activity  `json:"activities"`
	Since        *int         `json:"since"`
	ClientStatus ClientStatus `json:"client_status"`
}

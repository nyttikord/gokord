package types

// Activity is the type of status.Activity.
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-types
type Activity int

const (
	ActivityGame      Activity = 0
	ActivityStreaming Activity = 1
	ActivityListening Activity = 2
	ActivityWatching  Activity = 3
	ActivityCustom    Activity = 4
	ActivityCompeting Activity = 5
)

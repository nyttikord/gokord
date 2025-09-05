package types

// Activity is the type of Activity (see Activity* consts) in the Activity struct
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-types
type Activity int

// Valid Activity values
const (
	ActivityGame      Activity = 0
	ActivityStreaming Activity = 1
	ActivityListening Activity = 2
	ActivityWatching  Activity = 3
	ActivityCustom    Activity = 4
	ActivityCompeting Activity = 5
)

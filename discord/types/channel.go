package types

// Channel is the type of channel.Channel
type Channel int

// Block contains known Channel type values
const (
	ChannelGuildText          Channel = 0
	ChannelDM                 Channel = 1
	ChannelGuildVoice         Channel = 2
	ChannelGroupDM            Channel = 3
	ChannelGuildCategory      Channel = 4
	ChannelGuildNews          Channel = 5
	ChannelGuildStore         Channel = 6
	ChannelGuildNewsThread    Channel = 10
	ChannelGuildPublicThread  Channel = 11
	ChannelGuildPrivateThread Channel = 12
	ChannelGuildStageVoice    Channel = 13
	ChannelGuildDirectory     Channel = 14
	ChannelGuildForum         Channel = 15
	ChannelGuildMedia         Channel = 16
)

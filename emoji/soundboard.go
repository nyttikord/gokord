package emoji

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/user"
)

// AnimationType is an int indicating the type of animation.
// https://discord.com/developers/docs/events/gateway-events#voice-channel-effect-send-animation-types
type AnimationType int

const (
	AnimationTypePremium AnimationType = 0
	AnimationTypeBasic   AnimationType = 1
)

// SoundboardSound represents a sound on the soundboard.
type SoundboardSound struct {
	Name    string `json:"name"`
	SoundID uint64 `json:"sound_id,string"`
	// Volume of the sound (0-1)
	Volume float64 `json:"volume"`
	// ID of the [Emoji] for this sound.
	EmojiID uint64 `json:"emoji_id,string"`
	// The unicode character of this sound's standard [Emoji].
	EmojiName string `json:"emoji_name"`
	// The ID of the [guild.Guild] that owns this sound, if it is a [guild.Guild] sound
	GuildID uint64 `json:"guild_id,string"`
	// Whether this sound is available for use, may be false due to loss of server boosts.
	Available bool `json:"available"`
	// The [user.User] that created this sound, if not default sound.
	User *user.User `json:"user"`
}

// SoundboardSoundSend is used to send a sound from the soundboard in a given [guild.Guild].
type SoundboardSoundSend struct {
	// The ID of the [SoundboardSound] to send
	SoundID uint64 `json:"sound_id,string"`
	// Guild ID of the [SoundboardSound] to send, if it is a [guild.Guild] sound.
	// Required to send a sound from another [guild.Guild].
	GuildID uint64 `json:"guild_id,omitempty,string"`
}

// SendSoundboardSound in a specific [channel.Channel].
func SendSoundboardSound(channelID uint64, data SoundboardSoundSend) Empty {
	return WrapAsEmpty(NewSimple(http.MethodPost, discord.EndpointChannelSoundboardSoundSend(channelID)).WithData(data))
}

// ListDefaultSoundboardSounds returns all default [SoundboardSound]s.
func ListDefaultSoundboardSounds() Request[[]*SoundboardSound] {
	return NewData[[]*SoundboardSound](http.MethodGet, discord.EndpointSoundboardSounds)
}

// ListGuildSoundboardSounds returns all [SoundboardSound]s in the given [guild.Guild].
func ListGuildSoundboardSounds(guildID uint64) Request[[]*SoundboardSound] {
	return NewData[[]*SoundboardSound](http.MethodGet, discord.EndpointGuildSoundboardSounds(guildID))
}

// GetGuildSoundboardSound in the given [guild.Guild].
func GetGuildSoundboardSound(guildID, soundID uint64) Request[*SoundboardSound] {
	return NewData[*SoundboardSound](http.MethodGet, discord.EndpointGuildSoundboardSound(guildID, soundID)).
		WithBucketID(discord.EndpointGuildSoundboardSounds(guildID))
}

// SoundboardSoundParams is used to create or to edit a [SoundboardSound].
type SoundboardSoundParams struct {
	Name string `json:"name"`
	// The sound file to play, base64 encoded .mp3 or .ogg file, data URI format.
	// e.g. "data:audio/ogg;base64,SUQzBAAAAAAAI1RTU..."
	//
	// Soundboard sounds have a max file size of 512kb and a max duration of 5.2 seconds.
	//
	// Required when creating a new song, but cannot be used in an edit.
	Sound string `json:"sound,omitempty"`
	// Volume of the sound (0-1).
	Volume float64 `json:"volume,omitempty"`
	// ID of the custom [emoji.Emoji] for this sound.
	EmojiID uint64 `json:"emoji_id,omitempty,string"`
	// The unicode character of this sound's standard emoji.
	EmojiName string `json:"emoji_name,omitempty"`
}

// CreateGuildSoundboardSound in the given [guild.Guild].
func CreateGuildSoundboardSound(guildID uint64, data SoundboardSoundParams) Request[*SoundboardSound] {
	return NewData[*SoundboardSound](http.MethodPost, discord.EndpointGuildSoundboardSounds(guildID)).
		WithData(data)
}

// EditGuildSoundboardSound and returns updated [SoundboardSound].
func EditGuildSoundboardSound(guildID uint64, data SoundboardSoundParams) Request[*SoundboardSound] {
	return NewData[*SoundboardSound](http.MethodPatch, discord.EndpointGuildSoundboardSounds(guildID)).
		WithData(data)
}

// DeleteGuildSoundboardSound in the given [guild.Guild].
func DeleteGuildSoundboardSound(guildID, soundID uint64) Empty {
	return WrapAsEmpty(NewSimple(http.MethodDelete, discord.EndpointGuildSoundboardSound(guildID, soundID)).
		WithBucketID(discord.EndpointGuildSoundboardSounds(guildID)))
}

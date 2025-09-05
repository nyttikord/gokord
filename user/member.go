package user

import (
	"github.com/nyttikord/gokord/endpoints"
	"time"
)

// MemberFlags represent flags of a guild member.
// https://discord.com/developers/docs/resources/guild#guild-member-object-guild-member-flags
type MemberFlags int

// Block containing known MemberFlags values.
const (
	// MemberFlagDidRejoin indicates whether the Member has left and rejoined the guild.
	MemberFlagDidRejoin MemberFlags = 1 << 0
	// MemberFlagCompletedOnboarding indicates whether the Member has completed onboarding.
	MemberFlagCompletedOnboarding MemberFlags = 1 << 1
	// MemberFlagBypassesVerification indicates whether the Member is exempt from guild verification requirements.
	MemberFlagBypassesVerification MemberFlags = 1 << 2
	// MemberFlagStartedOnboarding indicates whether the Member has started onboarding.
	MemberFlagStartedOnboarding MemberFlags = 1 << 3
)

// A Member stores user information for Guild members. A guild
// member represents a certain user's presence in a guild.
type Member struct {
	// The guild ID on which the member exists.
	GuildID string `json:"guild_id"`

	// The time at which the member joined the guild.
	JoinedAt time.Time `json:"joined_at"`

	// The nickname of the member, if they have one.
	Nick string `json:"nick"`

	// Whether the member is deafened at a guild level.
	Deaf bool `json:"deaf"`

	// Whether the member is muted at a guild level.
	Mute bool `json:"mute"`

	// The hash of the avatar for the guild member, if any.
	Avatar string `json:"avatar"`

	// The hash of the banner for the guild member, if any.
	Banner string `json:"banner"`

	// The underlying user on which the member is based.
	User *User `json:"user"`

	// A list of IDs of the roles which are possessed by the member.
	Roles []string `json:"roles"`

	// When the user used their Nitro boost on the server
	PremiumSince *time.Time `json:"premium_since"`

	// The flags of this member. This is a combination of bit masks; the presence of a certain
	// flag can be checked by performing a bitwise AND between this int and the flag.
	Flags MemberFlags `json:"flags"`

	// Is true while the member hasn't accepted the membership screen.
	Pending bool `json:"pending"`

	// Total permissions of the member in the channel, including overrides, returned when in the interaction object.
	Permissions int64 `json:"permissions,string"`

	// The time at which the member's timeout will expire.
	// Time in the past or nil if the user is not timed out.
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until"`
}

// Mention creates a member mention
func (m *Member) Mention() string {
	return "<@!" + m.User.ID + ">"
}

// AvatarURL returns the URL of the member's avatar
//
//	size:    The size of the user's avatar as a power of two
//	         if size is an empty string, no size parameter will
//	         be added to the URL.
func (m *Member) AvatarURL(size string) string {
	if m.Avatar == "" {
		return m.User.AvatarURL(size)
	}
	// The default/empty avatar case should be handled by the above condition
	return avatarURL(m.Avatar, "", endpoints.EndpointGuildMemberAvatar(m.GuildID, m.User.ID, m.Avatar),
		endpoints.EndpointGuildMemberAvatarAnimated(m.GuildID, m.User.ID, m.Avatar), size)

}

// BannerURL returns the URL of the member's banner image.
//
//	size:    The size of the desired banner image as a power of two
//	         Image size can be any power of two between 16 and 4096.
func (m *Member) BannerURL(size string) string {
	if m.Banner == "" {
		return m.User.BannerURL(size)
	}
	return bannerURL(
		m.Banner,
		endpoints.EndpointGuildMemberBanner(m.GuildID, m.User.ID, m.Banner),
		endpoints.EndpointGuildMemberBannerAnimated(m.GuildID, m.User.ID, m.Banner),
		size,
	)
}

// DisplayName returns the member's guild nickname if they have one,
// otherwise it returns their discord display name.
func (m *Member) DisplayName() string {
	if m.Nick != "" {
		return m.Nick
	}
	return m.User.DisplayName()
}

// A VoiceState stores the voice states of Guilds
type VoiceState struct {
	GuildID                 string     `json:"guild_id"`
	ChannelID               string     `json:"channel_id"`
	UserID                  string     `json:"user_id"`
	Member                  *Member    `json:"member"`
	SessionID               string     `json:"session_id"`
	Deaf                    bool       `json:"deaf"`
	Mute                    bool       `json:"mute"`
	SelfDeaf                bool       `json:"self_deaf"`
	SelfMute                bool       `json:"self_mute"`
	SelfStream              bool       `json:"self_stream"`
	SelfVideo               bool       `json:"self_video"`
	Suppress                bool       `json:"suppress"`
	RequestToSpeakTimestamp *time.Time `json:"request_to_speak_timestamp"`
}

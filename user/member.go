package user

import (
	"time"

	"github.com/nyttikord/gokord/discord"
)

// MemberFlags represent flags of a guild.Guild Member.
// https://discord.com/developers/docs/resources/guild#guild-member-object-guild-member-flags
type MemberFlags int

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

// Member stores user information for a guild.Guild member.
// This represents a certain user.User's presence in a guild.Guild.
type Member struct {
	// The GuildID on which the Member exists.
	GuildID string `json:"guild_id"`

	// The time at which the Member joined the guild.Guild.
	JoinedAt time.Time `json:"joined_at"`

	// The nickname of the Member, if they have one.
	Nick string `json:"nick"`

	// Whether the Member is deafened at a guild.Guild level.
	Deaf bool `json:"deaf"`

	// Whether the Member is muted at a guild.Guild level.
	Mute bool `json:"mute"`

	// The hash of the Avatar for the guild.Guild Member, if any.
	Avatar string `json:"avatar"`

	// The hash of the Banner for the guild.Guild Member, if any.
	Banner string `json:"banner"`

	// The underlying user.User on which the Member is based.
	User *User `json:"user"`

	// A list of IDs of the Roles which are possessed by the Member.
	Roles []string `json:"roles"`

	// Time since the Member used their Nitro boost on the guild.Guild.
	PremiumSince *time.Time `json:"premium_since"`

	// The flags of this member.
	// This is a combination of bit masks; the presence of a certain flag can be checked by performing a bitwise AND
	// between this int and the flag.
	Flags MemberFlags `json:"flags"`

	// Is true while the Member hasn't accepted the membership screen.
	Pending bool `json:"pending"`

	// Total Permissions of the Member in the channel, including overrides (only returned from an interaction.Interaction).
	Permissions int64 `json:"permissions,string"`

	// The time at which the Member's timeout will expire.
	// Time in the past or nil if the Member is not timed out.
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until"`
}

// Mention creates a Member mention.
func (m *Member) Mention() string {
	return "<@!" + m.User.ID + ">"
}

// AvatarURL returns the URL of the Member.Avatar
//
// size is the size of the avatar as a power of two if size is an empty string, no size parameter will be added to the
// URL (between 16 and 4096).
func (m *Member) AvatarURL(size string) string {
	if m.Avatar == "" {
		return m.User.AvatarURL(size)
	}
	// The default/empty avatar case should be handled by the above condition
	return discord.AvatarURL(m.Avatar, "", discord.EndpointGuildMemberAvatar(m.GuildID, m.User.ID, m.Avatar),
		discord.EndpointGuildMemberAvatarAnimated(m.GuildID, m.User.ID, m.Avatar), size)

}

// BannerURL returns the URL of the Member.Banner.
//
// size is the size of the banner as a power of two if size is an empty string, no size parameter will be added to the
// URL (between 16 and 4096).
func (m *Member) BannerURL(size string) string {
	if m.Banner == "" {
		return m.User.BannerURL(size)
	}
	return discord.BannerURL(
		m.Banner,
		discord.EndpointGuildMemberBanner(m.GuildID, m.User.ID, m.Banner),
		discord.EndpointGuildMemberBannerAnimated(m.GuildID, m.User.ID, m.Banner),
		size,
	)
}

// DisplayName returns the member's guild nickname if they have one, otherwise it returns their discord display name.
func (m *Member) DisplayName() string {
	if m.Nick != "" {
		return m.Nick
	}
	return m.User.DisplayName()
}

// VoiceState stores the voice states of guild.Guild.
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

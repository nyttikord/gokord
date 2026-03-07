// Package invite contains the Invite... and that's all...
package invite

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/nyttikord/gokord/application"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/internal/structs"
	"github.com/nyttikord/gokord/user"
)

// Invite stores all data related to a specific Discord [guild.Guild] or [channel.Channel] invite.
type Invite struct {
	Type      types.Invite     `json:"type"`
	Guild     *guild.Guild     `json:"guild"`
	Channel   *channel.Channel `json:"channel"`
	Inviter   *user.User       `json:"inviter"`
	Code      string           `json:"code"`
	CreatedAt time.Time        `json:"created_at"`
	MaxAge    int              `json:"max_age"`
	Uses      int              `json:"uses"`
	MaxUses   int              `json:"max_uses"`
	Revoked   bool             `json:"revoked"`
	Temporary bool             `json:"temporary"`
	Unique    bool             `json:"unique"`
	// See [Invite.TargetUser] and [Invite.TargetApplication].
	TargetType types.InviteTarget `json:"target_type"`
	// TargetUser is the [user.User] streaming displayed for this [Invite].
	// Requires [Invite.TargetType] to be [types.InviteTargetStream].
	// Set [Invite.TargetUserID] when creating an [Invite] to use this feature.
	TargetUser *user.User `json:"target_user,omitempty"`
	// TargetApplication is the embedded [application.Application] to open for this [Invite].
	// Requires [Invite.TargetType] to be [types.InviteTargetEmbeddedApplication].
	// Set [Invite.TargetApplicationID] when creating an [Invite] to use this feature.
	TargetApplication *application.Application `json:"target_application,omitempty"`
	// TargetUsersFile is a CSV with a single column of [user.User] IDs for all the [user.User] able to accept this
	// [Invite].
	// Does not work with a [channel.Channel] [Invite].
	TargetUsersFile []byte `json:"target_users_file,omitempty"`
	// Roles are the [guild.Role] given when the [user.User] joins the [guild.Guild].
	// Does not work with a [channel.Channel] [Invite].
	Roles []uint64 `json:"-"`

	// will only be filled when using [GetWithCounts].
	ApproximatePresenceCount int `json:"approximate_presence_count"`
	ApproximateMemberCount   int `json:"approximate_member_count"`

	ExpiresAt *time.Time `json:"expires_at"`
}

func (i *Invite) MarshalJSON() ([]byte, error) {
	type t Invite
	v := struct {
		t
		Roles []string `json:"role_ids,omitempty"`
	}{t(*i), structs.UintsToSnowflakes(i.Roles)}
	return json.Marshal(v)
}

func (i *Invite) UnmarshalJSON(data []byte) error {
	type t Invite
	var v struct {
		t
		Roles []string `json:"role_ids,omitempty"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*i = Invite(v.t)
	i.Roles = structs.SnowflakesToUints(v.Roles)
	return nil
}

// The Discord's documentation does not yet provide complete information.
// Check https://discord.com/developers/docs/resources/invite#get-target-users-job-status for more information.
type TargetUsersJobStatus struct {
	Status         types.TargetUsersJobStatus `json:"status"`
	TotalUsers     uint                       `json:"total_users"`
	ProcessedUsers uint                       `json:"processed_users"`
	CreatedAt      time.Time                  `json:"created_at"`
	CompletedAt    time.Time                  `json:"completed_at"`
	ErrorMessage   string                     `json:"error_message"`
}

// Get the [Invite].
func Get(inviteID string) Request[*Invite] {
	return GetComplex(inviteID, "", false, false)
}

// GetWithCounts returns the [Invite] including approximate [user.Member] counts.
func GetWithCounts(inviteID string) Request[*Invite] {
	return GetComplex(inviteID, "", true, false)
}

// GetComplex returns the [Invite] with the given ID including specified fields.
//
// If specified, it includes specified [guild.ScheduledEvent] with guildScheduledEventID.
// withCounts indicates whether to include approximate [user.Member] counts.
// withExpiration indicates whether to include expiration time.
func GetComplex(inviteID, guildScheduledEventID string, withCounts, withExpiration bool) Request[*Invite] {
	endpoint := discord.EndpointInvite(inviteID)
	v := url.Values{}
	if guildScheduledEventID != "" {
		v.Set("guild_scheduled_event_id", guildScheduledEventID)
	}
	if withCounts {
		v.Set("with_counts", "true")
	}
	if withExpiration {
		v.Set("with_expiration", "true")
	}

	if len(v) != 0 {
		endpoint += "?" + v.Encode()
	}

	return NewData[*Invite](http.MethodGet, endpoint).
		WithBucketID(discord.EndpointInvite(""))
}

// Delete an existing [Invite].
func Delete(inviteID string) Request[*Invite] {
	return NewData[*Invite](http.MethodDelete, discord.EndpointInvite(inviteID)).
		WithBucketID(discord.EndpointInvite(""))
}

// Accept an [Invite].
func Accept(inviteID string) Request[*Invite] {
	return NewData[*Invite](http.MethodPut, discord.EndpointInvite(inviteID)).
		WithBucketID(discord.EndpointInvite(""))
}

// GetTargetUsers returns a CSV with a single column Users containing the [user.User] IDs targetted by the [Invite].
func GetTargetUsers(inviteID string) Request[[]byte] {
	return NewSimple(http.MethodPut, discord.EndpointInviteTargetUsers(inviteID)).
		WithBucketID(discord.EndpointInvite(""))
}

// UpdateTargetUsers updates the [user.User] allowed to see and accept this
// See [Invite.TargetUsers].
/*func (r Requester) UpdateTargetUsers(ctx context.Context, inviteID string, csvFile string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		ctx,
		http.MethodGet,
		discord.EndpointInviteTargetUsers(inviteID)+"?target_users_file="+url.PathEscape(csvFile),
		nil,
		discord.EndpointInvite(""),
		options...,
	)
	return err
}*/

// The Discord's documentation does not yet provide complete information.
// Check https://discord.com/developers/docs/resources/invite#get-target-users-job-status for more information.
func GetTargetUsersJobStatus(inviteID string) Request[*TargetUsersJobStatus] {
	return NewData[*TargetUsersJobStatus](http.MethodPut, discord.EndpointInviteTargetUsersJobStatus(inviteID)).
		WithBucketID(discord.EndpointInvite(""))
}

// List returns all [Invite] for the given [channel.Channel].
func List(channelID uint64) Request[[]*Invite] {
	return NewData[[]*Invite](http.MethodDelete, discord.EndpointChannelInvites(channelID))
}

// Create a new [Invite] for the given [channel.Channel].
//
// NOTE: [Invite] must have MaxAge, MaxUses and Temporary.
func Create(channelID uint64, i Invite) Request[*Invite] {
	uID := uint64(0)
	if i.TargetUser != nil {
		uID = i.TargetUser.ID
	}
	appID := uint64(0)
	if i.TargetApplication != nil {
		appID = i.TargetApplication.ID
	}
	data := struct {
		MaxAge            int                `json:"max_age"`
		MaxUses           int                `json:"max_uses"`
		Temporary         bool               `json:"temporary"`
		Unique            bool               `json:"unique"`
		TargetType        types.InviteTarget `json:"target_type"`
		TargetUser        uint64             `json:"target_user_id,string"`
		TargetApplication uint64             `json:"target_application_id,string"`
		//TargerUsers       []byte             `json:"target_users_file,omitempty"`
		Roles []string `json:"role_ids,omitempty,string"`
	}{i.MaxAge, i.MaxUses, i.Temporary, i.Unique, i.TargetType, uID, appID, structs.UintsToSnowflakes(i.Roles)}

	req := NewData[*Invite](http.MethodPost, discord.EndpointChannelInvites(channelID)).
		WithData(data)

	if len(i.TargetUsersFile) > 0 {
		return WrapWarn(req, "InviteCreate does not support yet TargetUsersFile")
	}

	return req
}

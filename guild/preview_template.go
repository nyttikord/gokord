package guild

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/user"
)

// Preview holds data related to a specific public Discord [Guild], even if the [user.User] is not in the [Guild].
type Preview struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
	// The hash of the [Guild]'s Icon.
	//
	// Use [Preview.IconURL] to retrieve the icon itself.
	Icon string `json:"icon"`
	// The hash of the Guild's Splash.
	Splash string `json:"splash"`
	// The hash of the Guild's DiscoverySplash.
	DiscoverySplash string `json:"discovery_splash"`
	// A list of the custom [emoji.Emoji]s present in the [Guild].
	Emojis []*emoji.Emoji `json:"emojis"`
	// The list of enabled [Guild] Features
	Features []string `json:"features"`
	// Approximate number of [user.Member]s in this [Guild].
	//
	// NOTE: this field is only filled when using gokord.Session.GetWithCounts.
	ApproximateMemberCount int `json:"approximate_member_count"`
	// Approximate number of non-offline [user.Member]s in this [Guild].
	//
	// NOTE: this field is only filled when using GetWithCounts.
	ApproximatePresenceCount int    `json:"approximate_presence_count"`
	Description              string `json:"description"`
}

// IconURL returns an URL to the [Preview.Icon].
//
// size is the size of the desired icon image as a power of two.
// It can be any power of two between 16 and 4096.
func (g *Preview) IconURL(size string) string {
	return discord.IconURL(g.Icon, discord.EndpointGuildIcon(g.ID, g.Icon), discord.EndpointGuildIconAnimated(g.ID, g.Icon), size)
}

// A Template represents a replicable template for [Guild] creation.
type Template struct {
	// The unique Code for the [Guild] [Template].
	Code        string  `json:"code"`
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	// The number of times this [Template] has been used.
	UsageCount int `json:"usage_count"`
	// The ID of the [user.User] who created the [Template].
	CreatorID uint64 `json:"creator_id,string"`
	// The [user.User] who created the [Template].
	Creator *user.User `json:"creator"`
	// The timestamp of when the [Template] was created.
	CreatedAt time.Time `json:"created_at"`
	// The timestamp of when the [Template] was last synced.
	UpdatedAt time.Time `json:"updated_at"`
	// The ID of the [Guild] the [Template] was based on.
	SourceGuildID uint64 `json:"source_guild_id,string"`
	// The [Guild] snapshot this [Template] contains.
	SerializedSourceGuild *Guild `json:"serialized_source_guild"`
	// Whether the [Template] has unsynced changes.
	IsDirty bool `json:"is_dirty"`
}

// TemplateParams stores the data needed to create or update a [Template].
type TemplateParams struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// A UserGuild holds a brief version of a [Guild].
type UserGuild struct {
	ID          uint64    `json:"id,string"`
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	Owner       bool      `json:"owner"`
	Permissions int64     `json:"permissions,string"`
	Features    []Feature `json:"features"`
	// Approximate number of [user.Member]s in this [Guild].
	//
	// NOTE: this field is only filled when withCounts is true.
	ApproximateMemberCount int `json:"approximate_member_count"`
	// Approximate number of non-offline [user.Member]s in this [Guild].
	//
	// NOTE: this field is only filled when withCounts is true.
	ApproximatePresenceCount int `json:"approximate_presence_count"`
}

// GetPreview returns the [Preview] for the given public [Guild] guildID.
func GetPreview(guildID uint64) Request[*Preview] {
	return NewData[*Preview](http.MethodGet, discord.EndpointGuildPreview(guildID))
}

// GetTemplate for the given code.
func GetTemplate(code string) Request[*Template] {
	return NewData[*Template](http.MethodGet, discord.EndpointGuildTemplate(code)).
		WithBucketID(discord.EndpointGuildTemplate(""))
}

// CreateWithTemplate a [Guild] based on a [Template].
//
// code is the [Template.Code].
// name is the [Guild.Name] (2-100 characters).
// icon is the base64 encoded 128x128 image for the [Guild.Icon].
func CreateWithTemplate(templateCode, name, icon string) Request[*Guild] {
	data := struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	}{name, icon}

	return NewData[*Guild](http.MethodPost, discord.EndpointGuildTemplate(templateCode)).
		WithBucketID(discord.EndpointGuildTemplate("")).WithData(data)
}

// ListTemplates returns every [Template] of the given [Guild].
func ListTemplates(guildID uint64) Request[[]*Template] {
	return NewData[[]*Template](http.MethodGet, discord.EndpointGuildTemplates(guildID))
}

// CreateTemplate for the given [Guild].
func CreateTemplate(guildID uint64, data *TemplateParams) Request[*Template] {
	return NewData[*Template](http.MethodPost, discord.EndpointGuildTemplates(guildID)).
		WithData(data)
}

// SyncTemplate to the [Guild]'s current state.
//
// code is [Template.Code].
func SyncTemplate(guildID uint64, code string) Empty {
	req := NewSimple(http.MethodPut, discord.EndpointGuildTemplateSync(guildID, code)).
		WithBucketID(discord.EndpointGuildTemplates(guildID))
	return WrapAsEmpty(req)
}

// EditTemplate metadata of the given [Guild].
func EditTemplate(guildID uint64, code string, data *TemplateParams) Request[*Template] {
	return NewData[*Template](http.MethodPatch, discord.EndpointGuildTemplateSync(guildID, code)).
		WithBucketID(discord.EndpointGuildTemplates(guildID)).WithData(data)
}

// DeleteTemplate of the given [Guild].
func DeleteTemplate(guildID uint64, code string) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointGuildTemplateSync(guildID, code)).
		WithBucketID(discord.EndpointGuildTemplates(guildID))
	return WrapAsEmpty(req)
}

// ListUserGuilds returns an array of [UserGuild] for all [Guild]s.
//
// limit is the number of [Guild]s that can be returned (max 200).
// If beforeID is set, it will return all [Guild]s before this ID.
// If afterID is set, it will return all [Guild]s after this ID.
// Set withCounts to true if you want to include approximate [user.Member] and [status.Presence] counts.
func ListUserGuilds(limit int, beforeID, afterID string, withCounts bool) Request[[]*UserGuild] {
	v := url.Values{}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if withCounts {
		v.Set("with_counts", "true")
	}

	uri := discord.EndpointUserGuilds(0)

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	return NewData[[]*UserGuild](http.MethodGet, uri).
		WithBucketID(discord.EndpointUserGuilds(0))
}

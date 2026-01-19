package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	. "github.com/nyttikord/gokord/guild"
)

// Template returns a guild.Template for the given code.
func (r Requester) Template(code string) Request[*Template] {
	return NewSimpleData[*Template](
		r, http.MethodGet, discord.EndpointGuildTemplate(code),
	).WithBucketID(discord.EndpointGuildTemplate(""))
}

// CreateWithTemplate creates a guild.Guild based on a guild.Template.
//
// code is the Code of the guild.Template.
// name is the name of the guild.Guild (2-100 characters).
// icon is the base64 encoded 128x128 image for the Guild icon.
func (r Requester) CreateWithTemplate(templateCode, name, icon string) Request[*Guild] {
	data := struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	}{name, icon}

	return NewSimpleData[*Guild](
		r, http.MethodPost, discord.EndpointGuildTemplate(templateCode),
	).WithBucketID(discord.EndpointGuildTemplate("")).WithData(data)
}

// Templates returns every guild.Template of the given guild.Guild.
func (r Requester) Templates(guildID string) Request[[]*Template] {
	return NewSimpleData[[]*Template](
		r, http.MethodGet, discord.EndpointGuildTemplates(guildID),
	)
}

// TemplateCreate creates a guild.Template for the given guild.Guild.
func (r Requester) TemplateCreate(guildID string, data *TemplateParams) Request[*Template] {
	return NewSimpleData[*Template](
		r, http.MethodPost, discord.EndpointGuildTemplates(guildID),
	).WithData(data)
}

// TemplateSync syncs the guild.Template to the guild.Guild's current state.
//
// code is the code of the guild.Template.
func (r Requester) TemplateSync(guildID, code string) EmptyRequest {
	req := NewSimple(
		r, http.MethodPut, discord.EndpointGuildTemplateSync(guildID, code),
	).WithBucketID(discord.EndpointGuildTemplates(guildID))
	return WrapAsEmpty(req)
}

// TemplateEdit modifies the guild.Template's metadata of the given guild.Guild.
func (r Requester) TemplateEdit(guildID, code string, data *TemplateParams) Request[*Template] {
	return NewSimpleData[*Template](
		r, http.MethodPatch, discord.EndpointGuildTemplateSync(guildID, code),
	).WithBucketID(discord.EndpointGuildTemplates(guildID)).WithData(data)
}

// TemplateDelete deletes the guild.Template of the given guild.Guild.
func (r Requester) TemplateDelete(guildID, code string) EmptyRequest {
	req := NewSimple(
		r, http.MethodDelete, discord.EndpointGuildTemplateSync(guildID, code),
	).WithBucketID(discord.EndpointGuildTemplates(guildID))
	return WrapAsEmpty(req)
}

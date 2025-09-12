package guildapi

import (
	"net/http"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/guild"
)

// Template returns a guild.Template for the given code.
func (r Requester) Template(code string, options ...discord.RequestOption) (*guild.Template, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuildTemplate(code), nil, options...)
	if err != nil {
		return nil, err
	}

	var t guild.Template
	return &t, r.Unmarshal(body, &t)
}

// CreateWithTemplate creates a guild.Guild based on a guild.Template.
//
// code is the Code of the guild.Template.
// name is the name of the guild.Guild (2-100 characters).
// icon is the base64 encoded 128x128 image for the guild.Guild icon.
func (r Requester) CreateWithTemplate(templateCode, name, icon string, options ...discord.RequestOption) (*guild.Guild, error) {
	data := struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	}{name, icon}

	body, err := r.Request(http.MethodPost, discord.EndpointGuildTemplate(templateCode), data, options...)
	if err != nil {
		return nil, err
	}

	var g guild.Guild
	return &g, r.Unmarshal(body, &g)
}

// Templates returns every guild.Template of the given guild.Guild.
func (r Requester) Templates(guildID string, options ...discord.RequestOption) ([]*guild.Template, error) {
	body, err := r.Request(http.MethodGet, discord.EndpointGuildTemplates(guildID), nil, options...)
	if err != nil {
		return nil, err
	}

	var t []*guild.Template
	return t, r.Unmarshal(body, &t)
}

// TemplateCreate creates a guild.Template for the guild.Guild.
func (r Requester) TemplateCreate(guildID string, data *guild.TemplateParams, options ...discord.RequestOption) (*guild.Template, error) {
	body, err := r.Request(http.MethodPost, discord.EndpointGuildTemplates(guildID), data, options...)
	if err != nil {
		return nil, err
	}

	var t guild.Template
	return &t, r.Unmarshal(body, &t)
}

// TemplateSync syncs the guild.Template to the guild.Guild's current state
//
// code is the code of the guild.Template.
func (r Requester) TemplateSync(guildID, code string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodPut,
		discord.EndpointGuildTemplateSync(guildID, code),
		nil,
		discord.EndpointGuildTemplateSync(guildID, ""),
		options...,
	)
	return err
}

// TemplateEdit modifies the guild.Template's metadata of the given guild.Guild.
//
// code is the code of the guild.Template.
func (r Requester) TemplateEdit(guildID, code string, data *guild.TemplateParams, options ...discord.RequestOption) (*guild.Template, error) {
	body, err := r.RequestWithBucketID(
		http.MethodPatch,
		discord.EndpointGuildTemplateSync(guildID, code),
		data,
		discord.EndpointGuildTemplateSync(guildID, ""),
		options...,
	)
	if err != nil {
		return nil, err
	}

	var t guild.Template
	return &t, r.Unmarshal(body, &t)
}

// TemplateDelete deletes the guild.Template of the given guild.Guild.
//
// code is the code of the guild.Template.
func (r Requester) TemplateDelete(guildID, templateCode string, options ...discord.RequestOption) error {
	_, err := r.RequestWithBucketID(
		http.MethodDelete,
		discord.EndpointGuildTemplateSync(guildID, templateCode),
		nil,
		discord.EndpointGuildTemplateSync(guildID, ""),
		options...,
	)
	return err
}

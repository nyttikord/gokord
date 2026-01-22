package interaction

import (
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/component"
	"github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
)

func NewDeferredResponse() *Response {
	return &Response{Type: types.InteractionResponseDeferredChannelMessageWithSource}
}

// MessageResponse is a text response to an interaction.
// It helps creating a Response or a channel.WebhookEdit.
// See ModalResponse to create modal.
// See NewSimpleResponse to create one.
type MessageResponse struct {
	components []component.Message
	res        Response
}

// NewMessageResponse creates a new SimpleResponse.
func NewMessageResponse() *MessageResponse {
	return &MessageResponse{res: Response{Data: new(ResponseData)}}
}

func (r *MessageResponse) Response() *Response {
	r.res.Data.Components = make([]component.Component, len(r.components))
	for i, v := range r.components {
		r.res.Data.Components[i] = v
	}
	return &r.res
}

func (r *MessageResponse) WebhookEdit() *channel.WebhookEdit {
	data := r.res.Data
	var e channel.WebhookEdit
	e.Content = &data.Content
	e.Embeds = &data.Embeds
	e.Files = data.Files
	e.Components = &r.components
	return &e
}

func (r *MessageResponse) IsEphemeral() *MessageResponse {
	r.res.Data.Flags |= channel.MessageFlagsEphemeral
	return r
}

func (r *MessageResponse) IsComponentsV2() *MessageResponse {
	r.res.Data.Flags |= channel.MessageFlagsIsComponentsV2
	return r
}

func (r *MessageResponse) Message(s string) *MessageResponse {
	r.res.Data.Content = s
	return r
}

func (r *MessageResponse) AddEmbed(e *channel.MessageEmbed) *MessageResponse {
	r.res.Data.Embeds = append(r.res.Data.Embeds, e)
	return r
}

func (r *MessageResponse) AddComponent(c component.Message) *MessageResponse {
	r.components = append(r.components, c)
	return r
}

func (r *MessageResponse) AddFile(f *request.File) *MessageResponse {
	r.res.Data.Files = append(r.res.Data.Files, f)
	return r
}

// ModalResponse is a Modal response to an interaction.
// It helps creating a Response.
// See SimpleResponse to create a text response.
// See NewModalResponse to create one.
type ModalResponse struct {
	res Response
}

// NewModalResponse creates a new ModalResponse.
func NewModalResponse() *ModalResponse {
	return &ModalResponse{res: Response{Data: new(ResponseData)}}
}

func (r *ModalResponse) Response() *Response {
	return &r.res
}

func (r *ModalResponse) Title(s string) *ModalResponse {
	r.res.Data.Title = s
	return r
}

func (r *ModalResponse) CustomID(s string) *ModalResponse {
	r.res.Data.CustomID = s
	return r
}

func (r *ModalResponse) AddComponent(c component.Modal) *ModalResponse {
	r.res.Data.Components = append(r.res.Data.Components, c)
	return r
}

package interaction

import (
	"encoding/json"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/component"
	"github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/user"
)

// MessageComponent is an Interaction provoked by message component.
type MessageComponent struct {
	*Interaction
	Data *MessageComponentData
}

// MessageComponentData is helper function to assert the inner Data to MessageComponentData.
// Make sure to check that the Type of the interaction is types.InteractionMessageComponent before calling.
func (i *Interaction) MessageComponent() *MessageComponent {
	if i.Type != types.InteractionMessageComponent {
		panic("MessageComponent() called on interaction of type " + i.Type.String())
	}
	return &MessageComponent{Interaction: i, Data: i.Data.(*MessageComponentData)}
}

// ApplicationCommand is an Interaction provoked by application commands.
type ApplicationCommand struct {
	*Interaction
	Data *CommandInteractionData
}

// CommandData is helper function to assert the inner Data to CommandInteractionData.
// Make sure to check that the Type of the interaction is types.InteractionApplicationCommand before calling.
func (i *Interaction) Command() *ApplicationCommand {
	if i.Type != types.InteractionApplicationCommand && i.Type != types.InteractionApplicationCommandAutocomplete {
		panic("Command() called on interaction of type " + i.Type.String())
	}
	return &ApplicationCommand{Interaction: i, Data: i.Data.(*CommandInteractionData)}
}

// InteractionCommand is an Interaction provoked by modal.
type ModalSubmit struct {
	*Interaction
	Data *ModalSubmitData
}

// ModalSubmitData is helper function to assert the inner Data to ModalSubmitData.
// Make sure to check that the Type of the interaction is types.InteractionModalSubmit before calling.
func (i *Interaction) ModalSubmit() *ModalSubmit {
	if i.Type != types.InteractionModalSubmit {
		panic("ModalSubmit() called on interaction of type " + i.Type.String())
	}
	return &ModalSubmit{Interaction: i, Data: i.Data.(*ModalSubmitData)}
}

// GetUser returns the user.User of the Interaction.
func (i *Interaction) GetUser() *user.User {
	if i.Member == nil {
		return i.User
	}
	return i.Member.User
}

// Data is a common interface for all types of interaction data.
type Data interface {
	Type() types.Interaction
}

// MessageComponentData contains the data of component.Message Interaction.
type MessageComponentData struct {
	CustomID      string                       `json:"custom_id"`
	ComponentType types.Component              `json:"component_type"`
	Resolved      MessageComponentDataResolved `json:"resolved"`

	// Note: Only filled when ComponentType is types.SelectMenu.
	// Otherwise, is nil.
	Values []string `json:"values"`
}

// MessageComponentDataResolved contains the resolved data of selected option.
type MessageComponentDataResolved struct {
	Users    map[string]*user.User       `json:"users"`
	Members  map[string]*user.Member     `json:"members"`
	Roles    map[string]*guild.Role      `json:"roles"`
	Channels map[string]*channel.Channel `json:"channels"`
}

// Type returns the type of interaction data.
func (*MessageComponentData) Type() types.Interaction {
	return types.InteractionMessageComponent
}

// ModalSubmitData contains the data of modal submit Interaction.
type ModalSubmitData struct {
	CustomID   string            `json:"custom_id"`
	Components []component.Modal `json:"-"`
}

// Type returns the type of interaction data.
func (*ModalSubmitData) Type() types.Interaction {
	return types.InteractionModalSubmit
}

func (d *ModalSubmitData) UnmarshalJSON(data []byte) error {
	type t ModalSubmitData
	var v struct {
		t
		RawComponents []component.Unmarshaler `json:"components"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*d = ModalSubmitData(v.t)
	d.Components = make([]component.Modal, len(v.RawComponents))
	for i, r := range v.RawComponents {
		d.Components[i] = r.Component.(component.Modal)
	}
	return nil
}

// Response represents a response for an Interaction event.
type Response struct {
	Type types.InteractionResponse `json:"type,omitempty"`
	Data *ResponseData             `json:"data,omitempty"`
}

// ResponseData is response data for an Interaction.
type ResponseData struct {
	TTS             bool                            `json:"tts"`
	Content         string                          `json:"content"`
	Components      []component.Component           `json:"components"`
	Embeds          []*channel.MessageEmbed         `json:"embeds"`
	AllowedMentions *channel.MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Files           []*request.File                 `json:"-"`
	Attachments     *[]*channel.MessageAttachment   `json:"attachments,omitempty"`
	Poll            *channel.Poll                   `json:"poll,omitempty"`

	// NOTE: only channel.MessageFlagsSuppressEmbeds and channel.MessageFlagsEphemeral can be set.
	Flags channel.MessageFlags `json:"flags,omitempty"`

	// NOTE: autocomplete Interaction only.
	Choices []*CommandOptionChoice `json:"choices,omitempty"`

	// NOTE: modal Interaction only.
	CustomID string `json:"custom_id,omitempty"`
	// NOTE: modal Interaction only.
	Title string `json:"title,omitempty"`
}

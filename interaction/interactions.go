package interaction

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/component"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/guild"
	"github.com/nyttikord/gokord/premium"
	"github.com/nyttikord/gokord/user"
)

// Deadline is the time allowed to respond to an Interaction.
const Deadline = time.Second * 3

// Interaction represents data of an interaction.
type Interaction struct {
	ID        string            `json:"id"`
	AppID     string            `json:"application_id"`
	Type      types.Interaction `json:"type"`
	Data      Data              `json:"data"`
	GuildID   string            `json:"guild_id"`
	ChannelID string            `json:"channel_id"`

	// The Message on which Interaction was used.
	//
	// NOTE: this field is only filled when a component.Button click triggered the Interaction.
	// Otherwise, it will be nil.
	Message *channel.Message `json:"message"`

	// Bitwise set of permissions the app or bot has within the channel.Channel the Interaction was sent from.
	AppPermissions int64 `json:"app_permissions,string"`

	// The Member who invoked this Interaction.
	//
	// NOTE: this field is only filled when the slash Command was invoked in a guild.Guild;
	// if it was invoked in a DM, the User field will be filled instead.
	// Make sure to check for nil before using this field.
	//
	// See GetUser to directly get a nil-safe user.User.
	Member *user.Member `json:"member"`
	// The user.User who invoked this Interaction.
	//
	// NOTE: this field is only filled when the slash Command was invoked in a DM;
	// if it was invoked in a guild.Guild, the Member field will be filled instead.
	// Make sure to check for nil before using this field.
	//
	// See GetUser to directly get a nil-safe user.User.
	User *user.User `json:"user"`

	// The user.User's discord client discord.Locale.
	Locale discord.Locale `json:"locale"`
	// The guild.Guild's discord.Locale.
	// This defaults to discord.LocaleEnglishUS
	//
	// Note: this field is only filled when the Interaction was invoked in a guild.Guild.
	GuildLocale *discord.Locale `json:"guild_locale"`

	Context                      types.InteractionContext            `json:"context"`
	AuthorizingIntegrationOwners map[types.IntegrationInstall]string `json:"authorizing_integration_owners"`

	Token   string `json:"token"`
	Version int    `json:"version"`

	// Any entitlements for the invoking user.User, representing access to premium SKUs.
	//
	// NOTE: this field is only filled in monetized apps.
	Entitlements []*premium.Entitlement `json:"entitlements"`
}

// UnmarshalJSON is a method for unmarshalling JSON object to Interaction.
func (i *Interaction) UnmarshalJSON(raw []byte) error {
	var tmp struct {
		Interaction
		Data json.RawMessage `json:"data"`
	}
	err := json.Unmarshal(raw, &tmp)
	if err != nil {
		return err
	}

	*i = tmp.Interaction

	switch tmp.Type {
	case types.InteractionApplicationCommand, types.InteractionApplicationCommandAutocomplete:
		v := CommandInteractionData{}
		err = json.Unmarshal(tmp.Data, &v)
		if err != nil {
			return err
		}
		i.Data = &v
	case types.InteractionMessageComponent:
		v := MessageComponentData{}
		err = json.Unmarshal(tmp.Data, &v)
		if err != nil {
			return err
		}
		i.Data = &v
	case types.InteractionModalSubmit:
		v := ModalSubmitData{}
		err = json.Unmarshal(tmp.Data, &v)
		if err != nil {
			return err
		}
		i.Data = &v
	}
	return nil
}

// MessageComponentData is helper function to assert the inner Data to MessageComponentData.
// Make sure to check that the Type of the interaction is types.InteractionMessageComponent before calling.
func (i *Interaction) MessageComponentData() *MessageComponentData {
	if i.Type != types.InteractionMessageComponent {
		panic("MessageComponentData called on interaction of type " + i.Type.String())
	}
	return i.Data.(*MessageComponentData)
}

// CommandData is helper function to assert the inner Data to CommandInteractionData.
// Make sure to check that the Type of the interaction is types.InteractionApplicationCommand before calling.
func (i *Interaction) CommandData() *CommandInteractionData {
	if i.Type != types.InteractionApplicationCommand && i.Type != types.InteractionApplicationCommandAutocomplete {
		panic("CommandData called on interaction of type " + i.Type.String())
	}
	return i.Data.(*CommandInteractionData)
}

// ModalSubmitData is helper function to assert the inner Data to ModalSubmitData.
// Make sure to check that the Type of the interaction is types.InteractionModalSubmit before calling.
func (i *Interaction) ModalSubmitData() *ModalSubmitData {
	if i.Type != types.InteractionModalSubmit {
		panic("ModalSubmitData called on interaction of type " + i.Type.String())
	}
	return i.Data.(*ModalSubmitData)
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
	Files           []*channel.File                 `json:"-"`
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

// VerifyInteraction implements channel.Message verification of the Discord interactions API signing algorithm, as
// documented here:
// https://discord.com/developers/docs/interactions/receiving-and-responding#security-and-authorization
func VerifyInteraction(r *http.Request, key ed25519.PublicKey) bool {
	var msg bytes.Buffer

	signature := r.Header.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize {
		return false
	}

	timestamp := r.Header.Get("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	msg.WriteString(timestamp)

	defer r.Body.Close()
	var body bytes.Buffer

	// at the end of the function, copy the original body back into the request
	defer func() {
		r.Body = io.NopCloser(&body)
	}()

	// copy body into buffers
	_, err = io.Copy(&msg, io.TeeReader(r.Body, &body))
	if err != nil {
		return false
	}

	return ed25519.Verify(key, msg.Bytes(), sig)
}

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
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/premium"
	"github.com/nyttikord/gokord/user"
)

const (
	// Deadline is the time allowed to acknowledge an Interaction (like by sending a response or a deferred).
	Deadline = 3 * time.Second
	// DeadlineDeferred is the time allowed to use an Interaction token.
	DeadlineDeferred = 15 * time.Minute
)

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

	// Any entitlements for the invoking user.User, representing access to premium.SKU.
	//
	// NOTE: this field is only filled in monetized apps.
	Entitlements []*premium.Entitlement `json:"entitlements"`
}

// UnmarshalJSON is a method for unmarshalling JSON object to Interaction.
func (i *Interaction) UnmarshalJSON(raw []byte) error {
	type in Interaction
	var tmp struct {
		in
		Data json.RawMessage `json:"data"`
	}
	err := json.Unmarshal(raw, &tmp)
	if err != nil {
		return err
	}

	*i = Interaction(tmp.in)

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

// Package interactionhandler contains utilities helping handling interactions.
// It provides a higher API than the traditional interaction and interactionapi packages.
//
// You can register custom interaction handlers via the Manager (available with gokord.Session InteractionManager()
// method).
// An handler is linked with one interaction type and its unique identifier (like its custom id or its name).
// The context received is automatically cancelled after interaction.Deadline if nothing is sent.
// When you send something, this time is delayed to interaction.DeadlineDeferred.
// The context is always cancelled when your function returns.
package interactionhandler

import (
	"context"
	"time"

	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	. "github.com/nyttikord/gokord/interaction"
)

// Handler is for handling any interaction.Interaction types.
type Handler[T any] func(context.Context, bot.Session, T)
type handlers[T any] map[string]Handler[T]

func getCommandHandlers(ctx context.Context) handlers[*ApplicationCommand] {
	raw := ctx.Value(discord.ContextCommandHandlers)
	if raw == nil {
		return nil
	}
	return raw.(handlers[*ApplicationCommand])
}

func getMessageComponentHandlers(ctx context.Context) handlers[*MessageComponent] {
	raw := ctx.Value(discord.ContextMessageComponentHandlers)
	if raw == nil {
		return nil
	}
	return raw.(handlers[*MessageComponent])
}

func getModalSubmitHandlers(ctx context.Context) handlers[*ModalSubmit] {
	raw := ctx.Value(discord.ContextModalSubmitHandlers)
	if raw == nil {
		return nil
	}
	return raw.(handlers[*ModalSubmit])
}

// Handle handles event.InteractionCreate and redirects them.
func Handle(ctx context.Context, s bot.Session, i *event.InteractionCreate) {
	ctx, cancel := context.WithTimeout(ctx, DeadlineDeferred)
	defer cancel()
	logger := bot.Logger(ctx)
	logger = logger.With("interaction_type", i.Type)
	logger = logger.With("interaction_app", i.AppID)
	logger = logger.With("interaction_guild", i.GuildID)
	ctx = bot.SetLogger(ctx, logger)
	// using a buffered channel to avoid goroutines lock
	// this shouldn't happened because when context is cancelled, nothing must be sent through this channel
	responsec := make(chan struct{}, 1)
	ctx = context.WithValue(ctx, discord.ContextInteractionResponse, responsec)

	go func() {
		select {
		case <-time.After(Deadline):
			logger.Warn("context deadline reached")
			cancel()
		case <-ctx.Done():
		case <-responsec:
		}
	}()

	switch i.Type {
	case types.InteractionApplicationCommand:
		handleCommand(ctx, s, i.Interaction)
	case types.InteractionMessageComponent:
		handleMessageComponent(ctx, s, i.Interaction)
	case types.InteractionModalSubmit:
		handleModalSubmit(ctx, s, i.Interaction)
	default:
		logger.Debug("interaction not supported by general handler")
		return
	}
}

func handleCommand(ctx context.Context, s bot.Session, i *Interaction) {
	cmd := i.Command()
	handlers := getCommandHandlers(ctx)
	h, ok := handlers[cmd.Data.Name]
	if !ok {
		bot.Logger(ctx).Debug("command not found in handlers")
		return
	}
	h(ctx, s, cmd)
}

func handleMessageComponent(ctx context.Context, s bot.Session, i *Interaction) {
	msg := i.MessageComponent()
	handlers := getMessageComponentHandlers(ctx)
	h, ok := handlers[msg.Data.CustomID]
	if !ok {
		bot.Logger(ctx).Debug("message component not found in handlers")
		return
	}
	h(ctx, s, msg)
}

func handleModalSubmit(ctx context.Context, s bot.Session, i *Interaction) {
	modal := i.ModalSubmit()
	handlers := getModalSubmitHandlers(ctx)
	h, ok := handlers[modal.Data.CustomID]
	if !ok {
		bot.Logger(ctx).Debug("modal submit not found in handlers")
		return
	}
	h(ctx, s, modal)
}

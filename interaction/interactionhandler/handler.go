// Package interactionhandler contains utilities helping handling interactions.
// It provides a higher API than the traditional interaction and interactionapi packages.
package interactionhandler

import (
	"context"

	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	. "github.com/nyttikord/gokord/interaction"
)

// CommandHandler is for handling types.InteractionApplicationCommand.
type CommandHandler func(context.Context, bot.Session, *Interaction, *CommandInteractionData)
type commandHandlers map[string]CommandHandler

func getCommandHandlers(ctx context.Context) commandHandlers {
	raw := ctx.Value(discord.ContextCommandHandlers)
	if raw == nil {
		return nil
	}
	return raw.(commandHandlers)
}

// MessageComponentHandler is for handling types.InteractionMessageComponent.
type MessageComponentHandler func(context.Context, bot.Session, *Interaction, *MessageComponentData)
type messageComponentHandlers map[string]MessageComponentHandler

func getMessageComponentHandlers(ctx context.Context) messageComponentHandlers {
	raw := ctx.Value(discord.ContextMessageComponentHandlers)
	if raw == nil {
		return nil
	}
	return raw.(messageComponentHandlers)
}

// MessageComponentHandler is for handling types.InteractionModalSubmit.
type ModalSubmitHandler func(context.Context, bot.Session, *Interaction, *ModalSubmitData)
type modalSubmitHandlers map[string]ModalSubmitHandler

func getModalSubmitHandlers(ctx context.Context) modalSubmitHandlers {
	raw := ctx.Value(discord.ContextModalSubmitHandlers)
	if raw == nil {
		return nil
	}
	return raw.(modalSubmitHandlers)
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
	data := i.CommandData()
	handlers := getCommandHandlers(ctx)
	h, ok := handlers[data.Name]
	if !ok {
		bot.Logger(ctx).Debug("command not found in handlers")
	}
	h(ctx, s, i, data)
}

func handleMessageComponent(ctx context.Context, s bot.Session, i *Interaction) {
	data := i.MessageComponentData()
	handlers := getMessageComponentHandlers(ctx)
	h, ok := handlers[data.CustomID]
	if !ok {
		bot.Logger(ctx).Debug("message component not found in handlers")
	}
	h(ctx, s, i, data)
}

func handleModalSubmit(ctx context.Context, s bot.Session, i *Interaction) {
	data := i.ModalSubmitData()
	handlers := getModalSubmitHandlers(ctx)
	h, ok := handlers[data.CustomID]
	if !ok {
		bot.Logger(ctx).Debug("modal submit not found in handlers")
	}
	h(ctx, s, i, data)
}

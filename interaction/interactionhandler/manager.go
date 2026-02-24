package interactionhandler

import (
	"context"
	"log/slog"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/interaction"
	"github.com/nyttikord/gokord/logger"
)

// Manager manages CommandHandler, MessageComponentHandler, and ModalSubmitHandler.
type Manager struct {
	logger                   *slog.Logger
	commandHandlers          handlers[*ApplicationCommand]
	messageComponentHandlers handlers[*MessageComponent]
	modalSubmitHandlers      handlers[*ModalSubmit]
	rawHandlers              []Handler[*Interaction]
}

func NewManager(logger *slog.Logger) *Manager {
	return &Manager{
		logger:                   logger.With("module", "interaction"),
		commandHandlers:          make(handlers[*ApplicationCommand]),
		messageComponentHandlers: make(handlers[*MessageComponent]),
		modalSubmitHandlers:      make(handlers[*ModalSubmit]),
	}
}

// HandleRaw receives every interaction after specific handlers.
func (m *Manager) HandleRaw(h Handler[*Interaction]) {
	m.rawHandlers = append(m.rawHandlers, h)
}

// HandleCommand with given name.
// Returns a function that cancel the handle.
func (m *Manager) HandleCommand(name string, h Handler[*ApplicationCommand]) func() {
	if _, ok := m.commandHandlers[name]; ok {
		m.logger.WarnContext(
			logger.NewContext(context.Background(), 1),
			"command handler already registered, skipping",
			"name", name,
		)
		return func() {}
	}
	m.commandHandlers[name] = h
	return func() {
		delete(m.commandHandlers, name)
	}
}

// HandleMessageComponent with given customID.
// Returns a function that cancel the handle.
func (m *Manager) HandleMessageComponent(customID string, h Handler[*MessageComponent]) func() {
	if _, ok := m.messageComponentHandlers[customID]; ok {
		m.logger.WarnContext(
			logger.NewContext(context.Background(), 1),
			"message component handler already registered, skipping",
			"customID", customID,
		)
		return func() {}
	}
	m.messageComponentHandlers[customID] = h
	return func() {
		delete(m.messageComponentHandlers, customID)
	}
}

// HandleModalSubmit with given customID.
// Returns a function that cancel the handle.
func (m *Manager) HandleModalSubmit(customID string, h Handler[*ModalSubmit]) func() {
	if _, ok := m.modalSubmitHandlers[customID]; ok {
		m.logger.WarnContext(
			logger.NewContext(context.Background(), 1),
			"modal submit already registered, skipping",
			"customID", customID,
		)
		return func() {}
	}
	m.modalSubmitHandlers[customID] = h
	return func() {
		delete(m.modalSubmitHandlers, customID)
	}
}

// Context returns the context associated with the manager.
func (m *Manager) setContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, discord.ContextCommandHandlers, m.commandHandlers)
	ctx = context.WithValue(ctx, discord.ContextMessageComponentHandlers, m.messageComponentHandlers)
	ctx = context.WithValue(ctx, discord.ContextModalSubmitHandlers, m.modalSubmitHandlers)
	return ctx
}

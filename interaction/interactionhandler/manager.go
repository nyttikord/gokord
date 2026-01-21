package interactionhandler

import (
	"context"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/interaction"
)

// Manager manages CommandHandler, MessageComponentHandler, and ModalSubmitHandler.
type Manager struct {
	commandHandlers          handlers[*ApplicationCommand]
	messageComponentHandlers handlers[*MessageComponent]
	modalSubmitHandlers      handlers[*ModalSubmit]
}

func NewManager() *Manager {
	return &Manager{
		commandHandlers:          make(handlers[*ApplicationCommand]),
		messageComponentHandlers: make(handlers[*MessageComponent]),
		modalSubmitHandlers:      make(handlers[*ModalSubmit]),
	}
}

// HandleCommand with given name.
// Returns a function that cancel the handle.
func (m *Manager) HandleCommand(name string, h Handler[*ApplicationCommand]) func() {
	m.commandHandlers[name] = Handler[*ApplicationCommand](h)
	return func() {
		delete(m.commandHandlers, name)
	}
}

// HandleMessageComponent with given customID.
// Returns a function that cancel the handle.
func (m *Manager) HandleMessageComponent(customID string, h Handler[*MessageComponent]) func() {
	m.messageComponentHandlers[customID] = h
	return func() {
		delete(m.messageComponentHandlers, customID)
	}
}

// HandleModalSubmit with given customID.
// Returns a function that cancel the handle.
func (m *Manager) HandleModalSubmit(customID string, h Handler[*ModalSubmit]) func() {
	m.modalSubmitHandlers[customID] = h
	return func() {
		delete(m.modalSubmitHandlers, customID)
	}
}

// Context returns the context associated with the manager.
func (m *Manager) Context(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, discord.ContextCommandHandlers, m.commandHandlers)
	ctx = context.WithValue(ctx, discord.ContextMessageComponentHandlers, m.messageComponentHandlers)
	ctx = context.WithValue(ctx, discord.ContextModalSubmitHandlers, m.modalSubmitHandlers)
	return ctx
}

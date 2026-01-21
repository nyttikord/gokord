package discord

type ContextKey uint8

const (
	// bot package
	ContextLogger  ContextKey = 0
	ContextSession ContextKey = 1
	// interaction package
	ContextCommandHandlers          ContextKey = 2
	ContextMessageComponentHandlers ContextKey = 3
	ContextModalSubmitHandlers      ContextKey = 4
)

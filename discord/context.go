package discord

import (
	"context"
	"log/slog"
)

type ContextKey uint8

const (
	contextLoggerKey ContextKey = iota
	// bot package
	ContextSession
	// interaction package
	ContextCommandHandlers
	ContextMessageComponentHandlers
	ContextModalSubmitHandlers
	ContextInteractionResponse
	// request package
	ContextREST
)

// ContextLogger returns the [slog.Logger] of the current context.
func ContextLogger(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value(contextLoggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return l
}

// ContextLogger sets the [slog.Logger] of the current context.
func SetContextLogger(parent context.Context, l *slog.Logger) context.Context {
	return context.WithValue(parent, contextLoggerKey, l)
}

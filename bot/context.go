package bot

import (
	"context"
	"log/slog"

	"github.com/nyttikord/gokord/discord"
)

func Logger(ctx context.Context) *slog.Logger {
	return ctx.Value(discord.ContextLogger).(*slog.Logger)
}

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, discord.ContextLogger, logger)
}

func CreateContext(ctx context.Context, logger *slog.Logger, s Session) context.Context {
	ctx = SetLogger(ctx, logger)
	ctx = context.WithValue(ctx, discord.ContextSession, s)
	return ctx
}

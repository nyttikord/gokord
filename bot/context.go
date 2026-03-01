package bot

import (
	"context"
	"log/slog"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/request"
)

func Logger(ctx context.Context) *slog.Logger {
	return ctx.Value(discord.ContextLogger).(*slog.Logger)
}

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, discord.ContextLogger, logger)
}

func NewContext(ctx context.Context, logger *slog.Logger, s Session, rest request.REST) context.Context {
	ctx = SetLogger(ctx, logger)
	ctx = context.WithValue(ctx, discord.ContextSession, s)
	ctx = context.WithValue(ctx, discord.ContextREST, rest)
	return ctx
}

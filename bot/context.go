package bot

import (
	"context"
	"log/slog"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/request"
)

func NewContext(ctx context.Context, logger *slog.Logger, s Session, rest request.REST) context.Context {
	ctx = discord.SetContextLogger(ctx, logger)
	ctx = context.WithValue(ctx, discord.ContextSession, s)
	ctx = context.WithValue(ctx, discord.ContextREST, rest)
	return ctx
}

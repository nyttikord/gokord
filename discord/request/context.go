package request

import (
	"context"
	"log/slog"

	"github.com/nyttikord/gokord/discord"
)

func getREST(ctx context.Context) REST {
	return ctx.Value(discord.ContextREST).(REST)
}

func getLogger(ctx context.Context) *slog.Logger {
	return ctx.Value(discord.ContextLogger).(*slog.Logger).With("module", "request")
}

// Unmarshal performs a specific [json.Unmarshal] handling custom API errors.
func Unmarshal(ctx context.Context, b []byte, target any) error {
	return getREST(ctx).Unmarshal(b, target)
}

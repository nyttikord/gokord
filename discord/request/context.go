package request

import (
	"context"

	"github.com/nyttikord/gokord/discord"
)

func getREST(ctx context.Context) REST {
	return ctx.Value(discord.ContextREST).(REST)
}

// Unmarshal performs a specific [json.Unmarshal] handling custom API errors.
func Unmarshal(ctx context.Context, b []byte, target any) error {
	return getREST(ctx).Unmarshal(b, target)
}

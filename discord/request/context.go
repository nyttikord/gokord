package request

import (
	"context"

	"github.com/nyttikord/gokord/discord"
)

func ContextREST(ctx context.Context) REST {
	return ctx.Value(discord.ContextREST).(REST)
}

func SetContextREST(parent context.Context, r REST) context.Context {
	return context.WithValue(parent, discord.ContextREST, r)
}

// Unmarshal performs a specific [json.Unmarshal] handling custom API errors.
func Unmarshal(ctx context.Context, b []byte, target any) error {
	return ContextREST(ctx).Unmarshal(b, target)
}

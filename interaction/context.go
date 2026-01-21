package interaction

import (
	"context"

	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/request"
)

// getter is a trick to get needed methods on bot.Session via a context
type getters interface {
	ChannelGetter() channelGetter
	RolesGetter() rolesGetter
	UserGetter() userGetter
}

func loadGetters(ctx context.Context) getters {
	return ctx.Value(discord.ContextSession).(getters)
}

func loadRolesGetter(ctx context.Context) rolesGetter {
	return loadGetters(ctx).RolesGetter()
}

func loadChannelGetter(ctx context.Context) channelGetter {
	return loadGetters(ctx).ChannelGetter()
}

func loadUserGetter(ctx context.Context) userGetter {
	return loadGetters(ctx).UserGetter()
}

// getResponseChannel returns a channel that must be called when a response is send.
func getResponseChannel(ctx context.Context) chan<- struct{} {
	raw := ctx.Value(discord.ContextInteractionResponse)
	if raw == nil {
		return nil
	}
	return raw.(chan struct{})
}

type ResponseRequest[T any] struct {
	request.Request[T]
}

func WrapRequestAsResponse[T any](r request.Request[T]) ResponseRequest[T] {
	return ResponseRequest[T]{r}
}

func (r ResponseRequest[T]) Do(ctx context.Context) (T, error) {
	v, err := r.Request.Do(ctx)
	if err != nil {
		return v, err
	}
	responsec := getResponseChannel(ctx)
	if responsec != nil {
		responsec <- struct{}{}
	}
	return v, nil
}

type ResponseEmptyRequest struct {
	request.Empty
}

func WrapEmptyRequestAsResponse(r request.Empty) ResponseEmptyRequest {
	return ResponseEmptyRequest{r}
}

func (r ResponseEmptyRequest) Do(ctx context.Context) error {
	err := r.Empty.Do(ctx)
	if err != nil {
		return err
	}
	responsec := getResponseChannel(ctx)
	if responsec != nil {
		responsec <- struct{}{}
	}
	return nil
}

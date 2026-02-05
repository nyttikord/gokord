package interaction

import (
	"context"

	"github.com/nyttikord/gokord/discord"
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

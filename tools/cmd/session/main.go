package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/event"
)

var (
	token   string
	timeout int
)

func init() {
	flag.StringVar(&token, "token", os.Getenv("DG_TOKEN"), "token of the bot (required)")
	flag.IntVar(&timeout, "timeout", timeout, "timeout in minutes")
}

func main() {
	flag.Parse()
	if token == "" {
		println("token not set, use '-token' to set it")
		os.Exit(1)
		return
	}
	var ctx context.Context
	if timeout <= 0 {
		ctx = context.Background()
	} else {
		var cancel func()
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Minute)
		defer cancel()
	}
	dg := gokord.NewWithLogLevel("Bot "+token, slog.LevelDebug)
	dg.EventManager().AddHandler(func(ctx context.Context, s bot.Session, r *event.Ready) {
		bot.Logger(ctx).Info("bot ready")
		s.BotAPI().UpdateGameStatus(ctx, 0, "testing!")
		for _, g := range r.Guilds {
			m, err := s.GuildAPI().Member(g.ID, r.User.ID).Do(ctx)
			if err != nil {
				panic(err)
			}
			bot.Logger(ctx).Info("Who am I?", "guild", g.ID, "nick", m.DisplayName())
		}
	})
	err := dg.OpenAndBlock(ctx)
	if err != nil && !errors.Is(err, ctx.Err()) {
		panic(err)
	}
}

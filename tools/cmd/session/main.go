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
	flag.IntVar(&timeout, "timeout", 0, "timeout in minutes")
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
	dg.SyncEvents = true
	dg.EventManager().AddHandler(func(ctx context.Context, s bot.Session, r *event.Ready) {
		s.Logger().Info("bot ready")
		s.BotAPI().UpdateGameStatus(ctx, 0, "testing!")
	})
	err := dg.OpenAndBlock(ctx)
	if err != nil && !errors.Is(err, ctx.Err()) {
		panic(err)
	}
}

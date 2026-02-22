package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/interaction"

	"github.com/nyttikord/gokord"
)

func handleEcho(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	opts := i.OptionMap()
	builder := new(strings.Builder)
	if v, ok := opts["author"]; ok && v.BoolValue() {
		author := i.GetUser()
		builder.WriteString("**" + author.String() + "** says: ")
	}
	builder.WriteString(opts["message"].StringValue())

	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: builder.String(),
		},
	}).Do(ctx)

	if err != nil {
		log.Panicf("could not respond to interaction: %s", err)
	}
}

var commands = []*interaction.Command{
	{
		Name:        "echo",
		Description: "Say something through a bot",
		Options: []*interaction.CommandOption{
			{
				Name:        "message",
				Description: "Contents of the message",
				Type:        types.CommandOptionString,
				Required:    true,
			},
			{
				Name:        "author",
				Description: "Whether to prepend message's author",
				Type:        types.CommandOptionBoolean,
			},
		},
	},
}

var (
	Token = flag.String("token", "", "Bot authentication token")
	App   = flag.String("app", "", "Application ID")
	Guild = flag.String("guild", "", "Application ID")
)

func main() {
	flag.Parse()
	if *App == "" {
		log.Fatal("application id is not set")
	}

	dg := gokord.New("Bot " + *Token)
	dg.InteractionManager().HandleCommand("echo", handleEcho)

	dg.EventManager().AddHandler(func(_ context.Context, s bot.Session, r *event.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})

	_, err := dg.InteractionAPI().CommandBulkOverwrite(*App, *Guild, commands).Do(context.Background())
	if err != nil {
		log.Fatalf("could not register commands: %s", err)
	}

	err = dg.Open(context.Background())
	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = dg.Close(context.Background())
	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}

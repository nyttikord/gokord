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
	"github.com/nyttikord/gokord/user"

	"github.com/nyttikord/gokord"
)

type optionMap = map[string]*interaction.CommandInteractionDataOption

func parseOptions(options []*interaction.CommandInteractionDataOption) (om optionMap) {
	om = make(optionMap)
	for _, opt := range options {
		om[opt.Name] = opt
	}
	return
}

func interactionAuthor(i *interaction.Interaction) *user.User {
	if i.Member != nil {
		return i.Member.User
	}
	return i.User
}

func handleEcho(ctx context.Context, s bot.Session, i *event.InteractionCreate, opts optionMap) {
	builder := new(strings.Builder)
	if v, ok := opts["author"]; ok && v.BoolValue() {
		author := interactionAuthor(i.Interaction)
		builder.WriteString("**" + author.String() + "** says: ")
	}
	builder.WriteString(opts["message"].StringValue())

	err := s.InteractionAPI().Respond(ctx, i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: builder.String(),
		},
	})

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

	session := gokord.New("Bot " + *Token)

	session.EventManager().AddHandler(func(ctx context.Context, s bot.Session, i *event.InteractionCreate) {
		if i.Type != types.InteractionApplicationCommand {
			return
		}

		data := i.CommandData()
		if data.Name != "echo" {
			return
		}

		handleEcho(ctx, s, i, parseOptions(data.Options))
	})

	session.EventManager().AddHandler(func(_ context.Context, s bot.Session, r *event.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})

	_, err := session.InteractionAPI().CommandBulkOverwrite(context.Background(), *App, *Guild, commands)
	if err != nil {
		log.Fatalf("could not register commands: %s", err)
	}

	err = session.Open(context.Background())
	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = session.Close(context.Background())
	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/nyttikord/gokord/discord/types"
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

func handleEcho(s *gokord.Session, i *gokord.InteractionCreate, opts optionMap) {
	builder := new(strings.Builder)
	if v, ok := opts["author"]; ok && v.BoolValue() {
		author := interactionAuthor(i.Interaction)
		builder.WriteString("**" + author.String() + "** says: ")
	}
	builder.WriteString(opts["message"].StringValue())

	err := s.InteractionAPI().Respond(i.Interaction, &interaction.InteractionResponse{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.InteractionResponseData{
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
				Type:        types.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "author",
				Description: "Whether to prepend message's author",
				Type:        types.ApplicationCommandOptionBoolean,
			},
		},
	},
}

var (
	Token = flag.String("token", "", "Bot authentication token")
	App   = flag.String("app", "", "Get ID")
	Guild = flag.String("guild", "", "Get ID")
)

func main() {
	flag.Parse()
	if *App == "" {
		log.Fatal("application id is not set")
	}

	session := gokord.New("Bot " + *Token)

	session.AddHandler(func(s *gokord.Session, i *gokord.InteractionCreate) {
		if i.Type != types.InteractionApplicationCommand {
			return
		}

		data := i.ApplicationCommandData()
		if data.Name != "echo" {
			return
		}

		handleEcho(s, i, parseOptions(data.Options))
	})

	session.AddHandler(func(s *gokord.Session, r *gokord.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})

	_, err := session.InteractionAPI().CommandBulkOverwrite(*App, *Guild, commands)
	if err != nil {
		log.Fatalf("could not register commands: %s", err)
	}

	err = session.Open()
	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = session.Close()
	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}

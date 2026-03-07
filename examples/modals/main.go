package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/component"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/interaction"

	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/nyttikord/gokord"
)

// Bot parameters
var (
	GuildID        = flag.Uint("guild", 0, "Test guild ID")
	BotToken       = flag.String("token", "", "Bot access token")
	AppID          = flag.Uint("app", 0, "Get ID")
	Cleanup        = flag.Bool("cleanup", true, "Cleanup of commands")
	ResultsChannel = flag.Uint("results", 0, "Get where send survey results to")
)

var s *gokord.Session

func init() {
	flag.Parse()
}

func init() {
	s = gokord.New("Bot " + *BotToken)
}

var (
	commands = []interaction.Command{
		{
			Name:        "modals-survey",
			Description: "Take a survey about modals",
		},
	}
)

func handleModalsSurvey(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	err := interaction.Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseModal,
		Data: &interaction.ResponseData{
			CustomID: fmt.Sprintf("modals_survey_%d", i.Interaction.Member.User.ID),
			Title:    "Modals survey",
			Components: []component.Component{
				&component.Label{
					Label: "What is your opinion on them?",
					Component: &component.TextInput{
						CustomID:    "opinion",
						Style:       component.TextInputShort,
						Placeholder: "Don't be shy, share your opinion with us",
						Required:    true,
						MaxLength:   300,
						MinLength:   10,
					},
				},
				&component.Label{
					Label: "What would you suggest to improve them?",
					Component: &component.TextInput{
						CustomID: "suggestions",

						Style:     component.TextInputParagraph,
						Required:  false,
						MaxLength: 2000,
					},
				},
			},
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}

func main() {
	s.EventManager().AddHandler(func(_ context.Context, s bot.Session, r *event.Ready) {
		log.Println("Bot is up!")
	})

	s.InteractionManager().HandleCommand("modals-survey", handleModalsSurvey)
	s.InteractionManager().HandleRaw(func(ctx context.Context, s bot.Session, i *interaction.Interaction) {
		if i.Type != types.InteractionModalSubmit {
			return
		}
		data := i.ModalSubmit().Data

		if !strings.HasPrefix(data.CustomID, "modals_survey") {
			return
		}

		err := interaction.Respond(i, &interaction.Response{
			Type: types.InteractionResponseChannelMessageWithSource,
			Data: &interaction.ResponseData{
				Content: "Thank you for taking your time to fill this survey",
				Flags:   channel.MessageFlagsEphemeral,
			},
		}).Do(ctx)
		if err != nil {
			panic(err)
		}

		userid := strings.Split(data.CustomID, "_")[2]
		_, err = channel.SendMessage(uint64(*ResultsChannel), fmt.Sprintf(
			"Feedback received. From <@%s>\n\n**Opinion**:\n%s\n\n**Suggestions**:\n%s",
			userid,
			data.Components[0].(*component.Label).Component.(*component.TextInput).Value,
			data.Components[1].(*component.Label).Component.(*component.TextInput).Value,
		)).Do(ctx)
		if err != nil {
			panic(err)
		}
	})

	cmdIDs := make(map[uint64]string, len(commands))

	ctx := s.NewContext(context.Background())

	for _, cmd := range commands {
		rcmd, err := interaction.CreateCommand(uint64(*AppID), uint64(*GuildID), &cmd).Do(ctx)
		if err != nil {
			log.Fatalf("Cannot create slash command %q: %v", cmd.Name, err)
		}

		cmdIDs[rcmd.ID] = rcmd.Name
	}

	err := s.Open(context.Background())
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

	if !*Cleanup {
		return
	}

	for id, name := range cmdIDs {
		err := interaction.DeleteCommand(uint64(*AppID), uint64(*GuildID), id).Do(ctx)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}
}

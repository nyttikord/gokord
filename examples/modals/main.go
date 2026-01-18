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
	GuildID        = flag.String("guild", "", "Test guild ID")
	BotToken       = flag.String("token", "", "Bot access token")
	AppID          = flag.String("app", "", "Get ID")
	Cleanup        = flag.Bool("cleanup", true, "Cleanup of commands")
	ResultsChannel = flag.String("results", "", "Get where send survey results to")
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
	commandsHandlers = map[string]func(ctx context.Context, s bot.Session, i *event.InteractionCreate){
		"modals-survey": func(ctx context.Context, s bot.Session, i *event.InteractionCreate) {
			err := s.InteractionAPI().Respond(ctx, i.Interaction, &interaction.Response{
				Type: types.InteractionResponseModal,
				Data: &interaction.ResponseData{
					CustomID: "modals_survey_" + i.Interaction.Member.User.ID,
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
			})
			if err != nil {
				panic(err)
			}
		},
	}
)

func main() {
	s.EventManager().AddHandler(func(_ context.Context, s bot.Session, r *event.Ready) {
		log.Println("Bot is up!")
	})

	s.EventManager().AddHandler(func(ctx context.Context, s bot.Session, i *event.InteractionCreate) {
		switch i.Type {
		case types.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.CommandData().Name]; ok {
				h(ctx, s, i)
			}
		case types.InteractionModalSubmit:
			err := s.InteractionAPI().Respond(ctx, i.Interaction, &interaction.Response{
				Type: types.InteractionResponseChannelMessageWithSource,
				Data: &interaction.ResponseData{
					Content: "Thank you for taking your time to fill this survey",
					Flags:   channel.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
			data := i.ModalSubmitData()

			if !strings.HasPrefix(data.CustomID, "modals_survey") {
				return
			}

			userid := strings.Split(data.CustomID, "_")[2]
			_, err = s.ChannelAPI().MessageSend(ctx, *ResultsChannel, fmt.Sprintf(
				"Feedback received. From <@%s>\n\n**Opinion**:\n%s\n\n**Suggestions**:\n%s",
				userid,
				data.Components[0].(*component.Label).Component.(*component.TextInput).Value,
				data.Components[1].(*component.Label).Component.(*component.TextInput).Value,
			))
			if err != nil {
				panic(err)
			}
		}
	})

	cmdIDs := make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := s.InteractionAPI().CommandCreate(context.Background(), *AppID, *GuildID, &cmd)
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
		err := s.InteractionAPI().CommandDelete(context.Background(), *AppID, *GuildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}
}

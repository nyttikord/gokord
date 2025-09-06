package main

import (
	"flag"
	"fmt"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/component"
	"github.com/nyttikord/gokord/interactions"

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
	AppID          = flag.String("app", "", "Application ID")
	Cleanup        = flag.Bool("cleanup", true, "Cleanup of commands")
	ResultsChannel = flag.String("results", "", "Channel where send survey results to")
)

var s *gokord.Session

func init() {
	flag.Parse()
}

func init() {
	var err error
	s, err = gokord.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []interactions.Command{
		{
			Name:        "modals-survey",
			Description: "Take a survey about modals",
		},
	}
	commandsHandlers = map[string]func(s *gokord.Session, i *gokord.InteractionCreate){
		"modals-survey": func(s *gokord.Session, i *gokord.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &interactions.InteractionResponse{
				Type: gokord.InteractionResponseModal,
				Data: &interactions.InteractionResponseData{
					CustomID: "modals_survey_" + i.Interaction.Member.User.ID,
					Title:    "Modals survey",
					Components: []component.Message{
						component.ActionsRow{
							Components: []component.Message{
								component.TextInput{
									CustomID:    "opinion",
									Label:       "What is your opinion on them?",
									Style:       component.TextInputShort,
									Placeholder: "Don't be shy, share your opinion with us",
									Required:    true,
									MaxLength:   300,
									MinLength:   10,
								},
							},
						},
						component.ActionsRow{
							Components: []component.Message{
								component.TextInput{
									CustomID:  "suggestions",
									Label:     "What would you suggest to improve them?",
									Style:     component.TextInputParagraph,
									Required:  false,
									MaxLength: 2000,
								},
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
	s.AddHandler(func(s *gokord.Session, r *gokord.Ready) {
		log.Println("Bot is up!")
	})

	s.AddHandler(func(s *gokord.Session, i *gokord.InteractionCreate) {
		switch i.Type {
		case gokord.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case gokord.InteractionModalSubmit:
			err := s.InteractionRespond(i.Interaction, &interactions.InteractionResponse{
				Type: gokord.InteractionResponseChannelMessageWithSource,
				Data: &interactions.InteractionResponseData{
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
			_, err = s.ChannelMessageSend(*ResultsChannel, fmt.Sprintf(
				"Feedback received. From <@%s>\n\n**Opinion**:\n%s\n\n**Suggestions**:\n%s",
				userid,
				data.Components[0].(*component.ActionsRow).Components[0].(*component.TextInput).Value,
				data.Components[1].(*component.ActionsRow).Components[0].(*component.TextInput).Value,
			))
			if err != nil {
				panic(err)
			}
		}
	})

	cmdIDs := make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := s.ApplicationCommandCreate(*AppID, *GuildID, &cmd)
		if err != nil {
			log.Fatalf("Cannot create slash command %q: %v", cmd.Name, err)
		}

		cmdIDs[rcmd.ID] = rcmd.Name
	}

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

	if !*Cleanup {
		return
	}

	for id, name := range cmdIDs {
		err := s.ApplicationCommandDelete(*AppID, *GuildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}

}

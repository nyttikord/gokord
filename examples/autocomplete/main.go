package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/interaction"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *gokord.Session

func init() { flag.Parse() }

func init() {
	s = gokord.New("Bot " + *BotToken)
}

var (
	commands = []*interaction.Command{
		{
			Name:        "single-autocomplete",
			Description: "Showcase of single autocomplete option",
			Type:        types.CommandChat,
			Options: []*interaction.CommandOption{
				{
					Name:         "autocomplete-option",
					Description:  "Autocomplete option",
					Type:         types.CommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "multi-autocomplete",
			Description: "Showcase of multiple autocomplete option",
			Type:        types.CommandChat,
			Options: []*interaction.CommandOption{
				{
					Name:         "autocomplete-option-1",
					Description:  "Autocomplete option 1",
					Type:         types.CommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "autocomplete-option-2",
					Description:  "Autocomplete option 2",
					Type:         types.CommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s bot.Session, i *event.InteractionCreate){
		"single-autocomplete": func(s bot.Session, i *event.InteractionCreate) {
			switch i.Type {
			case types.InteractionApplicationCommand:
				data := i.CommandData()
				err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
					Type: types.InteractionResponseChannelMessageWithSource,
					Data: &interaction.ResponseData{
						Content: fmt.Sprintf(
							"You picked %q autocompletion",
							// Autocompleted options do not affect usual flow of handling application command. They are ordinary options at this stage
							data.Options[0].StringValue(),
						),
					},
				})
				if err != nil {
					panic(err)
				}
			// Autocomplete options introduce a new interaction type (8) for returning custom autocomplete results.
			case types.InteractionApplicationCommandAutocomplete:
				data := i.CommandData()
				choices := []*interaction.CommandOptionChoice{
					{
						Name:  "Autocomplete",
						Value: "autocomplete",
					},
					{
						Name:  "Autocomplete is best!",
						Value: "autocomplete_is_best",
					},
					{
						Name:  "Choice 3",
						Value: "choice3",
					},
					{
						Name:  "Choice 4",
						Value: "choice4",
					},
					{
						Name:  "Choice 5",
						Value: "choice5",
					},
					// And so on, up to 25 choices
				}

				if data.Options[0].StringValue() != "" {
					choices = append(choices, &interaction.CommandOptionChoice{
						Name:  data.Options[0].StringValue(), // To get user input you just get value of the autocomplete option.
						Value: "choice_custom",
					})
				}

				err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
					Type: types.InteractionApplicationCommandAutocompleteResult,
					Data: &interaction.ResponseData{
						Choices: choices, // This is basically the whole purpose of autocomplete interaction - return custom options to the user.
					},
				})
				if err != nil {
					panic(err)
				}
			}
		},
		"multi-autocomplete": func(s bot.Session, i *event.InteractionCreate) {
			switch i.Type {
			case types.InteractionApplicationCommand:
				data := i.CommandData()
				err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
					Type: types.InteractionResponseChannelMessageWithSource,
					Data: &interaction.ResponseData{
						Content: fmt.Sprintf(
							"Option 1: %s\nOption 2: %s",
							data.Options[0].StringValue(),
							data.Options[1].StringValue(),
						),
					},
				})
				if err != nil {
					panic(err)
				}
			case types.InteractionApplicationCommandAutocomplete:
				data := i.CommandData()
				var choices []*interaction.CommandOptionChoice
				switch {
				// In this case there are multiple autocomplete options. The Focused field shows which option user is focused on.
				case data.Options[0].Focused:
					choices = []*interaction.CommandOptionChoice{
						{
							Name:  "Autocomplete 4 first option",
							Value: "autocomplete_default",
						},
						{
							Name:  "Choice 3",
							Value: "choice3",
						},
						{
							Name:  "Choice 4",
							Value: "choice4",
						},
						{
							Name:  "Choice 5",
							Value: "choice5",
						},
					}
					if data.Options[0].StringValue() != "" {
						choices = append(choices, &interaction.CommandOptionChoice{
							Name:  data.Options[0].StringValue(),
							Value: "choice_custom",
						})
					}

				case data.Options[1].Focused:
					choices = []*interaction.CommandOptionChoice{
						{
							Name:  "Autocomplete 4 second option",
							Value: "autocomplete_1_default",
						},
						{
							Name:  "Choice 3.1",
							Value: "choice3_1",
						},
						{
							Name:  "Choice 4.1",
							Value: "choice4_1",
						},
						{
							Name:  "Choice 5.1",
							Value: "choice5_1",
						},
					}
					if data.Options[1].StringValue() != "" {
						choices = append(choices, &interaction.CommandOptionChoice{
							Name:  data.Options[1].StringValue(),
							Value: "choice_custom_2",
						})
					}
				}

				err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
					Type: types.InteractionApplicationCommandAutocompleteResult,
					Data: &interaction.ResponseData{
						Choices: choices,
					},
				})
				if err != nil {
					panic(err)
				}
			}
		},
	}
)

func main() {
	s.EventManager().AddHandler(func(s bot.Session, r *event.Ready) { log.Println("Bot is up!") })
	s.EventManager().AddHandler(func(s bot.Session, i *event.InteractionCreate) {
		if h, ok := commandHandlers[i.CommandData().Name]; ok {
			h(s, i)
		}
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	createdCommands, err := s.InteractionAPI().CommandBulkOverwrite(s.SessionState().User().ID, *GuildID, commands)

	if err != nil {
		log.Fatalf("Cannot register commands: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")

	if *RemoveCommands {
		for _, cmd := range createdCommands {
			err := s.InteractionAPI().CommandDelete(s.SessionState().User().ID, *GuildID, cmd.ID)
			if err != nil {
				log.Fatalf("Cannot delete %q command: %v", cmd.Name, err)
			}
		}
	}
}

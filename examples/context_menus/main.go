package main

import (
	"flag"
	"fmt"

	"github.com/nyttikord/gokord/channel"
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
	GuildID  = flag.String("guild", "", "Test guild ID")
	BotToken = flag.String("token", "", "Bot access token")
	AppID    = flag.String("app", "", "Application ID")
	Cleanup  = flag.Bool("cleanup", true, "Cleanup of commands")
)

var s *gokord.Session

func init() { flag.Parse() }

func init() {
	s = gokord.New("Bot " + *BotToken)
}

func searchLink(message, format, sep string) string {
	return fmt.Sprintf(format, strings.Join(
		strings.Split(
			message,
			" ",
		),
		sep,
	))
}

var (
	commands = []interaction.Command{
		{
			Name: "rickroll-em",
			Type: types.CommandUser,
		},
		{
			Name: "google-it",
			Type: types.CommandMessage,
		},
		{
			Name: "stackoverflow-it",
			Type: types.CommandMessage,
		},
		{
			Name: "godoc-it",
			Type: types.CommandMessage,
		},
		{
			Name: "discordjs-it",
			Type: types.CommandMessage,
		},
		{
			Name: "discordpy-it",
			Type: types.CommandMessage,
		},
	}
	commandsHandlers = map[string]func(s event.Session, i *event.InteractionCreate){
		"rickroll-em": func(s event.Session, i *event.InteractionCreate) {
			err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
				Type: types.InteractionResponseChannelMessageWithSource,
				Data: &interaction.ResponseData{
					Content: "Operation rickroll has begun",
					Flags:   channel.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}

			ch, err := s.UserAPI().ChannelCreate(
				i.CommandData().TargetID,
			)
			if err != nil {
				_, err = s.InteractionAPI().FollowupMessageCreate(i.Interaction, true, &channel.WebhookParams{
					Content: fmt.Sprintf("Mission failed. Cannot send a message to this user: %q", err.Error()),
					Flags:   channel.MessageFlagsEphemeral,
				})
				if err != nil {
					panic(err)
				}
			}
			_, err = s.ChannelAPI().MessageSend(
				ch.ID,
				fmt.Sprintf("%s sent you this: https://youtu.be/dQw4w9WgXcQ", i.Member.Mention()),
			)
			if err != nil {
				panic(err)
			}
		},
		"google-it": func(s event.Session, i *event.InteractionCreate) {
			err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
				Type: types.InteractionResponseChannelMessageWithSource,
				Data: &interaction.ResponseData{
					Content: searchLink(
						i.CommandData().Resolved.Messages[i.CommandData().TargetID].Content,
						"https://google.com/search?q=%s", "+"),
					Flags: channel.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"stackoverflow-it": func(s event.Session, i *event.InteractionCreate) {
			err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
				Type: types.InteractionResponseChannelMessageWithSource,
				Data: &interaction.ResponseData{
					Content: searchLink(
						i.CommandData().Resolved.Messages[i.CommandData().TargetID].Content,
						"https://stackoverflow.com/search?q=%s", "+"),
					Flags: channel.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"godoc-it": func(s event.Session, i *event.InteractionCreate) {
			err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
				Type: types.InteractionResponseChannelMessageWithSource,
				Data: &interaction.ResponseData{
					Content: searchLink(
						i.CommandData().Resolved.Messages[i.CommandData().TargetID].Content,
						"https://pkg.go.dev/search?q=%s", "+"),
					Flags: channel.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"discordjs-it": func(s event.Session, i *event.InteractionCreate) {
			err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
				Type: types.InteractionResponseChannelMessageWithSource,
				Data: &interaction.ResponseData{
					Content: searchLink(
						i.CommandData().Resolved.Messages[i.CommandData().TargetID].Content,
						"https://discord.js.org/#/docs/main/stable/search?query=%s", "+"),
					Flags: channel.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"discordpy-it": func(s event.Session, i *event.InteractionCreate) {
			err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
				Type: types.InteractionResponseChannelMessageWithSource,
				Data: &interaction.ResponseData{
					Content: searchLink(
						i.CommandData().Resolved.Messages[i.CommandData().TargetID].Content,
						"https://discordpy.readthedocs.io/en/stable/search.html?q=%s", "+"),
					Flags: channel.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
)

func main() {
	s.EventManager().AddHandler(func(s event.Session, r *event.Ready) {
		log.Println("Bot is up!")
	})

	s.EventManager().AddHandler(func(s event.Session, i *event.InteractionCreate) {
		if h, ok := commandsHandlers[i.CommandData().Name]; ok {
			h(s, i)
		}
	})

	cmdIDs := make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := s.InteractionAPI().CommandCreate(*AppID, *GuildID, &cmd)
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
		err := s.InteractionAPI().CommandDelete(*AppID, *GuildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}

}

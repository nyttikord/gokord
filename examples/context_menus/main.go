package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/nyttikord/gokord/bot"
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
		strings.Split(message, " "), sep,
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
)

func main() {
	s.EventManager().AddHandler(func(_ context.Context, s bot.Session, r *event.Ready) {
		log.Println("Bot is up!")
	})

	s.InteractionManager().HandleCommand("rickroll-em", rickrollEm)
	s.InteractionManager().HandleCommand("google-it", googleIt)
	s.InteractionManager().HandleCommand("stackoverflow-it", stackoverflowIt)
	s.InteractionManager().HandleCommand("godoc-it", godocIt)
	s.InteractionManager().HandleCommand("discordjs-it", djsIt)
	s.InteractionManager().HandleCommand("discordpy-it", dpyIt)

	cmdIDs := make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := s.InteractionAPI().CommandCreate(*AppID, *GuildID, &cmd).Do(context.Background())
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
		err := s.InteractionAPI().CommandDelete(*AppID, *GuildID, id).Do(context.Background())
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}

}
func rickrollEm(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: "Operation rickroll has begun",
			Flags:   channel.MessageFlagsEphemeral,
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}

	ch, err := s.UserAPI().ChannelCreate(i.Data.TargetID).Do(ctx)
	if err != nil {
		_, err = s.InteractionAPI().FollowupMessageCreate(i.Interaction, true, &channel.WebhookParams{
			Content: fmt.Sprintf("Mission failed. Cannot send a message to this user: %q", err.Error()),
			Flags:   channel.MessageFlagsEphemeral,
		}).Do(ctx)
		if err != nil {
			panic(err)
		}
	}
	_, err = s.ChannelAPI().MessageSend(
		ch.ID,
		fmt.Sprintf("%s sent you this: https://youtu.be/dQw4w9WgXcQ", i.Member.Mention()),
	).Do(ctx)
	if err != nil {
		panic(err)
	}
}

func googleIt(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: searchLink(
				i.Data.Resolved.Messages[i.Data.TargetID].Content,
				"https://google.com/search?q=%s", "+"),
			Flags: channel.MessageFlagsEphemeral,
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}
func stackoverflowIt(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: searchLink(
				i.Data.Resolved.Messages[i.Data.TargetID].Content,
				"https://stackoverflow.com/search?q=%s", "+"),
			Flags: channel.MessageFlagsEphemeral,
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}
func godocIt(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: searchLink(
				i.Data.Resolved.Messages[i.Data.TargetID].Content,
				"https://pkg.go.dev/search?q=%s", "+"),
			Flags: channel.MessageFlagsEphemeral,
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}
func djsIt(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: searchLink(
				i.Data.Resolved.Messages[i.Data.TargetID].Content,
				"https://discord.js.org/#/docs/main/stable/search?query=%s", "+"),
			Flags: channel.MessageFlagsEphemeral,
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}
func dpyIt(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: searchLink(
				i.Data.Resolved.Messages[i.Data.TargetID].Content,
				"https://discordpy.readthedocs.io/en/stable/search.html?q=%s", "+"),
			Flags: channel.MessageFlagsEphemeral,
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}

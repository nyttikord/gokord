package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/event"
)

// Flags
var (
	GuildID        = flag.Uint("guild", 0, "Test guild ID")
	StageChannelID = flag.Uint("stage", 0, "Test stage channel ID")
	BotToken       = flag.String("token", "", "Bot token")
)

func init() { flag.Parse() }

// To be correctly used, the bot needs to be in a guild.
// All actions must be done on a stage channel event
func main() {
	s := gokord.New("Bot " + *BotToken)
	s.EventManager().AddHandler(func(_ context.Context, s bot.Session, r *event.Ready) {
		fmt.Println("Bot is ready")
	})

	err := s.Open(context.Background())
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close(context.Background())

	ctx := s.NewContext(context.Background())

	// Create a new Stage instance on the previous channel
	si, err := channel.CreateStageInstance(&channel.StageInstanceParams{
		ChannelID:             uint64(*StageChannelID),
		Topic:                 "Amazing topic",
		PrivacyLevel:          channel.StageInstancePrivacyLevelGuildOnly,
		SendStartNotification: true,
	}).Do(ctx)
	if err != nil {
		log.Fatalf("Cannot create stage instance: %v", err)
	}
	log.Printf("Stage Instance %s has been successfully created", si.Topic)

	// Edit the stage instance with a new Topic
	si, err = channel.EditStageInstance(uint64(*StageChannelID), &channel.StageInstanceParams{
		Topic: "New amazing topic",
	}).Do(ctx)
	if err != nil {
		log.Fatalf("Cannot edit stage instance: %v", err)
	}
	log.Printf("Stage Instance %s has been successfully edited", si.Topic)

	time.Sleep(5 * time.Second)
	if err = channel.DeleteStageInstance(uint64(*StageChannelID)).Do(ctx); err != nil {
		log.Fatalf("Cannot delete stage instance: %v", err)
	}
	log.Printf("Stage Instance %s has been successfully deleted", si.Topic)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/guild"
)

// Flags
var (
	GuildID        = flag.String("guild", "", "Test guild ID")
	VoiceChannelID = flag.String("voice", "", "Test voice channel ID")
	BotToken       = flag.String("token", "", "Bot token")
)

func init() { flag.Parse() }

func main() {
	s := gokord.New("Bot " + *BotToken)
	s.EventManager().AddHandler(func(s *gokord.Session, r *event.Ready) {
		fmt.Println("Bot is ready")
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	event := createAmazingEvent(s)
	transformEventToExternalEvent(s, event)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

}

// Create a new event on guild
func createAmazingEvent(s *gokord.Session) *guild.ScheduledEvent {
	// Define the starting time (must be in future)
	startingTime := time.Now().Add(1 * time.Hour)
	// Define the ending time (must be after starting time)
	endingTime := startingTime.Add(30 * time.Minute)
	// Create the event
	scheduledEvent, err := s.GuildAPI().ScheduledEventCreate(*GuildID, &guild.ScheduledEventParams{
		Name:               "Amazing Event",
		Description:        "This event will start in 1 hour and last 30 minutes",
		ScheduledStartTime: &startingTime,
		ScheduledEndTime:   &endingTime,
		EntityType:         types.ScheduledEventEntityVoice,
		ChannelID:          *VoiceChannelID,
		PrivacyLevel:       guild.ScheduledEventPrivacyLevelGuildOnly,
	})
	if err != nil {
		log.Printf("Error creating scheduled event: %v", err)
		return nil
	}

	fmt.Println("Created scheduled event:", scheduledEvent.Name)
	return scheduledEvent
}

func transformEventToExternalEvent(s *gokord.Session, event *guild.ScheduledEvent) {
	scheduledEvent, err := s.GuildAPI().ScheduledEventEdit(*GuildID, event.ID, &guild.ScheduledEventParams{
		Name:       "Amazing Event @ Discord Website",
		EntityType: types.ScheduledEventEntityExternal,
		EntityMetadata: &guild.ScheduledEventEntityMetadata{
			Location: "https://discord.com",
		},
	})
	if err != nil {
		log.Printf("Error during transformation of scheduled voice event into external event: %v", err)
		return
	}

	fmt.Println("Created scheduled event:", scheduledEvent.Name)
}

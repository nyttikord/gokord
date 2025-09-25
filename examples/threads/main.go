package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/event"
)

// Flags
var (
	BotToken = flag.String("token", "", "Bot token")
)

const timeout time.Duration = time.Second * 10

var games map[string]time.Time = make(map[string]time.Time)

func init() { flag.Parse() }

func main() {
	s := gokord.New("Bot " + *BotToken)
	s.EventManager().AddHandler(func(s event.Session, r *event.Ready) {
		fmt.Println("Bot is ready")
	})
	s.EventManager().AddHandler(func(s event.Session, m *event.MessageCreate) {
		if strings.Contains(m.Content, "ping") {
			if ch, err := s.ChannelAPI().State.Channel(m.ChannelID); err != nil || !ch.IsThread() {
				thread, err := s.ChannelAPI().MessageThreadStartComplex(m.ChannelID, m.ID, &channel.ThreadStart{
					Name:                "Pong game with " + m.Author.Username,
					AutoArchiveDuration: 60,
					Invitable:           false,
					RateLimitPerUser:    10,
				})
				if err != nil {
					panic(err)
				}
				_, _ = s.ChannelAPI().MessageSend(thread.ID, "pong")
				m.ChannelID = thread.ID
			} else {
				_, _ = s.ChannelAPI().MessageSendReply(m.ChannelID, "pong", m.Reference())
			}
			games[m.ChannelID] = time.Now()
			<-time.After(timeout)
			if time.Since(games[m.ChannelID]) >= timeout {
				archived := true
				locked := true
				_, err := s.ChannelAPI().ChannelEdit(m.ChannelID, &channel.Edit{
					Archived: &archived,
					Locked:   &locked,
				})
				if err != nil {
					panic(err)
				}
			}
		}
	})
	s.Identify.Intents = discord.IntentsAllWithoutPrivileged

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

}

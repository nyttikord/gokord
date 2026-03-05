package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/bot"
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
	s.EventManager().AddHandler(func(_ context.Context, s bot.Session, r *event.Ready) {
		fmt.Println("Bot is ready")
	})
	s.EventManager().AddHandler(func(ctx context.Context, s bot.Session, m *event.MessageCreate) {
		if strings.Contains(m.Content, "ping") {
			if ch, err := s.ChannelState().GetChannel(m.ChannelID); err != nil || !ch.IsThread() {
				thread, err := channel.StartThreadMessageComplex(m.ChannelID, m.ID, &channel.ThreadStart{
					Name:                "Pong game with " + m.Author.Username,
					AutoArchiveDuration: 60,
					Invitable:           false,
					RateLimitPerUser:    10,
				}).Do(ctx)
				if err != nil {
					panic(err)
				}
				_, _ = channel.SendMessage(thread.ID, "pong").Do(ctx)
				m.ChannelID = thread.ID
			} else {
				_, _ = channel.SendMessageReply(m.ChannelID, "pong", m.Reference()).Do(ctx)
			}
			games[m.ChannelID] = time.Now()
			<-time.After(timeout)
			if time.Since(games[m.ChannelID]) >= timeout {
				archived := true
				locked := true
				_, err := channel.Edit(m.ChannelID, &channel.EditData{
					Archived: &archived,
					Locked:   &locked,
				}).Do(ctx)
				if err != nil {
					panic(err)
				}
			}
		}
	})
	s.Identify.Intents = discord.IntentsAllWithoutPrivileged

	err := s.Open(context.Background())
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/nyttikord/gokord"
)

// Command line flags
var (
	BotToken  = flag.String("token", "", "Bot authorization token")
	GuildID   = flag.String("guild", "", "ID of the testing guild")
	ChannelID = flag.String("channel", "", "ID of the testing channel")
)

func init() { flag.Parse() }

func main() {
	session, _ := gokord.New("Bot " + *BotToken)
	session.Identify.Intents |= gokord.IntentAutoModerationExecution
	session.Identify.Intents |= gokord.IntentMessageContent

	enabled := true
	rule, err := session.AutoModerationRuleCreate(*GuildID, &gokord.AutoModerationRule{
		Name:        "Auto Moderation example",
		EventType:   gokord.AutoModerationEventMessageSend,
		TriggerType: gokord.AutoModerationEventTriggerKeyword,
		TriggerMetadata: &gokord.AutoModerationTriggerMetadata{
			KeywordFilter: []string{"*cat*"},
			RegexPatterns: []string{"(c|b)at"},
		},

		Enabled: &enabled,
		Actions: []gokord.AutoModerationAction{
			{Type: gokord.AutoModerationRuleActionBlockMessage},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created the rule")
	defer session.AutoModerationRuleDelete(*GuildID, rule.ID)

	session.AddHandlerOnce(func(s *gokord.Session, e *gokord.AutoModerationActionExecution) {
		_, err = session.AutoModerationRuleEdit(*GuildID, rule.ID, &gokord.AutoModerationRule{
			TriggerMetadata: &gokord.AutoModerationTriggerMetadata{
				KeywordFilter: []string{"cat"},
			},
			Actions: []gokord.AutoModerationAction{
				{Type: gokord.AutoModerationRuleActionTimeout, Metadata: &gokord.AutoModerationActionMetadata{Duration: 60}},
				{Type: gokord.AutoModerationRuleActionSendAlertMessage, Metadata: &gokord.AutoModerationActionMetadata{
					ChannelID: e.ChannelID,
				}},
			},
		})
		if err != nil {
			session.AutoModerationRuleDelete(*GuildID, rule.ID)
			panic(err)
		}

		s.ChannelMessageSend(e.ChannelID, "Congratulations! You have just triggered an auto moderation rule.\n"+
			"The current trigger can match anywhere in the word, so even if you write the trigger word as a part of another word, it will still match.\n"+
			"The rule has now been changed, now the trigger matches only in the full words.\n"+
			"Additionally, when you send a message, an alert will be sent to this channel and you will be **timed out** for a minute.\n")

		var counter int
		var counterMutex sync.Mutex
		session.AddHandler(func(s *gokord.Session, e *gokord.AutoModerationActionExecution) {
			action := "unknown"
			switch e.Action.Type {
			case gokord.AutoModerationRuleActionBlockMessage:
				action = "block message"
			case gokord.AutoModerationRuleActionSendAlertMessage:
				action = "send alert message into <#" + e.Action.Metadata.ChannelID + ">"
			case gokord.AutoModerationRuleActionTimeout:
				action = "timeout"
			}

			counterMutex.Lock()
			counter++
			if counter == 1 {
				counterMutex.Unlock()
				s.ChannelMessageSend(e.ChannelID, "Nothing has changed, right? "+
					"Well, since separate gateway events are fired per each action (current is "+action+"), "+
					"you'll see a second message about an action pop up soon")
			} else if counter == 2 {
				counterMutex.Unlock()
				s.ChannelMessageSend(e.ChannelID, "Now the second ("+action+") action got executed.")
				s.ChannelMessageSend(e.ChannelID, "And... you've made it! That's the end of the example.\n"+
					"For more information about the automod and how to use it, "+
					"you can visit the official Discord docs: https://discord.dev/resources/auto-moderation or ask in our server: https://discord.gg/6dzbuDpSWY",
				)

				session.Close()
				session.AutoModerationRuleDelete(*GuildID, rule.ID)
				os.Exit(0)
			}
		})
	})

	err = session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

}

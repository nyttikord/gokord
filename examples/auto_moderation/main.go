package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/guild"
)

// Command line flags
var (
	BotToken  = flag.String("token", "", "Bot authorization token")
	GuildID   = flag.String("guild", "", "ID of the testing guild")
	ChannelID = flag.String("channel", "", "ID of the testing channel")
)

func init() { flag.Parse() }

func main() {
	session := gokord.New("Bot " + *BotToken)
	session.Identify.Intents |= discord.IntentAutoModerationExecution
	session.Identify.Intents |= discord.IntentMessageContent

	enabled := true
	rule, err := session.GuildAPI().AutoModerationRuleCreate(*GuildID, &guild.AutoModerationRule{
		Name:        "Auto Moderation example",
		EventType:   guild.AutoModerationRuleEventMessageSend,
		TriggerType: guild.AutoModerationRuleTriggerKeywordPreset,
		TriggerMetadata: &guild.AutoModerationTriggerMetadata{
			KeywordFilter: []string{"*cat*"},
			RegexPatterns: []string{"(c|b)at"},
		},

		Enabled: &enabled,
		Actions: []guild.AutoModerationAction{
			{Type: types.AutoModerationActionBlockMessage},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created the rule")
	defer session.GuildAPI().AutoModerationRuleDelete(*GuildID, rule.ID)

	session.EventManager().AddHandlerOnce(func(s event.Session, e *event.AutoModerationActionExecution) {
		_, err = session.GuildAPI().AutoModerationRuleEdit(*GuildID, rule.ID, &guild.AutoModerationRule{
			TriggerMetadata: &guild.AutoModerationTriggerMetadata{
				KeywordFilter: []string{"cat"},
			},
			Actions: []guild.AutoModerationAction{
				{Type: types.AutoModerationActionTimeout, Metadata: &guild.AutoModerationActionMetadata{Duration: 60}},
				{Type: types.AutoModerationActionSendAlertMessage, Metadata: &guild.AutoModerationActionMetadata{
					ChannelID: e.ChannelID,
				}},
			},
		})
		if err != nil {
			session.GuildAPI().AutoModerationRuleDelete(*GuildID, rule.ID)
			panic(err)
		}

		s.ChannelAPI().MessageSend(e.ChannelID, "Congratulations! You have just triggered an auto moderation rule.\n"+
			"The current trigger can match anywhere in the word, so even if you write the trigger word as a part of another word, it will still match.\n"+
			"The rule has now been changed, now the trigger matches only in the full words.\n"+
			"Additionally, when you send a message, an alert will be sent to this channel and you will be **timed out** for a minute.\n")

		var counter int
		var counterMutex sync.Mutex
		session.EventManager().AddHandler(func(s event.Session, e *event.AutoModerationActionExecution) {
			action := "unknown"
			switch e.Action.Type {
			case types.AutoModerationActionBlockMessage:
				action = "block message"
			case types.AutoModerationActionSendAlertMessage:
				action = "send alert message into <#" + e.Action.Metadata.ChannelID + ">"
			case types.AutoModerationActionTimeout:
				action = "timeout"
			}

			counterMutex.Lock()
			counter++
			if counter == 1 {
				counterMutex.Unlock()
				s.ChannelAPI().MessageSend(e.ChannelID, "Nothing has changed, right? "+
					"Well, since separate gateway events are fired per each action (current is "+action+"), "+
					"you'll see a second message about an action pop up soon")
			} else if counter == 2 {
				counterMutex.Unlock()
				s.ChannelAPI().MessageSend(e.ChannelID, "Now the second ("+action+") action got executed.")
				s.ChannelAPI().MessageSend(e.ChannelID, "And... you've made it! That's the end of the example.\n"+
					"For more information about the automod and how to use it, "+
					"you can visit the official Discord docs: https://discord.dev/resources/auto-moderation or ask in our server: https://discord.gg/6dzbuDpSWY",
				)

				session.Close()
				session.GuildAPI().AutoModerationRuleDelete(*GuildID, rule.ID)
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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/channel"
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
	dg := gokord.New("Bot " + *BotToken)
	dg.Identify.Intents |= discord.IntentAutoModerationExecution
	dg.Identify.Intents |= discord.IntentMessageContent

	ctx := dg.NewContext(context.Background())

	enabled := true
	rule, err := guild.CreateAutoModerationRule(*GuildID, &guild.AutoModerationRule{
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
	}).Do(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created the rule")
	defer guild.DeleteAutoModerationRule(*GuildID, rule.ID).Do(ctx)

	dg.EventManager().AddHandlerOnce(func(ctx context.Context, s bot.Session, e *event.AutoModerationActionExecution) {
		_, err = guild.EditAutoModerationRule(*GuildID, rule.ID, &guild.AutoModerationRule{
			TriggerMetadata: &guild.AutoModerationTriggerMetadata{
				KeywordFilter: []string{"cat"},
			},
			Actions: []guild.AutoModerationAction{
				{Type: types.AutoModerationActionTimeout, Metadata: &guild.AutoModerationActionMetadata{Duration: 60}},
				{Type: types.AutoModerationActionSendAlertMessage, Metadata: &guild.AutoModerationActionMetadata{
					ChannelID: e.ChannelID,
				}},
			},
		}).Do(ctx)
		if err != nil {
			guild.DeleteAutoModerationRule(*GuildID, rule.ID).Do(ctx)
			panic(err)
		}

		channel.SendMessage(e.ChannelID, "Congratulations! You have just triggered an auto moderation rule.\n"+
			"The current trigger can match anywhere in the word, so even if you write the trigger word as a part of another word, it will still match.\n"+
			"The rule has now been changed, now the trigger matches only in the full words.\n"+
			"Additionally, when you send a message, an alert will be sent to this channel and you will be **timed out** for a minute.\n",
		).Do(ctx)

		var counter int
		var counterMutex sync.Mutex
		dg.EventManager().AddHandler(func(s bot.Session, e *event.AutoModerationActionExecution) {
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
				channel.SendMessage(e.ChannelID, "Nothing has changed, right? "+
					"Well, since separate gateway events are fired per each action (current is "+action+"), "+
					"you'll see a second message about an action pop up soon").Do(ctx)
			} else if counter == 2 {
				counterMutex.Unlock()
				channel.SendMessage(e.ChannelID, "Now the second ("+action+") action got executed.").Do(ctx)
				channel.SendMessage(e.ChannelID, "And... you've made it! That's the end of the example.\n"+
					"For more information about the automod and how to use it, "+
					"you can visit the official Discord docs: https://discord.dev/resources/auto-moderation or ask in our server: https://discord.gg/6dzbuDpSWY",
				).Do(ctx)

				dg.Close(ctx)
				guild.DeleteAutoModerationRule(*GuildID, rule.ID).Do(ctx)
				os.Exit(0)
			}
		})
	})

	err = dg.Open(context.Background())
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer dg.Close(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}

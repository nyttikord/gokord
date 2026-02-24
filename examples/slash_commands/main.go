package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/interaction"

	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/nyttikord/gokord"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *gokord.Session

func init() { flag.Parse() }

var (
	integerOptionMinValue          = 1.0
	defaultMemberPermissions int64 = discord.PermissionManageGuild

	commands = []*interaction.Command{
		{
			Name: "basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
		{
			Name:                     "permission-overview",
			Description:              "Command for demonstration of default command permissions",
			DefaultMemberPermissions: &defaultMemberPermissions,
		},
		{
			Name:        "basic-command-with-files",
			Description: "Basic command with files",
		},
		{
			Name:        "localized-command",
			Description: "Localized command. Description and name may vary depending on the Language setting",
			NameLocalizations: &map[discord.Locale]string{
				discord.LocaleChineseCN: "本地化的命令",
			},
			DescriptionLocalizations: &map[discord.Locale]string{
				discord.LocaleChineseCN: "这是一个本地化的命令",
			},
			Options: []*interaction.CommandOption{
				{
					Name:        "localized-option",
					Description: "Localized option. Description and name may vary depending on the Language setting",
					NameLocalizations: map[discord.Locale]string{
						discord.LocaleChineseCN: "一个本地化的选项",
					},
					DescriptionLocalizations: map[discord.Locale]string{
						discord.LocaleChineseCN: "这是一个本地化的选项",
					},
					Type: types.CommandOptionInteger,
					Choices: []*interaction.CommandOptionChoice{
						{
							Name: "First",
							NameLocalizations: map[discord.Locale]string{
								discord.LocaleChineseCN: "一的",
							},
							Value: 1,
						},
						{
							Name: "Second",
							NameLocalizations: map[discord.Locale]string{
								discord.LocaleChineseCN: "二的",
							},
							Value: 2,
						},
					},
				},
			},
		},
		{
			Name:        "options",
			Description: "Command for demonstrating options",
			Options: []*interaction.CommandOption{

				{
					Type:        types.CommandOptionString,
					Name:        "string-option",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        types.CommandOptionInteger,
					Name:        "integer-option",
					Description: "Integer option",
					MinValue:    &integerOptionMinValue,
					MaxValue:    10,
					Required:    true,
				},
				{
					Type:        types.CommandOptionNumber,
					Name:        "number-option",
					Description: "Float option",
					MaxValue:    10.1,
					Required:    true,
				},
				{
					Type:        types.CommandOptionBoolean,
					Name:        "bool-option",
					Description: "Boolean option",
					Required:    true,
				},

				// Required options must be listed first since optional parameters
				// always come after when they're used.
				// The same concept applies to Discord's Slash-commands API

				{
					Type:        types.CommandOptionChannel,
					Name:        "channel-option",
					Description: "Application option",
					// Application type mask
					ChannelTypes: []types.Channel{
						types.ChannelGuildText,
						types.ChannelGuildVoice,
					},
					Required: false,
				},
				{
					Type:        types.CommandOptionUser,
					Name:        "user-option",
					Description: "Application option",
					Required:    false,
				},
				{
					Type:        types.CommandOptionRole,
					Name:        "role-option",
					Description: "Role option",
					Required:    false,
				},
			},
		},
		{
			Name:        "subcommands",
			Description: "Subcommands and command groups example",
			Options: []*interaction.CommandOption{
				// When a command has subcommands/subcommand groups
				// It must not have top-level options, they aren't accesible in the UI
				// in this case (at least not yet), so if a command has
				// subcommands/subcommand any groups registering top-level options
				// will cause the registration of the command to fail

				{
					Name:        "subcommand-group",
					Description: "Subcommands group",
					Options: []*interaction.CommandOption{
						// Also, subcommand groups aren't capable of
						// containing options, by the name of them, you can see
						// they can only contain subcommands
						{
							Name:        "nested-subcommand",
							Description: "Nested subcommand",
							Type:        types.CommandOptionSubCommand,
						},
					},
					Type: types.CommandOptionSubCommandGroup,
				},
				// Also, you can create both subcommand groups and subcommands
				// in the command at the same time. But, there's some limits to
				// nesting, count of subcommands (top level and nested) and options.
				// Read the intro of slash-commands docs on Discord dev portal
				// to get more information
				{
					Name:        "subcommand",
					Description: "Top-level subcommand",
					Type:        types.CommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "responses",
			Description: "Interaction responses testing initiative",
			Options: []*interaction.CommandOption{
				{
					Name:        "resp-type",
					Description: "Response type",
					Type:        types.CommandOptionInteger,
					Choices: []*interaction.CommandOptionChoice{
						{
							Name:  "Application message with source",
							Value: 4,
						},
						{
							Name:  "Deferred response With Source",
							Value: 5,
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "followups",
			Description: "Followup messages",
		},
	}
)

func main() {
	s := gokord.New("Bot " + *BotToken)

	s.EventManager().AddHandler(func(ctx context.Context, s bot.Session, r *event.Ready) {
		bot.Logger(ctx).Info("Bot connected", "as", s.SessionState().User().Username)
	})

	s.InteractionManager().HandleCommand("basic-command", basicCommand)
	s.InteractionManager().HandleCommand("basic-command-with-files", basicCommandWithFiles)
	s.InteractionManager().HandleCommand("followups", followups)
	s.InteractionManager().HandleCommand("options", options)
	s.InteractionManager().HandleCommand("localized-command", localizedCommand)
	s.InteractionManager().HandleCommand("permission-overview", permissionOverview)
	s.InteractionManager().HandleCommand("responses", responses)
	s.InteractionManager().HandleCommand("subcommands", subcommands)

	err := s.Open(context.Background())
	if err != nil {
		panic(err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*interaction.Command, len(commands))
	for i, v := range commands {
		cmd, err := s.InteractionAPI().CommandCreate(s.SessionState().User().ID, *GuildID, v).Do(context.Background())
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.Application.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.InteractionAPI().CommandDelete(s.SessionState().User().ID, *GuildID, v.ID).Do(context.Background())
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}

func basicCommand(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command",
		},
	}).Do(ctx)
}

func basicCommandWithFiles(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command with a file in the response",
			Files: []*request.File{
				{
					ContentType: "text/plain",
					Name:        "test.txt",
					Reader:      strings.NewReader("Hello Discord!!"),
				},
			},
		},
	}).Do(ctx)
}

func localizedCommand(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	responses := map[discord.Locale]string{
		discord.LocaleChineseCN: "你好！ 这是一个本地化的命令",
	}
	response := "Hi! This is a localized message"
	if r, ok := responses[i.Locale]; ok {
		response = r
	}
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: response,
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}

func options(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	// Access options in the order provided by the user.
	options := i.Data.Options

	// Or convert the slice into a map
	optionMap := make(map[string]*interaction.CommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// This example stores the provided arguments in an []interface{}
	// which will be used to format the bot's response
	margs := make([]any, 0, len(options))
	msgformat := "You learned how to use command options! " +
		"Take a look at the value(s) you entered:\n"

	// Application the value from the option map.
	// When the option exists, ok = true
	if option, ok := optionMap["string-option"]; ok {
		// Option values must be type asserted from interface{}.
		// Discordgo provides utility functions to make this simple.
		margs = append(margs, option.StringValue())
		msgformat += "> string-option: %s\n"
	}

	if opt, ok := optionMap["integer-option"]; ok {
		margs = append(margs, opt.IntValue())
		msgformat += "> integer-option: %d\n"
	}

	if opt, ok := optionMap["number-option"]; ok {
		margs = append(margs, opt.FloatValue())
		msgformat += "> number-option: %f\n"
	}

	if opt, ok := optionMap["bool-option"]; ok {
		margs = append(margs, opt.BoolValue())
		msgformat += "> bool-option: %v\n"
	}

	if opt, ok := optionMap["channel-option"]; ok {
		margs = append(margs, opt.ChannelValue(ctx, s.ChannelAPI().State).ID)
		msgformat += "> channel-option: <#%s>\n"
	}

	if opt, ok := optionMap["user-option"]; ok {
		margs = append(margs, opt.UserValue(ctx).ID)
		msgformat += "> user-option: <@%s>\n"
	}

	if opt, ok := optionMap["role-option"]; ok {
		margs = append(margs, opt.RoleValue(ctx, "", s.GuildAPI().State).ID)
		msgformat += "> role-option: <@&%s>\n"
	}

	s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		// Ignore type for now, they will be discussed in "responses"
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	}).Do(ctx)
}

func permissionOverview(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	perms, err := s.InteractionAPI().CommandPermissions(s.SessionState().User().ID, i.GuildID, i.Data.ID).Do(ctx)

	var restError *gokord.RESTError
	if errors.As(err, &restError) && restError.Message != nil && restError.Message.Code == discord.ErrCodeUnknownApplicationCommandPermissions {
		s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
			Type: types.InteractionResponseChannelMessageWithSource,
			Data: &interaction.ResponseData{
				Content: ":x: No permission overwrites",
			},
		}).Do(ctx)
		return
	} else if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	format := "- %s %s\n"

	channels := ""
	users := ""
	roles := ""

	for _, o := range perms.Permissions {
		emoji := "❌"
		if o.Permission {
			emoji = "☑"
		}

		switch o.Type {
		case types.CommandPermissionUser:
			users += fmt.Sprintf(format, emoji, "<@!"+o.ID+">")
		case types.CommandPermissionChannel:
			allChannels, _ := interaction.GuildAllChannelsID(i.GuildID)

			if o.ID == allChannels {
				channels += fmt.Sprintf(format, emoji, "All channels")
			} else {
				channels += fmt.Sprintf(format, emoji, "<#"+o.ID+">")
			}
		case types.CommandPermissionRole:
			if o.ID == i.GuildID {
				roles += fmt.Sprintf(format, emoji, "@everyone")
			} else {
				roles += fmt.Sprintf(format, emoji, "<@&"+o.ID+">")
			}
		}
	}

	s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Embeds: []*channel.MessageEmbed{
				{
					Title:       "Permissions overview",
					Description: "Overview of permissions for this command",
					Fields: []*channel.MessageEmbedField{
						{
							Name:  "Users",
							Value: users,
						},
						{
							Name:  "Channels",
							Value: channels,
						},
						{
							Name:  "Roles",
							Value: roles,
						},
					},
				},
			},
			AllowedMentions: &channel.MessageAllowedMentions{},
		},
	}).Do(ctx)
}

func subcommands(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	options := i.Data.Options
	content := ""

	// As you can see, names of subcommands (nested, top-level)
	// and subcommand groups are provided through the arguments.
	switch options[0].Name {
	case "subcommand":
		content = "The top-level subcommand is executed. Now try to execute the nested one."
	case "subcommand-group":
		options = options[0].Options
		switch options[0].Name {
		case "nested-subcommand":
			content = "Nice, now you know how to execute nested commands too"
		default:
			content = "Oops, something went wrong.\n" +
				"Hol' up, you aren't supposed to see this message."
		}
	}

	s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: content,
		},
	}).Do(ctx)
}

func responses(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	// Responses to a command are very important.
	// First of all, because you need to react to the interaction
	// by sending the response in 3 seconds after receiving, otherwise
	// interaction will be considered invalid and you can no longer
	// use the interaction token and ID for responding to the user's request

	content := ""
	// As you can see, the response type names used here are pretty self-explanatory,
	// but for those who want more information see the official documentation
	switch i.Data.Options[0].IntValue() {
	case int64(types.InteractionResponseChannelMessageWithSource):
		content =
			"You just responded to an interaction, sent a message and showed the original one. " +
				"Congratulations!"
		content +=
			"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
	default:
		err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
			Type: types.InteractionResponse(i.Data.Options[0].IntValue()),
		}).Do(ctx)
		if err != nil {
			s.InteractionAPI().FollowupMessageCreate(i.Interaction, true, &channel.WebhookParams{
				Content: "Something went wrong",
			}).Do(ctx)
		}
		return
	}

	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponse(i.Data.Options[0].IntValue()),
		Data: &interaction.ResponseData{
			Content: content,
		},
	}).Do(ctx)
	if err != nil {
		s.InteractionAPI().FollowupMessageCreate(i.Interaction, true, &channel.WebhookParams{
			Content: "Something went wrong",
		}).Do(ctx)
		return
	}
	time.AfterFunc(time.Second*5, func() {
		content := content + "\n\nWell, now you know how to create and edit responses. " +
			"But you still don't know how to delete them... so... wait 10 seconds and this " +
			"message will be deleted."
		_, err = s.InteractionAPI().ResponseEdit(i.Interaction, &channel.WebhookEdit{
			Content: &content,
		}).Do(ctx)
		if err != nil {
			s.InteractionAPI().FollowupMessageCreate(i.Interaction, true, &channel.WebhookParams{
				Content: "Something went wrong",
			}).Do(ctx)
			return
		}
		time.Sleep(time.Second * 10)
		s.InteractionAPI().ResponseDelete(i.Interaction).Do(ctx)
	})
}

func followups(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	// Followup messages are basically regular messages (you can create as many of them as you wish)
	// but work as they are created by webhooks and their functionality
	// is for handling additional messages after sending a response.

	s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:   channel.MessageFlagsEphemeral,
			Content: "Surprise!",
		},
	}).Do(ctx)
	msg, err := s.InteractionAPI().FollowupMessageCreate(i.Interaction, true, &channel.WebhookParams{
		Content: "Followup message has been created, after 5 seconds it will be edited",
	}).Do(ctx)
	if err != nil {
		s.InteractionAPI().FollowupMessageCreate(i.Interaction, true, &channel.WebhookParams{
			Content: "Something went wrong",
		}).Do(ctx)
		return
	}
	time.Sleep(time.Second * 5)

	content := "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted."
	s.InteractionAPI().FollowupMessageEdit(i.Interaction, msg.ID, &channel.WebhookEdit{
		Content: &content,
	}).Do(ctx)

	time.Sleep(time.Second * 10)

	s.InteractionAPI().FollowupMessageDelete(i.Interaction, msg.ID).Do(ctx)

	s.InteractionAPI().FollowupMessageCreate(i.Interaction, true, &channel.WebhookParams{
		Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
			"take a unicorn :unicorn: as reward!\n" +
			"Also, as bonus... look at the original interaction response :D",
	}).Do(ctx)
}

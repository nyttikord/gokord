package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/component"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/interaction"
)

// Bot parameters
var (
	GuildID  = flag.String("guild", "", "Test guild ID")
	BotToken = flag.String("token", "", "Bot access token")
	AppID    = flag.String("app", "", "Application ID")
)

func init() { flag.Parse() }

func main() {
	s := gokord.New("Bot " + *BotToken)

	s.InteractionManager().HandleMessageComponent("fd_no", fdNo)
	s.InteractionManager().HandleMessageComponent("fd_yes", fdYes)
	s.InteractionManager().HandleMessageComponent("select", selectFn)
	s.InteractionManager().HandleMessageComponent("stackoverflow_tags", stackoverflowTags)
	s.InteractionManager().HandleMessageComponent("channel_select", channelSelect)

	s.InteractionManager().HandleCommand("selects", selects)
	s.InteractionManager().HandleCommand("buttons", buttons)

	s.EventManager().AddHandler(func(ctx context.Context, s bot.Session, r *event.Ready) {
		bot.Logger(ctx).Info("Bot is up!")
	})

	_, err := s.InteractionAPI().CommandCreate(*AppID, *GuildID, &interaction.Command{
		Name:        "buttons",
		Description: "Test the buttons if you got courage",
	}).Do(context.Background())

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
	_, err = s.InteractionAPI().CommandCreate(*AppID, *GuildID, &interaction.Command{
		Name: "selects",
		Options: []*interaction.CommandOption{
			{
				Type:        types.CommandOptionSubCommand,
				Name:        "multi",
				Description: "Multi-item select menu",
			},
			{
				Type:        types.CommandOptionSubCommand,
				Name:        "single",
				Description: "Single-item select menu",
			},
			{
				Type:        types.CommandOptionSubCommand,
				Name:        "auto-populated",
				Description: "Automatically populated select menu, which lets you pick a member, channel or role",
			},
		},
		Description: "Lo and behold: dropdowns are coming",
	}).Do(context.Background())

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}

	err = s.OpenAndBlock(context.Background())
	if err != nil {
		panic(err)
	}
}

func fdNo(ctx context.Context, s bot.Session, i *interaction.MessageComponent) {
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: "Huh. I see, maybe some of these resources might help you?",
			Flags:   channel.MessageFlagsEphemeral,
			Components: []component.Component{
				&component.ActionsRow{
					Components: []component.Message{
						&component.Button{
							Emoji: &emoji.Component{
								Name: "📜",
							},
							Label: "Documentation",
							Style: component.ButtonStyleLink,
							URL:   "https://discord.com/developers/docs/interactions/message-components#buttons",
						},
						&component.Button{
							Emoji: &emoji.Component{
								Name: "🔧",
							},
							Label: "Discord developers",
							Style: component.ButtonStyleLink,
							URL:   "https://discord.gg/discord-developers",
						},
						&component.Button{
							Emoji: &emoji.Component{
								Name: "🦫",
							},
							Label: "Discord Gophers",
							Style: component.ButtonStyleLink,
							URL:   "https://discord.gg/7RuRrVHyXF",
						},
					},
				},
			},
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}

func fdYes(ctx context.Context, s bot.Session, i *interaction.MessageComponent) {
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: "Great! If you wanna know more or just have questions, feel free to visit Discord Devs and Discord Gophers server. " +
				"But now, when you know how buttons work, let's move onto select menus (execute `/selects single`)",
			Flags: channel.MessageFlagsEphemeral,
			Components: []component.Component{
				&component.ActionsRow{
					Components: []component.Message{
						&component.Button{
							Emoji: &emoji.Component{
								Name: "🔧",
							},
							Label: "Discord developers",
							Style: component.ButtonStyleLink,
							URL:   "https://discord.gg/discord-developers",
						},
						&component.Button{
							Emoji: &emoji.Component{
								Name: "🦫",
							},
							Label: "Discord Gophers",
							Style: component.ButtonStyleLink,
							URL:   "https://discord.gg/7RuRrVHyXF",
						},
					},
				},
			},
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}

func selectFn(ctx context.Context, s bot.Session, i *interaction.MessageComponent) {
	var response *interaction.Response

	data := i.Data

	switch data.Values[0] {
	case "go":
		response = &interaction.Response{
			Type: types.InteractionResponseChannelMessageWithSource,
			Data: &interaction.ResponseData{
				Content: "This is the way.",
				Flags:   channel.MessageFlagsEphemeral,
			},
		}
	default:
		response = &interaction.Response{
			Type: types.InteractionResponseChannelMessageWithSource,
			Data: &interaction.ResponseData{
				Content: "It is not the way to go.",
				Flags:   channel.MessageFlagsEphemeral,
			},
		}
	}
	err := s.InteractionAPI().Respond(i.Interaction, response).Do(ctx)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second) // Doing that so user won't see instant response.
	_, err = s.InteractionAPI().FollowupMessageCreate(i.Interaction, true, &channel.WebhookParams{
		Content: "Anyways, now when you know how to use single select menus, let's see how multi select menus work. " +
			"Try calling `/selects multi` command.",
		Flags: channel.MessageFlagsEphemeral,
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}

func stackoverflowTags(ctx context.Context, s bot.Session, i *interaction.MessageComponent) {
	data := i.Data

	const stackoverflowFormat = `https://stackoverflow.com/questions/tagged/%s`

	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: "Here is your stackoverflow URL: " + fmt.Sprintf(stackoverflowFormat, strings.Join(data.Values, "+")),
			Flags:   channel.MessageFlagsEphemeral,
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second) // Doing that so user won't see instant response.
	_, err = s.InteractionAPI().FollowupMessageCreate(i.Interaction, true, &channel.WebhookParams{
		Content: "But wait, there is more! You can also auto populate the select menu. Try executing `/selects auto-populated`.",
		Flags:   channel.MessageFlagsEphemeral,
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}

func channelSelect(ctx context.Context, s bot.Session, i *interaction.MessageComponent) {
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: "This is it. You've reached your destination. Your choice was <#" + i.Data.Values[0] + ">\n" +
				"If you want to know more, check out the links below",
			Components: []component.Component{
				&component.ActionsRow{
					Components: []component.Message{
						&component.Button{
							Emoji: &emoji.Component{
								Name: "📜",
							},
							Label: "Documentation",
							Style: component.ButtonStyleLink,
							URL:   "https://discord.com/developers/docs/interactions/message-components#select-menus",
						},
						&component.Button{
							Emoji: &emoji.Component{
								Name: "🔧",
							},
							Label: "Discord developers",
							Style: component.ButtonStyleLink,
							URL:   "https://discord.gg/discord-developers",
						},
						&component.Button{
							Emoji: &emoji.Component{
								Name: "🦫",
							},
							Label: "Discord Gophers",
							Style: component.ButtonStyleLink,
							URL:   "https://discord.gg/7RuRrVHyXF",
						},
					},
				},
			},

			Flags: channel.MessageFlagsEphemeral,
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}

func buttons(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	err := s.InteractionAPI().Respond(i.Interaction, &interaction.Response{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.ResponseData{
			Content: "Are you comfortable with buttons and other message components?",
			Flags:   channel.MessageFlagsEphemeral,
			// Buttons and other components are specified in Components field.
			Components: []component.Component{
				// ActionRow is a container of all buttons within the same row.
				&component.ActionsRow{
					Components: []component.Message{
						&component.Button{
							// Label is what the user will see on the button.
							Label: "Yes",
							// Style provides coloring of the button. There are not so many styles tho.
							Style: component.ButtonStyleSuccess,
							// Disabled allows bot to disable some buttons for users.
							Disabled: false,
							// CustomID is a thing telling Discord which data to send when this button will be pressed.
							CustomID: "fd_yes",
						},
						&component.Button{
							Label:    "No",
							Style:    component.ButtonStyleDanger,
							Disabled: false,
							CustomID: "fd_no",
						},
						&component.Button{
							Label:    "I don't know",
							Style:    component.ButtonStyleLink,
							Disabled: false,
							// Link buttons don't require CustomID and do not trigger the gateway/HTTP event
							URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
							Emoji: &emoji.Component{
								Name: "🤷",
							},
						},
					},
				},
				// The message may have multiple actions rows.
				&component.ActionsRow{
					Components: []component.Message{
						&component.Button{
							Label:    "Discord Developers server",
							Style:    component.ButtonStyleLink,
							Disabled: false,
							URL:      "https://discord.gg/discord-developers",
						},
					},
				},
			},
		},
	}).Do(ctx)
	if err != nil {
		panic(err)
	}
}

func selects(ctx context.Context, s bot.Session, i *interaction.ApplicationCommand) {
	var response *interaction.Response
	switch i.Data.Options[0].Name {
	case "single":
		response = &interaction.Response{
			Type: types.InteractionResponseChannelMessageWithSource,
			Data: &interaction.ResponseData{
				Content: "Now let's take a look on selects. This is single item select menu.",
				Flags:   channel.MessageFlagsEphemeral,
				Components: []component.Component{
					&component.ActionsRow{
						Components: []component.Message{
							&component.SelectMenu{
								// Select menu, as other components, must have a customID, so we set it to this value.
								CustomID:    "select",
								Placeholder: "Choose your favorite programming language 👇",
								Options: []component.SelectMenuOption{
									{
										Label: "Go",
										// As with components, this things must have their own unique "id" to identify which is which.
										// In this case such id is Value field.
										Value: "go",
										Emoji: &emoji.Component{
											Name: "🦦",
										},
										// You can also make it a default option, but in this case we won't.
										Default:     false,
										Description: "Go programming language",
									},
									{
										Label: "JS",
										Value: "js",
										Emoji: &emoji.Component{
											Name: "🟨",
										},
										Description: "JavaScript programming language",
									},
									{
										Label: "Python",
										Value: "py",
										Emoji: &emoji.Component{
											Name: "🐍",
										},
										Description: "Python programming language",
									},
								},
							},
						},
					},
				},
			},
		}
	case "multi":
		minValues := 1
		response = &interaction.Response{
			Type: types.InteractionResponseChannelMessageWithSource,
			Data: &interaction.ResponseData{
				Content: "Now let's see how the multi-item select menu works: " +
					"try generating your own stackoverflow search link",
				Flags: channel.MessageFlagsEphemeral,
				Components: []component.Component{
					&component.ActionsRow{
						Components: []component.Message{
							&component.SelectMenu{
								CustomID:    "stackoverflow_tags",
								Placeholder: "Select tags to search on StackOverflow",
								// This is where confusion comes from. If you don't specify these things you will get single item select.
								// These fields control the minimum and maximum amount of selected items.
								MinValues: &minValues,
								MaxValues: 3,
								Options: []component.SelectMenuOption{
									{
										Label:       "Go",
										Description: "Simple yet powerful programming language",
										Value:       "go",
										// Default works the same for multi-select menus.
										Default: false,
										Emoji: &emoji.Component{
											Name: "🦦",
										},
									},
									{
										Label:       "JS",
										Description: "Multiparadigm OOP language",
										Value:       "javascript",
										Emoji: &emoji.Component{
											Name: "🟨",
										},
									},
									{
										Label:       "Python",
										Description: "OOP prototyping programming language",
										Value:       "python",
										Emoji: &emoji.Component{
											Name: "🐍",
										},
									},
									{
										Label:       "Web",
										Description: "Web related technologies",
										Value:       "web",
										Emoji: &emoji.Component{
											Name: "🌐",
										},
									},
									{
										Label:       "Desktop",
										Description: "Desktop applications",
										Value:       "desktop",
										Emoji: &emoji.Component{
											Name: "💻",
										},
									},
								},
							},
						},
					},
				},
			},
		}
	case "auto-populated":
		response = &interaction.Response{
			Type: types.InteractionResponseChannelMessageWithSource,
			Data: &interaction.ResponseData{
				Content: "The tastiest things are left for the end. Meet auto populated select menus.\n" +
					"By setting `MenuType` on the select menu you can tell Discord to automatically populate the menu with entities of your choice: roles, members, channels. Try one below.",
				Flags: channel.MessageFlagsEphemeral,
				Components: []component.Component{
					&component.ActionsRow{
						Components: []component.Message{
							&component.SelectMenu{
								MenuType:     types.SelectMenuChannel,
								CustomID:     "channel_select",
								Placeholder:  "Pick your favorite channel!",
								ChannelTypes: []types.Channel{types.ChannelGuildText},
							},
						},
					},
				},
			},
		}
	}
	err := s.InteractionAPI().Respond(i.Interaction, response).Do(ctx)
	if err != nil {
		panic(err)
	}
}

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
)

// Bot parameters
var (
	GuildID  = flag.String("guild", "", "Test guild ID")
	BotToken = flag.String("token", "", "Bot access token")
	AppID    = flag.String("app", "", "Application ID")
)

var s *gokord.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = gokord.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

// Important note: call every command in order it's placed in the example.

var (
	componentsHandlers = map[string]func(s *gokord.Session, i *gokord.InteractionCreate){
		"fd_no": func(s *gokord.Session, i *gokord.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &gokord.InteractionResponse{
				Type: gokord.InteractionResponseChannelMessageWithSource,
				Data: &gokord.InteractionResponseData{
					Content: "Huh. I see, maybe some of these resources might help you?",
					Flags:   gokord.MessageFlagsEphemeral,
					Components: []gokord.MessageComponent{
						gokord.ActionsRow{
							Components: []gokord.MessageComponent{
								gokord.Button{
									Emoji: &gokord.ComponentEmoji{
										Name: "📜",
									},
									Label: "Documentation",
									Style: gokord.LinkButton,
									URL:   "https://discord.com/developers/docs/interactions/message-components#buttons",
								},
								gokord.Button{
									Emoji: &gokord.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: gokord.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								gokord.Button{
									Emoji: &gokord.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: gokord.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"fd_yes": func(s *gokord.Session, i *gokord.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &gokord.InteractionResponse{
				Type: gokord.InteractionResponseChannelMessageWithSource,
				Data: &gokord.InteractionResponseData{
					Content: "Great! If you wanna know more or just have questions, feel free to visit Discord Devs and Discord Gophers server. " +
						"But now, when you know how buttons work, let's move onto select menus (execute `/selects single`)",
					Flags: gokord.MessageFlagsEphemeral,
					Components: []gokord.MessageComponent{
						gokord.ActionsRow{
							Components: []gokord.MessageComponent{
								gokord.Button{
									Emoji: &gokord.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: gokord.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								gokord.Button{
									Emoji: &gokord.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: gokord.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"select": func(s *gokord.Session, i *gokord.InteractionCreate) {
			var response *gokord.InteractionResponse

			data := i.MessageComponentData()
			switch data.Values[0] {
			case "go":
				response = &gokord.InteractionResponse{
					Type: gokord.InteractionResponseChannelMessageWithSource,
					Data: &gokord.InteractionResponseData{
						Content: "This is the way.",
						Flags:   gokord.MessageFlagsEphemeral,
					},
				}
			default:
				response = &gokord.InteractionResponse{
					Type: gokord.InteractionResponseChannelMessageWithSource,
					Data: &gokord.InteractionResponseData{
						Content: "It is not the way to go.",
						Flags:   gokord.MessageFlagsEphemeral,
					},
				}
			}
			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second) // Doing that so user won't see instant response.
			_, err = s.FollowupMessageCreate(i.Interaction, true, &gokord.WebhookParams{
				Content: "Anyways, now when you know how to use single select menus, let's see how multi select menus work. " +
					"Try calling `/selects multi` command.",
				Flags: gokord.MessageFlagsEphemeral,
			})
			if err != nil {
				panic(err)
			}
		},
		"stackoverflow_tags": func(s *gokord.Session, i *gokord.InteractionCreate) {
			data := i.MessageComponentData()

			const stackoverflowFormat = `https://stackoverflow.com/questions/tagged/%s`

			err := s.InteractionRespond(i.Interaction, &gokord.InteractionResponse{
				Type: gokord.InteractionResponseChannelMessageWithSource,
				Data: &gokord.InteractionResponseData{
					Content: "Here is your stackoverflow URL: " + fmt.Sprintf(stackoverflowFormat, strings.Join(data.Values, "+")),
					Flags:   gokord.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second) // Doing that so user won't see instant response.
			_, err = s.FollowupMessageCreate(i.Interaction, true, &gokord.WebhookParams{
				Content: "But wait, there is more! You can also auto populate the select menu. Try executing `/selects auto-populated`.",
				Flags:   gokord.MessageFlagsEphemeral,
			})
			if err != nil {
				panic(err)
			}
		},
		"channel_select": func(s *gokord.Session, i *gokord.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &gokord.InteractionResponse{
				Type: gokord.InteractionResponseChannelMessageWithSource,
				Data: &gokord.InteractionResponseData{
					Content: "This is it. You've reached your destination. Your choice was <#" + i.MessageComponentData().Values[0] + ">\n" +
						"If you want to know more, check out the links below",
					Components: []gokord.MessageComponent{
						gokord.ActionsRow{
							Components: []gokord.MessageComponent{
								gokord.Button{
									Emoji: &gokord.ComponentEmoji{
										Name: "📜",
									},
									Label: "Documentation",
									Style: gokord.LinkButton,
									URL:   "https://discord.com/developers/docs/interactions/message-components#select-menus",
								},
								gokord.Button{
									Emoji: &gokord.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: gokord.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								gokord.Button{
									Emoji: &gokord.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: gokord.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},

					Flags: gokord.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
	commandsHandlers = map[string]func(s *gokord.Session, i *gokord.InteractionCreate){
		"buttons": func(s *gokord.Session, i *gokord.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &gokord.InteractionResponse{
				Type: gokord.InteractionResponseChannelMessageWithSource,
				Data: &gokord.InteractionResponseData{
					Content: "Are you comfortable with buttons and other message components?",
					Flags:   gokord.MessageFlagsEphemeral,
					// Buttons and other components are specified in Components field.
					Components: []gokord.MessageComponent{
						// ActionRow is a container of all buttons within the same row.
						gokord.ActionsRow{
							Components: []gokord.MessageComponent{
								gokord.Button{
									// Label is what the user will see on the button.
									Label: "Yes",
									// Style provides coloring of the button. There are not so many styles tho.
									Style: gokord.SuccessButton,
									// Disabled allows bot to disable some buttons for users.
									Disabled: false,
									// CustomID is a thing telling Discord which data to send when this button will be pressed.
									CustomID: "fd_yes",
								},
								gokord.Button{
									Label:    "No",
									Style:    gokord.DangerButton,
									Disabled: false,
									CustomID: "fd_no",
								},
								gokord.Button{
									Label:    "I don't know",
									Style:    gokord.LinkButton,
									Disabled: false,
									// Link buttons don't require CustomID and do not trigger the gateway/HTTP event
									URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
									Emoji: &gokord.ComponentEmoji{
										Name: "🤷",
									},
								},
							},
						},
						// The message may have multiple actions rows.
						gokord.ActionsRow{
							Components: []gokord.MessageComponent{
								gokord.Button{
									Label:    "Discord Developers server",
									Style:    gokord.LinkButton,
									Disabled: false,
									URL:      "https://discord.gg/discord-developers",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"selects": func(s *gokord.Session, i *gokord.InteractionCreate) {
			var response *gokord.InteractionResponse
			switch i.ApplicationCommandData().Options[0].Name {
			case "single":
				response = &gokord.InteractionResponse{
					Type: gokord.InteractionResponseChannelMessageWithSource,
					Data: &gokord.InteractionResponseData{
						Content: "Now let's take a look on selects. This is single item select menu.",
						Flags:   gokord.MessageFlagsEphemeral,
						Components: []gokord.MessageComponent{
							gokord.ActionsRow{
								Components: []gokord.MessageComponent{
									gokord.SelectMenu{
										// Select menu, as other components, must have a customID, so we set it to this value.
										CustomID:    "select",
										Placeholder: "Choose your favorite programming language 👇",
										Options: []gokord.SelectMenuOption{
											{
												Label: "Go",
												// As with components, this things must have their own unique "id" to identify which is which.
												// In this case such id is Value field.
												Value: "go",
												Emoji: &gokord.ComponentEmoji{
													Name: "🦦",
												},
												// You can also make it a default option, but in this case we won't.
												Default:     false,
												Description: "Go programming language",
											},
											{
												Label: "JS",
												Value: "js",
												Emoji: &gokord.ComponentEmoji{
													Name: "🟨",
												},
												Description: "JavaScript programming language",
											},
											{
												Label: "Python",
												Value: "py",
												Emoji: &gokord.ComponentEmoji{
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
				response = &gokord.InteractionResponse{
					Type: gokord.InteractionResponseChannelMessageWithSource,
					Data: &gokord.InteractionResponseData{
						Content: "Now let's see how the multi-item select menu works: " +
							"try generating your own stackoverflow search link",
						Flags: gokord.MessageFlagsEphemeral,
						Components: []gokord.MessageComponent{
							gokord.ActionsRow{
								Components: []gokord.MessageComponent{
									gokord.SelectMenu{
										CustomID:    "stackoverflow_tags",
										Placeholder: "Select tags to search on StackOverflow",
										// This is where confusion comes from. If you don't specify these things you will get single item select.
										// These fields control the minimum and maximum amount of selected items.
										MinValues: &minValues,
										MaxValues: 3,
										Options: []gokord.SelectMenuOption{
											{
												Label:       "Go",
												Description: "Simple yet powerful programming language",
												Value:       "go",
												// Default works the same for multi-select menus.
												Default: false,
												Emoji: &gokord.ComponentEmoji{
													Name: "🦦",
												},
											},
											{
												Label:       "JS",
												Description: "Multiparadigm OOP language",
												Value:       "javascript",
												Emoji: &gokord.ComponentEmoji{
													Name: "🟨",
												},
											},
											{
												Label:       "Python",
												Description: "OOP prototyping programming language",
												Value:       "python",
												Emoji: &gokord.ComponentEmoji{
													Name: "🐍",
												},
											},
											{
												Label:       "Web",
												Description: "Web related technologies",
												Value:       "web",
												Emoji: &gokord.ComponentEmoji{
													Name: "🌐",
												},
											},
											{
												Label:       "Desktop",
												Description: "Desktop applications",
												Value:       "desktop",
												Emoji: &gokord.ComponentEmoji{
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
				response = &gokord.InteractionResponse{
					Type: gokord.InteractionResponseChannelMessageWithSource,
					Data: &gokord.InteractionResponseData{
						Content: "The tastiest things are left for the end. Meet auto populated select menus.\n" +
							"By setting `MenuType` on the select menu you can tell Discord to automatically populate the menu with entities of your choice: roles, members, channels. Try one below.",
						Flags: gokord.MessageFlagsEphemeral,
						Components: []gokord.MessageComponent{
							gokord.ActionsRow{
								Components: []gokord.MessageComponent{
									gokord.SelectMenu{
										MenuType:     gokord.ChannelSelectMenu,
										CustomID:     "channel_select",
										Placeholder:  "Pick your favorite channel!",
										ChannelTypes: []gokord.ChannelType{gokord.ChannelTypeGuildText},
									},
								},
							},
						},
					},
				}
			}
			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				panic(err)
			}
		},
	}
)

func main() {
	s.AddHandler(func(s *gokord.Session, r *gokord.Ready) {
		log.Println("Bot is up!")
	})
	// Components are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *gokord.Session, i *gokord.InteractionCreate) {
		switch i.Type {
		case gokord.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case gokord.InteractionMessageComponent:

			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
	_, err := s.ApplicationCommandCreate(*AppID, *GuildID, &gokord.ApplicationCommand{
		Name:        "buttons",
		Description: "Test the buttons if you got courage",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
	_, err = s.ApplicationCommandCreate(*AppID, *GuildID, &gokord.ApplicationCommand{
		Name: "selects",
		Options: []*gokord.ApplicationCommandOption{
			{
				Type:        gokord.ApplicationCommandOptionSubCommand,
				Name:        "multi",
				Description: "Multi-item select menu",
			},
			{
				Type:        gokord.ApplicationCommandOptionSubCommand,
				Name:        "single",
				Description: "Single-item select menu",
			},
			{
				Type:        gokord.ApplicationCommandOptionSubCommand,
				Name:        "auto-populated",
				Description: "Automatically populated select menu, which lets you pick a member, channel or role",
			},
		},
		Description: "Lo and behold: dropdowns are coming",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}

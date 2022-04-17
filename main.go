package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken = flag.String("token", "", "Bot access token")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	integerOptionMinValue = 1.0

	objects = []string{
		"mimi", "momo", "mimo", "tree", "birb", "shiro",
	}

	commands = []*discordgo.ApplicationCommand{}
	// {
	// 	Name: "basic-command",
	// 	// All commands and options must have a description
	// 	// Commands/options without description will fail the registration
	// 	// of the command.
	// 	Description: "Basic command",
	// },
	// {
	// 	Name:        "basic-command-with-files",
	// 	Description: "Basic command with files",
	// },
	// {
	// 	Name:        "subcommands",
	// 	Description: "Subcommands and command groups example",
	// 	Options: []*discordgo.ApplicationCommandOption{
	// 		// When a command has subcommands/subcommand groups
	// 		// It must not have top-level options, they aren't accesible in the UI
	// 		// in this case (at least not yet), so if a command has
	// 		// subcommands/subcommand any groups registering top-level options
	// 		// will cause the registration of the command to fail
	// 		{
	// 			Name:        "subcommand-group",
	// 			Description: "Subcommands group",
	// 			Options: []*discordgo.ApplicationCommandOption{
	// 				// Also, subcommand groups aren't capable of
	// 				// containing options, by the name of them, you can see
	// 				// they can only contain subcommands
	// 				{
	// 					Name:        "nested-subcommand",
	// 					Description: "Nested subcommand",
	// 					Type:        discordgo.ApplicationCommandOptionSubCommand,
	// 				},
	// 			},
	// 			Type: discordgo.ApplicationCommandOptionSubCommandGroup,
	// 		},
	// 		// Also, you can create both subcommand groups and subcommands
	// 		// in the command at the same time. But, there's some limits to
	// 		// nesting, count of subcommands (top level and nested) and options.
	// 		// Read the intro of slash-commands docs on Discord dev portal
	// 		// to get more information
	// 		{
	// 			Name:        "subcommand",
	// 			Description: "Top-level subcommand",
	// 			Type:        discordgo.ApplicationCommandOptionSubCommand,
	// 		},
	// 	},
	// },
	// {
	// 	Name:        "responses",
	// 	Description: "Interaction responses testing initiative",
	// 	Options: []*discordgo.ApplicationCommandOption{
	// 		{
	// 			Name:        "resp-type",
	// 			Description: "Response type",
	// 			Type:        discordgo.ApplicationCommandOptionInteger,
	// 			Choices: []*discordgo.ApplicationCommandOptionChoice{
	// 				{
	// 					Name:  "Channel message with source",
	// 					Value: 4,
	// 				},
	// 				{
	// 					Name:  "Deferred response With Source",
	// 					Value: 5,
	// 				},
	// 			},
	// 			Required: true,
	// 		},
	// 	},
	// },
	// {
	// 	Name:        "followups",
	// 	Description: "Followup messages",
	// },
	// }

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
	// "basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: &discordgo.InteractionResponseData{
	// 			Content: "Hey there! Congratulations, you just executed your first slash command",
	// 		},
	// 	})
	// },
	// "basic-command-with-files": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: &discordgo.InteractionResponseData{
	// 			Content: "Hey there! Congratulations, you just executed your first slash command with a file in the response",
	// 			Files: []*discordgo.File{
	// 				{
	// 					ContentType: "text/plain",
	// 					Name:        "test.txt",
	// 					Reader:      strings.NewReader("Hello Discord!!"),
	// 				},
	// 			},
	// 		},
	// 	})
	// },
	// "options": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	// Access options in the order provided by the user.
	// 	options := i.ApplicationCommandData().Options

	// 	// Or convert the slice into a map
	// 	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	// 	for _, opt := range options {
	// 		optionMap[opt.Name] = opt
	// 	}

	// 	// This example stores the provided arguments in an []interface{}
	// 	// which will be used to format the bot's response
	// 	margs := make([]interface{}, 0, len(options))
	// 	msgformat := "You learned how to use command options! " +
	// 		"Take a look at the value(s) you entered:\n"

	// 	// Get the value from the option map.
	// 	// When the option exists, ok = true
	// 	if option, ok := optionMap["attachment-option"]; ok {
	// 		key := option.Value.(string)

	// 		attachment_url := i.Data.(discordgo.ApplicationCommandInteractionData).Resolved.Attachments[key].URL
	// 		fmt.Println(attachment_url)
	// 		msgformat += "> attachment-option: %s\n"
	// 	}

	// 	if opt, ok := optionMap["integer-option"]; ok {
	// 		margs = append(margs, opt.IntValue())
	// 		msgformat += "> integer-option: %d\n"
	// 	}

	// 	if opt, ok := optionMap["number-option"]; ok {
	// 		margs = append(margs, opt.FloatValue())
	// 		msgformat += "> number-option: %f\n"
	// 	}

	// 	if opt, ok := optionMap["bool-option"]; ok {
	// 		margs = append(margs, opt.BoolValue())
	// 		msgformat += "> bool-option: %v\n"
	// 	}

	// 	if opt, ok := optionMap["channel-option"]; ok {
	// 		margs = append(margs, opt.ChannelValue(nil).ID)
	// 		msgformat += "> channel-option: <#%s>\n"
	// 	}

	// 	if opt, ok := optionMap["user-option"]; ok {
	// 		margs = append(margs, opt.UserValue(nil).ID)
	// 		msgformat += "> user-option: <@%s>\n"
	// 	}

	// 	if opt, ok := optionMap["role-option"]; ok {
	// 		margs = append(margs, opt.RoleValue(nil, "").ID)
	// 		msgformat += "> role-option: <@&%s>\n"
	// 	}

	// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		// Ignore type for now, they will be discussed in "responses"
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: &discordgo.InteractionResponseData{
	// 			Content: fmt.Sprintf(
	// 				msgformat,
	// 				margs...,
	// 			),
	// 		},
	// 	})
	// },
	// "subcommands": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	options := i.ApplicationCommandData().Options
	// 	content := ""

	// 	// As you can see, names of subcommands (nested, top-level)
	// 	// and subcommand groups are provided through the arguments.
	// 	switch options[0].Name {
	// 	case "subcommand":
	// 		content = "The top-level subcommand is executed. Now try to execute the nested one."
	// 	case "subcommand-group":
	// 		options = options[0].Options
	// 		switch options[0].Name {
	// 		case "nested-subcommand":
	// 			content = "Nice, now you know how to execute nested commands too"
	// 		default:
	// 			content = "Oops, something went wrong.\n" +
	// 				"Hol' up, you aren't supposed to see this message."
	// 		}
	// 	}

	// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: &discordgo.InteractionResponseData{
	// 			Content: content,
	// 		},
	// 	})
	// },
	// "responses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	// Responses to a command are very important.
	// 	// First of all, because you need to react to the interaction
	// 	// by sending the response in 3 seconds after receiving, otherwise
	// 	// interaction will be considered invalid and you can no longer
	// 	// use the interaction token and ID for responding to the user's request

	// 	content := ""
	// 	// As you can see, the response type names used here are pretty self-explanatory,
	// 	// but for those who want more information see the official documentation
	// 	switch i.ApplicationCommandData().Options[0].IntValue() {
	// 	case int64(discordgo.InteractionResponseChannelMessageWithSource):
	// 		content =
	// 			"You just responded to an interaction, sent a message and showed the original one. " +
	// 				"Congratulations!"
	// 		content +=
	// 			"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
	// 	default:
	// 		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 			Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
	// 		})
	// 		if err != nil {
	// 			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
	// 				Content: "Something went wrong",
	// 			})
	// 		}
	// 		return
	// 	}

	// 	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
	// 		Data: &discordgo.InteractionResponseData{
	// 			Content: content,
	// 		},
	// 	})
	// 	if err != nil {
	// 		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
	// 			Content: "Something went wrong",
	// 		})
	// 		return
	// 	}
	// 	time.AfterFunc(time.Second*5, func() {
	// 		_, err = s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
	// 			Content: content + "\n\nWell, now you know how to create and edit responses. " +
	// 				"But you still don't know how to delete them... so... wait 10 seconds and this " +
	// 				"message will be deleted.",
	// 		})
	// 		if err != nil {
	// 			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
	// 				Content: "Something went wrong",
	// 			})
	// 			return
	// 		}
	// 		time.Sleep(time.Second * 10)
	// 		s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
	// 	})
	// },
	// "followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	// Followup messages are basically regular messages (you can create as many of them as you wish)
	// 	// but work as they are created by webhooks and their functionality
	// 	// is for handling additional messages after sending a response.

	// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: &discordgo.InteractionResponseData{
	// 			// Note: this isn't documented, but you can use that if you want to.
	// 			// This flag just allows you to create messages visible only for the caller of the command
	// 			// (user who triggered the command)
	// 			Flags:   1 << 6,
	// 			Content: "Surprise!",
	// 		},
	// 	})
	// 	msg, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
	// 		Content: "Followup message has been created, after 5 seconds it will be edited",
	// 	})
	// 	if err != nil {
	// 		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
	// 			Content: "Something went wrong",
	// 		})
	// 		return
	// 	}
	// 	time.Sleep(time.Second * 5)

	// 	s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msg.ID, &discordgo.WebhookEdit{
	// 		Content: "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted.",
	// 	})

	// 	time.Sleep(time.Second * 10)

	// 	s.FollowupMessageDelete(s.State.User.ID, i.Interaction, msg.ID)

	// 	s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
	// 		Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
	// 			"take a unicorn :unicorn: as reward!\n" +
	// 			"Also, as bonus... look at the original interaction response :D",
	// 	})
	// },
	// }
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	s.AddHandler(messageCreate)
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Populating command info...")

	commands = populateCommandInfo(commandHandlers, objects)

	log.Println("Adding commands...")

	// registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	// for i, v := range commands {
	// 	cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
	// 	if err != nil {
	// 		log.Panicf("Cannot create '%v' command: %v", v.Name, err)
	// 	}
	// 	registeredCommands[i] = cmd
	// }

	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, *GuildID, commands)

	if err != nil {
		log.Println("FAILED to bulk overwrite commands")
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "!listbirds" {
		s.ChannelMessageSend(m.ChannelID, "Placeholder")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "!honk" {
		s.ChannelMessageSend(m.ChannelID, "HOOOONK")
	}
}

func populateCommandInfo(commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate),
	objects []string) (commands []*discordgo.ApplicationCommand) {
	// Loop over every object
	for _, obj := range objects {
		log.Println("Creating command for " + obj)

		// create command object
		var command *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
			Name:        fmt.Sprintf("%s", obj),
			Description: fmt.Sprintf("Allows you to upload one %s", obj),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        obj + " to upload",
					Description: "upload a new " + obj,
					Required:    true,
				},
			},
		}

		// create commandHandler
		var commandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			option := i.ApplicationCommandData().Options[0]
			key := option.Value.(string)

			attachment_url := i.Data.(discordgo.ApplicationCommandInteractionData).Resolved.Attachments[key].URL

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "TODO upload picture with url " + attachment_url,
				},
			})
		}

		// Add json to commands and commandHandlers for the object
		commands = append(commands, command)
		commandHandlers[obj] = commandHandler
	}

	return commands
}

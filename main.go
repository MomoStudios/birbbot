package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	objects = []string{
		"mimi", "momo", "mimo", "tree", "birb", "shiro",
	}

	bucket = "www.momobot.net"

	commands = []*discordgo.ApplicationCommand{}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	momos []string
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

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	out, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &objects[1],
	})

	if err != nil {
		log.Printf("Unable to fetch momos with error %v\n", err)
	} else {
		momos = make([]string, len(out.Contents))

		for i, obj := range out.Contents {
			momos[i] = *obj.Key
		}
	}

	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, *GuildID, commands)

	if err != nil {
		log.Printf("FAILED to bulk overwrite commands because %v\n", err)
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

	if m.Content == "!momo" {
		if len(momos) == 0 {
			return
		}

		var n = rand.Intn(len(momos))

		s.ChannelMessageSend(m.ChannelID, momos[n])
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
					Name:        obj + "-to-upload",
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

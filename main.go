package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

	authorized_accounts = []string{
		"198514377421225984", //km42
		"222144155604877322", //kara
		"123583494365380608", //Nile
	}

	bucket = "www.momobot.net"

	commands = []*discordgo.ApplicationCommand{}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	momos []string

	client s3.Client
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	// s.AddHandler(messageCreate)
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
	client = *s3.NewFromConfig(cfg)

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

func populateCommandInfo(commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate),
	objects []string) (commands []*discordgo.ApplicationCommand) {
	// Loop over every object
	for _, obj := range objects {
		// command to upload the obj
		var key = "upload" + obj

		log.Println("Creating command for " + key)

		// create command object
		var command *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
			Name:        key,
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

		var commandHandler = create_upload_command_handler(obj)

		// Add json to commands and commandHandlers for the object
		commands = append(commands, command)
		commandHandlers[key] = commandHandler

		key = obj

		// command to fetch the obj

		log.Println("Creating command for " + key)

		// create command object
		command = &discordgo.ApplicationCommand{
			Name:        key,
			Description: fmt.Sprintf("Allows you to fetch one %s", obj),
			Options:     []*discordgo.ApplicationCommandOption{},
		}

		commandHandler = create_fetch_command_handler(obj)

		// Add json to commands and commandHandlers for the object
		commands = append(commands, command)
		commandHandlers[key] = commandHandler

	}

	return commands
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func create_url(key string) string {
	return fmt.Sprintf("https://s3.us-east-1.amazonaws.com/%s/%s", bucket, url.QueryEscape(key))
}

func create_fetch_command_handler(obj string) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		out, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: &bucket,
			Prefix: &obj,
		})

		if err != nil {
			log.Printf("Unable to fetch %ss with error %v\n", obj, err)
			// s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Failure to load your %s :[", obj))
			respondOrLog(s, i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Failure to load your %s :[", obj),
				},
			})
			return
		}

		momos = make([]string, len(out.Contents))

		for i, obj := range out.Contents {
			momos[i] = *obj.Key
		}

		if out.ContinuationToken != nil {
			log.Printf("There were more %s available!\n", obj)
		}

		if len(momos) == 0 {
			// s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Aint any %s here", obj))
			respondOrLog(s, i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Aint any %s here", obj),
				},
			})
			return
		}

		var n = rand.Intn(len(momos))

		// s.ChannelMessageSend(i.ChannelID, create_url(momos[n]))
		respondOrLog(s, i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: create_url(momos[n]),
			},
		})
	}
}

func create_upload_command_handler(obj string) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		id := i.Member.User.ID

		if !contains(authorized_accounts, id) {
			log.Printf("Attempted upload by UNAUTHORIZED account %s", id)

			respondOrLog(s, i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Nuh uh",
				},
			})

			return
		}

		option := i.ApplicationCommandData().Options[0]
		key := option.Value.(string)

		attachment_url := i.Data.(discordgo.ApplicationCommandInteractionData).Resolved.Attachments[key].URL

		// Download from URL into memory
		resp, err := http.Get(attachment_url)

		if err != nil {
			fmt.Printf("Failed to fetch from URL %s\n", attachment_url)
			return
		} else {
			// Create reader from image
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			reader := bytes.NewReader(buf.Bytes())

			var s3key = fmt.Sprintf("%s/%d.jpg", obj, rand.Uint64())

			// Put it
			_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    &s3key,
				Body:   reader,
				ACL:    types.ObjectCannedACLPublicRead,
			})

			if err != nil {
				fmt.Printf("Failed to upload to s3 with error %v\n", err)
				respondOrLog(s, i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Failed to upload your picture with err %v", err),
					},
				})
			} else {
				respondOrLog(s, i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Uploaded your picture! Now at URL %s", create_url(s3key)),
					},
				})
			}
		}
	}
}

func respondOrLog(s *discordgo.Session, i *discordgo.Interaction, resp *discordgo.InteractionResponse) {
	err := s.InteractionRespond(i, resp)

	if err != nil {
		log.Printf("Error responding: %v", err)
	}
}

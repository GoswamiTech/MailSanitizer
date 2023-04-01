package senitizer

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GoswamiTech/MailSenitizer/mail"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	pipeDelimeter = "|"
)

type GmailService struct {
	Srvc *gmail.Service
	Ctx  context.Context
}
type CommandLineArguments struct {
	GmailService
	CredentialFile string
	User           string
	TokenFile      string
	Label          string
	SpammersFile   string
	Threads        int16
	labelList      []string
}

func (args *CommandLineArguments) NewService() error {

	args.GmailService.Ctx = context.Background()
	b, err := os.ReadFile(args.CredentialFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := mail.GetClient(config, args.TokenFile)

	args.GmailService.Srvc, err = gmail.NewService(args.GmailService.Ctx, option.WithHTTPClient(client))
	return err
}

func (args *CommandLineArguments) GetLabel() {
	labels := mail.GetLabel(args.Srvc, args.User)

	for _, label := range labels {
		if args.Label == label.Name || args.Label == label.Id {
			args.Label = label.Id
			return
		}
	}

	fmt.Printf("Please enter  the index of Labels(folders) you want to clean.\n")

	for index, label := range labels {
		fmt.Printf("%d. %s\n", index, label.Name)
	}

	label := 1
	_, err := fmt.Scanln(&label)
	if err != nil {
		log.Fatalln("[ERROR] could not read user input for label")
	}

	args.Label = labels[label].Id
}

func (args *CommandLineArguments) DeleteEmails() {

	umlc := args.Srvc.Users.Messages.List(args.User)
	umlc.LabelIds(args.Label)
	parallel := make(chan bool, args.Threads)
	umlc.Pages(args.Ctx, func(lmr *gmail.ListMessagesResponse) error {
		for _, msg := range lmr.Messages {
			parallel <- true
			go func(message *gmail.Message) error {

				defer func() {
					<-parallel
				}()

				msg, err := args.Srvc.Users.Messages.Get(args.User, message.Id).Format("full").Do()
				if err != nil {
					fmt.Printf("msg error: %v\n", err)
					// <-parallel
					return err
				}
				if args.SpammersFile != "" {
					sender := getSender(msg.Payload.Headers)
					if !isSpammer(sender) {
						fmt.Printf("%s is not spammer \n", sender)
						// <-parallel
						return nil
					}
				}
				err = args.Srvc.Users.Messages.Delete(args.User, message.Id).Do()
				if err != nil {
					fmt.Printf("message: %s - not deleted: %v\n", msg.Id, err)
				}
				fmt.Printf("message: %s - deleted\n", msg.Id)
				// <-parallel
				return nil
			}(msg)
		}

		return nil
	})

	for i := 0; i < int(args.Threads); i++ {
		parallel <- true
	}
}

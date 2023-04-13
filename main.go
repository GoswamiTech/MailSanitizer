package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GoswamiTech/MailSanitizer/sanitizer"
)

const usage = ` usage:
./bin/MailSanitizer [-credential <credential file path>] [-label <label>] [-spammers <spammers file path>]

options:
	-credential     Optional 	oauth credential filen downloaded from google console
	-label		    Optional 	labelIds of labels/folder, separated by ","
	-token 			Optiobal	path of token file
	-spammers		Optional    path of spamers file. If not provided, all mails will be deleted into label
	-user			Optional 	emailId, default value is me
`

func main() {
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }

	// command line arguments
	credential := flag.String("credential", "", "oauth credential filen downloaded from google console")
	label := flag.String("label", "", "labelIds of labels/folder, separated by comma")
	token := flag.String("token", "token.json", "path of token file.")
	spammers := flag.String("spammers", "", "path of spamers file. If not provided, all mails will be deleted into label")
	user := flag.String("user", "me", "emailId, default value is me")
	threads := int16(*flag.Int64("threads", 10, "Number of threads to be allocated for the service"))

	flag.Parse()

	args := &sanitizer.CommandLineArguments{
		CredentialFile: *credential,
		User:           *user,
		TokenFile:      *token,
		Label:          *label,
		SpammersFile:   *spammers,
		Threads:        threads,
	}

	if _, err := os.Stat(*token); err != nil && *token != "token.json" {
		log.Printf("[WARNING] token file %s doesn't exist. Therefore you need to reauthorize.\n", *token)
	}
	if _, err := os.Stat(*credential); err != nil {
		log.Fatalf("\n[ERROR]credential file: %s doen't exist.", *credential)

	}

	if args.SpammersFile == "" {
		fmt.Println("\n[Warning] spammers file path is not provided.All messages will be deleted from selected folder. Please type Y|y to continue.")
		input := ""
		fmt.Scanln(&input)
		if input != "Y" && input != "y" {
			os.Exit(0)
			return
		}
	}
	if err := args.NewService(); err != nil {
		log.Fatalf("[ERROR] error while authentication %v", err)
	}
	args.GetLabel()
	args.DeleteEmails()
}

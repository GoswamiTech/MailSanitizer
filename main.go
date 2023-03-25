package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GoswamiTech/MailSenirtizer/senitizer"
)

const usage = ` usage:
	remove-spam [-credential <credential file path>] [-label <label>] [-spammers <spammers file path>]

options:
	-credential     Optional 	oauth credential filen downloaded from google console(default value is credentials.json)
	-label		    Optional 	labelIds of labels/folder, separated by ","
	-token 			Optiobal	path of token file (default value is token.json)
	-spammers		Optional    path of spamers file. If not provided, all mails will be deleted into label
	-user			Optional 	emailId, default value is me
`

func main() {

	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }

	// command line arguments
	credential := flag.String("credential", "credentials.json", "oauth credential filen downloaded from google console")
	label := flag.String("label", "", "labelIds of labels/folder, separated by comma")
	token := flag.String("token", "token.json", "path of token file (default value is token.json)")
	spammers := flag.String("spammers", "", "path of spamers file. If not provided, all mails will be deleted into label")
	user := flag.String("user", "me", "emailId, default value is me")
	threads := int16(*flag.Int64("threads", 10, "Number of threads to be allocated for the service"))

	flag.Parse()

	args := &senitizer.CommandLineArguments{
		CredentialFile: *credential,
		User:           *user,
		TokenFile:      *token,
		Label:          *label,
		SpammersFile:   *spammers,
		Threads:        threads,
	}

	if _, err := os.Stat(*credential); err != nil {
		fmt.Printf("credential file: %s doen't exist", *credential)
		os.Exit(1)
	}

	if args.SpammersFile == "" {
		fmt.Println("[Warning] spammers file path is not provided.All messages will be deleted from selected folder. Please type Y|y to continue.")
		input := ""
		fmt.Scanln(&input)
		if input != "y" || input != "y" {
			os.Exit(0)
			return
		}
	}
	args.NewService()
	args.GetLabel()
	args.DeleteEmails()
}

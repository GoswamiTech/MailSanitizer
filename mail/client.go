package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoswamiTech/MailSenitizer/tokenListener"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config, tokFile string) *http.Client {
	// The file tokFile(token.json) stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {

	oauthCodeChannel := make(chan string)
	// Start http server to read auth code
	tokenListener.Start(oauthCodeChannel)

	defer func() {
		tokenListener.Shutdown()
	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser,"+
		"then authoerize the service to perform read,delete messages\n%v\n", authURL)

	var authCode string

	select {
	case authCode = <-oauthCodeChannel:
		log.Println("Authcode recieved")
		tok, err := config.Exchange(context.TODO(), authCode)
		if err != nil {
			log.Fatalf("Unable to retrieve token from web: %v", err)
			return nil
		}
		return tok
	case <-time.After(time.Duration(time.Second * 300)):
		log.Fatalf("Unable to get oauth code after 300 second. Please retry and authorize again")
		return nil
	}
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetLabel(srv *gmail.Service, user string) []*gmail.Label {
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return nil
	}
	return r.Labels
}

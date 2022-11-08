package goc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const TOKEN_FILE = "/token.json"
const CREDENTIALS_FILE = "/credentials.json"

func GetClient() (*calendar.Service, oauth2.TokenSource) {
	config := getConfig()
	tok := getToken(config)

	ctx := context.Background()
	source := config.TokenSource(ctx, tok)
	client := oauth2.NewClient(ctx, source)

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("unable to retrieve client: %v", err)
	}
	return srv, source
}

func getConfig() *oauth2.Config {
	path := getSharedPath()
	b := getCredentials(path)

	// Need to delete token.json on a scope change
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("unable to parse client secret file to config: %v", err)
	}
	return config
}

func getToken(config *oauth2.Config) *oauth2.Token {
	path := getSharedPath()
	tokFile := path + TOKEN_FILE
	tok := getTokenFromFile(tokFile)
	if tok == nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)

		fmt.Println("Run setup again to select calendar")
		os.Exit(0)
	}
	return tok
}

func updateToken(source oauth2.TokenSource) {
	path := getSharedPath()
	tokFile := path + TOKEN_FILE
	tok := getTokenFromFile(tokFile)
	if tok == nil {
		log.Printf("unable to read token from file")
	}

	sourceToken, err := source.Token()
	if err != nil {
		log.Printf("unable to get token from source: %v", err)
	}
	if sourceToken == nil {
		log.Println("source token is nil")
	}

	saveToken(tokFile, sourceToken)
}

func getCredentials(path string) []byte {
	b, err := os.ReadFile(path + CREDENTIALS_FILE)
	if err != nil {
		writeCredentialsInstructionsAndExit(path)
	}
	return b
}

func writeCredentialsInstructionsAndExit(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Setup and download the Google credentials.json file here: https://console.cloud.google.com/apis/credentials")
	fmt.Println("- Create a new project")
	fmt.Println("- Setup OAuth consent screen")
	fmt.Println("- Fill out all required fields")
	fmt.Println("- For scopes, add `calendarlist.readonly` and `calendar.events`")
	fmt.Println("- For test users, add your own email")
	fmt.Println("- Click on credentials and create credentials then select OAuth client ID")
	fmt.Println("  - Set application type to Desktop app and follow the steps")
	fmt.Println("  - Choose download JSON after creating the credential")
	fmt.Println("  - Rename the file you downloaded to `credentials.json`")
	fmt.Println("  - Move this file into `" + path + "`")

	fmt.Println()
	fmt.Println("- Run goc setup to continue")

	os.Exit(0)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the link below and follow the steps.\n\n%v\n\nPaste the 'code' value from the localhost-URL here: ", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("unable to scan input: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("unable to retrieve token from web: %v", err)
	}
	return tok
}

func getTokenFromFile(file string) *oauth2.Token {
	f, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return nil
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("unable to save oauth tokens: %v", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("unable to encode token: %v", err)
	}
}

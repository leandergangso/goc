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

func GetClient() *calendar.Service {
	config := getConfig()
	tok := getToken(config)

	client := config.Client(context.Background(), tok)
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	return srv
}

func getConfig() *oauth2.Config {
	b := getCredentials()

	// Need to delete token.json on a scope change
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

func getToken(config *oauth2.Config) *oauth2.Token {
	tokFile := SHARED_PATH + "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)

		fmt.Println("Run setup again to select calendar")
		os.Exit(0)
	}

	source := config.TokenSource(context.Background(), tok)
	ntok, err := source.Token()
	if err != nil {
		log.Fatalf("Unable to get token source: %v", err)
	}

	if tok.RefreshToken != ntok.RefreshToken {
		saveToken(tokFile, ntok)
	}

	return tok
}

func getCredentials() []byte {
	b, err := os.ReadFile(SHARED_PATH + "credentials.json")
	if err != nil {
		writeCredentialsInstructionsAndExit()
	}
	return b
}

func writeCredentialsInstructionsAndExit() {
	if _, err := os.Stat(SHARED_PATH); os.IsNotExist(err) {
		if err := os.Mkdir(SHARED_PATH, os.ModePerm); err != nil {
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
	fmt.Println("  - Move this file into `" + SHARED_PATH + "`")

	fmt.Println()
	fmt.Println("- Run goc setup to continue")

	os.Exit(0)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the link below and follow the steps.\n\n%v\n\nPaste the 'code' value from the localhost-URL here: ", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

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

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving new tokens to file: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to save oauth tokens: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

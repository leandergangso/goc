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

func GetClient() (*calendar.Service, oauth2.TokenSource, error) {
	sharedPath, err := getSharedPath()
	if err != nil {
		return nil, nil, err
	}

	config := getConfig(sharedPath)
	tok, err := getToken(config)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	source := config.TokenSource(ctx, tok)
	client := oauth2.NewClient(ctx, source)
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to retrieve client: %v", err)
	}

	return srv, source, nil
}

func getConfig(path string) *oauth2.Config {
	b := getCredentials(path)

	// Need to delete token.json on a scope change
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

func getToken(config *oauth2.Config) (*oauth2.Token, error) {
	sharedPath, err := getSharedPath()
	if err != nil {
		return nil, err
	}

	tokFile := sharedPath + TOKEN_FILE
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		saveToken(tokFile, tok)

		fmt.Println("Run setup again to select calendar")
		os.Exit(0)
	}

	return tok, nil
}

func updateToken(source oauth2.TokenSource) error {
	sharedPath, err := getSharedPath()
	if err != nil {
		return err
	}

	tokFile := sharedPath + TOKEN_FILE
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return err
	}

	sourceToken, err := source.Token()
	if err != nil {
		return err
	}

	if tok.RefreshToken != sourceToken.RefreshToken || tok.AccessToken != sourceToken.AccessToken {
		saveToken(tokFile, sourceToken)
	}
	return nil
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

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the link below and follow the steps.\n\n%v\n\nPaste the 'code' value from the localhost-URL here: ", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return tok, nil
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
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("unable to save oauth tokens: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

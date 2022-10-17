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
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
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
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	return b
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the link bellow and follow the steps.\n\n%v\n\nPaste the 'code value' from the localhost-URL here: ", authURL)

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

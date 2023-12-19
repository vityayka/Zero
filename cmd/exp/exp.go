package main

import (
	"context"
	"fmt"
	"github.com/lpernett/godotenv"
	"golang.org/x/oauth2"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_APP_KEY"),
		ClientSecret: os.Getenv("DROPBOX_APP_SECRET"),
		Scopes:       []string{"files.content.read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
	}

	// use PKCE to protect against CSRF attacks
	// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
	verifier := oauth2.GenerateVerifier()

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL(
		"state",
		oauth2.SetAuthURLParam("token_access_type", "offline"),
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)
	fmt.Printf("Visit the URL for the auth dialog: %v \n", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	resp, err := client.Post("https://api.dropboxapi.com/2/files/list_folder", "application/json", strings.NewReader(`{
		"path": ""
	}`))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}

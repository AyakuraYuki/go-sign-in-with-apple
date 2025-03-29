package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AyakuraYuki/go-sign-in-with-apple/apple"
)

func main() {
	// some user identifiers from team A
	someSubs := []string{
		"the_team_scope_user_identifier_1",
		"the_team_scope_user_identifier_2",
		"the_team_scope_user_identifier_3",
		"the_team_scope_user_identifier_4",
	}

	// setup recipient's credentials
	recipientAuthKey := apple.AuthKey{
		KeyID:    "OP12QR34ST",
		ClientID: "com.recipient.bundleid",
		TeamID:   "UV56WX78YZ",
		SigningKey: `-----BEGIN RSA PUBLIC KEY-----
YOUR_P8_PRIVATE_KEY
-----END RSA PUBLIC KEY-----`,
	}

	// generate the client_secret for accessing Apple's validation API
	clientSecret, err := apple.GenerateClientSecret(recipientAuthKey)
	if err != nil {
		log.Fatalf("Error generating client secret: %v", err)
	}

	// create a new Sign in with Apple client
	client, err := apple.NewClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	// obtain the user access token
	tokenRsp, err := client.ObtainMigrationAccessToken(
		ctx,
		recipientAuthKey.ClientID,
		clientSecret)
	if err != nil {
		log.Fatalf("Failed to obtain access_token: %v", err)
	}

	// transfer and exchange users from the old team to recipient team
	for _, sub := range someSubs {
		// generate the transfer identifier
		transferSub, err := client.GenerateTransferSub(
			ctx,
			recipientAuthKey.ClientID,
			recipientAuthKey.TeamID,
			clientSecret,
			tokenRsp.AccessToken,
			sub)
		if err != nil {
			log.Printf("Failed to generate transfer sub [%s]: %v", sub, err)
			continue
		}

		// <it's up to you to check if transferSub is empty>

		// exchange identifiers
		exchangedUser, err := client.ExchangeIdentifier(
			ctx,
			recipientAuthKey.ClientID,
			clientSecret,
			tokenRsp.AccessToken,
			transferSub)
		if err != nil {
			log.Printf("Failed to exchange identifier [%s]: %v", sub, err)
			continue
		}

		// do what ever you want with the new user identifiers
		fmt.Printf("new sub: %s\n", exchangedUser.Sub)
	}
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AyakuraYuki/go-sign-in-with-apple/apple"
)

func main() {
	// setup your credentials
	authKey := apple.AuthKey{
		KeyID:    "AB12CD34EF",
		ClientID: "com.yourapp.bundleid",
		TeamID:   "GH56IJ78KL",
		SigningKey: `-----BEGIN RSA PUBLIC KEY-----
YOUR_P8_PRIVATE_KEY
-----END RSA PUBLIC KEY-----`,
	}

	// generate the client_secret for accessing Apple's validation API
	clientSecret, err := apple.GenerateClientSecret(authKey)
	if err != nil {
		log.Fatalf("Error generating client secret: %v", err)
	}

	// create a new Sign in with Apple client
	client, err := apple.NewClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// do the validation
	rsp, err := client.RevokeRefreshToken(
		context.Background(),
		authKey.ClientID,
		clientSecret,
		"the_refresh_token_to_revoke",
	)
	if err != nil {
		log.Fatalf("Error validating: %v", err)
	}

	// Voila!!
	fmt.Printf("token revoked: (%s) %s", rsp.Error, rsp.ErrorDescription)
}

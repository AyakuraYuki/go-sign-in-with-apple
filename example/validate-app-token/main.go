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
	rsp, err := client.ValidateAppToken(
		context.Background(),
		authKey.ClientID,
		clientSecret,
		"the_authorization_code_to_validate",
	)
	if err != nil {
		log.Fatalf("Error validating: %v", err)
	}

	// verifying the ID token signature
	pass, token, err := client.VerifyTokenSignature(rsp.IDToken)
	if err != nil {
		log.Fatalf("Error verifying signature: %v", err)
	}
	if !pass {
		log.Fatalf("Error verifying signature, wrong id_token")
	}

	// get the unique user identifier
	unique, err := apple.GetUniqueID(token)
	if err != nil {
		log.Fatalf("Error getting unique id: %v", err)
	}

	// get the email
	email, emailVerified, isPrivateEmail, ok := apple.GetEmail(token)
	if !ok {
		log.Fatalf("Failed to get email")
	}

	// Voila!!
	fmt.Println(unique)
	fmt.Println(email)
	fmt.Println(emailVerified)
	fmt.Println(isPrivateEmail)
}

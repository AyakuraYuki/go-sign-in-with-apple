package apple

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AuthKey represents an authentication key configuration required for
// Apple Developer API access.
//
// These credentials can be created and managed in the [Certificates,
// Identifiers & Profiles > Keys] section of the Apple Developer Portal.
//
// Here is an example of a valid AuthKey configuration:
//
//	KeyID: "AB12CD34EF"
//	ClientID: "com.example.app"
//	TeamID: "GH56IJ78KL"
//	SigningKey: "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"
type AuthKey struct {
	// KeyID is the 10-character identifier assigned by Apple to the
	// authentication key, also known as the "Key Identifier" in the Apple
	// Developer Portal.
	KeyID string `json:"key_id"`

	// ClientID is your Apple Services ID, usually it is your application's
	// bundle id.
	// Example: com.example.app
	ClientID string `json:"client_id"`

	// TeamID is the 10-character Apple Developer Team identifier.
	// Visible in the top-right corner of the Apple Developer Portal.
	TeamID string `json:"team_id"`

	// SigningKey contains the PEM-encoded private key material for JWT
	// signing.
	//
	// This sensitive value should be stored securely and never commited to
	// Version Control System.
	SigningKey string `json:"-"`
}

// GenerateClientSecret generates the client_secret used to make request to
// the Sign in with Apple REST API. A client_secret expires after 6 months.
func GenerateClientSecret(authKey AuthKey) (string, error) {
	block, _ := pem.Decode([]byte(authKey.SigningKey))
	if block == nil {
		return "", errors.New("failed to decode signing private key")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	now := time.Now()

	claims := &jwt.RegisteredClaims{
		Issuer:    authKey.TeamID,
		Subject:   authKey.ClientID,
		Audience:  jwt.ClaimStrings{"https://appleid.apple.com"},
		ExpiresAt: &jwt.NumericDate{Time: now.Add(180*24*time.Hour - time.Second)}, // 6 months
		IssuedAt:  &jwt.NumericDate{Time: now},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["alg"] = "ES256"
	token.Header["kid"] = authKey.KeyID

	return token.SignedString(privKey)
}

package apple

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
)

const (
	baseURL = `https://appleid.apple.com`

	applePublicKeyURI = `/auth/keys`
	validationURI     = `/auth/token`
	revokeURI         = `/auth/revoke`
	userMigrationURI  = `/auth/usermigrationinfo`

	headerAuthorization = `Authorization`
	headerContentType   = `Content-Type`
	headerUserAgent     = `User-Agent`

	headerValueUserAgent   = `go-sign-in-with-apple`
	headerValueContentType = `application/x-www-form-urlencoded`
	headerValueAccept      = `application/json`
)

type Client interface {
	// VerifyTokenSignature for verifying the ID token signature
	//
	// Ref: https://developer.apple.com/documentation/sign_in_with_apple/processing-changes-for-sign-in-with-apple-accounts#Decode-and-validate-the-notifications
	VerifyTokenSignature(idToken string) (pass bool, token *jwt.Token, err error)

	// ValidateAppToken sends the validation request and gets TokenResponse
	//
	// @param clientID: The identifier (App ID or Services ID) for your app.
	// The identifier must not include your Team ID, to help prevent the
	// possibility of exposing sensitive data to the end user.
	// (see AuthKey.ClientID)
	//
	// @param clientSecret: A secret JSON Web Token, generated by the developer,
	// that uses the Sign in with Apple private key associated with your
	// developer account. Authorization code and refresh token validation
	// requests require this parameter. Use GenerateClientSecret function to
	// create this token.
	//
	// @param code: The authorization code received in an authorization response
	// sent to your app. The code is single-use only and valid for five minutes.
	ValidateAppToken(ctx context.Context, clientID, clientSecret, code string) (*TokenResponse, error)

	// ValidateWebToken sends the validation request and gets TokenResponse
	//
	// @param clientID: The identifier (App ID or Services ID) for your app.
	// The identifier must not include your Team ID, to help prevent the
	// possibility of exposing sensitive data to the end user.
	// (see AuthKey.ClientID)
	//
	// @param clientSecret: A secret JSON Web Token, generated by the developer,
	// that uses the Sign in with Apple private key associated with your
	// developer account. Authorization code and refresh token validation
	// requests require this parameter. Use GenerateClientSecret function to
	// create this token.
	//
	// @param code: The authorization code received in an authorization response
	// sent to your app. The code is single-use only and valid for five minutes.
	//
	// @param redirectURL: The destination URI provided in the authorization
	// request when authorizing a user with your app, if applicable. The URI
	// must use the HTTPS protocol, include a domain name, and can’t contain
	// an IP address or localhost.
	ValidateWebToken(ctx context.Context, clientID, clientSecret, code, redirectURL string) (*TokenResponse, error)

	// ValidateRefreshToken sends the validation request and gets TokenResponse
	//
	// @param clientID: The identifier (App ID or Services ID) for your app.
	// The identifier must not include your Team ID, to help prevent the
	// possibility of exposing sensitive data to the end user.
	// (see AuthKey.ClientID)
	//
	// @param clientSecret: A secret JSON Web Token, generated by the developer,
	// that uses the Sign in with Apple private key associated with your
	// developer account. Authorization code and refresh token validation
	// requests require this parameter. Use GenerateClientSecret function to
	// create this token.
	//
	// @param refreshToken: The refresh token received from the validation
	// server during an authorization request.
	ValidateRefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*TokenResponse, error)

	// RevokeAccessToken revokes the access_token
	//
	// @param clientID: The identifier (App ID or Services ID) for your app.
	// The identifier must match the value provided during the authorization
	// request for the user’s information. Additionally, the identifier must
	// not include your Team ID, to help mitigate the possibility of exposing
	// sensitive data to the end user.
	//
	// @param clientSecret: A secret JSON Web Token (JWT) that uses the Sign in
	// with Apple private key associated with your developer account. Use
	// GenerateClientSecret function to create this token.
	//
	// @param accessToken: The user access token intended to be revoked. The user
	// session associated with the token provided is revoked if the request is
	// successful.
	RevokeAccessToken(ctx context.Context, clientID, clientSecret, accessToken string) (*RevokeResponse, error)

	// RevokeRefreshToken revokes the refresh_token
	//
	// @param clientID: The identifier (App ID or Services ID) for your app.
	// The identifier must match the value provided during the authorization
	// request for the user’s information. Additionally, the identifier must
	// not include your Team ID, to help mitigate the possibility of exposing
	// sensitive data to the end user.
	//
	// @param clientSecret: A secret JSON Web Token (JWT) that uses the Sign in
	// with Apple private key associated with your developer account. Use
	// GenerateClientSecret function to create this token.
	//
	// @param refreshToken: The user refresh token intended to be revoked. The user
	// session associated with the token provided is revoked if the request is
	// successful.
	RevokeRefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*RevokeResponse, error)

	// ObtainMigrationAccessToken generates an access_key for migrating users
	//
	// In order to transfer your users, you must obtain their user access token
	// and generate a transfer identifier. You normally obtain the user access
	// token when your user signs in, or when you validate a stored refresh token.
	//
	// @param clientID: The identifier (App ID or Services ID) for the transferring
	// app. The identifier must not include your Team ID, to help mitigate the
	// possibility of exposing sensitive data to the end user.
	//
	// @param clientSecret: The client secret of the transferring team, represented
	// as a JSON Web Token (JWT). The JWT payload should contain a `sub` claim that
	// matches the transferring app’s bundle ID or associated Services ID.
	//
	// Ref: https://developer.apple.com/documentation/sign_in_with_apple/transferring-your-apps-and-users-to-another-team#Obtain-the-user-access-token
	ObtainMigrationAccessToken(ctx context.Context, clientID, clientSecret string) (rsp *TokenResponse, err error)

	// GenerateTransferSub generates a transfer identifier that can transfer user
	// from last team to a new recipient
	//
	// After the recipient team accepts the transfer, you have 60 days to generate
	// transfer identifiers for the client.
	//
	// @param clientID: The identifier (App ID or Services ID) for the transferring
	// app. The identifier must not include your Team ID, to help mitigate the
	// possibility of exposing sensitive data to the end user.
	//
	// @param clientSecret: The client secret of the transferring team, represented
	// as a JSON Web Token (JWT). The JWT payload should contain a `sub` claim that
	// matches the transferring app’s bundle ID or associated Services ID.
	//
	// @param recipientTeamID: The Team ID of the recipient team to which you
	// transfer the application.
	//
	// @param accessToken: The migration access_token obtained by recipient team.
	//
	// @param sub: The team-scoped user identifier that Apple provides.
	//
	// Ref: https://developer.apple.com/documentation/sign_in_with_apple/transferring-your-apps-and-users-to-another-team#Generate-the-transfer-identifier
	GenerateTransferSub(ctx context.Context, clientID, recipientTeamID, clientSecret, accessToken, sub string) (transferSub string, err error)

	// ExchangeIdentifier exchanges identifier of the user in recipient team
	//
	// @param clientID: The identifier (App ID or Services ID) for the transferring
	// app. The identifier must not include your Team ID, to help mitigate the
	// possibility of exposing sensitive data to the end user.
	//
	// @param clientSecret: The client secret of the transferring team, represented
	// as a JSON Web Token (JWT). The JWT payload should contain a `sub` claim that
	// matches the transferring app’s bundle ID or associated Services ID.
	//
	// @param accessToken: The migration access_token obtained by recipient team.
	//
	// @param transferSub: The transfer identifier that you obtained from the
	// sending team.
	//
	// Ref: https://developer.apple.com/documentation/sign_in_with_apple/bringing-new-apps-and-users-into-your-team#Exchange-identifiers
	ExchangeIdentifier(ctx context.Context, clientID, clientSecret, accessToken, transferSub string) (rsp *ExchangeIdentifierResponse, err error)
}

type client struct {
	client *resty.Client

	ticker         *time.Ticker
	pubkey         *JWKSet
	pubkeyUpdateAt time.Time

	onUpdatePubkeyFailed func()

	closed atomic.Bool
	stop   chan bool
}

func NewClient(opts ...Option) (Client, error) {
	c := &client{
		ticker:               time.NewTicker(32 * time.Minute),
		stop:                 make(chan bool),
		onUpdatePubkeyFailed: func() {},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.client = resty.New()
	c.client.SetTimeout(30 * time.Second)
	c.client.SetBaseURL(baseURL)
	c.client.SetRetryCount(0) // no retry
	c.client.SetHeader(headerContentType, headerValueContentType)
	c.client.SetHeader(headerUserAgent, headerValueUserAgent)

	// fetch Apple's public key
	if err := c.fetchApplePublicKey(); err != nil {
		return nil, fmt.Errorf("cannot create Sign in with Apple client cause error when fetching Apple's public key: %w", err)
	}
	if c.pubkey == nil || len(c.pubkey.Keys) == 0 {
		return nil, errors.New("cannot create Sign in with Apple client caused Apple's public key not found")
	}

	// start Apple's public key updater
	go c.startUpdater()

	return c, nil
}

func (c *client) Close() {
	c.closed.Store(true)
	c.stop <- true
	c.ticker.Stop()
}

func (c *client) startUpdater() {
	for {
		select {
		case <-c.ticker.C:
			for i := 0; i < 3; i++ {
				_ = c.fetchApplePublicKey()
				if c.pubkey != nil && len(c.pubkey.Keys) > 0 {
					break
				}
			}
			if c.pubkey == nil || len(c.pubkey.Keys) == 0 {
				c.onUpdatePubkeyFailed()
			}

		case stop := <-c.stop:
			if stop {
				return
			}
		}
	}
}

// fetch Apple's public key for verifying token signature
func (c *client) fetchApplePublicKey() (err error) {
	set := &JWKSet{Keys: make([]*Keys, 0)}
	_, err = c.client.R().SetResult(set).Get(applePublicKeyURI)
	if set != nil && len(set.Keys) > 0 {
		c.pubkey = set
		c.pubkeyUpdateAt = time.Now()
	}
	return err
}

func (c *client) loadApplePublicKey(keyID string) (pubkey *rsa.PublicKey, err error) {
	for _, key := range c.pubkey.Keys {
		if key.KID == keyID {
			n, _ := base64.StdEncoding.DecodeString(key.N)
			e, _ := base64.StdEncoding.DecodeString(key.E)
			modulus := new(big.Int).SetBytes(n)
			publicExponent := int(new(big.Int).SetBytes(e).Int64()) // ensure there's no overflow
			pubkey = &rsa.PublicKey{N: modulus, E: publicExponent}
			return pubkey, nil
		}
	}
	return nil, errors.New("missing Apple's public key")
}

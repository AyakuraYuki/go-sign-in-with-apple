package apple

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func (c *client) VerifyTokenSignature(idToken string) (pass bool, token *jwt.Token, err error) {
	if idToken == "" {
		return false, nil, errors.New("idToken is required, must not be empty")
	}

	token, _, err = jwt.NewParser().ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return false, token, err
	}

	keyID := token.Header["kid"].(string)
	audience := token.Header["aud"].([]string)[0]
	subject := token.Header["sub"].(string)

	pubkey, err := c.loadApplePublicKey(keyID)
	if err != nil {
		return false, token, err
	}

	claims := jwt.MapClaims{
		"iss": "https://appleid.apple.com",
		"aud": audience,
		"sub": subject,
	}

	token, err = jwt.ParseWithClaims(idToken, claims, func(_ *jwt.Token) (interface{}, error) {
		return pubkey, nil
	})
	if err != nil {
		return false, token, err
	}
	if token == nil {
		return false, token, errors.New("token is nil")
	}

	return true, token, nil
}

func (c *client) ValidateAppToken(ctx context.Context, clientID, clientSecret, code string) (rsp *TokenResponse, err error) {
	form := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
	}
	return c.doRequestValidation(ctx, form)
}

func (c *client) ValidateWebToken(ctx context.Context, clientID, clientSecret, code, redirectURL string) (rsp *TokenResponse, err error) {
	formData := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"redirect_uri":  redirectURL,
		"grant_type":    "authorization_code",
	}
	return c.doRequestValidation(ctx, formData)
}

func (c *client) ValidateRefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (rsp *TokenResponse, err error) {
	formData := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}
	return c.doRequestValidation(ctx, formData)
}

func (c *client) doRequestValidation(ctx context.Context, formData map[string]string) (rsp *TokenResponse, err error) {
	rsp = &TokenResponse{}

	_, err = c.client.R().
		SetContext(ctx).
		SetFormData(formData).
		SetResult(rsp).
		Post(validationURI)
	if err != nil {
		return nil, err
	}

	if rsp.Error != "" {
		return nil, fmt.Errorf("error %q: %s", rsp.Error, rsp.ErrorDescription)
	}

	return rsp, nil
}

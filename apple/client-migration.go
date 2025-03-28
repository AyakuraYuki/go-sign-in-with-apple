package apple

import (
	"context"
	"errors"
)

func (c *client) ObtainMigrationAccessToken(ctx context.Context, clientID, clientSecret string) (rsp *TokenResponse, err error) {
	formData := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"scope":         "user.migration",
		"grant_type":    "client_credentials",
	}
	return c.doRequestValidation(ctx, formData)
}

func (c *client) GenerateTransferSub(ctx context.Context, clientID, recipientTeamID, clientSecret, accessToken, sub string) (transferSub string, err error) {
	formData := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"sub":           sub,
		"target":        recipientTeamID,
	}

	rsp := &GenerateTransferSubResponse{}

	_, err = c.client.R().
		SetContext(ctx).
		SetHeader(headerAuthorization, "Bearer "+accessToken).
		SetFormData(formData).
		SetResult(rsp).
		Post(userMigrationURI)
	if err != nil {
		return "", err
	}
	if rsp == nil {
		return "", errors.New("no response")
	}

	return rsp.TransferSub, nil
}

func (c *client) ExchangeIdentifier(ctx context.Context, clientID, clientSecret, accessToken, transferSub string) (rsp *ExchangeIdentifierResponse, err error) {
	formData := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"transfer_sub":  transferSub,
	}

	rsp = &ExchangeIdentifierResponse{}

	_, err = c.client.R().
		SetContext(ctx).
		SetHeader(headerAuthorization, "Bearer "+accessToken).
		SetFormData(formData).
		SetResult(rsp).
		Post(userMigrationURI)
	if err != nil {
		return nil, err
	}
	if rsp == nil {
		return nil, errors.New("no response")
	}

	return rsp, nil
}

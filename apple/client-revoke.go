package apple

import (
	"context"
	"fmt"
)

func (c *client) RevokeAccessToken(ctx context.Context, clientID, clientSecret, accessToken string) (rsp *RevokeResponse, err error) {
	formData := map[string]string{
		"client_id":       clientID,
		"client_secret":   clientSecret,
		"token":           accessToken,
		"token_type_hint": "access_token",
	}
	return c.doRequestRevoke(ctx, formData)
}

func (c *client) RevokeRefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (rsp *RevokeResponse, err error) {
	formData := map[string]string{
		"client_id":       clientID,
		"client_secret":   clientSecret,
		"token":           refreshToken,
		"token_type_hint": "refresh_token",
	}
	return c.doRequestRevoke(ctx, formData)
}

func (c *client) doRequestRevoke(ctx context.Context, formData map[string]string) (rsp *RevokeResponse, err error) {
	rsp = &RevokeResponse{}

	_, err = c.client.R().
		SetContext(ctx).
		SetFormData(formData).
		SetResult(rsp).
		Post(revokeURI)
	if err != nil {
		return nil, err
	}

	if rsp.Error != "" {
		return nil, fmt.Errorf("error %q: %s", rsp.Error, rsp.ErrorDescription)
	}

	return rsp, nil
}

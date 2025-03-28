package apple

// Keys is an object that defines a single JSON Web Key.
type Keys struct {
	KTY string `json:"kty"` // The key type parameter setting. You must set to "RSA".
	KID string `json:"kid"` // A 10-character identifier key, obtained from your developer account.
	USE string `json:"use"` // The intended use for the public key.
	ALG string `json:"alg"` // The encryption algorithm used to encrypt the token.
	N   string `json:"n"`   // The modulus value for the RSA public key.
	E   string `json:"e"`   // The exponent value for the RSA public key.
}

// JWKSet is a set of JSON Web Key objects.
type JWKSet struct {
	Keys []*Keys `json:"keys"` // An array that contains JSON Web Key objects.
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`  // (Reserved for future use) A token used to access allowed data. Currently, no data set has been defined for access.
	TokenType    string `json:"token_type"`    // The type of access token. It will always be bearer.
	ExpiresIn    int    `json:"expires_in"`    // The amount of time, in seconds, before the access token expires.
	RefreshToken string `json:"refresh_token"` // The refresh token used to regenerate new access tokens. Store this token securely on your server.
	IDToken      string `json:"id_token"`      // A JSON Web Token that contains the user’s identity information.

	// A string that describes the reason for the unsuccessful request.
	// The string consists of a single allowed value.
	//
	// Possible Values:
	// 	invalid_request, invalid_client, invalid_grant, unauthorized_client, unsupported_grant_type, invalid_scope
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"` // More detailed precision about the current error.
}

type RevokeResponse struct {
	// A string that describes the reason for the unsuccessful request.
	// The string consists of a single allowed value.
	//
	// Possible Values:
	// 	invalid_request, invalid_client, invalid_grant, unauthorized_client, unsupported_grant_type, invalid_scope
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"` // More detailed precision about the current error.
}

type GenerateTransferSubResponse struct {
	// TransferSub is the transfer identifier for the user that you send to the recipient team
	TransferSub string `json:"transfer_sub"`
}

type ExchangeIdentifierResponse struct {
	// Sub is the recipient team-scoped identifier for the user.
	// This value is the same as the sub in the ID token issued during user sign-in.
	Sub string `json:"sub"`

	// Email is the private email address specific to the recipient team.
	// This attribute returns only if the user utilized a private email address with the transferred application.
	Email string `json:"email"`

	// IsPrivateEmail specifies if the email address provided is the private mail relay address.
	IsPrivateEmail bool `json:"is_private_email"`
}

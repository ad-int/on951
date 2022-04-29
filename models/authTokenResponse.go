package models

type AuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
}

func (auth *AuthTokenResponse) GetAuthorizationString() string {
	return auth.TokenType + " " + auth.AccessToken
}

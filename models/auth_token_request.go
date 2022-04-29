package models

type AuthTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

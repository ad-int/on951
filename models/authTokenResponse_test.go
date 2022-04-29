package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAuthorizationString(t *testing.T) {
	auth := AuthTokenResponse{
		TokenType:   "Bearer",
		AccessToken: "123",
	}
	assert.Equal(t, "Bearer 123", auth.GetAuthorizationString())
}

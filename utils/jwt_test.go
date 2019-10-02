package utils

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/assert.v1"
)

func TestNewJWTToken(t *testing.T) {
	token := NewJWTToken()
	claims := token.Claims.(jwt.MapClaims)

	assert.Equal(t, "auth", claims["iss"])
}

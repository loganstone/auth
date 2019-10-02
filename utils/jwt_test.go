package utils

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/assert.v1"
)

const (
	testEmail = "test@email.com"
)

func TestNewJWTToken(t *testing.T) {
	token := NewJWTToken(1)
	claims := token.Claims.(jwt.MapClaims)

	assert.Equal(t, "auth", claims["iss"])
	assert.Equal(t, "", claims["sub"])
	assert.Equal(t, "", claims["aud"])
}

func TestParseToken(t *testing.T) {
	token := NewJWTToken(5)
	signupToken, err := token.Signup(testEmail)
	assert.Equal(t, err, nil)

	tokenData, err := ParseJWTToken(signupToken)
	assert.Equal(t, err, nil)
	assert.Equal(t, testEmail, tokenData["aud"])
	assert.Equal(t, "Signup", tokenData["sub"])
}

func TestParseTokenWithExpired(t *testing.T) {
	token := NewJWTToken(-1)
	signupToken, err := token.Signup(testEmail)
	assert.Equal(t, err, nil)

	_, err = ParseJWTToken(signupToken)
	ve, ok := err.(*jwt.ValidationError)
	assert.Equal(t, ok, true)
	assert.Equal(t, jwt.ValidationErrorExpired, ve.Errors)
}

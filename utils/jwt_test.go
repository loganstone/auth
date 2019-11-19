package utils

import (
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

const (
	testEmailFmt  = "test-%s@email.com"
	testPassword  = "ok1234"
	testSecretkey = "thisissecertkey"
)

func getTestEmail() string {
	return fmt.Sprintf(testEmailFmt, uuid.New().String())
}

func TestNewJWTToken(t *testing.T) {
	token := NewJWTToken(1)
	assert.Equal(t, jwt.MapClaims{}, token.Claims)
}

func TestParseToken(t *testing.T) {
	testEmail := getTestEmail()
	token := NewJWTToken(5)
	signupToken, err := token.Signup(testEmail, testSecretkey)
	assert.Nil(t, err)

	signupClaims, err := ParseJWTSignupToken(signupToken, testSecretkey)
	assert.Nil(t, err)
	assert.Equal(t, testEmail, signupClaims.Email)
	assert.Equal(t, Signup, signupClaims.Subject)

	var userID uint = 1
	userEmail := getTestEmail()

	sessionToken, err := token.Session(userID, userEmail, testSecretkey)
	assert.Nil(t, err)

	sessionClaims, err := ParseJWTSessionToken(sessionToken, testSecretkey)
	assert.Nil(t, err)
	assert.Equal(t, Session, sessionClaims.Subject)

	assert.Equal(t, userEmail, sessionClaims.UserEmail)
	assert.Equal(t, userID, sessionClaims.UserID)
}

func TestParseTokenWithExpired(t *testing.T) {
	testEmail := getTestEmail()
	token := NewJWTToken(-1)
	signupToken, err := token.Signup(testEmail, testSecretkey)
	assert.Nil(t, err)

	_, err = ParseJWTSignupToken(signupToken, testSecretkey)
	ve, ok := err.(*jwt.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, jwt.ValidationErrorExpired, ve.Errors)
}

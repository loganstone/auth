package utils

import (
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"gopkg.in/go-playground/assert.v1"
)

const (
	testEmailFmt = "test-%s@email.com"
	testPassword = "ok1234"
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
	signupToken, err := token.Signup(testEmail)
	assert.Equal(t, err, nil)

	signupClaims, err := ParseJWTSignupToken(signupToken)
	assert.Equal(t, err, nil)
	assert.Equal(t, testEmail, signupClaims.Email)
	assert.Equal(t, Signup, signupClaims.Subject)

	var userID uint = 1
	userEmail := getTestEmail()

	sessionToken, err := token.Session(userID, userEmail)
	assert.Equal(t, err, nil)

	sessionClaims, err := ParseJWTSessionToken(sessionToken)
	assert.Equal(t, err, nil)
	assert.Equal(t, Session, sessionClaims.Subject)

	assert.Equal(t, userEmail, sessionClaims.UserEmail)
	assert.Equal(t, userID, sessionClaims.UserID)
}

func TestParseTokenWithExpired(t *testing.T) {
	testEmail := getTestEmail()
	token := NewJWTToken(-1)
	signupToken, err := token.Signup(testEmail)
	assert.Equal(t, err, nil)

	_, err = ParseJWTSignupToken(signupToken)
	ve, ok := err.(*jwt.ValidationError)
	assert.Equal(t, ok, true)
	assert.Equal(t, jwt.ValidationErrorExpired, ve.Errors)
}

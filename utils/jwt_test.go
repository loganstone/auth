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

func testEmail() string {
	return fmt.Sprintf(testEmailFmt, uuid.New().String())
}

func TestNewJWT(t *testing.T) {
	token := NewJWT(1)
	assert.Equal(t, jwt.MapClaims{}, token.Claims)
}

func TestParseJWT(t *testing.T) {
	email := testEmail()
	token := NewJWT(5)
	signupToken, err := token.Signup(email, testSecretkey)
	assert.Nil(t, err)

	signupClaims, err := ParseSignupJWT(signupToken, testSecretkey)
	assert.Nil(t, err)
	assert.Equal(t, email, signupClaims.Email)
	assert.Equal(t, Signup, signupClaims.Subject)

	var userID uint = 1
	userEmail := testEmail()

	sessionToken, err := token.Session(userID, userEmail, testSecretkey)
	assert.Nil(t, err)

	sessionClaims, err := ParseSessionJWT(sessionToken, testSecretkey)
	assert.Nil(t, err)
	assert.Equal(t, Session, sessionClaims.Subject)

	assert.Equal(t, userEmail, sessionClaims.UserEmail)
	assert.Equal(t, userID, sessionClaims.UserID)
}

func TestParseJWTWithExpired(t *testing.T) {
	email := testEmail()
	token := NewJWT(-1)
	signupToken, err := token.Signup(email, testSecretkey)
	assert.Nil(t, err)

	_, err = ParseSignupJWT(signupToken, testSecretkey)
	ve, ok := err.(*jwt.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, jwt.ValidationErrorExpired, ve.Errors)
}

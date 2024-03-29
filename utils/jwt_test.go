package utils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

const (
	testEmailFmt  = "test-%s@email.com"
	testPassword  = "ok1234"
	testSecretkey = "thisissecertkey"
	testIssuer    = "testIssuer"
)

func testEmail() string {
	return fmt.Sprintf(testEmailFmt, uuid.New().String())
}

func TestNewJWT(t *testing.T) {
	token := NewJWT(1)
	assert.Equal(t, jwt.MapClaims{}, token.Claims)
}

func TestJWTPaeseError(t *testing.T) {
	const fnTest = "Test"
	const signedString = "testSignedString"
	const errMessage = "test error"
	expected := fmt.Sprintf("utils,%s: %s - '%s'", fnTest, errMessage, signedString)
	err := &JWTParseError{fnTest, signedString, errors.New(errMessage)}
	assert.Equal(t, expected, err.Error())
}

func TestParseJWT(t *testing.T) {
	email := testEmail()
	token := NewJWT(5)
	signupToken, err := token.Signup(email, testSecretkey, testIssuer)
	assert.NoError(t, err)

	signupClaims, err := ParseSignupJWT(signupToken, testSecretkey)
	assert.NoError(t, err)
	assert.Equal(t, email, signupClaims.Email)
	assert.Equal(t, Signup, signupClaims.Subject)

	var userID uint = 1
	userEmail := testEmail()

	sessionToken, err := token.Session(userID, userEmail, testSecretkey, testIssuer)
	assert.NoError(t, err)

	sessionClaims, err := ParseSessionJWT(sessionToken, testSecretkey)
	assert.NoError(t, err)
	assert.Equal(t, Session, sessionClaims.Subject)

	assert.Equal(t, userEmail, sessionClaims.UserEmail)
	assert.Equal(t, userID, sessionClaims.UserID)
}

func TestParseJWTWithExpired(t *testing.T) {
	email := testEmail()
	token := NewJWT(-1)
	signupToken, err := token.Signup(email, testSecretkey, testIssuer)
	assert.NoError(t, err)

	_, err = ParseSignupJWT(signupToken, testSecretkey)
	ve, ok := err.(*jwt.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, jwt.ValidationErrorExpired, ve.Errors)

	var userID uint = 1
	userEmail := testEmail()

	sessionToken, err := token.Session(userID, userEmail, testSecretkey, testIssuer)
	assert.NoError(t, err)

	_, err = ParseSessionJWT(sessionToken, testSecretkey)
	ve, ok = err.(*jwt.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, jwt.ValidationErrorExpired, ve.Errors)
}

func TestParseJWTWithBadMethod(t *testing.T) {
	// reference - https://github.com/golang-jwt/jwt/blob/main/ecdsa_test.go#L21
	ecdsa256Token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJmb28iOiJiYXIifQ.feG39E-bn8HXAKhzDZq7yEAPWYDhZlwTn3sePJnU9VrGMmwdXAIEyoOnrjreYlVM_Z4N13eK9-TmMTWyfKJtHQ"
	expectedError := fmt.Sprintf("utils,parseWithClaims: unexpected signing method 'ES256' - '%v'", ecdsa256Token)

	_, err := ParseSignupJWT(ecdsa256Token, testSecretkey)
	assert.EqualError(t, err, expectedError)
}

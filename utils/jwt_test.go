package utils

import (
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"gopkg.in/go-playground/assert.v1"
)

const (
	testEmailFmt = "test-%s@email.com"
	testPassword = "ok1234"
)

func getTestEmail() string {
	return fmt.Sprintf(testEmailFmt, uuid.New().String())
}

func newTestUser() *models.User {
	con := db.Connection()
	defer con.Close()

	user := models.User{
		Email:    getTestEmail(),
		Password: testPassword,
	}

	_ = user.SetPassword()
	_ = db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	})
	return &user
}

func TestNewJWTToken(t *testing.T) {
	token := NewJWTToken(1)
	claims := token.Claims.(jwt.MapClaims)

	assert.Equal(t, "auth", claims["iss"])
	assert.Equal(t, "", claims["sub"])
	assert.Equal(t, "", claims["aud"])
}

func TestParseToken(t *testing.T) {
	testEmail := getTestEmail()
	token := NewJWTToken(5)
	signupToken, err := token.Signup(testEmail)
	assert.Equal(t, err, nil)

	tokenData, err := ParseJWTToken(signupToken)
	assert.Equal(t, err, nil)
	assert.Equal(t, testEmail, tokenData["aud"])
	assert.Equal(t, "Signup", tokenData["sub"])

	user := newTestUser()
	sessionToken, err := token.Session(user)
	assert.Equal(t, err, nil)

	tokenData, err = ParseJWTToken(sessionToken)
	assert.Equal(t, err, nil)
	assert.Equal(t, user.Email, tokenData["aud"])
	assert.Equal(t, "Authorization", tokenData["sub"])

	sessionUser, ok := tokenData["user"].(map[string]interface{})
	assert.Equal(t, ok, true)
	assert.Equal(t, user.Email, sessionUser["email"])
}

func TestParseTokenWithExpired(t *testing.T) {
	testEmail := getTestEmail()
	token := NewJWTToken(-1)
	signupToken, err := token.Signup(testEmail)
	assert.Equal(t, err, nil)

	_, err = ParseJWTToken(signupToken)
	ve, ok := err.(*jwt.ValidationError)
	assert.Equal(t, ok, true)
	assert.Equal(t, jwt.ValidationErrorExpired, ve.Errors)
}

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/utils"
)

const (
	testEmailFmt = "test-%s@email.com"
	testPassword = "ok12345678"
)

func getTestEmail() string {
	return fmt.Sprintf(testEmailFmt, uuid.New().String())
}

func TestUser(t *testing.T) {
	conf := configs.App()
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errPayload := createNewUser(&user)
	assert.Equal(t, errPayload == nil, true)

	router := NewTest()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Equal(t, err, nil)

	token := utils.NewJWTToken(10)
	sessionToken, err := token.Session(user.ID, user.Email, conf.JWTSigninKey)
	assert.Equal(t, err, nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionToken))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody models.JSONUser
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, email, resBody.Email)
	assert.Equal(t, int64(0), resBody.OTPConfirmedAt)
}

func TestUserWithNonexistentEmail(t *testing.T) {
	conf := configs.App()
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errPayload := createNewUser(&user)
	assert.Equal(t, errPayload == nil, true)

	nonexistentEmail := getTestEmail()

	router := NewTest()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", nonexistentEmail)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Equal(t, err, nil)

	token := utils.NewJWTToken(10)
	sessionToken, err := token.Session(user.ID, user.Email, conf.JWTSigninKey)
	assert.Equal(t, err, nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionToken))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteUser(t *testing.T) {
	conf := configs.App()
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errPayload := createNewUser(&user)
	assert.Equal(t, errPayload == nil, true)

	router := NewTest()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Equal(t, err, nil)

	token := utils.NewJWTToken(10)
	sessionToken, err := token.Session(user.ID, user.Email, conf.JWTSigninKey)
	assert.Equal(t, err, nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionToken))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUserAsOtherUser(t *testing.T) {
	conf := configs.App()
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errPayload := createNewUser(&user)
	assert.Equal(t, errPayload == nil, true)

	router := NewTest()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Equal(t, err, nil)

	token := utils.NewJWTToken(10)

	otherUser := models.User{
		Email:    getTestEmail(),
		Password: testPassword,
	}
	errPayload = createNewUser(&otherUser)
	assert.Equal(t, errPayload == nil, true)
	sessionToken, err := token.Session(otherUser.ID, otherUser.Email, conf.JWTSigninKey)
	assert.Equal(t, err, nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionToken))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteUserAsAdmin(t *testing.T) {
	conf := configs.App()
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errPayload := createNewUser(&user)
	assert.Equal(t, errPayload == nil, true)

	router := NewTest()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Equal(t, err, nil)

	token := utils.NewJWTToken(10)

	admin := models.User{
		Email:    getTestEmail(),
		Password: testPassword,
		IsAdmin:  true,
	}
	errPayload = createNewUser(&admin)
	assert.Equal(t, errPayload == nil, true)
	sessionToken, err := token.Session(admin.ID, admin.Email, conf.JWTSigninKey)
	assert.Equal(t, err, nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionToken))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

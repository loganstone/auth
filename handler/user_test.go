package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/db/models"
)

func TestUser(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(&user)
	assert.Equal(t, errRes.ErrorCode, 0)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody models.JSONUser
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, email, resBody.Email)
	assert.Equal(t, int64(0), resBody.OTPConfirmedAt)
}

func TestUserWithNonexistentEmail(t *testing.T) {
	email := getTestEmail()
	admin := models.User{
		Email:    email,
		Password: testPassword,
		IsAdmin:  true,
	}
	errRes := createNewUser(&admin)
	assert.Equal(t, errRes.ErrorCode, 0)

	nonexistentEmail := getTestEmail()

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s", nonexistentEmail)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &admin)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteUser(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(&user)
	assert.Equal(t, errRes.ErrorCode, 0)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUserAsOtherUser(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(&user)
	assert.Equal(t, errRes.ErrorCode, 0)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	otherUser := models.User{
		Email:    getTestEmail(),
		Password: testPassword,
	}
	errRes = createNewUser(&otherUser)
	assert.Equal(t, errRes.ErrorCode, 0)

	setSessionTokenInReqHeaderForTest(req, &otherUser)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteUserAsAdmin(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(&user)
	assert.Equal(t, errRes.ErrorCode, 0)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s", email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	admin := models.User{
		Email:    getTestEmail(),
		Password: testPassword,
		IsAdmin:  true,
	}
	errRes = createNewUser(&admin)
	assert.Equal(t, errRes.ErrorCode, 0)

	setSessionTokenInReqHeaderForTest(req, &admin)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

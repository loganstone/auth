package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/payload"
)

const (
	changedPassword = "changedPassw0rd"
)

func TestUser(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", user.Email)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody db.JSONUser
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, user.Email, resBody.Email)
	assert.Equal(t, int64(0), resBody.OTPConfirmedAt)
}

func TestUserWithNonexistentEmail(t *testing.T) {
	admin, err := testAdmin(testDBCon)
	assert.Nil(t, err)

	nonexistentEmail := testEmail()

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s", nonexistentEmail)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, admin)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var errRes payload.ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)
	assert.Equal(t, payload.ErrorCodeNotFoundUser, errRes.ErrorCode)
}

func TestDeleteUser(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", user.Email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUserAsOtherUser(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", user.Email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	otherUser := db.User{
		Email:    testEmail(),
		Password: testPassword,
	}
	err = otherUser.Create(testDBCon)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &otherUser)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteUserAsAdmin(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s", user.Email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	admin, err := testAdmin(testDBCon)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, admin)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestChangePassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	reqBody := ChangePasswordParam{
		CurrentPassword: user.Password,
		Password:        changedPassword,
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/password", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Signin
	signinReqBody := SigninParam{
		Email:    user.Email,
		Password: changedPassword,
	}
	body, err = json.Marshal(signinReqBody)
	assert.Nil(t, err)

	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resBody SiginResponse
	json.NewDecoder(w.Body).Decode(&resBody)
	assert.Equal(t, signinReqBody.Email, resBody.User.Email)
	assert.NotEqual(t, "", resBody.Token)
}

func TestChangePasswordWithIncorrectCurrentPassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	reqBody := ChangePasswordParam{
		CurrentPassword: "incorrectcurrentpassword",
		Password:        changedPassword,
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/password", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errRes payload.ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)
	assert.Equal(t, payload.ErrorCodeIncorrectPassword, errRes.ErrorCode)
}

func TestChangePasswordWithOutPassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	reqBody := ChangePasswordParam{
		CurrentPassword: user.Password,
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/password", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errRes payload.ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)
	assert.Equal(t, payload.ErrorCodeBindJSON, errRes.ErrorCode)
}

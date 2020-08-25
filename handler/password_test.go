package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	changedPassword            = "changedPassw0rd%"
	resetPasswordEmailSubject  = "[auth] Reset password."
	resetPasswordEmailBodyTmpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Reset password.</title>
</head>

<body>
    <p>Hi {{ .UserEmail}}. Do you want to reset password?</p>

    <p><a href="{{ .ResetURL }}">Reset Password</a></p>

    <p>If you don’t use this link within {{ .ExpireMin }} minutes, it will expire.</p>

    <p>Thanks,</p>
    <p>Your friends at {{ .Organization }}.</p>

    <p>If this wasn’t you, please ignore this email.</p>
</body>

</html>`
)

func TestChangePassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.NoError(t, err)

	reqBody := ChangePasswordParam{
		CurrentPassword: testPassword,
		Password:        changedPassword,
	}

	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := New()

	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/password", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.NoError(t, err)
	setAuthJWTForTest(req, user)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Signin
	signinReqBody := SigninParam{
		Email:    user.Email,
		Password: changedPassword,
	}
	body, err = json.Marshal(signinReqBody)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody SiginResponse
	err = json.NewDecoder(w.Body).Decode(&resBody)
	assert.NoError(t, err)
	assert.Equal(t, signinReqBody.Email, resBody.User.Email)
	assert.NotEqual(t, "", resBody.Token)

	w = httptest.NewRecorder()
	uri = fmt.Sprintf("/users/%s/password", testEmail())
	req, err = http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.NoError(t, err)
	setAuthJWTForTest(req, user)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestChangePasswordWithIncorrectCurrentPassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.NoError(t, err)

	reqBody := ChangePasswordParam{
		CurrentPassword: "incorrectcurrentpassword",
		Password:        changedPassword,
	}
	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := New()

	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/password", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.NoError(t, err)
	setAuthJWTForTest(req, user)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errRes ErrorCodeResponse
	err = json.NewDecoder(w.Body).Decode(&errRes)
	assert.NoError(t, err)
	assert.Equal(t, ErrorCodeIncorrectPassword, errRes.ErrorCode)
}

func TestChangePasswordWithoutPassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.NoError(t, err)

	reqBody := ChangePasswordParam{
		CurrentPassword: testPassword,
	}
	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := New()

	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/password", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.NoError(t, err)
	setAuthJWTForTest(req, user)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errRes ErrorCodeResponse
	err = json.NewDecoder(w.Body).Decode(&errRes)
	assert.NoError(t, err)
	assert.Equal(t, ErrorCodeBindJSON, errRes.ErrorCode)
}

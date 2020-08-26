package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/utils"
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

func TestSendResetPasswordEmail(t *testing.T) {
	conf := configs.App()
	user, err := testUser(testDBCon)
	assert.NoError(t, err)

	var emailBody bytes.Buffer
	data := ResetPasswordEmailData{
		UserEmail:    user.Email,
		ResetURL:     conf.ResetPasswordURL("token"),
		ExpireMin:    conf.ResetPasswordTokenExpire / oneMinuteSeconds,
		Organization: conf.Org,
	}

	emailTmpl, err := template.New("reset password email").Parse(resetPasswordEmailBodyTmpl)
	assert.NoError(t, err)
	err = emailTmpl.Execute(&emailBody, data)
	assert.NoError(t, err)

	ln, err := utils.NewLocalListener(utils.MockSMTPPort)
	assert.NoError(t, err)
	defer ln.Close()

	go func() {
		c, err := ln.Accept()
		if err != nil {
			t.Errorf("local listener accept: %v", err)
			return
		}
		defer c.Close()
		handler := utils.MockSMTPHandler{
			Con:     c,
			Name:    utils.NameFromEmail(user.Email),
			From:    conf.SupportEmail,
			To:      user.Email,
			Subject: resetPasswordEmailSubject,
			Body:    emailBody.String(),
		}
		if err := handler.Handle(); err != nil {
			t.Errorf("mock smtp handle error: %v", err)
		}
	}()
	configs.SetSMTPPort(utils.MockSMTPPort)

	reqBody := SendEmailParam{
		Email:   user.Email,
		Subject: resetPasswordEmailSubject,
		Body:    resetPasswordEmailBodyTmpl,
	}
	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := New()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/email/reset_password", bytes.NewReader(body))
	defer req.Body.Close()
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

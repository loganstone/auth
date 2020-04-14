package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/utils"
)

const (
	emailSubject = "[auth] Sign up for email address."
	emailBody    = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Please verify your email address.</title>
</head>

<body>
    <p>Hi. Do you want to create a new account?</p>

    <p>Help us secure your account by verifying your email address ({{ .UserEmail }})</p>

    <p><a href="{{ .SignupURL }}">Sign Up</a></p>

    <p>If you don’t use this link within {{ .ExpireMin }} minutes, it will expire.</p>

    <p>Thanks,</p>
    <p>Your friends at {{ .Organization }}.</p>

    <p>You’re receiving this email because you recently created a new account. If this wasn’t you, please ignore this email.</p>
</body>

</html>`
)

func TestSendVerificationEmail(t *testing.T) {
	conf := configs.App()
	reqBody := VerificationEmailParam{
		Email:   testEmail(),
		Subject: emailSubject,
		Body:    emailBody,
	}
	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup/email/verification", bytes.NewReader(body))
	defer req.Body.Close()
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody VerificationEmailResponseForTest
	err = json.NewDecoder(w.Body).Decode(&resBody)
	assert.NoError(t, err)

	claims, err := utils.ParseSignupJWT(resBody.SignupToken, conf.JWTSigninKey)
	assert.NoError(t, err)
	assert.Equal(t, reqBody.Email, claims.Email)
	assert.Equal(t, emailSubject, resBody.Subject)

	var expectedEmailBody bytes.Buffer
	emailTmpl, err := template.New("verification email").Parse(emailBody)
	assert.NoError(t, err)
	err = emailTmpl.Execute(&expectedEmailBody, resBody.VerificationEmailData)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmailBody.String(), resBody.Body)
}

func TestVerifySignupToken(t *testing.T) {
	conf := configs.App()
	email := testEmail()
	token := utils.NewJWT(conf.SignupTokenExpire)
	signupToken, err := token.Signup(email, conf.JWTSigninKey, conf.Org)
	assert.NoError(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/signup/email/verification/%s", signupToken)
	req, err := http.NewRequest("GET", uri, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]string
	err = json.NewDecoder(w.Body).Decode(&resBody)
	assert.NoError(t, err)

	assert.Equal(t, email, resBody["email"])
}

func TestVerifySignupTokenWithExpiredToken(t *testing.T) {
	conf := configs.App()
	email := testEmail()
	token := utils.NewJWT(-1)
	signupToken, err := token.Signup(email, conf.JWTSigninKey, conf.Org)
	assert.NoError(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/signup/email/verification/%s", signupToken)
	req, err := http.NewRequest("GET", uri, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resBody ErrorCodeResponse
	err = json.NewDecoder(w.Body).Decode(&resBody)
	assert.NoError(t, err)

	assert.Equal(t, ErrorCodeExpiredToken, resBody.ErrorCode)
}

func TestSignup(t *testing.T) {
	conf := configs.App()
	email := testEmail()
	token := utils.NewJWT(conf.SignupTokenExpire)
	signupToken, err := token.Signup(email, conf.JWTSigninKey, conf.Org)
	assert.NoError(t, err)

	reqBody := map[string]string{
		"token":    signupToken,
		"password": testPassword,
	}
	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup", bytes.NewReader(body))
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resBody db.User
	err = json.NewDecoder(w.Body).Decode(&resBody)
	assert.NoError(t, err)

	assert.Equal(t, email, resBody.Email)
}

func TestSignupWithShortPassword(t *testing.T) {
	conf := configs.App()
	email := testEmail()
	token := utils.NewJWT(conf.SignupTokenExpire)
	signupToken, err := token.Signup(email, conf.JWTSigninKey, conf.Org)
	assert.NoError(t, err)

	reqBody := map[string]string{
		"token":    signupToken,
		"password": "ok1234",
	}
	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup", bytes.NewReader(body))
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errRes ErrorCodeResponse
	err = json.NewDecoder(w.Body).Decode(&errRes)
	assert.NoError(t, err)

	assert.Equal(t, ErrorCodeInvalidPassword, errRes.ErrorCode)
}

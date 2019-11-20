package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

func TestSendVerificationEmail(t *testing.T) {
	conf := configs.App()
	reqBody := map[string]string{
		"email": getTestEmail(),
	}
	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup/email/verification", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.NotEqual(t, resBody["token"], "")
	token := resBody["token"]

	claims, err := utils.ParseJWTSignupToken(token, conf.JWTSigninKey)
	assert.Nil(t, err)

	assert.Equal(t, reqBody["email"], claims.Email)
}

func TestVerifySignupToken(t *testing.T) {
	conf := configs.App()
	email := getTestEmail()
	token := utils.NewJWTToken(conf.SignupTokenExpire)
	signupToken, err := token.Signup(email, conf.JWTSigninKey)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/signup/email/verification/%s", signupToken)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, email, resBody["email"])
}

func TestSignup(t *testing.T) {
	conf := configs.App()
	email := getTestEmail()
	token := utils.NewJWTToken(conf.SignupTokenExpire)
	signupToken, err := token.Signup(email, conf.JWTSigninKey)
	assert.Nil(t, err)

	reqBody := map[string]string{
		"token":    signupToken,
		"password": testPassword,
	}
	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup", bytes.NewReader(body))
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, email, resBody["email"])
}

func TestSignupWithShortPassword(t *testing.T) {
	conf := configs.App()
	email := getTestEmail()
	token := utils.NewJWTToken(conf.SignupTokenExpire)
	signupToken, err := token.Signup(email, conf.JWTSigninKey)
	assert.Nil(t, err)

	reqBody := map[string]string{
		"token":    signupToken,
		"password": "ok1234",
	}
	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup", bytes.NewReader(body))
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resBody payload.ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, payload.ErrorCodeBindJSON, resBody.ErrorCode)
}

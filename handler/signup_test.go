package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/utils"
	"github.com/stretchr/testify/assert"
)

func TestSendVerificationEmail(t *testing.T) {
	reqBody := map[string]string{
		"email": getTestEmail(),
	}
	body, err := json.Marshal(reqBody)
	assert.Equal(t, err, nil)

	router := NewTest()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup/email/verification", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Equal(t, err, nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.NotEqual(t, resBody["token"], "")
	token := resBody["token"]

	decodedToken, err := utils.ParseJWTToken(token)
	assert.Equal(t, err, nil)

	assert.Equal(t, reqBody["email"], decodedToken["aud"])
}

func TestVerifySignupToken(t *testing.T) {
	email := getTestEmail()
	token := utils.NewJWTToken(configs.App().SignupTokenExpire)
	signupToken, err := token.Signup(email)
	assert.Equal(t, err, nil)

	router := NewTest()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/signup/email/verification/%s", signupToken)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Equal(t, err, nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, email, resBody["email"])
}

func TestSignup(t *testing.T) {
	email := getTestEmail()
	token := utils.NewJWTToken(configs.App().SignupTokenExpire)
	signupToken, err := token.Signup(email)
	assert.Equal(t, err, nil)

	reqBody := map[string]string{
		"token":    signupToken,
		"password": testPassword,
	}
	body, err := json.Marshal(reqBody)
	assert.Equal(t, err, nil)

	router := NewTest()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup", bytes.NewReader(body))
	assert.Equal(t, err, nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, email, resBody["email"])
}

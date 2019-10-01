package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	decodedToken, err := utils.Load(token)
	assert.Equal(t, err, nil)

	var tokenData TokenData
	err = json.Unmarshal(decodedToken, &tokenData)
	assert.Equal(t, err, nil)

	assert.Equal(t, reqBody["email"], tokenData.Email)
}

func TestVerifySignupToken(t *testing.T) {
	email := getTestEmail()
	v, err := json.Marshal(TokenData{
		Email:     email,
		ExpiredAt: time.Now().Unix() + configs.App().SignupTokenExpire,
	})
	assert.Equal(t, err, nil)

	token, err := utils.Sign(v)
	assert.Equal(t, err, nil)

	router := NewTest()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/signup/email/verification/%s", token)
	req, err := http.NewRequest("GET", uri, nil)

	assert.Equal(t, err, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, email, resBody["email"])
}

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/db/models"
)

func TestGenerateOTP(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(&user)
	assert.Equal(t, errRes.ErrorCode, 0)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/otp", email)
	req, err := http.NewRequest("POST", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	ok := reloadUser(&user)
	assert.True(t, ok)

	assert.Equal(t, user.OTPSecretKey, resBody["secert_key"])
	otpLink, _ := user.OTPProvisioningURI()
	assert.Equal(t, otpLink, resBody["key_uri"])
}

func TestConfirmOTP(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(&user)
	assert.Equal(t, errRes.ErrorCode, 0)

	router := New()

	// Generate
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/otp", email)
	req, err := http.NewRequest("POST", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	ok := reloadUser(&user)
	assert.True(t, ok)

	// Confirm
	totp, err := user.GetTOTP()
	assert.Nil(t, err)
	reqBody := map[string]string{
		"otp": totp.Now(),
	}
	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	w = httptest.NewRecorder()
	uri = fmt.Sprintf("/users/%s/otp", email)
	req, err = http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody []string
	json.NewDecoder(w.Body).Decode(&resBody)

	ok = reloadUser(&user)
	assert.True(t, ok)

	assert.NotEqual(t, user.OTPConfirmedAt, 0)
	assert.Equal(t, len(resBody), 10)
	var prev string
	for _, code := range resBody {
		assert.NotEqual(t, prev, code)
		prev = code
	}

	assert.True(t, user.ConfirmedOTP())
}

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loganstone/auth/db/models"
	"github.com/stretchr/testify/assert"
)

func TestSignin(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(&user)
	assert.Equal(t, errRes.ErrorCode, 0)

	reqBody := map[string]string{
		"email":    user.Email,
		"password": user.Password,
	}
	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody SiginResponse
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, reqBody["email"], resBody.User.Email)
	assert.NotEqual(t, "", resBody.Token)
}

func TestSigninWithOTP(t *testing.T) {
	user, _ := generateOTP(t)
	user, _ = confirmOTP(t, user)

	totp, err := user.GetTOTP()
	assert.Nil(t, err)

	reqBody := map[string]string{
		"email":    user.Email,
		"password": testPassword,
		"otp":      totp.Now(),
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSigninWithBackupCode(t *testing.T) {
	user, _ := generateOTP(t)
	user, _ = confirmOTP(t, user)

	var backupCodes []string
	err := json.Unmarshal(user.OTPBackupCodes, &backupCodes)
	assert.Nil(t, err)

	reqBody := map[string]string{
		"email":    user.Email,
		"password": testPassword,
		"otp":      backupCodes[0],
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	user = getUserByEmailForTest(user.Email)
	err = json.Unmarshal(user.OTPBackupCodes, &backupCodes)
	assert.Nil(t, err)
	assert.Equal(t, 9, len(backupCodes))
}

func TestSigninWithBackupCodes(t *testing.T) {
	user, _ := generateOTP(t)
	user, _ = confirmOTP(t, user)

	var backupCodes []string
	err := json.Unmarshal(user.OTPBackupCodes, &backupCodes)
	assert.Nil(t, err)

	for _, code := range backupCodes {
		reqBody := map[string]string{
			"email":    user.Email,
			"password": testPassword,
			"otp":      code,
		}

		body, err := json.Marshal(reqBody)
		assert.Nil(t, err)

		router := New()
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/signin", bytes.NewReader(body))
		defer req.Body.Close()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
	user = getUserByEmailForTest(user.Email)
	err = json.Unmarshal(user.OTPBackupCodes, &backupCodes)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(backupCodes))

	reqBody := map[string]string{
		"email":    user.Email,
		"password": testPassword,
	}
	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSigninWithOutOTP(t *testing.T) {
	user, _ := generateOTP(t)
	user, _ = confirmOTP(t, user)

	reqBody := map[string]string{
		"email":    user.Email,
		"password": testPassword,
	}
	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

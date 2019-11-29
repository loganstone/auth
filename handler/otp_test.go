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

func TestGenerateOTP(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/otp", user.Email)
	req, err := http.NewRequest("POST", uri, nil)
	assert.Nil(t, err)

	setAuthJWTForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	user, err = user.Fetch(testDBCon)
	assert.Nil(t, err)
	otpLink, err := user.OTPProvisioningURI()
	assert.Nil(t, err)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)
	assert.Equal(t, user.OTPSecretKey, resBody["secert_key"])
	assert.Equal(t, otpLink, resBody["key_uri"])
}

func TestConfirmOTP(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	totp, err := user.GetTOTP()
	assert.Nil(t, err)
	reqBody := map[string]string{
		"otp": totp.Now(),
	}
	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/otp", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setAuthJWTForTest(req, user)

	router := New()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string][]string
	json.NewDecoder(w.Body).Decode(&resBody)

	user, err = user.Fetch(testDBCon)
	assert.Nil(t, err)

	assert.NotEqual(t, 0, user.OTPConfirmedAt)
	assert.Equal(t, len(resBody["otp_backup_codes"]), 10)
	var prev string
	for _, code := range resBody["otp_backup_codes"] {
		assert.NotEqual(t, prev, code)
		prev = code
	}

	assert.True(t, user.ConfirmedOTP())
}

func TestResetOTP(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	// Reset
	assert.Nil(t, err)
	reqBody := map[string]string{
		"backup_code": user.OTPBackupCodes.Get()[0],
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/otp", user.Email)
	req, err := http.NewRequest("DELETE", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setAuthJWTForTest(req, user)

	router := New()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	user, err = user.Fetch(testDBCon)
	assert.Nil(t, err)
	assert.False(t, user.ConfirmedOTP())
	assert.Nil(t, user.OTPConfirmedAt)
	assert.Nil(t, user.OTPBackupCodes)
}

func TestResetOTPAsAdmin(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	// Reset - Admin
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s/otp", user.Email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	admin, err := testAdmin(testDBCon)
	assert.Nil(t, err)

	setAuthJWTForTest(req, admin)

	router := New()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	user, err = user.Fetch(testDBCon)
	assert.Nil(t, err)
	assert.False(t, user.ConfirmedOTP())
	assert.Nil(t, user.OTPConfirmedAt)
	assert.Nil(t, user.OTPBackupCodes)
}

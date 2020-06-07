package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loganstone/auth/configs"
	"github.com/stretchr/testify/assert"
)

func TestGenerateOTP(t *testing.T) {
	conf := configs.App()
	user, err := testUser(testDBCon)
	assert.NoError(t, err)

	router := New()

	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/otp", user.Email)
	req, err := http.NewRequest("POST", uri, nil)
	assert.NoError(t, err)
	setAuthJWTForTest(req, user)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	user, err = user.Fetch(testDBCon)
	assert.NoError(t, err)
	otpLink, err := user.OTPProvisioningURI(conf.Org)
	assert.NoError(t, err)

	var resBody map[string]string
	err = json.NewDecoder(w.Body).Decode(&resBody)
	assert.NoError(t, err)
	assert.Equal(t, user.OTPSecretKey, resBody["secert_key"])
	assert.Equal(t, otpLink, resBody["key_uri"])
}

func TestConfirmOTP(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.NoError(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)
	totp, err := user.TOTP()
	assert.NoError(t, err)
	reqBody := map[string]string{
		"otp": totp.Now(),
	}
	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := New()

	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/otp", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.NoError(t, err)
	setAuthJWTForTest(req, user)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string][]string
	err = json.NewDecoder(w.Body).Decode(&resBody)
	assert.NoError(t, err)

	user, err = user.Fetch(testDBCon)
	assert.NoError(t, err)

	assert.NotEqual(t, 0, user.OTPConfirmedAt)
	assert.Equal(t, len(resBody["otp_backup_codes"]), 10)
	var prev string
	for _, code := range resBody["otp_backup_codes"] {
		assert.NotEqual(t, prev, code)
		prev = code
	}

	assert.True(t, user.ConfirmedOTP())
}

func TestConfirmOTPWithoutOTPSecretKey(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.NoError(t, err)

	reqBody := ConfirmOTPParam{OTP: "111111"}
	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := New()

	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/otp", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.NoError(t, err)
	setAuthJWTForTest(req, user)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	var resBody ErrorCodeResponse
	err = json.NewDecoder(w.Body).Decode(&resBody)
	assert.NoError(t, err)

	assert.Equal(t, resBody.Links[0].Rel, "otp.generate")
	assert.Equal(t, resBody.Links[0].Method, "POST")
	assert.Equal(t, resBody.Links[0].Href, fmt.Sprintf("/%s/otp", user.Email))
}

func TestResetOTP(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.NoError(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	// Reset
	assert.NoError(t, err)
	reqBody := map[string]string{
		"backup_code": user.OTPBackupCodes.Value()[0],
	}

	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := New()

	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/otp", testEmail())
	req, err := http.NewRequest("DELETE", uri, bytes.NewReader(body))
	assert.NoError(t, err)
	setAuthJWTForTest(req, user)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	w = httptest.NewRecorder()
	uri = fmt.Sprintf("/users/%s/otp", user.Email)
	req, err = http.NewRequest("DELETE", uri, bytes.NewReader(body))
	assert.NoError(t, err)
	setAuthJWTForTest(req, user)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	user, err = user.Fetch(testDBCon)
	assert.NoError(t, err)
	assert.False(t, user.ConfirmedOTP())
	assert.Nil(t, user.OTPConfirmedAt)
	assert.Nil(t, user.OTPBackupCodes)
}

func TestResetOTPAsAdmin(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.NoError(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	router := New()

	// Reset - Admin
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s/otp", user.Email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.NoError(t, err)
	admin, err := testAdmin(testDBCon)
	assert.NoError(t, err)
	setAuthJWTForTest(req, admin)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	user, err = user.Fetch(testDBCon)
	assert.NoError(t, err)
	assert.False(t, user.ConfirmedOTP())
	assert.Nil(t, user.OTPConfirmedAt)
	assert.Nil(t, user.OTPBackupCodes)
}

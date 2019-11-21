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

func generateOTP(t *testing.T) (*models.User, map[string]string) {
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

	return getUserByEmailForTest(user.Email), resBody
}

func confirmOTP(t *testing.T, user *models.User) (*models.User, map[string][]string) {
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

	setSessionTokenInReqHeaderForTest(req, user)

	router := New()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string][]string
	json.NewDecoder(w.Body).Decode(&resBody)

	return getUserByEmailForTest(user.Email), resBody
}

func TestGenerateOTP(t *testing.T) {
	user, resBody := generateOTP(t)
	assert.Equal(t, user.OTPSecretKey, resBody["secert_key"])
	otpLink, _ := user.OTPProvisioningURI()
	assert.Equal(t, otpLink, resBody["key_uri"])
}

func TestConfirmOTP(t *testing.T) {
	user, _ := generateOTP(t)
	user, resBody := confirmOTP(t, user)
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
	user, _ := generateOTP(t)
	user, _ = confirmOTP(t, user)

	// Reset
	var backupCodes []string
	err := json.Unmarshal(user.OTPBackupCodes, &backupCodes)
	assert.Nil(t, err)
	reqBody := map[string]string{
		"backup_code": backupCodes[0],
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/otp", user.Email)
	req, err := http.NewRequest("DELETE", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, user)

	router := New()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	user = getUserByEmailForTest(user.Email)
	assert.False(t, user.ConfirmedOTP())
	assert.Nil(t, user.OTPConfirmedAt)
	assert.True(t, user.OTPBackupCodes.IsNull())
}

func TestResetOTPAsAdmin(t *testing.T) {
	user, _ := generateOTP(t)
	user, _ = confirmOTP(t, user)

	// Reset - Admin
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s/otp", user.Email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	admin := models.User{
		Email:    getTestEmail(),
		Password: testPassword,
		IsAdmin:  true,
	}
	errRes := createNewUser(&admin)
	assert.Equal(t, errRes.ErrorCode, 0)

	setSessionTokenInReqHeaderForTest(req, &admin)

	router := New()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	user = getUserByEmailForTest(user.Email)
	assert.False(t, user.ConfirmedOTP())
	assert.Nil(t, user.OTPConfirmedAt)
	assert.True(t, user.OTPBackupCodes.IsNull())
}

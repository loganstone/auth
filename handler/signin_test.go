package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loganstone/auth/db/models"
	"github.com/loganstone/auth/payload"
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

	reqBody := SigninParam{
		Email:    user.Email,
		Password: user.Password,
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

	assert.Equal(t, reqBody.Email, resBody.User.Email)
	assert.NotEqual(t, "", resBody.Token)
}

func TestSigninWithWrongPassword(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(&user)
	assert.Equal(t, errRes.ErrorCode, 0)

	reqBody := SigninParam{
		Email:    user.Email,
		Password: "wrongpassword",
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

	var resBody payload.ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, payload.ErrorCodeIncorrectPassword, resBody.ErrorCode)
}

func TestSigninWithOutEmail(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(&user)
	assert.Equal(t, errRes.ErrorCode, 0)

	reqBody := SigninParam{
		Password: user.Password,
	}
	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resBody payload.ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, payload.ErrorCodeBindJSON, resBody.ErrorCode)
}

func TestSigninWithOutPassword(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(&user)
	assert.Equal(t, errRes.ErrorCode, 0)

	reqBody := SigninParam{
		Email: user.Email,
	}
	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resBody payload.ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, payload.ErrorCodeBindJSON, resBody.ErrorCode)
}

func TestSigninWithOTP(t *testing.T) {
	con := GetDBConnection()
	defer con.Close()

	email := getTestEmail()
	user := &models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(user)
	assert.Equal(t, errRes.ErrorCode, 0)

	_, errCodeRes := generateOTP(con, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(con, user)
	assert.Nil(t, errCodeRes)

	totp, err := user.GetTOTP()
	assert.Nil(t, err)

	reqBody := SigninParam{
		Email:    user.Email,
		Password: testPassword,
		OTP:      totp.Now(),
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
	con := GetDBConnection()
	defer con.Close()

	email := getTestEmail()
	user := &models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(user)
	assert.Equal(t, errRes.ErrorCode, 0)

	_, errCodeRes := generateOTP(con, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(con, user)
	assert.Nil(t, errCodeRes)

	var backupCodes []string
	err := json.Unmarshal(user.OTPBackupCodes, &backupCodes)
	assert.Nil(t, err)

	reqBody := SigninParam{
		Email:    user.Email,
		Password: testPassword,
		OTP:      backupCodes[0],
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

func TestSigninWithAllBackupCodes(t *testing.T) {
	con := GetDBConnection()
	defer con.Close()

	email := getTestEmail()
	user := &models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(user)
	assert.Equal(t, errRes.ErrorCode, 0)

	_, errCodeRes := generateOTP(con, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(con, user)
	assert.Nil(t, errCodeRes)

	var backupCodes []string
	err := json.Unmarshal(user.OTPBackupCodes, &backupCodes)
	assert.Nil(t, err)

	// NOTE(hs.lee): 모든 백업 코드 소모
	for _, code := range backupCodes {
		reqBody := SigninParam{
			Email:    user.Email,
			Password: testPassword,
			OTP:      code,
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

	reqBody := SigninParam{
		Email:    user.Email,
		Password: testPassword,
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
	con := GetDBConnection()
	defer con.Close()

	email := getTestEmail()
	user := &models.User{
		Email:    email,
		Password: testPassword,
	}
	errRes := createNewUser(user)
	assert.Equal(t, errRes.ErrorCode, 0)

	_, errCodeRes := generateOTP(con, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(con, user)
	assert.Nil(t, errCodeRes)

	reqBody := SigninParam{
		Email:    user.Email,
		Password: testPassword,
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

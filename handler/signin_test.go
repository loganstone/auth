package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignin(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

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
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody SiginResponse
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, reqBody.Email, resBody.User.Email)
	assert.NotEqual(t, "", resBody.Token)
}

func TestSigninWithWrongPassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

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

	var errRes ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)

	assert.Equal(t, ErrorCodeIncorrectPassword, errRes.ErrorCode)
}

func TestSigninWithoutEmail(t *testing.T) {
	_, err := testUser(testDBCon)
	assert.Nil(t, err)

	reqBody := SigninParam{
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
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errRes ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)

	assert.Equal(t, ErrorCodeBindJSON, errRes.ErrorCode)
}

func TestSigninWithoutPassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

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

	var errRes ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)

	assert.Equal(t, ErrorCodeBindJSON, errRes.ErrorCode)
}

func TestSigninWithOTP(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	totp, err := user.TOTP()
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
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	reqBody := SigninParam{
		Email:    user.Email,
		Password: testPassword,
		OTP:      user.OTPBackupCodes.Value()[0],
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

	user, err = user.Fetch(testDBCon)
	assert.Nil(t, err)
	assert.Equal(t, 9, len(user.OTPBackupCodes.Value()))
}

func TestSigninWithIncorrectOTP(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	reqBody := SigninParam{
		Email:    user.Email,
		Password: testPassword,
		OTP:      "111111", // incorrect otp
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

	var errRes ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)

	assert.Equal(t, ErrorCodeIncorrectOTP, errRes.ErrorCode)
}

func TestSigninWithAllBackupCodes(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	// NOTE(hs.lee): 모든 백업 코드 소모
	for _, code := range user.OTPBackupCodes.Value() {
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

	user, err = user.Fetch(testDBCon)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(user.OTPBackupCodes.Value()))

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

	var errRes ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)

	assert.Equal(t, ErrorCodeRequireVerifyOTP, errRes.ErrorCode)
}

func TestSigninWithoutOTP(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	_, errCodeRes := generateOTP(testDBCon, user)
	assert.Nil(t, errCodeRes)

	errCodeRes = confirmOTP(testDBCon, user)
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

	var errRes ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)

	assert.Equal(t, ErrorCodeRequireVerifyOTP, errRes.ErrorCode)
}

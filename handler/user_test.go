package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/db"
)

func TestUser(t *testing.T) {
	con := GetDBConnection()
	defer con.Close()

	email := getTestEmail()
	user := db.User{
		Email:    email,
		Password: testPassword,
	}
	err := user.Create(con)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody db.JSONUser
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, email, resBody.Email)
	assert.Equal(t, int64(0), resBody.OTPConfirmedAt)
}

func TestUserWithNonexistentEmail(t *testing.T) {
	con := GetDBConnection()
	defer con.Close()

	email := getTestEmail()
	admin := db.User{
		Email:    email,
		Password: testPassword,
		IsAdmin:  true,
	}
	err := admin.Create(con)
	assert.Nil(t, err)

	nonexistentEmail := getTestEmail()

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s", nonexistentEmail)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &admin)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteUser(t *testing.T) {
	con := GetDBConnection()
	defer con.Close()

	email := getTestEmail()
	user := db.User{
		Email:    email,
		Password: testPassword,
	}
	err := user.Create(con)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUserAsOtherUser(t *testing.T) {
	con := GetDBConnection()
	defer con.Close()

	email := getTestEmail()
	user := db.User{
		Email:    email,
		Password: testPassword,
	}
	err := user.Create(con)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	otherUser := db.User{
		Email:    getTestEmail(),
		Password: testPassword,
	}
	err = otherUser.Create(con)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &otherUser)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteUserAsAdmin(t *testing.T) {
	con := GetDBConnection()
	defer con.Close()

	email := getTestEmail()
	user := db.User{
		Email:    email,
		Password: testPassword,
	}
	err := user.Create(con)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s", email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	admin := db.User{
		Email:    getTestEmail(),
		Password: testPassword,
		IsAdmin:  true,
	}
	err = admin.Create(con)
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &admin)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestChangePassword(t *testing.T) {
	con := GetDBConnection()
	defer con.Close()

	changedPassword := "changedPassw0rd"

	email := getTestEmail()
	user := db.User{
		Email:    email,
		Password: testPassword,
	}
	err := user.Create(con)
	assert.Nil(t, err)

	reqBody := ChangePasswordParam{
		CurrentPassword: user.Password,
		Password:        changedPassword,
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/password", email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setSessionTokenInReqHeaderForTest(req, &user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Signin
	signinReqBody := SigninParam{
		Email:    user.Email,
		Password: changedPassword,
	}
	body, err = json.Marshal(signinReqBody)
	assert.Nil(t, err)

	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resBody SiginResponse
	json.NewDecoder(w.Body).Decode(&resBody)
	assert.Equal(t, signinReqBody.Email, resBody.User.Email)
	assert.NotEqual(t, "", resBody.Token)
}

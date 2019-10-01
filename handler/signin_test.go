package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loganstone/auth/models"
	"github.com/stretchr/testify/assert"
)

func TestSignin(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
	}
	errPayload := createNewUser(&user)
	assert.Equal(t, errPayload == nil, true)

	reqBody := map[string]string{
		"email":    user.Email,
		"password": user.Password,
	}
	body, err := json.Marshal(reqBody)
	assert.Equal(t, err, nil)

	router := NewTest()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Equal(t, err, nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, reqBody["email"], resBody["email"])
}

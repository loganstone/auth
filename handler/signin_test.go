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
	errPayload := createNewUser(&user)
	assert.Nil(t, errPayload)

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

	var resBody map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resBody)
	resUser, ok := resBody["user"].(map[string]interface{})
	assert.True(t, ok)

	token, ok := resBody["token"].(string)
	assert.True(t, ok)

	assert.Equal(t, reqBody["email"], resUser["email"])
	assert.NotEqual(t, "", token)
}

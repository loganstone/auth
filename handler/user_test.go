package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/loganstone/auth/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	reqBody := map[string]string{
		"email":    fmt.Sprintf("test-%s@email.com", uuid.New().String()),
		"password": "ok1234",
	}
	body, err := json.Marshal(reqBody)

	assert.Equal(t, err, nil)

	router := NewTest()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/users", bytes.NewReader(body))
	defer req.Body.Close()

	assert.Equal(t, err, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, reqBody["email"], resBody["email"])
}

func TestUser(t *testing.T) {
	email := fmt.Sprintf("test-%s@email.com", uuid.New().String())
	user := models.User{
		Email:    email,
		Password: "ok1234",
	}
	errPayload := createNewUser(&user)
	assert.Equal(t, errPayload == nil, true)

	router := NewTest()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", email)
	req, err := http.NewRequest("GET", uri, nil)

	assert.Equal(t, err, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, email, resBody["email"])
}

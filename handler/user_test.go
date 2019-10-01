package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/loganstone/auth/models"
	"github.com/stretchr/testify/assert"
)

const (
	testEmailFmt = "test-%s@email.com"
	testPassword = "ok1234"
)

func getTestEmail() string {
	return fmt.Sprintf(testEmailFmt, uuid.New().String())
}

func TestUser(t *testing.T) {
	email := getTestEmail()
	user := models.User{
		Email:    email,
		Password: testPassword,
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

func TestUserWithNonexistentEmail(t *testing.T) {
	nonexistentEmail := getTestEmail()

	router := NewTest()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", nonexistentEmail)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Equal(t, err, nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

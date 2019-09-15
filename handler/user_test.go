package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	user := map[string]string{
		"email":    fmt.Sprintf("test-%s@email.com", uuid.New().String()),
		"password": "ok1234",
	}
	reqBody, err := json.Marshal(user)

	assert.Equal(t, err, nil)

	router := New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

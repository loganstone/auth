package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/types"
	"github.com/loganstone/auth/validator"
)

var (
	emailFmt = "test_%s@mail.com"
	password = "password"
)

func TestCreateUser(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = validator.New()

	params := types.AddUserParams{
		Email:    fmt.Sprintf(emailFmt, uuid.New().String()),
		Password: password,
	}
	jsonBytes, _ := json.Marshal(params)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, CreateUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var jsonUser models.JSONUser
		decoder := json.NewDecoder(rec.Body)
		if assert.NoError(t, decoder.Decode(&jsonUser)) {
			assert.Equal(t, jsonUser.Email, params.Email)
		}
	}
}

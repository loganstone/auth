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

const NumberOfUsersToCreate = 10

func TestCreateUser(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = validator.New()

	params := types.AddUserParams{
		Email:    fmt.Sprintf(testEmailFmt, uuid.New().String()),
		Password: testPassword,
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

func TestUsers(t *testing.T) {
	// Setup
	for i := 0; i < NumberOfUsersToCreate; i++ {
		email := fmt.Sprintf(testEmailFmt, uuid.New().String())
		_, err := SetUpNewTestUser(email, testPassword)
		assert.Nil(t, err)
	}

	e := echo.New()
	e.Validator = validator.New()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, Users(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestUser(t *testing.T) {
	// Setup
	email := fmt.Sprintf(testEmailFmt, uuid.New().String())
	_, err := SetUpNewTestUser(email, testPassword)
	assert.Nil(t, err)

	e := echo.New()
	e.Validator = validator.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/users/:email")
	c.SetParamNames("email")
	c.SetParamValues(email)

	// Assertions
	if assert.NoError(t, User(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var jsonUser models.JSONUser
		decoder := json.NewDecoder(rec.Body)
		if assert.NoError(t, decoder.Decode(&jsonUser)) {
			assert.Equal(t, jsonUser.Email, email)
		}
	}
}
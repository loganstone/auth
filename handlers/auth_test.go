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

func TestSignin(t *testing.T) {
	// Setup
	email := fmt.Sprintf(testEmailFmt, uuid.New().String())
	_, err := SetUpNewTestUser(email, testPassword)
	assert.Nil(t, err)

	e := echo.New()
	e.Validator = validator.New()

	params := types.SigninParams{
		Email:    email,
		Password: testPassword,
	}
	jsonBytes, _ := json.Marshal(params)

	req := httptest.NewRequest(
		http.MethodPost, "/auth/signin", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, Signin(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var jsonUser models.JSONUser
		decoder := json.NewDecoder(rec.Body)
		if assert.NoError(t, decoder.Decode(&jsonUser)) {
			assert.Equal(t, jsonUser.Email, params.Email)
		}
	}
}

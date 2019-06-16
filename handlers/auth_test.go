package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/types"
	"github.com/loganstone/auth/validator"
)

func TestAuthenticate(t *testing.T) {

	// Setup
	con := db.Connection()
	defer con.Close()
	email := fmt.Sprintf(emailFmt, uuid.New().String())
	u := models.User{Email: email}
	u.SetPassword(password)
	assert.Nil(t, db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(&u).Error
	}))

	e := echo.New()
	e.Validator = validator.New()

	params := types.AuthenticateParams{
		Email:    email,
		Password: password,
	}
	jsonBytes, _ := json.Marshal(params)

	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, Authenticate(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var jsonUser models.JSONUser
		decoder := json.NewDecoder(rec.Body)
		if assert.NoError(t, decoder.Decode(&jsonUser)) {
			assert.Equal(t, jsonUser.Email, params.Email)
		}
	}
}

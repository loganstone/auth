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

const NumberOfUsersToCreate = 100

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

func TestUsers(t *testing.T) {
	// Setup
	con := db.Connection()
	defer con.Close()
	users := make([]models.User, NumberOfUsersToCreate)
	for i := 0; i > NumberOfUsersToCreate; i++ {
		users[i] = models.User{
			Email: fmt.Sprintf(emailFmt, uuid.New().String()),
		}
		users[i].SetPassword(password)
	}

	assert.Nil(t, db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(users).Error
	}))

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

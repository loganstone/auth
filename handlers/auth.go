package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/types"
)

// Authenticate ...
func Authenticate(c echo.Context) error {
	con := db.Connection()
	defer con.Close()

	authParams := new(types.AuthenticateParams)
	if err := c.Bind(authParams); err != nil {
		return err
	}
	if err := c.Validate(authParams); err != nil {
		return err
	}

	user := models.User{}
	con.Where("email = ? AND deleted_at is NULL", authParams.Email).First(&user)
	if !user.VerifyPassword(authParams.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect Password")
	}

	return c.JSON(http.StatusOK, user)
}

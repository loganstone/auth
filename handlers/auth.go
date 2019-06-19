package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/types"
)

// Signin .
func Signin(c echo.Context) error {
	con := db.Connection()
	defer con.Close()

	params := new(types.SigninParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	if err := c.Validate(params); err != nil {
		return err
	}

	user := models.User{Email: params.Email}
	if con.Where(&user).First(&user).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "User Not Found")
	}
	if !user.VerifyPassword(params.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect Password")
	}

	return c.JSON(http.StatusOK, user)
}

package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/response"
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
		return response.ValidateError(
			c, http.StatusBadRequest, err.Error())
	}

	user := models.User{Email: params.Email}
	if con.Where(&user).First(&user).RecordNotFound() {
		return response.NotFoundUser(c)
	}
	if !user.VerifyPassword(params.Password) {
		return c.JSON(http.StatusUnauthorized,
			types.Error{
				ErrorCode: types.NotFoundUser,
				Message:   "incorrect Password",
			})
	}

	return c.JSON(http.StatusOK, user)
}

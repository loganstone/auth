package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/response"
	"github.com/loganstone/auth/types"
)

// Users ...
func Users(c echo.Context) error {
	con := db.Connection()
	defer con.Close()
	users := []models.User{}
	con.Find(&users)
	return c.JSON(http.StatusOK, users)
}

// User ...
func User(c echo.Context) error {
	con := db.Connection()
	defer con.Close()

	email := c.Param("email")
	user := models.User{Email: email}
	if con.Where(&user).First(&user).RecordNotFound() {
		return response.NotFoundUser(c)
	}

	return c.JSON(http.StatusOK, user)
}

// CreateUser ...
func CreateUser(c echo.Context) error {
	con := db.Connection()
	defer con.Close()

	params := new(types.AddUserParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	if err := c.Validate(params); err != nil {
		return response.ValidateError(
			c, http.StatusBadRequest, err.Error())
	}

	user := models.User{Email: params.Email}
	if !con.Where(&user).First(&user).RecordNotFound() {
		return c.JSON(http.StatusBadRequest,
			types.Error{
				ErrorCode: types.UserAlreadyExists,
				Message:   "user already exists",
			})
	}

	user.SetPassword(params.Password)
	err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			types.Error{
				ErrorCode: types.DBTransactionError,
				Message:   err.Error(),
			})
	}

	return c.JSON(http.StatusCreated, user)
}

// DeleteUser ...
func DeleteUser(c echo.Context) error {
	con := db.Connection()
	defer con.Close()

	email := c.Param("email")
	user := models.User{Email: email}
	if con.Where(&user).First(&user).RecordNotFound() {
		return c.NoContent(http.StatusNoContent)
	}

	err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Delete(&user).Error
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			types.Error{
				ErrorCode: types.DBTransactionError,
				Message:   err.Error(),
			})
	}

	return c.NoContent(http.StatusNoContent)
}

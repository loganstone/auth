package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
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
		return c.JSON(http.StatusNotFound,
			types.Error{
				ErrorCode: types.NotFoundUser,
				Message:   "not such user",
			})
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
		return c.JSON(http.StatusBadRequest,
			types.Error{
				ErrorCode: types.ValidateError,
				Message:   err.Error(),
			})
	}

	user := models.User{Email: params.Email}
	if !con.Where(&user).First(&user).RecordNotFound() {
		return echo.NewHTTPError(http.StatusBadRequest, "User Already Exists")
	}

	user.SetPassword(params.Password)
	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	}); err != nil {
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

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Delete(&user).Error
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}

	return c.NoContent(http.StatusNoContent)
}

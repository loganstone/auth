package handlers

import (
	"net/http"
	"strconv"

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
	user := models.User{}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong User ID")
	}

	if con.First(&user, id).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "User Not Found")
	}
	return c.JSON(http.StatusOK, user)
}

// AddUser ...
func AddUser(c echo.Context) error {
	con := db.Connection()
	defer con.Close()

	userParams := new(types.AddUserParams)
	if err := c.Bind(userParams); err != nil {
		return err
	}
	if err := c.Validate(userParams); err != nil {
		return err
	}

	user := models.User{Email: userParams.Email}
	user.SetPassword(userParams.Password)

	if err := createUser(con, &user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "User creation failed")
	}

	return c.JSON(http.StatusCreated, user)
}

func createUser(db *gorm.DB, u *models.User) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

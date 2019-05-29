package handlers

import (
	"net/http"
	"strconv"

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

	con.First(&user, id)
	if user.ID == 0 {
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

	con.Create(&user)
	return c.JSON(http.StatusCreated, user)
}

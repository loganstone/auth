package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
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
	id := c.Param("id")
	con.First(&user, id)
	if user.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, user)
}

// AddUser ...
func AddUser(c echo.Context) error {
	con := db.Connection()
	defer con.Close()

	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return err
	}
	if err := c.Validate(user); err != nil {
		return err
	}

	con.Create(&user)
	return c.JSON(http.StatusCreated, user)
}

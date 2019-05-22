package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
)

// User ...
func User(c echo.Context) error {
	con := db.Connection()
	defer con.Close()
	user := models.User{}
	con.First(&user)
	return c.JSON(http.StatusOK, user)
}

// AddUser ...
func AddUser(c echo.Context) error {
	con := db.Connection()
	defer con.Close()
	user := models.User{Email: "test@email.com"}
	con.Create(&user)
	return c.JSON(http.StatusCreated, user)
}

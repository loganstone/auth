package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/handlers"
)

func main() {
	db.AutoMigrate()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/user", handlers.User)
	e.POST("/user", handlers.AddUser)

	// Start server
	e.Logger.Fatal(e.Start(":9090"))
}

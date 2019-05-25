package main

import (
	"flag"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/handlers"
)

var (
	portToListen = flag.Int("p", configs.DefaultPort, "port to listen on")
)

func main() {
	flag.Parse()

	db.Sync()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/user", handlers.User)
	e.POST("/user", handlers.AddUser)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *portToListen)))
}

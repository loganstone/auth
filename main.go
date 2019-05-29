package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/handlers"
	"github.com/loganstone/auth/validator"
)

var (
	portToListen = flag.Int("p", configs.DefaultPort, "port to listen on")
)

func main() {
	flag.Parse()

	db.Sync()

	// Echo instance
	e := echo.New()

	e.Validator = validator.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users", handlers.Users)
	e.GET("/users/:id", handlers.User)
	e.POST("/users", handlers.AddUser)
	e.POST("/signin", handlers.Authenticate)

	// Debug uri - /debug/pprof/
	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *portToListen)))
}

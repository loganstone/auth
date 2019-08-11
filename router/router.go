package router

import (
	"github.com/labstack/echo/v4"
	"github.com/loganstone/auth/handlers"
)

// Init .
func Init(e *echo.Echo) {
	e.GET("/users", handlers.Users)
	e.GET("/users/:email", handlers.User)
	e.POST("/users", handlers.CreateUser)
	e.DELETE("/users/:email", handlers.DeleteUser)

	auth := e.Group("/auth")
	auth.POST("/signin", handlers.Signin)
}

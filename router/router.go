package router

import (
	"github.com/labstack/echo/v4"
	"github.com/loganstone/auth/handlers"
)

// Init .
func Init(e *echo.Echo) {
	users := e.Group("/users")
	users.GET("", handlers.Users)
	users.GET("/:email", handlers.User)
	users.POST("", handlers.CreateUser)
	users.DELETE("/:email", handlers.DeleteUser)

	auth := e.Group("/auth")
	auth.POST("/signin", handlers.Signin)

	test := e.Group("/test")
	test.GET("/send/email/:email", handlers.SendEmailForUser)
}

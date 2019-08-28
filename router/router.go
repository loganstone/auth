package router

import (
	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/handler"
)

// New .
func New() *gin.Engine {
	router := gin.Default()

	users := router.Group("/users")
	{
		users.GET("", handler.Users)
		users.GET("/:email", handler.User)
		users.POST("", handler.CreateUser)
		users.DELETE("/:email", handler.DeleteUser)
	}

	router.POST("signin", handler.Signin)

	return router
}

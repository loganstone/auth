package handler

import (
	"github.com/gin-gonic/gin"
)

func bind(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("", Users)
		users.GET("/:email", User)
		users.POST("", CreateUser)
		users.DELETE("/:email", DeleteUser)
	}
	r.POST("/signin", Signin)
}

package handler

import (
	"github.com/gin-gonic/gin"
)

func bind(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("", Users)
		users.GET("/:email", User)
		users.DELETE("/:email", DeleteUser)
	}

	signup := r.Group("/signup")
	{
		signup.GET("/email/verification/:token", VerifySignupToken)
		signup.POST("/email/verification", SendVerificationEmail)
		signup.POST("", Signup)
	}

	r.POST("/signin", Signin)
}

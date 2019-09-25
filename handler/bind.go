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

	signup := r.Group("/signup")
	{
		signup.GET("/email/verification/:token", VerifySignupToken)
		signup.POST("/email/verification", SendVerificationEmail)
		// TODO(hs.lee): 토큰을 받아 검증 후 User 생성으로 변경.
		signup.POST("", CreateUser)
	}

	r.POST("/signin", Signin)
}

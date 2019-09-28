package handler

import (
	"github.com/gin-gonic/gin"
)

func bind(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("", Users)
		users.GET("/:email", User)
		// TODO(hs.lee):
		// SignUp 테스트 코드가 작성되면 삭제 한다.
		users.POST("", CreateUser)
		users.DELETE("/:email", DeleteUser)
	}

	signup := r.Group("/signup")
	{
		signup.GET("/email/verification/:token", VerifySignupToken)
		signup.POST("/email/verification", SendVerificationEmail)
		signup.POST("", SignUp)
	}

	r.POST("/signin", Signin)
}

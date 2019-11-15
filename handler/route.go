package handler

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/middleware"
)

func bind(r *gin.Engine) {
	users := r.Group("/users")
	users.Use(middleware.Authorize())
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

// New .
func New() *gin.Engine {
	router := gin.New()
	router.Use(middleware.LogFormat())
	router.Use(middleware.RequestID())
	router.Use(gin.Recovery())

	bind(router)

	if gin.Mode() == gin.DebugMode {
		pprof.Register(router)
	}

	return router
}

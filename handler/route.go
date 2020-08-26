package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
)

func bind(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(Authorize())
	admin.Use(AuthorizedUserIsAdmin())
	{
		users := admin.Group("users")
		users.GET("", Users)
		users.GET("/:email", User)
		users.DELETE("/:email", DeleteUser)
		users.DELETE("/:email/otp", ResetOTP)
	}

	users := r.Group("/users")
	users.Use(Authorize())
	users.Use(RequesterIsAuthorizedUser())
	{
		users.GET("/:email", User)
		users.DELETE("/:email", DeleteUser)
		users.PUT("/:email/password", ChangePassword)

		users.POST("/:email/otp", GenerateOTP)
		users.PUT("/:email/otp", ConfirmOTP)
		users.DELETE("/:email/otp", ResetOTP)

		users.PUT("/:email/session", RenewSession)
	}

	signup := r.Group("/signup")
	{
		signup.GET("/email/verification/:token", VerifySignupToken)
		signup.POST("/email/verification", SendVerificationEmail)
		signup.POST("", Signup)
	}

	r.POST("/email/reset_password", SendResetPasswordEmail)
	r.POST("/signin", Signin)
}

// New .
func New() http.Handler {
	mode := configs.Mode()
	gin.SetMode(mode)

	if mode == configs.DebugMode {
		fmt.Printf(
			`[INFO] running in "debug" mode. "%s" is overwrite "%s", ignore GIN-debug message below.
`, configs.EnvMode, gin.EnvGinMode)
	}

	router := gin.New()
	if mode != configs.TestMode {
		router.Use(LogFormat())
		router.Use(RequestID())
		router.Use(gin.Recovery())
	}

	router.Use(DBConnection())
	bind(router)

	if mode == configs.DebugMode {
		pprof.Register(router)
	}

	return router
}

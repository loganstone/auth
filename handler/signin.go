package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

// SiginResponse .
type SiginResponse struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
}

// SigninParam .
type SigninParam struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	OTP      string `json:"otp"`
}

// Signin .
func Signin(c *gin.Context) {
	conf := configs.App()
	con := GetDBConnection()
	defer con.Close()

	var params SigninParam

	if err := c.ShouldBindJSON(&params); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBindJSON(err.Error()))
		return
	}

	var user db.User
	if con.Where("email = ?", params.Email).First(&user).RecordNotFound() {
		c.AbortWithStatusJSON(
			http.StatusNotFound, payload.NotFoundUser())
		return
	}

	user.Password = params.Password

	if !user.VerifyPassword() {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			payload.ErrorResponse(
				payload.ErrorCodeIncorrectPassword,
				"incorrect Password"))
		return
	}

	if user.ConfirmedOTP() {
		if params.OTP == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				payload.ErrorRequireVerifyOTP())
			return
		}

		if !user.VerifyOTP(params.OTP) {
			if !user.VerifyOTPBackupCode(params.OTP) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					payload.ErrorIncorrectOTP())
				return
			}
			err := user.RemoveOTPBackupCode(params.OTP)
			if err != nil {
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					payload.ErrorResponse(
						payload.ErrorCodeRemoveOTPBackupCode, err.Error()))
				return
			}

			if err := user.Save(con); err != nil {
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					payload.ErrorDBTransaction(err.Error()))
				return
			}
		}
	}

	token := utils.NewJWTToken(conf.SessionTokenExpire)
	sessionToken, err := token.Session(user.ID, user.Email, conf.JWTSigninKey)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorSignJWTToken(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SiginResponse{User: user, Token: sessionToken})
}

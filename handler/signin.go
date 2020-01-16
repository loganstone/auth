package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

// SigninParam .
type SigninParam struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	OTP      string `json:"otp"`
}

// SiginResponse .
type SiginResponse struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
}

// Signin .
func Signin(c *gin.Context) {
	conf := configs.App()
	con := DBConnOrAbort(c)
	if con == nil {
		return
	}

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

	if !user.VerifyPassword(params.Password) {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			payload.ErrorIncorrectPassword())
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
			if _, ok := user.OTPBackupCodes.In(params.OTP); !ok {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					payload.ErrorIncorrectOTP())
				return
			}

			// 백업코드 확인은 성공 했으니,
			// 삭제를 실패해도 Signin 은 그대로 진행.
			message := "fail delete backup code '%s', error '%s'"
			ok, err := user.OTPBackupCodes.Del(params.OTP)
			if err != nil {
				log.Printf(message, params.OTP, err.Error())
			}

			if ok {
				if err := user.Save(con); err != nil {
					log.Printf(message, params.OTP, err.Error())
				}
			}
		}
	}

	token := utils.NewJWT(conf.SessionTokenExpire)
	sessionToken, err := token.Session(user.ID, user.Email, conf.JWTSigninKey)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorSignJWTToken(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SiginResponse{User: user, Token: sessionToken})
}

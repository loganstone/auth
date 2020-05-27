package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
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
			NewErrResWithErr(ErrorCodeBindJSON, err))
		return
	}

	user := findUserByEmailOrAbort(
		params.Email, c, con, http.StatusBadRequest)
	if user == nil {
		return
	}

	if !user.VerifyPassword(params.Password) {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			NewErrRes(ErrorCodeIncorrectPassword))
		return
	}

	if user.ConfirmedOTP() {
		if params.OTP == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				NewErrRes(ErrorCodeRequireVerifyOTP))
			return
		}

		if !user.VerifyOTP(params.OTP) {
			if _, ok := user.OTPBackupCodes.In(params.OTP); !ok {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					NewErrRes(ErrorCodeIncorrectOTP))
				return
			}

			// 백업코드 확인은 성공 했으니,
			// 삭제를 실패해도 Signin 은 그대로 진행.
			message := "failed delete backup code '%s', error '%s'"
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
	sessionToken, err := token.Session(user.ID, user.Email, conf.JWTSigninKey, conf.Org)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeSignJWT, err))
		return
	}
	c.JSON(http.StatusOK, SiginResponse{User: *user, Token: sessionToken})
}

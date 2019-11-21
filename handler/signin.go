package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db/models"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

// SiginResponse .
type SiginResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

// Signin .
func Signin(c *gin.Context) {
	conf := configs.App()
	con := GetDBConnection()
	defer con.Close()

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBindJSON(err.Error()))
		return
	}

	if con.Where(&user).First(&user).RecordNotFound() {
		c.AbortWithStatusJSON(
			http.StatusNotFound, payload.NotFoundUser())
		return
	}

	if !user.VerifyPassword() {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			payload.ErrorResponse(
				payload.ErrorCodeIncorrectPassword,
				"incorrect Password"))
		return
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

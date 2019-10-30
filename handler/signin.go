package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

// Signin .
func Signin(c *gin.Context) {
	con := db.Connection()
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
			payload.ErrorWithCode(
				payload.ErrorCodeIncorrectPassword,
				"incorrect Password"))
		return
	}

	token := utils.NewJWTToken(configs.App().SessionTokenExpire)
	sessionToken, err := token.Session(user.ID, user.Email)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorSignJWTToken(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user, "token": sessionToken})
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/payload"
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

	c.JSON(http.StatusOK, user)
}

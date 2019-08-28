package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/response"
	"github.com/loganstone/auth/types"
)

// Signin .
func Signin(c *gin.Context) {
	con := db.Connection()
	defer con.Close()

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, response.BindJSONError(err.Error()))
		return
	}

	if con.Where(&user).First(&user).RecordNotFound() {
		c.JSON(http.StatusNotFound, response.NotFoundUser())
		return
	}

	if !user.VerifyPassword() {
		c.JSON(http.StatusUnauthorized,
			response.ErrorCode(types.IncorrectPassword, "incorrect Password"))
		return
	}

	c.JSON(http.StatusOK, user)
}

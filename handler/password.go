package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loganstone/auth/db"
)

// ChangePasswordParam .
type ChangePasswordParam struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// ChangePassword .
func ChangePassword(c *gin.Context) {
	con := DBConnOrAbort(c)
	if con == nil {
		return
	}

	var param ChangePasswordParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			NewErrResWithErr(ErrorCodeBindJSON, err))
		return
	}

	user := findUserByEmailOrAbort(
		c.Param("email"), c, con, http.StatusNotFound)
	if user == nil {
		return
	}

	if !user.VerifyPassword(param.CurrentPassword) {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			NewErrRes(ErrorCodeIncorrectPassword))
		return
	}

	err := user.SetPassword(param.Password)
	if err != nil {
		httpStatusCode := http.StatusInternalServerError
		errRes := NewErrResWithErr(ErrorCodeSetPassword, err)
		if errors.Is(err, db.ErrorInvalidPassword) {
			httpStatusCode = http.StatusBadRequest
			errRes = NewErrResWithErr(ErrorCodeInvalidPassword, err)
		}
		c.AbortWithStatusJSON(httpStatusCode, errRes)
		return
	}

	err = user.Save(con)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeDBTransaction, err))
		return
	}

	c.Status(http.StatusOK)
}

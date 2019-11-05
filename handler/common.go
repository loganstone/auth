package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/loganstone/auth/models"
)

var (
	errEmptySessionUser = errors.New("'SessionUser' empty")
	errWrongSessionUser = errors.New("'SessionUser' not 'models.User' type")
)

// GetLoginUser .
func GetLoginUser(c *gin.Context) (loginUser models.User, err error) {
	sessionUser, ok := c.Get("SessionUser")
	if !ok {
		err = errEmptySessionUser
		return
	}
	loginUser, ok = sessionUser.(models.User)
	if !ok {
		err = errWrongSessionUser
		return
	}
	return
}

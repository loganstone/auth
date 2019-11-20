package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/db/models"
)

var (
	errEmptySessionUser = errors.New("'SessionUser' empty")
	errWrongSessionUser = errors.New("'SessionUser' not 'models.User' type")
)

// ErrorCodeResponse .
type ErrorCodeResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

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

// GetDBConnection .
func GetDBConnection() *gorm.DB {
	dbConf := configs.DB()
	return db.Connection(dbConf.ConnectionString(), dbConf.Echo)
}

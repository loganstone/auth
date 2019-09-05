package handler

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/payload"
)

const (
	defaultPageSize = "20"
)

func createNewUser(user *models.User) (errPayload gin.H) {
	con := db.Connection()
	defer con.Close()

	if !con.Where(&user).First(&user).RecordNotFound() {
		errPayload = payload.UserAlreadyExists()
		return
	}

	if err := user.SetPassword(); err != nil {
		errPayload = payload.ErrorSetPassword(err.Error())
		return
	}

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	}); err != nil {
		errPayload = payload.ErrorDBTransaction(err.Error())
		return
	}
	return
}

var (
	errPageType      = errors.New("'page' must be integer")
	errPageValue     = errors.New("'page' must not be less than zero")
	errPageSizeType  = errors.New("'page_size' must be integer")
	errPageSizeValue = errors.New("'page_size' must not be less than one")
)

// Page .
func Page(c *gin.Context) (int, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		return 0, errPageType
	}

	if page < 0 {
		return 0, errPageValue
	}

	return page, nil
}

// PageSize .
func PageSize(c *gin.Context) (int, error) {
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", defaultPageSize))
	if err != nil {
		return 0, errPageSizeType
	}

	if pageSize < 1 {
		return 0, errPageSizeValue
	}

	return pageSize, nil
}

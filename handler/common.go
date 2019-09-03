package handler

import (
	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/payload"
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

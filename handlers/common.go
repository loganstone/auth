package handlers

import (
	"github.com/jinzhu/gorm"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
)

var (
	testEmailFmt = "test_%s@mail.com"
	testPassword = "password"
)

// SetUpNewTestUser .
func SetUpNewTestUser(email string, pw string) (*models.User, error) {
	con := db.Connection()
	defer con.Close()
	u := models.User{Email: email}
	u.SetPassword(pw)
	err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(&u).Error
	})
	if err != nil {
		return nil, err
	}
	return &u, err
}

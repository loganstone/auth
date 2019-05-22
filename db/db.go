package db

import (
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" //
	"github.com/loganstone/auth/models"
)

var db *gorm.DB
var err error

// AutoMigrate ...
func AutoMigrate() {
	db := Connection()
	db.AutoMigrate(&models.User{})
	defer db.Close()
}

// Connection ..
func Connection() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("DB Connection Error")
	}
	return db
}

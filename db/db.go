package db

import (
	"fmt"
	"log"
	"os"

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
	const errMsgFmt = "'%s' environment variable is required"

	id, ok := os.LookupEnv("AUTH_DB_ID")
	if !ok {
		log.Fatalf(errMsgFmt, "AUTH_DB_ID")
	}

	pw, ok := os.LookupEnv("AUTH_DB_PW")
	if !ok {
		log.Fatalf(errMsgFmt, "AUTH_DB_PW")
	}

	name, ok := os.LookupEnv("AUTH_DB_NAME")
	if !ok {
		log.Fatalf(errMsgFmt, "AUTH_DB_NAME")
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8mb4&parseTime=True&loc=Local", id, pw, name))
	if err != nil {
		panic("DB Connection Error")
	}
	return db
}

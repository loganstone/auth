package db

import (
	"fmt"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" //
	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/models"
)

const connOpt = "charset=utf8mb4&parseTime=True&loc=Local"

// Sync ...
func Sync() {
	db := Connection()
	db.AutoMigrate(&models.User{})
	defer db.Close()
}

// Connection ..
func Connection() *gorm.DB {
	conf := configs.DB()
	connectionString := fmt.Sprintf("%s:%s@/%s?%s", conf.ID, conf.PW, conf.Name, connOpt)
	db, err := gorm.Open("mysql", connectionString)
	db.LogMode(conf.Echo)
	if err != nil {
		panic("DB Connection Error")
	}
	return db
}

// InTransaction ...
type InTransaction func(tx *gorm.DB) error

// DoInTransaction ...
func DoInTransaction(db *gorm.DB, fn InTransaction) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

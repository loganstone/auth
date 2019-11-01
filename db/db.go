package db

import (
	"log"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" // driver
)

// Connection .
func Connection(option string, echo bool) *gorm.DB {
	db, err := gorm.Open("mysql", option)
	db.LogMode(echo)
	if err != nil {
		log.Panicln("DB Connection Error")
	}
	return db
}

// InTransaction .
type InTransaction func(tx *gorm.DB) error

// DoInTransaction .
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

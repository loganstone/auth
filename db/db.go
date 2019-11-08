package db

import (
	"database/sql"
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

// ResetTestDB .
func ResetTestDB(option string) {
	db, err := sql.Open("mysql", option)
	if err != nil {
		log.Fatal("db connection failed")
	}
	defer db.Close()

	_, err = db.Exec(
		"DROP DATABASE IF EXISTS auth_test")
	if err != nil {
		log.Fatal("drop 'auth_test' database failed")
	}
	_, err = db.Exec(
		"CREATE DATABASE auth_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	if err != nil {
		log.Fatal("create 'auth_test' database failed")
	}
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

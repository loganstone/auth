package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" // driver
)

// IDField ...
type IDField struct {
	ID uint `gorm:"primary_key"`
}

// DateTimeFields ...
type DateTimeFields struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Sync .
func Sync(option string, echo bool) {
	con := Connection(option, echo)
	defer con.Close()
	con.AutoMigrate(&User{})
}

// Connection .
func Connection(option string, echo bool) *gorm.DB {
	db, err := gorm.Open("mysql", option)
	if err != nil {
		log.Panicln("DB Connection Error")
	}
	db.LogMode(echo)
	return db
}

// Reset .
func Reset(option string, dbname string) {
	db, err := sql.Open("mysql", option)
	if err != nil {
		log.Fatal("db connection failed")
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS " + dbname)
	if err != nil {
		log.Fatalf("drop '%s' database failed\n", dbname)
	}
	_, err = db.Exec(
		"CREATE DATABASE " + dbname + " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	if err != nil {
		log.Fatalf("create '%s' database failed\n", dbname)
	}
}

// Do .
type Do func(tx *gorm.DB) error

// DoInTransaction .
func DoInTransaction(db *gorm.DB, fn Do) error {
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

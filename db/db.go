package db

import (
	"database/sql"
	"fmt"
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

// SyncModels .
func SyncModels(option string, echo bool) (*gorm.DB, error) {
	const maxWait = 1000
	con, err := Connection(option, echo)
	if err != nil {
		return nil, err
	}
	con.AutoMigrate(&User{})

	wait := 0
	for wait < maxWait {
		if con.HasTable("users") {
			break
		}
		time.Sleep(time.Microsecond * 1)
		wait++
	}
	return con, nil
}

// Connection .
func Connection(dsn string, echo bool) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dsn)
	if err == nil {
		db.LogMode(echo)
	}
	return db, err
}

// Reset .
func Reset(dsn string, dbname string) error {
	db, err := sql.Open("mysql", dsn)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("db connection failed")
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS " + dbname)
	if err != nil {
		return fmt.Errorf("drop '%s' database failed", dbname)
	}

	_, err = db.Exec(
		"CREATE DATABASE " + dbname + " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	if err != nil {
		return fmt.Errorf("create '%s' database failed", dbname)
	}

	return nil
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

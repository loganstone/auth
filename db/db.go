package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" // driver
)

// IDField is primary key definition.
type IDField struct {
	ID uint `gorm:"primary_key"`
}

// DateTimeFields is base columns definition.
type DateTimeFields struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// SyncModels .
func SyncModels(dataSourceName string, echo bool) (*gorm.DB, error) {
	const maxWait = 1000
	con, err := Connection(dataSourceName, echo)
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

	if wait == maxWait {
		return nil, errors.New("sync models failed with timeout")
	}
	return con, nil
}

// Connection .
func Connection(dataSourceName string, echo bool) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	db.LogMode(echo)
	return db, nil
}

// Reset .
func Reset(dataSourceName string, dbname string) error {
	dataSourceName = strings.Split(dataSourceName, dbname)[0]
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("db connection failed: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS " + dbname)
	if err != nil {
		return fmt.Errorf("drop '%s' database failed: %w", dbname, err)
	}

	_, err = db.Exec(
		"CREATE DATABASE " + dbname + " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	if err != nil {
		return fmt.Errorf("create '%s' database failed: %w", dbname, err)
	}

	return nil
}

// Do .
type Do func(tx *gorm.DB) error

// Transaction .
func Transaction(db *gorm.DB, do Do) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := do(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

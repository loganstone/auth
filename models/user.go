package models

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

// User .
type User struct {
	IDField
	Email          string `gorm:"index;not null" binding:"required,email"`
	Password       string `gorm:"-" binding:"required"`
	HashedPassword string `gorm:"not null"`
	DateTimeFields
}

// JSONUser .
type JSONUser struct {
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// SetPassword .
func (u *User) SetPassword() error {
	passwordBytes := []byte(u.Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.HashedPassword = string(hashedBytes[:])
	return nil
}

// VerifyPassword .
func (u *User) VerifyPassword() bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(u.HashedPassword), []byte(u.Password))
	return err == nil
}

// MarshalJSON .
func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&JSONUser{
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	})
}

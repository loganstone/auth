package models

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

// User ..
type User struct {
	IDField
	Email          string `gorm:"index;not null"`
	HashedPassword string `gorm:"not null"`
	DateTimeFields
}

// JSONUser .
type JSONUser struct {
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// SetPassword ...
func (u *User) SetPassword(password string) error {
	passwordBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.HashedPassword = string(hashedBytes[:])
	return nil
}

// VerifyPassword ...
func (u *User) VerifyPassword(password string) bool {
	incoming := []byte(password)
	existing := []byte(u.HashedPassword)
	err := bcrypt.CompareHashAndPassword(existing, incoming)
	return err == nil
}

// MarshalJSON ...
func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&JSONUser{
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	})
}

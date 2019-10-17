package models

import (
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User .
type User struct {
	IDField
	Email          string `gorm:"index;not null" binding:"required,email"`
	Password       string `gorm:"-" binding:"required,gte=10,alphanum"`
	HashedPassword string `gorm:"not null"`

	OTPSecretKey   string `gorm:"size:16"`
	OTPBackupCodes JSON
	OTPConfirmedAt time.Time

	DateTimeFields
}

// JSONUser .
type JSONUser struct {
	Email          string `json:"email"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
	OTPConfirmedAt int64  `json:"otp_confirmed_at"`
}

// SetPassword .
func (u *User) SetPassword() error {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(u.Password), bcrypt.DefaultCost)
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
	user := &JSONUser{
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
	if !u.OTPConfirmedAt.IsZero() {
		user.OTPConfirmedAt = u.OTPConfirmedAt.Unix()
	}
	return json.Marshal(user)
}

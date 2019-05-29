package models

import "golang.org/x/crypto/bcrypt"

// User ..
type User struct {
	CommonFields
	Email          string `gorm:"index;not null" json:"email"`
	HashedPassword string `gorm:"not null" json:"-"`
}

// SetPassword ...
func (u *User) SetPassword(password string) error {
	passwordBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
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

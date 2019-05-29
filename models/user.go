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

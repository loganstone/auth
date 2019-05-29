package models

// User ..
type User struct {
	CommonFields
	Email          string `gorm:"index;not null" json:"email"`
	HashedPassword string `gorm:"_" json:"-"`
}

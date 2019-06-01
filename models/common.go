package models

import (
	"strconv"
	"time"
)

// JSONTime ...
type JSONTime time.Time

// CommonFields ...
type CommonFields struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt JSONTime   `json:"created_at"`
	UpdatedAt JSONTime   `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// MarshalJSON ...
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := time.Time(t).Unix()
	return []byte(strconv.FormatInt(stamp, 10)), nil
}

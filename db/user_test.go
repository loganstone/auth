package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testEmailFmt = "test-%s@email.com"
	testPassword = "Ok1234567!"
)

func TestSetPassword(t *testing.T) {
	tables := []struct {
		Password string
		Err      error
	}{
		{testPassword, nil},
		{testPassword + "more", nil},
		{"", ErrorInvalidPassword},
		{"Ok123456!", ErrorInvalidPassword},
		{"ok12345678", ErrorInvalidPassword},
		{"OK12345678", ErrorInvalidPassword},
		{"okabcdefgh", ErrorInvalidPassword},
		{"1234567890", ErrorInvalidPassword},
	}

	for _, v := range tables {
		u := User{}
		err := u.SetPassword(v.Password)
		assert.Equal(t, v.Err, err)
	}
}

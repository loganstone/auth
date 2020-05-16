package db

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testEmailFmt = "test-%s@email.com"
	testPassword = "Ok1234567!"
)

func TestCodes(t *testing.T) {
	var codes Codes
	data := []string{"111111", "222222"}
	err := codes.Set(data)
	assert.NoError(t, err)
}

func TestSetPassword(t *testing.T) {
	tables := []struct {
		Password string
		Err      error
	}{
		{testPassword, nil},
		{testPassword + "more", nil},
		{"", ErrorInvalidPassword},
		{"Ok123456!", ErrorInvalidPassword},
		{"Ok12345678", ErrorInvalidPassword},
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

func TestVerifyPassword(t *testing.T) {
	u := User{}
	err := u.SetPassword(testPassword)
	assert.NoError(t, err)
	ok := u.VerifyPassword(testPassword)
	assert.True(t, ok)
}

func TestMarshalJSON(t *testing.T) {
	const zeroUnix = -62135596800
	now := time.Now()
	email := fmt.Sprintf(testEmailFmt, "test")
	expected := fmt.Sprintf(`{"email":"%s","is_admin":false,"created_at":%d,"updated_at":%d,"deleted_at":%d,"otp_confirmed_at":%d}`,
		email, zeroUnix, zeroUnix, now.Unix(), now.Unix())
	u := User{
		Email:          email,
		OTPConfirmedAt: &now,
		DateTimeFields: DateTimeFields{DeletedAt: &now},
	}

	v, err := json.Marshal(u)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(v))
}

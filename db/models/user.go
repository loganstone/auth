package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/xlzd/gotp"
	"golang.org/x/crypto/bcrypt"
)

var (
	errEmptyOTPSecretKey = errors.New("empty 'OTPSecretKey'")
)

// User .
type User struct {
	IDField
	Email          string `gorm:"index;not null" binding:"required,email"`
	Password       string `gorm:"-" binding:"required,gte=10,alphanum"`
	HashedPassword string `gorm:"not null"`
	IsAdmin        bool   `gorm:"default:false"`

	OTPSecretKey   string `gorm:"size:16"`
	OTPBackupCodes JSON
	OTPConfirmedAt *time.Time

	DateTimeFields
}

// JSONUser .
type JSONUser struct {
	Email          string `json:"email"`
	IsAdmin        bool   `json:"is_admin"`
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
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
	if u.OTPConfirmedAt != nil {
		user.OTPConfirmedAt = u.OTPConfirmedAt.Unix()
	}
	return json.Marshal(user)
}

// GetTOTP .
func (u *User) GetTOTP() (*gotp.TOTP, error) {
	if u.OTPSecretKey == "" {
		return nil, errEmptyOTPSecretKey
	}
	return gotp.NewDefaultTOTP(u.OTPSecretKey), nil
}

// GenerateOTPSecretKey .
func (u *User) GenerateOTPSecretKey() {
	// TODO: export config
	secretLength := 16
	u.OTPSecretKey = gotp.RandomSecret(secretLength)
}

// VerifyOTP .
func (u *User) VerifyOTP(otp string) bool {
	totp, err := u.GetTOTP()
	if err != nil {
		return false
	}
	return totp.Verify(otp, int(time.Now().Unix()))
}

// OTPProvisioningURI .
func (u *User) OTPProvisioningURI() (string, error) {
	totp, err := u.GetTOTP()
	if err != nil {
		return "", err
	}
	// TODO: export config - accountName, issuerName
	return totp.ProvisioningUri("demoAccountName", "issuerName"), nil
}

// ConfirmOTP .
func (u *User) ConfirmOTP() {
	now := time.Now()
	u.OTPConfirmedAt = &now
}

// ResetOTP .
func (u *User) ResetOTP() {
	u.OTPSecretKey = ""
	u.OTPConfirmedAt = nil
}

// ConfirmedOTP .
func (u *User) ConfirmedOTP() bool {
	return u.OTPConfirmedAt != nil
}

// SetOTPBackupCodes .
func (u *User) SetOTPBackupCodes(codes []byte) error {
	return u.OTPBackupCodes.Scan(codes)
}

// VerifyOTPBackupCode .
func (u *User) VerifyOTPBackupCode(code string) bool {
	// TODO(hs.lee): 구현 필요
	return false
}

// RemoveOTPBackupCode .
func (u *User) RemoveOTPBackupCode(code string) bool {
	// TODO(hs.lee): 구현 필요
	return false
}

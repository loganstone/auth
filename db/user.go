package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/xlzd/gotp"
	"golang.org/x/crypto/bcrypt"
)

const (
	failCreateUserMessage = "fail create user '%s': %w"
)

var (
	// ErrorUserAlreadyExists .
	ErrorUserAlreadyExists = errors.New("user already exists")
	// ErrorFailSetPassword .
	ErrorFailSetPassword   = errors.New("fail set password")
	errEmptyOTPSecretKey   = errors.New("empty 'OTPSecretKey'")
	errEmptyOTPBackupCodes = errors.New("empty 'OTPBackupCodes'")
)

// User .
type User struct {
	IDField
	Email          string `gorm:"index;not null" binding:"required,email"`
	Password       string `gorm:"-" binding:"required,gte=10,alphanum"`
	HashedPassword string `gorm:"not null"`
	IsAdmin        bool   `gorm:"default:false"`

	OTPSecretKey   string `gorm:"size:16"`
	OTPBackupCodes []byte
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
		return ErrorFailSetPassword
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
	u.OTPBackupCodes = nil
}

// ConfirmedOTP .
func (u *User) ConfirmedOTP() bool {
	return u.OTPConfirmedAt != nil
}

// GetOTPBackupCodes .
func (u *User) GetOTPBackupCodes() []string {
	var backupCodes []string
	err := json.Unmarshal(u.OTPBackupCodes, &backupCodes)
	if err != nil {
		return nil
	}
	return backupCodes
}

// SetOTPBackupCodes .
func (u *User) SetOTPBackupCodes(codes []string) error {
	result, err := json.Marshal(codes)
	if err != nil {
		return err
	}

	u.OTPBackupCodes = result
	return nil
}

// VerifyOTPBackupCode .
func (u *User) VerifyOTPBackupCode(code string) bool {
	if backupCodes := u.GetOTPBackupCodes(); backupCodes != nil {
		for _, v := range backupCodes {
			if v == code {
				return true
			}
		}
	}
	return false
}

// RemoveOTPBackupCode .
func (u *User) RemoveOTPBackupCode(code string) error {
	backupCodes := u.GetOTPBackupCodes()
	if backupCodes == nil {
		return errEmptyOTPBackupCodes
	}

	var target int
	for i, v := range backupCodes {
		if v == code {
			target = i
		}
	}

	backupCodes = append(backupCodes[:target], backupCodes[target+1:]...)
	return u.SetOTPBackupCodes(backupCodes)
}

// Create .
func (u *User) Create(con *gorm.DB) error {
	if u.Email == "" {
	}
	if !con.Where("email = ?", u.Email).First(u).RecordNotFound() {
		return fmt.Errorf(failCreateUserMessage, u.Email, ErrorUserAlreadyExists)
	}

	if err := u.SetPassword(); err != nil {
		return fmt.Errorf(failCreateUserMessage, u.Email, err)
	}

	if err := DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(u).Error
	}); err != nil {
		return fmt.Errorf(failCreateUserMessage, u.Email, err)
	}
	return nil
}

// Save .
func (u *User) Save(con *gorm.DB) error {
	if err := DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Save(u).Error
	}); err != nil {
		return err
	}
	return nil
}

// Delete .
func (u *User) Delete(con *gorm.DB) error {
	if err := DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Delete(u).Error
	}); err != nil {
		return err
	}
	return nil
}

// Fetch .
func (u *User) Fetch(con *gorm.DB) (*User, error) {
	user := &User{}
	if con.Where("email = ? ", u.Email).First(user).RecordNotFound() {
		return nil, fmt.Errorf("user not found: '%s'", u.Email)
	}
	return user, nil
}

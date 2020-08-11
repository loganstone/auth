package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/xlzd/gotp"
	"golang.org/x/crypto/bcrypt"
)

const (
	failedCreateUserMessage = "failed create user '%s': %w"
	passwordMinimumLen      = 10
)

var (
	// ErrorNoEmail .
	ErrorNoEmail = errors.New("no email")
	// ErrorUserAlreadyExists .
	ErrorUserAlreadyExists = errors.New("user already exists")
	// ErrorFailedSetPassword .
	ErrorFailedSetPassword = errors.New("failed set password")
	// ErrorInvalidPassword .
	ErrorInvalidPassword = errors.New("invalid password")
	errEmptyOTPSecretKey = errors.New("empty 'OTPSecretKey'")
)

var (
	digitRegExp  = regexp.MustCompile(`[[:digit:]]+`)
	upperRegExp  = regexp.MustCompile(`[[:upper:]]+`)
	lowerRegExp  = regexp.MustCompile(`[[:lower:]]+`)
	noWordRegExp = regexp.MustCompile(`[\W]+`)
)

// Codes .
type Codes []byte

// Value .
func (c Codes) Value() []string {
	var result []string
	err := json.Unmarshal(c, &result)
	if err != nil {
		return nil
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// Set .
func (c *Codes) Set(codes []string) error {
	result, err := json.Marshal(codes)
	if err != nil {
		return err
	}

	*c = result
	return nil
}

// In .
func (c Codes) In(code string) (int, bool) {
	if codes := c.Value(); codes != nil {
		for i, v := range codes {
			if v == code {
				return i, true
			}
		}
	}
	return 0, false
}

// Del .
func (c *Codes) Del(code string) (bool, error) {
	codes := c.Value()
	if codes == nil {
		return true, nil
	}

	i, ok := c.In(code)
	if !ok {
		return true, nil
	}

	codes = append(codes[:i], codes[i+1:]...)
	err := c.Set(codes)
	if err != nil {
		return false, err
	}
	return true, nil
}

// User is user ORM.
type User struct {
	IDField
	Email          string `gorm:"index;not null" binding:"required,email"`
	HashedPassword string `gorm:"not null"`
	IsAdmin        bool   `gorm:"default:false"`

	OTPSecretKey   string `gorm:"size:16"`
	OTPBackupCodes Codes
	OTPConfirmedAt *time.Time

	DateTimeFields
}

// JSONUser .
type JSONUser struct {
	Email          string `json:"email"`
	IsAdmin        bool   `json:"is_admin"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
	DeletedAt      *int64 `json:"deleted_at"`
	OTPConfirmedAt *int64 `json:"otp_confirmed_at"`
}

// SetPassword .
func (u *User) SetPassword(password string) error {
	// TODO(logan): 에러를 세분화
	if password == "" {
		return ErrorInvalidPassword
	}

	if len(password) < passwordMinimumLen {
		return ErrorInvalidPassword
	}

	ok := digitRegExp.MatchString(password)
	if !ok {
		return ErrorInvalidPassword
	}

	ok = upperRegExp.MatchString(password)
	if !ok {
		return ErrorInvalidPassword
	}

	ok = lowerRegExp.MatchString(password)
	if !ok {
		return ErrorInvalidPassword
	}

	ok = noWordRegExp.MatchString(password)
	if !ok {
		return ErrorInvalidPassword
	}

	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrorFailedSetPassword
	}

	u.HashedPassword = string(hashedBytes[:])
	return nil
}

// VerifyPassword .
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(u.HashedPassword), []byte(password))
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
	if u.DeletedAt != nil {
		ts := u.DeletedAt.Unix()
		user.DeletedAt = &ts
	}
	if u.OTPConfirmedAt != nil {
		ts := u.OTPConfirmedAt.Unix()
		user.OTPConfirmedAt = &ts
	}
	return json.Marshal(user)
}

// GenerateOTPSecretKey .
func (u *User) GenerateOTPSecretKey(secretKeyLen int) {
	u.OTPSecretKey = gotp.RandomSecret(secretKeyLen)
}

// TOTP .
func (u *User) TOTP() (*gotp.TOTP, error) {
	if u.OTPSecretKey == "" {
		return nil, errEmptyOTPSecretKey
	}
	return gotp.NewDefaultTOTP(u.OTPSecretKey), nil
}

// VerifyOTP .
func (u *User) VerifyOTP(otp string) bool {
	totp, err := u.TOTP()
	if err != nil {
		return false
	}
	return totp.Verify(otp, int(time.Now().Unix()))
}

// OTPProvisioningURI .
func (u *User) OTPProvisioningURI(issuer string) (string, error) {
	totp, err := u.TOTP()
	if err != nil {
		return "", err
	}
	return totp.ProvisioningUri(u.Email, issuer), nil
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

// Create .
func (u *User) Create(con *gorm.DB, password string) error {
	if u.Email == "" {
		return fmt.Errorf(
			failedCreateUserMessage, u.Email, ErrorNoEmail)
	}

	if !con.Where("email = ?", u.Email).First(u).RecordNotFound() {
		return fmt.Errorf(
			failedCreateUserMessage, u.Email, ErrorUserAlreadyExists)
	}

	if err := u.SetPassword(password); err != nil {
		return fmt.Errorf(failedCreateUserMessage, u.Email, err)
	}

	do := func(tx *gorm.DB) error {
		return tx.Create(u).Error
	}
	if err := Transaction(con, do); err != nil {
		return fmt.Errorf(failedCreateUserMessage, u.Email, err)
	}
	return nil
}

// Save .
func (u *User) Save(con *gorm.DB) error {
	do := func(tx *gorm.DB) error {
		return tx.Save(u).Error
	}
	if err := Transaction(con, do); err != nil {
		return err
	}
	return nil
}

// Delete .
func (u *User) Delete(con *gorm.DB) error {
	do := func(tx *gorm.DB) error {
		return tx.Delete(u).Error
	}
	if err := Transaction(con, do); err != nil {
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

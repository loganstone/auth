package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// const .
const (
	Signup        = "Signup"
	Session       = "Session"
	ResetPasswrod = "ResetPasswrod"
)

// SessionUser .
type SessionUser struct {
	UserID    uint
	UserEmail string
}

// Token .
type Token struct {
	expireAfterSec time.Duration
	jwt.Token
}

// SignupClaims .
type SignupClaims struct {
	Email string
	jwt.StandardClaims
}

// SessionClaims .
type SessionClaims struct {
	SessionUser
	jwt.StandardClaims
}

// ResetPasswrodClaims .
type ResetPasswrodClaims struct {
	Email           string
	PasswrodResetTs int
	jwt.StandardClaims
}

// JWTParseError .
type JWTParseError struct {
	Func         string
	SignedString string
	Err          error
}

func (e *JWTParseError) Error() string {
	return fmt.Sprintf("utils,%s: %s - '%s'", e.Func, e.Err.Error(), e.SignedString)
}

// NewJWT .
func NewJWT(expireAfterSec int) *Token {
	return &Token{
		time.Duration(expireAfterSec),
		*jwt.New(jwt.SigningMethodHS256),
	}
}

func newStandardClaims(
	subject, audience, issuer string,
	expireAfterSec time.Duration,
	notBeFore int64) *jwt.StandardClaims {

	now := time.Now()
	return &jwt.StandardClaims{
		Audience:  audience,
		ExpiresAt: now.Add(time.Second * expireAfterSec).Unix(),
		Id:        uuid.New().String(),
		IssuedAt:  now.Unix(),
		Issuer:    issuer,
		NotBefore: notBeFore,
		Subject:   subject,
	}
}

// Signup .
func (t *Token) Signup(email, secretkey, issuer string) (string, error) {
	t.Claims = SignupClaims{
		email,
		*newStandardClaims(Signup, email, issuer, t.expireAfterSec, 0),
	}
	return t.SignedString([]byte(secretkey))
}

// Session .
func (t *Token) Session(userID uint, userEmail, secretkey, issuer string) (string, error) {
	t.Claims = SessionClaims{
		SessionUser{UserID: userID, UserEmail: userEmail},
		*newStandardClaims(Session, userEmail, issuer, t.expireAfterSec, 0),
	}
	return t.SignedString([]byte(secretkey))
}

// ResetPasswrod .
func (t *Token) ResetPasswrod(email string, passwordResetTs int, secretkey, issuer string) (string, error) {
	t.Claims = ResetPasswrodClaims{
		email,
		passwordResetTs,
		*newStandardClaims(ResetPasswrod, email, issuer, t.expireAfterSec, 0),
	}
	return t.SignedString([]byte(secretkey))
}

func parseWithClaims(signedString, secretkey string, claims jwt.Claims) (*jwt.Token, error) {
	const fnName = "parseWithClaims"
	return jwt.ParseWithClaims(
		signedString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err := fmt.Errorf("unexpected signing method '%v'", token.Header["alg"])
				return nil, &JWTParseError{fnName, signedString, err}
			}
			return []byte(secretkey), nil
		})
}

// ParseSignupJWT .
func ParseSignupJWT(signedString, secretkey string) (*SignupClaims, error) {
	token, err := parseWithClaims(signedString, secretkey, &SignupClaims{})
	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(*SignupClaims)
	return claims, nil
}

// ParseSessionJWT .
func ParseSessionJWT(signedString, secretkey string) (*SessionClaims, error) {
	token, err := parseWithClaims(signedString, secretkey, &SessionClaims{})
	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(*SessionClaims)
	return claims, nil
}

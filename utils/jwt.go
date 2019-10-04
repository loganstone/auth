package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/models"
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

// NewJWTToken .
func NewJWTToken(expireAfterSec int) *Token {
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
func (t *Token) Signup(email string) (string, error) {
	t.Claims = SignupClaims{
		email,
		*newStandardClaims("Signup", email, "auth", t.expireAfterSec, 0),
	}
	return t.SignedString(configs.App().JWTSigninKey)
}

// Session .
func (t *Token) Session(user *models.User) (string, error) {
	// TODO(hs.lee):
	// expireAfterSec 을 session 용으로 변경한다
	t.Claims = SessionClaims{
		SessionUser{UserID: user.ID, UserEmail: user.Email},
		*newStandardClaims("Authorization", user.Email, "auth", t.expireAfterSec, 0),
	}
	return t.SignedString(configs.App().JWTSigninKey)
}

// ParseJWTSignupToken .
func ParseJWTSignupToken(signedString string) (*SignupClaims, error) {
	// TODO(hs.lee):
	// ParseJWTSessionToken 과 코드가 중복, 중복 제거 필요
	token, err := jwt.ParseWithClaims(
		signedString,
		&SignupClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return configs.App().JWTSigninKey, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SignupClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

// ParseJWTSessionToken .
func ParseJWTSessionToken(signedString string) (*SessionClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedString,
		&SessionClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return configs.App().JWTSigninKey, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SessionClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

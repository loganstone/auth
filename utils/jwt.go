package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/loganstone/auth/configs"
)

// const .
const (
	Signup  = "Signup"
	Session = "Session"
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
		*newStandardClaims(Signup, email, "auth", t.expireAfterSec, 0),
	}
	return t.SignedString(configs.App().JWTSigninKey)
}

// Session .
func (t *Token) Session(userID uint, userEmail string) (string, error) {
	t.Claims = SessionClaims{
		SessionUser{UserID: userID, UserEmail: userEmail},
		*newStandardClaims(Session, userEmail, "auth", t.expireAfterSec, 0),
	}
	return t.SignedString(configs.App().JWTSigninKey)
}

func parseWithClaims(signedString string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		signedString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: '%v'", token.Header["alg"])
			}
			return configs.App().JWTSigninKey, nil
		})
}

// ParseJWTSignupToken .
func ParseJWTSignupToken(signedString string) (*SignupClaims, error) {
	token, err := parseWithClaims(signedString, &SignupClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignupClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token: '%s'", signedString)
	}

	return claims, nil
}

// ParseJWTSessionToken .
func ParseJWTSessionToken(signedString string) (*SessionClaims, error) {
	token, err := parseWithClaims(signedString, &SessionClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SessionClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token: '%s'", signedString)
	}

	return claims, nil
}

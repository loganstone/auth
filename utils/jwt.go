package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/loganstone/auth/configs"
)

type Token struct {
	jwt.Token
}

func NewJWTToken() *Token {
	claims := jwt.MapClaims{
		"aud": "audience",
		"exp": time.Now().Add(time.Minute * time.Duration(5)).Unix(),
		"jti": "reqeustID",
		"iat": time.Now().Unix(),
		"iss": "auth",
		"nbf": 0,
		"sub": "title",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	return &Token{
		*token,
	}
}

func (t *Token) Sign() (string, error) {
	return t.SignedString(configs.App().JWTSigninKey)
}

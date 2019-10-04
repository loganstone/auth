package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/models"
)

type Token struct {
	jwt.Token
}

func NewJWTToken(expireAfterSec int) *Token {
	now := time.Now()
	claims := jwt.MapClaims{
		"aud": "",
		"exp": now.Add(time.Second * time.Duration(expireAfterSec)).Unix(),
		"jti": uuid.New().String(),
		"iat": now.Unix(),
		"iss": "auth", // TODO(hs.lee): 설정 가능하게 변경.
		"nbf": 0,
		"sub": "",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return &Token{
		*token,
	}
}

func (t *Token) Signup(email string) (string, error) {
	claims := t.Claims.(jwt.MapClaims)
	claims["sub"] = "Signup"
	claims["aud"] = email
	claims["email"] = email
	return t.SignedString(configs.App().JWTSigninKey)
}

func (t *Token) Session(user *models.User) (string, error) {
	claims := t.Claims.(jwt.MapClaims)
	claims["sub"] = "Authorization"
	claims["aud"] = user.Email
	claims["user"] = user.ToMap()
	return t.SignedString(configs.App().JWTSigninKey)
}

func ParseJWTToken(signedString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return configs.App().JWTSigninKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

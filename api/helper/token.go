package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenHelper interface {
	Generate(userID uint, duration time.Duration) (string, error)
}

type tokenHelper struct {
	Secret string
}

func NewTokenHelper(secret string) TokenHelper {
	return &tokenHelper{
		Secret: secret,
	}
}

func (th *tokenHelper) Generate(userID uint, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(th.Secret))
}

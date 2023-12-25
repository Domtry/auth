package utils

import (
	"auth/model"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ExtractToken(token string) string {
	parts := strings.Split(token, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}

func ParseToken(tokenString string) (claims *model.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("my_secret_key"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}

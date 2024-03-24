package utils

import (
	"auth/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
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

func VerifyPermission(ctx echo.Context, value string) error {
	roleCtx := ctx.Get("roles")

	if roleCtx == nil {
		return fmt.Errorf("User has not permission")
	}

	roles := roleCtx.(JsonRoleItem)
	hasPermit := false

	for i := 0; i < len(roles.Permissions); i++ {
		permission := roles.Permissions[i]
		if permission == value {
			hasPermit = true
			break
		}
	}

	if !hasPermit {
		return fmt.Errorf("User has not permission")
	}

	return nil
}

func OtpGenerator(n int) string {
	var number = []byte("0123456789")
	b := make([]byte, n)
	for i := range b {
		b[i] = number[rand.Intn(len(number))]
	}

	return string(b)
}

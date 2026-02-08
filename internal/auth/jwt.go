package auth

import (
	"errors"
	"time"

	"github.com/AlikhanF2006/Final_project/configs"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

func GenerateToken(userID int) (string, error) {
	secret := configs.AppConfig.Auth.JWTSecret
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func ParseToken(tokenStr string) (int, error) {
	secret := configs.AppConfig.Auth.JWTSecret
	tkn, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil || !tkn.Valid {
		return 0, ErrInvalidToken
	}
	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}
	uidFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}
	return int(uidFloat), nil
}

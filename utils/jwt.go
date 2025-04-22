package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(userID uuid.UUID, name string, email string, isRefresh bool) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expMinute := os.Getenv("JWT_EXPIRE_MINUTES")

	if isRefresh {
		secret = os.Getenv("JWT_REFRESH_SECRET")
		expMinute = os.Getenv("JWT_REFRESH_EXPIRE_MINUTES")
	}

	duration, _ := time.ParseDuration(expMinute + "m")

	claims := jwt.MapClaims{
		"id": userID,
		"name":    name,
		"email":   email,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyToken(tokenStr string, isRefresh bool) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if isRefresh {
		secret = os.Getenv("JWT_REFRESH_SECRET")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GenerateToken(id string) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:        id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	userToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return userToken, nil
}

func GetClaims(c echo.Context) (string, error) {
	token := c.Get("user").(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token")
	}
	return claims["jti"].(string), nil
}

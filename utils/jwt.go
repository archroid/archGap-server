package utils

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"os"
)

var secretKey = os.Getenv("JWT_SECRET_KEY")

func CreateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseJWT(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, err
	}

	userID := uint(claims["user_id"].(float64))
	return userID, nil
}

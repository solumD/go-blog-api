package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken генерирует, подписывает и возвращает jwt-токен
func GenerateToken(login string, secret string) (string, error) {
	payload := jwt.MapClaims{
		"sub": login,
		"exp": time.Now().Add(time.Hour * 4).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GetTokenClaims читает jwt-токен и возвращает его payload
func GetTokenClaims(secret, tokenString string) (jwt.MapClaims, error) {
	signature := []byte(secret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signature, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%v", err)
	}

	return claims, nil
}

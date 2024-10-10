package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func ValidateToken(tokenString string, tokenType string) (*jwt.Token, jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if tokenType == "access_token" {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		}
		if tokenType == "refresh_token" {
			return []byte(os.Getenv("REFRESH_SECRET")), nil
		}

		return nil, fmt.Errorf("invalid token type: %s", tokenType)
	})

	if err != nil || !token.Valid {
		return nil, nil, fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return nil, nil, fmt.Errorf("token expired")
	}

	return token, claims, nil
}

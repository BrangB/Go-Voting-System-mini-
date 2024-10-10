package utils

import (
	"os"
	"time"

	"github.com/brangb/go_voting_system/models"
	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"User_ID": user.ID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // Access token expires in 15 minutes
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"User_ID": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // Refresh token expires in 30 days
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

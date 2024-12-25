package middlewares

import (
	"time"

	"github.com/NutsBalls/Nexus/utils"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(userID uint, username string, secretKey string) (string, error) {
	claims := &utils.JWTCustomClaims{
		ID:       userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

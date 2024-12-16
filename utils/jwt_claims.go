package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTCustomClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func CreateJWTToken(id uint, email string, secret []byte, expirationTime int64) (string, error) {
	claims := &JWTCustomClaims{
		ID:    id,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateJWTToken(tokenString string, secret []byte) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func IsTokenExpired(claims *JWTCustomClaims) bool {
	return time.Now().Unix() > claims.ExpiresAt
}

func RefreshJWTToken(oldToken *jwt.Token, secret []byte, newExpirationTime int64) (string, error) {
	claims := oldToken.Claims.(*JWTCustomClaims)
	claims.ExpiresAt = newExpirationTime

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString(secret)
}

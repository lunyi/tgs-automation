package main

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	jwt.StandardClaims
}

// generateToken generates a JWT token
func GenerateToken() (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

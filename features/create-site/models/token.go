package models

import "github.com/golang-jwt/jwt"

var JwtKey = []byte("secret_key")

type Claims struct {
	APIKey string `json:"api_key"`
	jwt.StandardClaims
}

package jwttoken

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var JwtKey = []byte("secret_key")

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
	tokenString, err := token.SignedString(JwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func TokenHandler(c *gin.Context) {
	tokenString, err := GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

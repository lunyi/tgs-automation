package main

import (
	"net/http"
	"os"
	"time"

	"tgs-automation/features/create-site/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func validateAPIKey(apiKey string) bool {
	return apiKey == os.Getenv("API_KEY")
}

func TokenHandler(c *gin.Context) {
	apiKey := c.GetHeader("x-api-key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API key is required"})
		return
	}

	if !validateAPIKey(apiKey) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	tokenString, err := generateToken(apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// generateToken generates a JWT token
func generateToken(apiKey string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &models.Claims{
		APIKey: apiKey,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(models.JwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

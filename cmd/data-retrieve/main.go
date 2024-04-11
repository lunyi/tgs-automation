package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Public endpoint
	router.GET("/token", tokenHandler)

	// Protected endpoint
	router.GET("/album", AuthMiddleware(), albumHandler)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

// tokenHandler creates and sends a new JWT token
func tokenHandler(c *gin.Context) {
	tokenString, err := GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// albumHandler sends a protected resource
func albumHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "secret album data"})
}

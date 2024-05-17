package main

import (
	"fmt"
	"net/http"

	"tgs-automation/features/create-site/sites"

	_ "tgs-automation/features/data-retrieve-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/iris-contrib/swagger/swaggerFiles"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()
	router.GET("/healthz", healthCheckHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/token", tokenHandler)
	router.POST("/site", AuthMiddleware(), sites.CreateSiteHandler)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
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

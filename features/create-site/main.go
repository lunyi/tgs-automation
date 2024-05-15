package main

import (
	"fmt"
	"net/http"

	"tgs-automation/features/create-site/sites"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"

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
	router.POST("/site", AuthMiddleware(), createLobbyHandler)

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

func createLobbyHandler(c *gin.Context) {
	config := util.GetConfig()
	request, err := parseRequestBody(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	httpCode, response := sites.GetLobbyInfo(request, config)

	if httpCode != http.StatusOK {
		c.JSON(httpCode, response["lobby"])
		return
	}

	lobbyInfo, err := sites.ConvertToLobbyInfo(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not convert lobby info", "details": err.Error()})
		return
	}

	httpCode, response = sites.RunKubectlApply(lobbyInfo, request)
	if httpCode != http.StatusOK {
		c.JSON(httpCode, response)
		return
	}

	c.JSON(httpCode, response)
}

func parseRequestBody(c *gin.Context) (sites.CreateSiteRequest, error) {
	var request sites.CreateSiteRequest
	if err := c.BindJSON(&request); err != nil {
		log.LogError(fmt.Sprintf("Error parsing request body: %v", err))
		return request, err
	}
	log.LogInfo(fmt.Sprintf("Request: %v", request))
	return request, nil
}

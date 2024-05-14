package main

import (
	"fmt"
	"net/http"

	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"

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
	router.POST("/site", AuthMiddleware(), getLobbyInfo)

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

func getLobbyInfo(c *gin.Context) {
	config := util.GetConfig()

	// Parse JSON body into struct
	var request CreateSiteRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Fetch Docker image
	dockerhubService := NewDockerImageService(config.Dockerhub)
	image, err := dockerhubService.FetchDockerImage(request.lobbyTemplate)

	if err != nil {
		log.LogError(fmt.Sprintf("Error fetching docker image %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get docker image", "details": err.Error()})
		return
	}
	fmt.Println("Image:", image)

	// Get brand ID
	brandId, err := postgresql.GetBrandId(config.CreateSiteDb, "MOPH")
	if err != nil {
		log.LogError(fmt.Sprintf("Error getting brand id %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get brand id", "details": err.Error()})
	}
	fmt.Println("Brand ID:", brandId)

	// Get brand token
	token, err := getBrandToken(brandId, "staging", config.ApiUrl.BrandCert)
	if err != nil {
		log.LogError(fmt.Sprintf("Error getting brand token %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get brand token", "details": err.Error()})
	}

	lobby := &LobbyInfo{
		BrandToken:  token,
		DockerImage: image}
	c.JSON(http.StatusOK, gin.H{"lobby": lobby})
}

type CreateSiteRequest struct {
	BrandCode     string `json:"brandCode"`
	lobbyTemplate string `json:"lobbyTemplate"`
	Domain        string `json:"domain"`
	NameSpace     string `json:"namespace"`
}

type LobbyInfo struct {
	BrandToken  string
	DockerImage string
}

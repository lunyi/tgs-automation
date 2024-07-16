package main

import (
	"fmt"
	"net/http"

	jwttoken "tgs-automation/internal/jwt_token"
	"tgs-automation/internal/log"
	middleware "tgs-automation/internal/middleware"
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
	router.GET("/token", jwttoken.TokenHandler)
	router.POST("/player_adjust", middleware.AuthMiddleware(), getPlayersAdjustAmount)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

// albumHandler sends a protected resource
func getPlayersAdjustAmount(c *gin.Context) {
	// Load configuration
	config := util.GetConfig()

	// Initialize the database connection or service interface
	app := postgresql.NewGetPlayersAdjustAmountInterface(config.Postgresql)
	defer app.Close()

	// Parse JSON body into struct
	var requestData GetDataRequest
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	log.LogInfo(fmt.Sprintf("Request data: %+v", requestData))

	// Call GetData method with parameters from the JSON body
	data, err := app.GetData(requestData.BrandCode, requestData.StartDate, requestData.EndDate, requestData.TransType)

	log.LogInfo(fmt.Sprintf("Data: %+v", data))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get data", "details": err.Error()})
		return
	}

	// Send the data as part of the JSON response
	c.JSON(http.StatusOK, gin.H{"data": data})
}

type GetDataRequest struct {
	BrandCode string `json:"brandCode"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	TransType int    `json:"transType"`
}

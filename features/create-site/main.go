package main

import (
	"fmt"
	"net/http"
	"tgs-automation/features/create-site/controllers/sites"
	"tgs-automation/features/create-site/middleware"
	"tgs-automation/internal/log"

	_ "tgs-automation/features/data-retrieve-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/iris-contrib/swagger/swaggerFiles"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1", "10.139.0.0/16"})
	router.GET("/healthz", healthCheckHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/token", TokenHandler)
	router.POST("/site", middleware.AuthMiddleware(), sites.CreateSiteHandler)

	err := server.ListenAndServe()
	if err != nil {
		log.LogFatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

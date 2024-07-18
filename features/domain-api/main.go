package main

import (
	"fmt"
	"net/http"
	_ "github.com/go-programming-tour-book/blog-service/docs"
	jwttoken "tgs-automation/internal/jwt_token"
	middleware "tgs-automation/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/iris-contrib/swagger/swaggerFiles"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()
	router.GET("/healthz", healthCheckHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/token", jwttoken.TokenHandler)
	router.GET("/checknameserver", middleware.AuthMiddleware(), GetNameServer)
	router.PUT("/changenameserver", middleware.AuthMiddleware(), ChangeNameServer)
	router.GET("/domainprice", middleware.AuthMiddleware(), CheckDomainPrice)
	router.POST("/createdomain", middleware.AuthMiddleware(), CreateDomain)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

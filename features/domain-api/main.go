package main

import (
	"fmt"
	"net/http"
	jwttoken "tgs-automation/internal/jwt_token"
	"tgs-automation/internal/log"
	middleware "tgs-automation/internal/middleware"
	"tgs-automation/pkg/telegram"

	"github.com/gin-gonic/gin"
	"github.com/iris-contrib/swagger/swaggerFiles"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()
	router.GET("/healthz", healthCheckHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/token", jwttoken.TokenHandler)
	router.GET("/checknameserver", middleware.AuthMiddleware(), getNameServer)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

func getNameServer(c *gin.Context) {
	var request GetNameServerRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	log.LogInfo(fmt.Sprintf("Request data: %+v", request))

	targetNameServer, err := GetTargetNameServers(request.Domain)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get target nameserver", "details": err.Error()})
		return
	}

	originNameServer, err := GetOriginalNameServer(request.Domain)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get original nameserver", "details": err.Error()})
		return
	}

	message := fmt.Sprintf("Domain: %s\nNameServers: %s\nOriginal Nameservers: %s", request.Domain, targetNameServer, originNameServer)

	err = telegram.SendMessageWithChatId(message, request.ChatId)
	if err != nil {
		fmt.Println("Failed to send Telegram message:", err)
	}

	c.JSON(http.StatusOK, gin.H{"data": message})
}

type GetNameServerRequest struct {
	Domain string `form:"domain" json:"domain" binding:"required"`
	ChatId string `form:"chatid" json:"chatid" binding:"required"`
}

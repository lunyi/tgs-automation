package main

import (
	"fmt"
	"net/http"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/namecheap"
	"tgs-automation/pkg/telegram"

	"github.com/gin-gonic/gin"
)

type GetDomainPriceRequest struct {
	Domain string `form:"domain" json:"domain" binding:"required"`
	ChatId string `form:"chatid" json:"chatid" binding:"required"`
}

func CheckDomainPrice(c *gin.Context) {
	var request GetDomainPriceRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	config := util.GetConfig()
	isAvailable, err := namecheap.CheckDomainAvailable(request.Domain, config.Namecheap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking domain availability", "details": err.Error()})
		return
	}

	if !isAvailable {
		message := fmt.Sprintf("[Get Domain Price]\ndomain: %s\ndomain 已經被使用，無法提供資訊.", request.Domain)
		c.JSON(http.StatusBadRequest, gin.H{"info": message})
		telegram.SendMessageWithChatId(message, request.ChatId)
		return
	}

	domainPriceResponse, err := namecheap.CheckDomainPrice(request.Domain, config.Namecheap)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !domainPriceResponse.CanRegister {
		message := fmt.Sprintf("[Get Domain Price]\ndomain: %s\ndomain 無法註冊.", request.Domain)
		c.JSON(http.StatusConflict, gin.H{"info": message})
		telegram.SendMessageWithChatId(message, request.ChatId)
		return
	}

	message := fmt.Sprintf("[Get Domain Price]\ndomain: %s\nRegular Price: %s\nPromotion Price: %s", request.Domain, domainPriceResponse.RegularPrice, domainPriceResponse.Price)
	c.JSON(http.StatusBadRequest, gin.H{"info": message})
	telegram.SendMessageWithChatId(message, request.ChatId)
}

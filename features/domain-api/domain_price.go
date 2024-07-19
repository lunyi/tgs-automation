package main

import (
	"fmt"
	"net/http"
	"reflect"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/namecheap"
	"tgs-automation/pkg/telegram"

	"github.com/gin-gonic/gin"
)

type GetDomainPriceRequest struct {
	Domain string `form:"domain" json:"domain" binding:"required"`
	ChatId string `form:"chatid" json:"chatid" binding:"required"`
}

func printAllFields(c *gin.Context) {
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Check if the field is exported
		if fieldType.PkgPath != "" {
			continue
		}

		fmt.Printf("%s: %v\n", fieldType.Name, field.Interface())
	}
}

// Get the price of the domain
// @Summary      Get the price of the domain
// @Tags         domain
// @Description  Get the price of the domain
// @Accept       json
// @Produce      json
// @Param        domain  query  string  true  "domain"
// @Param        chatid  query  string  true  "chatid"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Router       /domain/price [get]
func CheckDomainPrice(c *gin.Context) {
	printAllFields(c)
	var request GetDomainPriceRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	log.LogInfo(fmt.Sprintf("Request data: %+v", request))

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

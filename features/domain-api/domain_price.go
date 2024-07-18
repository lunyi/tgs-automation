package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/telegram"

	"github.com/gin-gonic/gin"
)

type GetDomainPriceRequest struct {
	Domain string `form:"domain" json:"domain" binding:"required"`
	ChatId string `form:"chatid" json:"chatid" binding:"required"`
}

type checkDomainPriceResponse struct {
	RegularPrice string
	Price        string
	CanRegister  bool
}

type CheckDomainApiResponse struct {
	CommandResponse struct {
		DomainCheckResult struct {
			Available string `xml:"Available,attr"`
		} `xml:"DomainCheckResult"`
		UserGetPricingResult struct {
			ProductType struct {
				ProductCategory []struct {
					Name    string `xml:"Name,attr"`
					Product struct {
						Price []struct {
							Duration     string `xml:"Duration,attr"`
							RegularPrice string `xml:"RegularPrice,attr"`
							Price        string `xml:"Price,attr"`
						} `xml:"Price"`
					} `xml:"Product"`
				} `xml:"ProductCategory"`
			} `xml:"ProductType"`
		} `xml:"UserGetPricingResult"`
	} `xml:"CommandResponse"`
}

var (
	Address      = "Sec4WenxinRdBeitunDist"
	NameCheapUrl = "https://api.namecheap.com/xml.response?"
)

func CheckDomainPrice(c *gin.Context) {
	var request GetDomainPriceRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	config := util.GetConfig()
	isAvailable, err := checkDomainAvailable(request.Domain, config.Namecheap)
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

	domainPriceResponse, err := checkDomainPrice(request.Domain, config.Namecheap)
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

func checkDomainPrice(domain string, config util.NamecheapConfig) (*checkDomainPriceResponse, error) {
	array := strings.Split(domain, ".")
	TLD := array[1]

	checkDomainPriceUrl := fmt.Sprintf("%s&ApiUser=%s&ApiKey=%s&UserName=%s&Command=namecheap.users.getPricing&ClientIp=%s&ProductCategory=register&ProductName=%s&ProductType=DOMAIN",
		NameCheapUrl, config.NamecheapUsername, config.NamecheapApiKey, config.NamecheapUsername, config.NamecheapClientIp, TLD)

	fmt.Println("checkDomainPriceUrl=", checkDomainPriceUrl)

	// Get domain price
	response, err := http.Get(checkDomainPriceUrl)
	if err != nil {
		return nil, fmt.Errorf("Error getting domain price: %s", err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %s", err)

	}
	var apiResponse CheckDomainApiResponse
	err = xml.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling XML: %s", err)

	}

	count := len(apiResponse.CommandResponse.UserGetPricingResult.ProductType.ProductCategory)
	fmt.Println("Count:", count)

	if count != 0 {
		regularPrice := apiResponse.CommandResponse.UserGetPricingResult.ProductType.ProductCategory[0].Product.Price[0].RegularPrice
		price := apiResponse.CommandResponse.UserGetPricingResult.ProductType.ProductCategory[0].Product.Price[0].Price
		return &checkDomainPriceResponse{
			CanRegister:  true,
			RegularPrice: regularPrice,
			Price:        price,
		}, nil
	}

	return &checkDomainPriceResponse{
		CanRegister: false,
	}, nil
}

func checkDomainAvailable(domain string, config util.NamecheapConfig) (bool, error) {
	checkDomainAvailableUrl := fmt.Sprintf("%s&ApiUser=%s&ApiKey=%s&UserName=%s&Command=namecheap.domains.check&ClientIp=%s&DomainList=%s",
		NameCheapUrl, config.NamecheapUsername, config.NamecheapApiKey, config.NamecheapUsername, config.NamecheapClientIp, domain)

	// Check if domain is available
	responseDomainCheck, err := http.Get(checkDomainAvailableUrl)
	if err != nil {
		return false, fmt.Errorf("Error checking domain availability: %s", err)
	}
	defer responseDomainCheck.Body.Close()

	bodyDomainCheck, err := ioutil.ReadAll(responseDomainCheck.Body)
	if err != nil {
		return false, fmt.Errorf("Error reading response body: %s", err)
	}

	var apiResponse CheckDomainApiResponse
	err = xml.Unmarshal(bodyDomainCheck, &apiResponse)
	if err != nil {
		return false, fmt.Errorf("Error unmarshalling XML: %s", err)
	}

	available := apiResponse.CommandResponse.DomainCheckResult.Available
	fmt.Println("Available:", available)

	if available == "false" {
		return false, nil
	}
	return true, nil
}

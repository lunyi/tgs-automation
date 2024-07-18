package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/telegram"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"gopkg.in/xmlpath.v2"
)

const (
	Email        = "ly.lester@rayprosoft.com"
	NameServers  = "micah.ns.cloudflare.com,ulla.ns.cloudflare.com"
	Address      = "Sec4WenxinRdBeitunDist"
	NameCheapUrl = "https://sandbox.namecheap.com/xml.response?"
)

type CreateDomainRequest struct {
	Domain string
	ChatId string
}

// CreateDomain
// @Summary create new domain on namecheap
// @Tags domain
// @Description create new domain on namecheap
// @Accept json
// @Produce json
// @Param domain query string true "domain"
// @Param chatid query string true "chatid"
// @Success 200 {object} map[string]interface{} "success"
// @Failure 400 {object} map[string]interface{} "error"
// @Router /domain [post]
func CreateDomain(c *gin.Context) {
	var request CreateDomainRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	config := util.GetConfig()

	promotionCode := getCouponCode()
	domainResponse, err := createDomain(request.Domain, promotionCode, config.Namecheap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balanceResponse, err := getBalance(config.Namecheap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status, _ := parseXml(domainResponse, "//ApiResponse/@Status")
	errorMsg, _ := parseXml(domainResponse, "//ApiResponse/Errors/Error/text()")
	chargeAmount, _ := parseXml(domainResponse, "//ApiResponse/CommandResponse/DomainCreateResult/@ChargedAmount")
	currency, _ := parseXml(balanceResponse, "//ApiResponse/CommandResponse/UserGetBalancesResult/@Currency")
	balance, _ := parseXml(balanceResponse, "//ApiResponse/CommandResponse/UserGetBalancesResult/@AvailableBalance")

	var message string
	if status == "OK" {
		message = fmt.Sprintf("建立Domain: %s 成功\n優惠碼: %s\n費用: %s\n餘額: %s %s", request.Domain, promotionCode, chargeAmount, balance, currency)
	} else {
		message = fmt.Sprintf("建立Domain: %s 失敗\n優惠碼: %s\n原因: %s\n餘額: %s %s", request.Domain, promotionCode, errorMsg, balance, currency)
	}

	telegram.SendMessageWithChatId(message, request.ChatId)

	c.JSON(http.StatusOK, gin.H{"status": status, "message": message})
}

func getCouponCode() string {
	c := colly.NewCollector()

	var couponCode string
	c.OnHTML("button", func(e *colly.HTMLElement) {
		if e.Attr("class") == "button" {
			couponCode = e.Text
		}
	})

	c.Visit("https://www.namecheap.com/promos/coupons/")
	return couponCode
}

func createDomain(domainName, promotionCode string, config util.NamecheapConfig) (string, error) {
	url := fmt.Sprintf("%s&ApiUser=%s&ApiKey=%s&UserName=%s&DomainName=%s&Command=namecheap.domains.create&ClientIp=%s&Years=1&AuxBillingFirstName=Mark&AuxBillingLastName=Wu&AuxBillingAddress1=%s&AuxBillingStateProvince=TW&AuxBillingPostalCode=406&AuxBillingCountry=TW&AuxBillingPhone=+886.6613102107&AuxBillingEmailAddress=%s&AuxBillingOrganizationName=Raypro&AuxBillingCity=TC&TechFirstName=Mark&TechLastName=Wu&TechAddress1=%s&TechStateProvince=TW&TechPostalCode=90045&TechCountry=TW&TechPhone=+886.6613102107&TechEmailAddress=%s&TechOrganizationName=Raypro&TechCity=TW&AdminFirstName=Mark&AdminLastName=Wu&AdminAddress1=%s&AdminStateProvince=CA&AdminPostalCode=9004&AdminCountry=US&AdminPhone=+886.6613102107&AdminEmailAddress=%s&AdminOrganizationName=Raypro&AdminCity=CA&RegistrantFirstName=Mark&RegistrantLastName=Wu&RegistrantAddress1=%s&RegistrantStateProvince=TW&RegistrantPostalCode=406&RegistrantCountry=TW&RegistrantPhone=+886.6613102107&RegistrantEmailAddress=%s&RegistrantOrganizationName=Raypro&RegistrantCity=TW&Nameservers=%s&PromotionCode=%s",
		NameCheapUrl, config.NamecheapUsername, config.NamecheapApiKey, config.NamecheapUsername, domainName, config.NamecheapClientIp, Address, Email, Address, Email, Address, Email, Address, Email, NameServers, promotionCode)

	resp, err := http.Post(url, "application/xml", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getBalance(config util.NamecheapConfig) (string, error) {
	url := fmt.Sprintf("%s&ApiUser=%s&ApiKey=%s&UserName=%s&Command=namecheap.users.getBalances&ClientIp=%s",
		NameCheapUrl,
		config.NamecheapUsername,
		config.NamecheapApiKey,
		config.NamecheapUsername,
		config.NamecheapClientIp)

	resp, err := http.Post(url, "application/xml", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseXml(xmlStr string, xpathExpr string) (string, error) {
	root, err := xmlpath.Parse(strings.NewReader(xmlStr))
	if err != nil {
		return "", err
	}

	path := xmlpath.MustCompile(xpathExpr)
	if value, ok := path.String(root); ok {
		return value, nil
	}

	return "", fmt.Errorf("no match found")
}

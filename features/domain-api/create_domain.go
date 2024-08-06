package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/namecheap"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/xmlpath.v2"
)

type CreateDomainRequest struct {
	Domain string `form:"domain" json:"domain" binding:"required"`
	ChatId string `form:"chatid" json:"chatid" binding:"required"`
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
func CreateDomainHandler(
	api namecheap.NamecheapApi,
	natsSvc util.NatsPublisherService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateDomainRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		promotionCode, err := api.GetCouponCode(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"get coupon code error": err.Error()})
			return
		}

		domainResponse, err := api.CreateDomain(ctx, request.Domain, promotionCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"create domain error": err.Error()})
			return
		}

		balanceResponse, err := api.GetBalance(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"get balance error": err.Error()})
			return
		}

		status, _ := parseXml(domainResponse, "//ApiResponse/@Status")
		errorMsg, _ := parseXml(domainResponse, "//ApiResponse/Errors/Error/text()")
		chargeAmount, _ := parseXml(domainResponse, "//ApiResponse/CommandResponse/DomainCreateResult/@ChargedAmount")
		currency, _ := parseXml(balanceResponse, "//ApiResponse/CommandResponse/UserGetBalancesResult/@Currency")
		balance, _ := parseXml(balanceResponse, "//ApiResponse/CommandResponse/UserGetBalancesResult/@AvailableBalance")

		message := formatDomainMessage(status, request.Domain, promotionCode, chargeAmount, errorMsg, balance, currency)
		natsSvc.Publish(request.ChatId, message)

		c.JSON(http.StatusOK, gin.H{"status": status, "message": message})
	}
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

func formatDomainMessage(status, domain, promotionCode, chargeAmount, errorMsg, balance, currency string) string {
	if status == "OK" {
		return fmt.Sprintf("建立Domain: %s 成功\n優惠碼: %s\n費用: %s\n餘額: %s %s", domain, promotionCode, chargeAmount, balance, currency)
	}
	return fmt.Sprintf("建立Domain: %s 失敗\n優惠碼: %s\n原因: %s\n餘額: %s %s", domain, promotionCode, errorMsg, balance, currency)
}

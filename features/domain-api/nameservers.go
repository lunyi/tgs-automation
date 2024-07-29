package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/cloudflare"
	"tgs-automation/pkg/telegram"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type GetNameServerRequest struct {
	Domain string `form:"domain" json:"domain" binding:"required"`
	ChatId string `form:"chatid" json:"chatid" binding:"required"`
}

type UpdateNameServerRequest struct {
	Domain      string `form:"domain" json:"domain" binding:"required"`
	ChatId      string `form:"chatid" json:"chatid" binding:"required"`
	NameServers string `form:"nameservers" json:"nameservers" binding:"required"`
}

type ApiResponse struct {
	Status string `xml:"Status,attr"`
	Errors struct {
		Error string `xml:"Error"`
	} `xml:"Errors"`
}

// UpdateNameServer changes the name server information
// @Summary Update name server information
// @Tags NameServer
// @Description Update the name server information for a given domain
// @Accept json
// @Produce json
// @Param updateNameServerRequest body UpdateNameServerRequest true "Update name server request"
// @Success 200 {object} ApiResponse "Success"
// @Failure 400 {object} ApiResponse "Bad Request"
// @Security ApiKeyAuth
// @Router /nameservers [put]
func UpdateNameServer(c *gin.Context) {
	tracer := otel.Tracer("domain-api")
	ctx, span := tracer.Start(c.Request.Context(), "UpdateNameServer")
	defer func() {
		span.End()
		ctx.Done()
	}()

	printAllFields(c)
	var request UpdateNameServerRequest
	if err := c.BindJSON(&request); err != nil {
		span.RecordError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update nameserver request data", "details": err.Error()})
		return
	}

	config := util.GetConfig()
	apiUser := config.Namecheap.NamecheapUsername
	userName := config.Namecheap.NamecheapUsername
	apiKey := config.Namecheap.NamecheapApiKey
	clientIp := config.Namecheap.NamecheapClientIp
	nameCheapUrl := "https://api.namecheap.com/xml.response?"

	domainParts := strings.Split(request.Domain, ".")
	if len(domainParts) != 2 {
		fmt.Println("Invalid domain name format")
		span.RecordError(fmt.Errorf("Invalid domain name format"))
		return
	}
	sld := domainParts[0]
	tld := domainParts[1]

	urlParams := url.Values{}
	urlParams.Set("ApiUser", apiUser)
	urlParams.Set("ApiKey", apiKey)
	urlParams.Set("UserName", userName)
	urlParams.Set("DomainName", request.Domain)
	urlParams.Set("Command", "namecheap.domains.dns.setCustom")
	urlParams.Set("ClientIp", clientIp)
	urlParams.Set("SLD", sld)
	urlParams.Set("TLD", tld)
	urlParams.Set("Nameservers", request.NameServers)

	apiUrl := nameCheapUrl + urlParams.Encode()
	fmt.Println("Url=", apiUrl)

	span.AddEvent("Generated promotion code")

	resp, err := http.Get(apiUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error making API request:", "details": err.Error()})
		span.RecordError(fmt.Errorf("Error making API request: %v", err))
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response", "details": err.Error()})
		span.RecordError(fmt.Errorf("Error reading response: %v", err))
		return
	}

	var apiResponse ApiResponse
	if err := xml.Unmarshal(body, &apiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing XML:", "details": err.Error()})
		span.RecordError(fmt.Errorf("Error parsing XML: %v", err))
		return
	}

	var message string
	if strings.Contains(apiResponse.Status, "OK") {
		message = fmt.Sprintf("修改nameserver: %s 成功\nNameServers: %s", request.Domain, request.NameServers)
		span.AddEvent("Update nameserver successfully")
	} else {
		message = fmt.Sprintf("修改nameserver: %s 失敗\n原因: %s\nNameServers: %s", request.Domain, apiResponse.Errors.Error, request.NameServers)
		span.AddEvent("Update nameserver failed")
	}

	span.AddEvent("start to send message to telegram", trace.WithTimestamp(time.Now()))
	err = telegram.SendMessageWithChatId(message, request.ChatId)
	if err != nil {
		fmt.Println("Failed to send Telegram message:", err)
		recordError(span, "Failed to send Telegram message", err)
	}
	span.AddEvent("end to send message to telegram", trace.WithTimestamp(time.Now()))
	span.SetStatus(codes.Ok, "")
	c.JSON(http.StatusOK, gin.H{"data": message})
}

// GetNameServer retrieves the name server information
// @Summary Retrieve name server information
// @Tags NameServer
// @Description Retrieve the name server information for a given domain
// @Accept json
// @Produce json
// @Param domain query string true "Domain name"
// @Param chatid query string true "Chat ID"
// @Success 200 {object} ApiResponse "Success"
// @Failure 400 {object} ApiResponse "Bad Request"
// @Security ApiKeyAuth
// @Router /nameservers [get]
func GetNameServer(c *gin.Context) {
	tracer := otel.Tracer("domain-api")
	ctx, span := tracer.Start(c.Request.Context(), "GetNameServer")
	defer func() {
		span.End()
		ctx.Done()
	}()

	var request GetNameServerRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		recordError(span, "invalid request data", err)
		return
	}

	log.LogInfo(fmt.Sprintf("Request data: %+v", request))
	targetNameServer, err := cloudflare.GetTargetNameServers(request.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get target nameserver", "details": err.Error()})
		recordError(span, "could not get target nameserver", err)
		return
	}

	originNameServer, err := getOriginalNameServer(request.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get original nameserver", "details": err.Error()})
		recordError(span, "could not get original nameserver", err)
		return
	}

	message := fmt.Sprintf("Domain: %s\nNameServers: %s\nOriginal Nameservers: %s", request.Domain, targetNameServer, originNameServer)

	config := util.GetConfig()
	publish(request.ChatId, message, config.NatsUrl)

	// subCtx, subSpan := tracer.Start(ctx, "Send message to telegram")
	// err = telegram.SendMessageWithChatIdAndContext(subCtx, message, request.ChatId)
	// if err != nil {
	// 	fmt.Println("Failed to send Telegram message:", err)
	// }
	// subSpan.End()

	c.JSON(http.StatusOK, gin.H{"data": message})
}

type MessageData struct {
	ChatId  string `json:"chatid"`
	Message string `json:"message"`
}

func publish(chatId string, message string, natsUrl string) {
	msgContent := MessageData{
		ChatId:  chatId,
		Message: message,
	}

	msgBytes, err := json.Marshal(msgContent)
	if err != nil {
		log.LogError(fmt.Sprintln("error marshalling message: %v", err))
	}

	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.LogError(fmt.Sprintln("error connecting to NATS server: %v", err))
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.LogError(fmt.Sprintln("error creating JetStream context: %v", err))
	}

	_, err = js.Publish("telegram", msgBytes)
	if err != nil {
		log.LogError(fmt.Sprintln("error publishing message: %v", err))
	}
	fmt.Println("Message published: %w", msgContent)
}

func recordError(span trace.Span, message string, err error) {
	span.RecordError(fmt.Errorf("%v: %w", message, err))
	span.SetStatus(codes.Error, message)
}

func getOriginalNameServer(domain string) (string, error) {
	var nsRecords []string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeNS)
	m.RecursionDesired = true

	c := new(dns.Client)
	in, _, err := c.Exchange(m, "8.8.8.8:53") // Using Google's public DNS server
	if err != nil {
		return "", err
	}

	if len(in.Answer) == 0 {
		return "", fmt.Errorf("no NS records found for domain %s", domain)
	}

	for _, ans := range in.Answer {
		if ns, ok := ans.(*dns.NS); ok {
			nsRecords = append(nsRecords, ns.Ns)
		}
	}

	return strings.Join(nsRecords, " "), nil
}

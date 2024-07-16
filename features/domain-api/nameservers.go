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
	"tgs-automation/pkg/telegram"

	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
)

type GetNameServerRequest struct {
	Domain string `form:"domain" json:"domain" binding:"required"`
	ChatId string `form:"chatid" json:"chatid" binding:"required"`
}

type ChangeNameServerRequest struct {
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

func ChangeNameServer(c *gin.Context) {
	var request ChangeNameServerRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
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

	resp, err := http.Get(apiUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error making API request:", "details": err.Error()})
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response", "details": err.Error()})
		return
	}

	var apiResponse ApiResponse
	if err := xml.Unmarshal(body, &apiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing XML:", "details": err.Error()})
		return
	}

	var message string
	if strings.Contains(apiResponse.Status, "OK") {
		message = fmt.Sprintf("修改nameserver: %s 成功\nNameServers: %s", request.Domain, request.NameServers)
	} else {
		message = fmt.Sprintf("修改nameserver: %s 失敗\n原因: %s\nNameServers: %s", request.Domain, apiResponse.Errors.Error, request.NameServers)
	}
	c.JSON(http.StatusOK, gin.H{"data": message})
}

func GetNameServer(c *gin.Context) {
	var request GetNameServerRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	log.LogInfo(fmt.Sprintf("Request data: %+v", request))

	targetNameServer, err := getTargetNameServers(request.Domain)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get target nameserver", "details": err.Error()})
		return
	}

	originNameServer, err := getOriginalNameServer(request.Domain)

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

func getTargetNameServers(domain string) (string, error) {
	config := util.GetConfig()
	// Get Cloudflare name servers
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?per_page=1000&name=%s", domain)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.CloudflareToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching Cloudflare data:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	nameServers := ""
	if res, ok := result["result"].([]interface{}); ok && len(res) > 0 {
		if ns, ok := res[0].(map[string]interface{})["name_servers"].([]interface{}); ok {
			var nsList []string
			for _, v := range ns {
				nsList = append(nsList, v.(string))
			}
			nameServers = strings.Join(nsList, " ")
		}
	}

	return nameServers, nil
}

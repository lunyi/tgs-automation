package namecheap

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"tgs-automation/internal/util"
)

type UpdateNameServerApiResponse struct {
	Status string `xml:"Status,attr"`
	Errors struct {
		Error string `xml:"Error"`
	} `xml:"Errors"`
}

func (n *NamecheapService) UpdateNameServer(domain string, nameServers string) (*UpdateNameServerApiResponse, error) {

	apiUrl, err := getUrl(domain, nameServers, n.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid domain name format: %v", err.Error())
	}

	fmt.Println("Url=", apiUrl)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %v", err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err.Error())
	}

	var apiResponse UpdateNameServerApiResponse
	if err := xml.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("eror parsing XML: %v", err.Error())
	}
	return &apiResponse, nil
}

func getUrl(domain string, nameServers string, config util.TgsConfig) (string, error) {
	domainParts := strings.Split(domain, ".")
	if len(domainParts) != 2 {
		return "", fmt.Errorf("invalid domain name format")
	}
	sld := domainParts[0]
	tld := domainParts[1]

	apiUser := config.Namecheap.NamecheapUsername
	userName := config.Namecheap.NamecheapUsername
	apiKey := config.Namecheap.NamecheapApiKey
	clientIp := config.Namecheap.NamecheapClientIp
	nameCheapUrl := config.Namecheap.NamecheapBaseUrl

	urlParams := url.Values{}
	urlParams.Set("ApiUser", apiUser)
	urlParams.Set("ApiKey", apiKey)
	urlParams.Set("UserName", userName)
	urlParams.Set("DomainName", domain)
	urlParams.Set("Command", "namecheap.domains.dns.setCustom")
	urlParams.Set("ClientIp", clientIp)
	urlParams.Set("SLD", sld)
	urlParams.Set("TLD", tld)
	urlParams.Set("Nameservers", nameServers)

	return nameCheapUrl + urlParams.Encode(), nil
}

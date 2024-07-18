package namecheap

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"tgs-automation/internal/util"
)

const (
	Email        = "ly.lester@rayprosoft.com"
	NameServers  = "micah.ns.cloudflare.com,ulla.ns.cloudflare.com"
	Address      = "Sec4WenxinRdBeitunDist"
	NameCheapUrl = "https://sandbox.namecheap.com/xml.response?"
)

func CheckDomainAvailable(domain string, config util.NamecheapConfig) (bool, error) {
	checkDomainAvailableUrl := fmt.Sprintf("%s&ApiUser=%s&ApiKey=%s&UserName=%s&Command=namecheap.domains.check&ClientIp=%s&DomainList=%s",
		NameCheapUrl, config.NamecheapUsername, config.NamecheapApiKey, config.NamecheapUsername, config.NamecheapClientIp, domain)

	// Check if domain is available
	responseDomainCheck, err := http.Get(checkDomainAvailableUrl)
	if err != nil {
		return false, fmt.Errorf("error checking domain availability: %s", err)
	}
	defer responseDomainCheck.Body.Close()

	bodyDomainCheck, err := ioutil.ReadAll(responseDomainCheck.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %s", err)
	}

	var apiResponse CheckDomainApiResponse
	err = xml.Unmarshal(bodyDomainCheck, &apiResponse)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling XML: %s", err)
	}

	available := apiResponse.CommandResponse.DomainCheckResult.Available
	fmt.Println("Available:", available)

	if available == "false" {
		return false, nil
	}
	return true, nil
}

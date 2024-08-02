package namecheap

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (n *NamecheapService) CheckDomainAvailable(ctx context.Context, domain string) (bool, error) {
	checkDomainAvailableUrl := fmt.Sprintf("%s&ApiUser=%s&ApiKey=%s&UserName=%s&Command=namecheap.domains.check&ClientIp=%s&DomainList=%s",
		n.Config.Namecheap.NamecheapBaseUrl, n.Config.Namecheap.NamecheapUsername, n.Config.Namecheap.NamecheapApiKey, n.Config.Namecheap.NamecheapUsername, n.Config.Namecheap.NamecheapClientIp, domain)

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

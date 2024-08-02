package namecheap

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (n *NamecheapService) CreateDomain(ctx context.Context, domainName string, promotionCode string) (string, error) {
	email := n.Config.Namecheap.NamecheapEmail
	address := n.Config.Namecheap.NamecheapAddress
	username := n.Config.Namecheap.NamecheapUsername
	url := fmt.Sprintf("%s&ApiUser=%s&ApiKey=%s&UserName=%s&DomainName=%s&Command=namecheap.domains.create&ClientIp=%s&Years=1&AuxBillingFirstName=Mark&AuxBillingLastName=Wu&AuxBillingAddress1=%s&AuxBillingStateProvince=TW&AuxBillingPostalCode=406&AuxBillingCountry=TW&AuxBillingPhone=+886.6613102107&AuxBillingEmailAddress=%s&AuxBillingOrganizationName=Raypro&AuxBillingCity=TC&TechFirstName=Mark&TechLastName=Wu&TechAddress1=%s&TechStateProvince=TW&TechPostalCode=90045&TechCountry=TW&TechPhone=+886.6613102107&TechEmailAddress=%s&TechOrganizationName=Raypro&TechCity=TW&AdminFirstName=Mark&AdminLastName=Wu&AdminAddress1=%s&AdminStateProvince=CA&AdminPostalCode=9004&AdminCountry=US&AdminPhone=+886.6613102107&AdminEmailAddress=%s&AdminOrganizationName=Raypro&AdminCity=CA&RegistrantFirstName=Mark&RegistrantLastName=Wu&RegistrantAddress1=%s&RegistrantStateProvince=TW&RegistrantPostalCode=406&RegistrantCountry=TW&RegistrantPhone=+886.6613102107&RegistrantEmailAddress=%s&RegistrantOrganizationName=Raypro&RegistrantCity=TW&Nameservers=%s&PromotionCode=%s",
		n.Config.Namecheap.NamecheapBaseUrl, username, n.Config.Namecheap.NamecheapApiKey, username, domainName, n.Config.Namecheap.NamecheapClientIp,
		address, email, address, email, address, email, address, email, n.Config.Namecheap.NamecheapNameServers, promotionCode)

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

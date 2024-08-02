package namecheap

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type CheckDomainPriceResponse struct {
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

func (n *NamecheapService) GetDomainPrice(ctx context.Context, domain string) (*CheckDomainPriceResponse, error) {
	array := strings.Split(domain, ".")
	TLD := array[1]

	checkDomainPriceUrl := fmt.Sprintf("%s&ApiUser=%s&ApiKey=%s&UserName=%s&Command=namecheap.users.getPricing&ClientIp=%s&ProductCategory=register&ProductName=%s&ProductType=DOMAIN",
		n.Config.Namecheap.NamecheapBaseUrl, n.Config.Namecheap.NamecheapUsername, n.Config.Namecheap.NamecheapApiKey, n.Config.Namecheap.NamecheapUsername, n.Config.Namecheap.NamecheapClientIp, TLD)

	fmt.Println("checkDomainPriceUrl=", checkDomainPriceUrl)

	// Get domain price
	response, err := http.Get(checkDomainPriceUrl)
	if err != nil {
		return nil, fmt.Errorf("error getting domain price: %s", err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)

	}
	var apiResponse CheckDomainApiResponse
	err = xml.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %s", err)

	}

	count := len(apiResponse.CommandResponse.UserGetPricingResult.ProductType.ProductCategory)
	fmt.Println("Count:", count)

	if count != 0 {
		regularPrice := apiResponse.CommandResponse.UserGetPricingResult.ProductType.ProductCategory[0].Product.Price[0].RegularPrice
		price := apiResponse.CommandResponse.UserGetPricingResult.ProductType.ProductCategory[0].Product.Price[0].Price
		return &CheckDomainPriceResponse{
			CanRegister:  true,
			RegularPrice: regularPrice,
			Price:        price,
		}, nil
	}

	return &CheckDomainPriceResponse{
		CanRegister: false,
	}, nil
}

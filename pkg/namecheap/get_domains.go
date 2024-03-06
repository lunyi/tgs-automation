package namecheap

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	apiUrl     = "https://api.namecheap.com/xml.response"
	command    = "namecheap.domains.getList"
	pageSize   = 100
	maxPages   = 5
	dateFormat = "01/02/2006"
)

type DomainGetListResult struct {
	Domains []Domain `xml:"Domain"`
}

type Paging struct {
	TotalItems  int `xml:"TotalItems"`
	CurrentPage int `xml:"CurrentPage"`
	PageSize    int `xml:"PageSize"`
}

type CommandResponse struct {
	DomainGetListResult DomainGetListResult `xml:"DomainGetListResult"`
	Paging              Paging              `xml:"Paging"`
}

type ApiResponse struct {
	XMLName          xml.Name        `xml:"ApiResponse"`
	Status           string          `xml:"Status,attr"`
	RequestedCommand string          `xml:"RequestedCommand"`
	CommandResponse  CommandResponse `xml:"CommandResponse"`
}

type Domain struct {
	ID         string `xml:"ID,attr"`
	Name       string `xml:"Name,attr"`
	User       string `xml:"User,attr"`
	Created    string `xml:"Created,attr"`
	Expires    string `xml:"Expires,attr"`
	IsExpired  string `xml:"IsExpired,attr"`
	IsLocked   string `xml:"IsLocked,attr"`
	AutoRenew  string `xml:"AutoRenew,attr"`
	WhoisGuard string `xml:"WhoisGuard,attr"`
	IsPremium  string `xml:"IsPremium,attr"`
	IsOurDNS   string `xml:"IsOurDNS,attr"`
}

type FilteredDomain struct {
	Name    string `xml:"Name,attr"`
	Created string `xml:"Created,attr"`
	Expires string `xml:"Expires,attr"`
}

type NamecheapAPI interface {
	GetExpiredDomains() ([]FilteredDomain, error)
}

type NamecheapClient struct {
	Config util.NamecheapConfig
}

func New(config util.NamecheapConfig) NamecheapAPI {
	return &NamecheapClient{
		Config: config,
	}
}

func (nc *NamecheapClient) GetExpiredDomains() ([]FilteredDomain, error) {
	userName := nc.Config.NamecheapUsername
	apiKey := nc.Config.NamecheapApiKey
	clinetIp := nc.Config.NamecheapClientIp

	url := fmt.Sprintf("%s?ApiUser=%s&ApiKey=%s&UserName=%s&Command=%s&ClientIp=%s&PageSize=%d&Page=", apiUrl, userName, apiKey, userName, command, clinetIp, pageSize)

	var domains []Domain
	index := 1
	for index < maxPages {
		fmt.Println(index)
		_domains, totolCount, shouldReturn, err := getDomains(url, index)

		if err != nil {
			return nil, err
		}
		if !shouldReturn {
			domains = append(domains, _domains...)
		}
		index++
		fmt.Println("totolCount:", totolCount)
		fmt.Println("index_pageSize:", index*pageSize)
		if totolCount > index*pageSize {
			break
		}
	}

	return filterDomainsWithExpired(domains), nil
}

func filterDomainsWithExpired(domains []Domain) []FilteredDomain {
	currentDate := time.Now()
	next30Days := currentDate.AddDate(0, 0, 30)
	var result []FilteredDomain
	for _, domain := range domains {
		expireDate, err := time.Parse(dateFormat, domain.Expires)
		if err != nil {
			fmt.Println("Error parsing expiration date:", err)
			continue
		}

		if expireDate.After(currentDate) && expireDate.Before(next30Days) {
			result = append(result, FilteredDomain{
				Name:    domain.Name,
				Created: domain.Created,
				Expires: domain.Expires,
			})
		}
	}

	for i, domain := range result {
		fmt.Printf("%d, Name: %s, Created: %s, Expires: %s\n", i+1, domain.Name, domain.Created, domain.Expires)
	}
	return result
}

func getDomains(url string, index int) ([]Domain, int, bool, error) {
	url = url + strconv.Itoa(index)
	fmt.Println("url:", url)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching domain list:", err)
		return []Domain{}, -1, true, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return []Domain{}, -1, true, err
	}

	log.LogInfo(string(body))
	var apiResponse ApiResponse
	if err := xml.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return []Domain{}, -1, true, err
	}
	return apiResponse.CommandResponse.DomainGetListResult.Domains, apiResponse.CommandResponse.Paging.TotalItems, false, nil
}

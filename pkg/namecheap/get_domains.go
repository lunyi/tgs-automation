package namecheap

import (
	"cdnetwork/internal/httpclient"
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"encoding/xml"
	"fmt"
	"io/ioutil"
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
	Config     util.NamecheapConfig
	HTTPClient httpclient.HttpClient
}

func New(config util.NamecheapConfig, httpClient httpclient.HttpClient) NamecheapAPI {
	return &NamecheapClient{
		Config:     config,
		HTTPClient: httpClient,
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
		_domains, totolCount, shouldReturn, err := getDomains(nc.HTTPClient, url, index)

		if err != nil {
			return filterDomainsWithExpired(domains), nil
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
	nextMonth := currentDate.AddDate(0, 1, 0)
	firstDayOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())
	lastDayOfNextMonth := firstDayOfNextMonth.AddDate(0, 1, 0)

	var result []FilteredDomain
	for _, domain := range domains {
		expireDate, err := time.Parse(dateFormat, domain.Expires)
		if err != nil {
			fmt.Println("Error parsing expiration date2:", err)
			continue
		}

		if expireDate.After(firstDayOfNextMonth) && expireDate.Before(lastDayOfNextMonth) {
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

func getDomains(httpclient httpclient.HttpClient, url string, index int) ([]Domain, int, bool, error) {
	url = url + strconv.Itoa(index)
	fmt.Println("url:", url)

	response, err := httpclient.Get(url)
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

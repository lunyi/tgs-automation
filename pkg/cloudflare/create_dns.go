package cloudflare

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"tgs-automation/internal/log"
)

func (svc *CloudflareService) CreateDNS(domain string) error {
	apiToken := svc.Config.CloudflareToken
	zoneName := extractDomain(domain)
	zoneID, err := getZoneId(apiToken, domain)
	if err != nil {
		return err
	}

	record := DNSRecord{
		Type:    "CNAME",
		Name:    zoneName,                         // Name of the record
		Content: svc.Config.CdnNetwork.DnsContent, // IP address or content of the record
		TTL:     1,
		Proxied: false, // Whether the record is proxied through Cloudflare
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	// API endpoint to create DNS record
	endpoint := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)

	resp, err := sendRequest(apiToken, "POST", endpoint, strings.NewReader(string(recordJSON)))
	if err != nil {
		return err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		log.LogError(fmt.Sprintf("Error: %s", resp.Status))
		return errors.New("http status is not 200")
	}

	fmt.Println("DNS record created successfully!")
	return nil
}

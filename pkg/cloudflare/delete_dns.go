package cloudflare

import (
	"cdnetwork/internal/util"
	"errors"
	"fmt"
	"net/http"
)

func DeleteDNS(domain string) error {
	config := util.GetConfig()

	apiToken := config.CloudflareToken
	zoneName := extractDomain(domain)
	zoneID, err := GetZoneId(zoneName)
	if err != nil {
		return err
	}

	// DNS record ID to be deleted
	recordID := "DNS_RECORD_ID_TO_DELETE"

	// API endpoint to delete DNS record
	apiEndpoint := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)

	// Create HTTP request
	req, err := http.NewRequest("DELETE", apiEndpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiToken)

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return errors.New("http status is not 200")
	}

	fmt.Println("DNS record deleted successfully!")
	return nil
}

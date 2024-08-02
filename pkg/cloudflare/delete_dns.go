package cloudflare

import (
	"errors"
	"fmt"
	"net/http"
)

func (svc *CloudflareService) DeleteDNS(domain string) error {
	apiToken := svc.Config.CloudflareToken
	zoneName := extractDomain(domain)
	zoneID, err := getZoneId(svc.Config.CloudflareToken, zoneName)
	if err != nil {
		return err
	}

	recordID := "DNS_RECORD_ID_TO_DELETE"
	apiEndpoint := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)

	req, err := http.NewRequest("DELETE", apiEndpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)

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

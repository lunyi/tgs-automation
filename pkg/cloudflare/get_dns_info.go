package cloudflare

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}
type DNSRecordsResponse struct {
	Result []DNSRecord `json:"result"`
}

func GetDnsInfo(dnsName string) {
	// Define your Cloudflare API token or key
	config := util.GetConfig()
	zoneName := extractDomain(dnsName)
	zoneId := GetZoneId(zoneName)
	fmt.Println(fmt.Sprintf("zoneId: %s", zoneId))

	// Get DNS info for the specified DNS record
	dnsRecords, err := getDNSInfo(config.CloudflareToken, zoneId, dnsName)

	if err != nil {
		fmt.Println("Error getting DNS info:", err)
		return
	}

	for _, record := range dnsRecords {
		fmt.Printf("ID: %s, Type: %s, Name: %s, Content: %s\n", record.ID, record.Type, record.Name, record.Content)
	}
}

func extractDomain(dns string) string {
	// Split the DNS record by dot (.)
	parts := strings.Split(dns, ".")

	// Extract the domain
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return dns
}

// Function to get DNS info for a specific DNS record within a zone
func getDNSInfo(apiToken, zoneID, dnsName string) ([]DNSRecord, error) {
	// Construct the endpoint URL
	endpoint := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s", zoneID, dnsName)

	resp, err := sendRequest(apiToken, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.LogError(fmt.Sprintf("Error reading response body:", err))
	}

	// Print response body
	log.LogInfo(string(body))
	defer resp.Body.Close()

	// Decode the JSON response
	var dnsResponse DNSRecordsResponse
	err = json.Unmarshal(body, &dnsResponse)
	if err != nil {
		log.LogError(fmt.Sprintln("Error reading response body:", err))
	}

	log.LogInfo(fmt.Sprintf("%+v\n", dnsResponse))
	// Return the DNS records
	return dnsResponse.Result, nil
}

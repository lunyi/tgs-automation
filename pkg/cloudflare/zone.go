package cloudflare

import (
	"encoding/json"
	"errors"
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
)

type ZoneResponse struct {
	Result []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

func GetZoneId(domain string) (string, error) {
	// Define your Cloudflare API token or key
	config := util.GetConfig()

	apiToken := config.CloudflareToken
	// Get Zone ID for the domain
	zoneID, err := getZoneID(apiToken, domain)
	if err != nil {
		log.LogError(fmt.Sprintln("Error getting Zone ID:", err))
		return "", err
	}

	return zoneID, nil
}

// Function to get Zone ID for the domain
func getZoneID(apiToken, domain string) (string, error) {
	// Construct the endpoint URL
	endpoint := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?name=%s", domain)

	// Create a new HTTP client
	resp, err := sendRequest(apiToken, "GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var zoneResponse ZoneResponse
	err = json.NewDecoder(resp.Body).Decode(&zoneResponse)
	if err != nil {
		return "", err
	}

	// Check if any result is found
	if len(zoneResponse.Result) == 0 {
		return "", errors.New("no zone found")
	}

	// Return the Zone ID
	return zoneResponse.Result[0].ID, nil
}

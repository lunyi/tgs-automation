package cloudflare

import (
	"cdnetwork/internal/log"
	"fmt"
	"net/http"
	"strings"
)

func sendRequest(apiToken, httpMethod, endpoint string, body *strings.Reader) (*http.Response, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest(httpMethod, endpoint, body)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error create request: %s, method: %s, err: %v", endpoint, httpMethod, err))
		return nil, err
	}

	// Add the necessary headers for Cloudflare API
	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error in send request: %s, method: %s, err: %v", endpoint, httpMethod, err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return nil, err
	}

	return resp, nil
}

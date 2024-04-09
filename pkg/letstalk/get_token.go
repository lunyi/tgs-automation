package letstalk

import (
	"bytes"
	"cdnetwork/internal/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetToken(config util.LetsTalkConfig) (string, error) {
	apiURL := "https://message.biatalk.cc/bot/v3/sign-token"

	payload := map[string]string{
		"id":    config.AccountId,
		"token": config.ApiKey,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return "", err
	}

	// Create a new POST request with the payload
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	fmt.Println("Response body:", string(body))

	var apiResponse TokenResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return "", err // Assuming your function returns a string and an error
	}

	// Now, apiResponse contains the unmarshalled data and you can use it
	fmt.Printf("Response Status: %d\n", apiResponse.Status)
	fmt.Printf("Message: %s\n", apiResponse.Message)
	fmt.Printf("Token: %s\n", apiResponse.Token)

	// Continue with your logic, possibly returning something based on apiResponse
	return apiResponse.Token, nil
}

type TokenResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

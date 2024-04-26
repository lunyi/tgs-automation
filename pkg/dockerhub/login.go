package dockerhub

// Path: pkg/dockerhub/login.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tgs-automation/internal/util"
)

// Struct to hold the login credentials
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct to parse the response which contains the token
type LoginResponse struct {
	Token string `json:"token"`
}

// Login attempts to authenticate a user and retrieve an authentication token.
func Login(config util.DockerhubConfig) (string, error) {
	credentials := LoginCredentials{
		Username: config.Username,
		Password: config.Password,
	}
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		return "", fmt.Errorf("marshaling credentials: %w", err)
	}

	url := fmt.Sprintf("%s/users/login", config.BaseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	return parseLoginResponse(resp)
}

// parseLoginResponse reads and unmarshals the response body into the LoginResponse struct.
func parseLoginResponse(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	var response LoginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("unmarshaling response: %w", err)
	}

	return response.Token, nil
}

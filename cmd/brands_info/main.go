package main

import (
	"bytes"
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"cdnetwork/pkg/postgresql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	config := util.GetConfig()
	app := postgresql.NewBrandsIncomeInterface(config.Postgresql)

	brands, err := app.GetBrandsIncome()
	if err != nil {
		panic(err)
	}

	message := ""

	for _, brand := range brands {
		message += "\n" + brand.PlatformCode + "\n" +
			"當日營收：" + brand.DailyRevenueUSD + "\n" +
			"當日訂單數量：" + brand.DailyOrderCount + "\n" +
			"當日活躍人數：" + fmt.Sprintf("%d", brand.ActiveUserCount) + "\n" +
			"當月營收：" + brand.CumulativeRevenueUSD + "\n\n"

	}
	log.LogInfo(message)

	token, err := getToken()
	if err != nil {
		log.LogInfo("Token:" + token)
	}
	var rooms []Room
	rooms, err = getRooms(token)
	if err == nil {
		for _, room := range rooms {
			log.LogInfo("Room:" + room.Title + " Token:" + room.Token)
		}
	}
}

func getToken() (string, error) {
	apiURL := "https://message.biatalk.cc/bot/v3/sign-token"

	// Create the request payload
	payload := map[string]string{
		"id":    "100110",
		"token": "UU8baU1KAUm1gqGbsXdclfaicu69Zvms",
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

	var apiResponse ApiResponse
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

func getRooms(token string) ([]Room, error) {
	url := "https://message.biatalk.cc/bot/v3/chatroom/metadata"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var apiResponse RoomApiResponse
	err = json.Unmarshal([]byte(body), &apiResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return nil, err // Assuming your function returns a string and an error
	}

	// Now, apiResponse contains the unmarshalled data and you can use it
	fmt.Printf("Response Status: %d\n", apiResponse.Status)
	fmt.Printf("Message: %s\n", apiResponse.Message)
	fmt.Printf("totalCount: %v\n", apiResponse.TotalCount)
	fmt.Printf("rooms: %v\n", apiResponse.Rooms)

	// Continue with your logic, possibly returning something based on apiResponse
	return apiResponse.Rooms, nil
}

type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}
type RoomApiResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	TotalCount int    `json:"totalCount"`
	Rooms      []Room `json:"rooms"`
}

// Room represents each room within the "rooms" array in the JSON response
type Room struct {
	Title string `json:"title"`
	Token string `json:"token"`
}

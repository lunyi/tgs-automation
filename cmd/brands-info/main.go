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

	message := "日期: " + fmt.Sprintf("%v", brands[0].Date.Format("2006-01-02")) + "<br>"

	for _, brand := range brands {
		message += "<br><b>[" + brand.PlatformCode + "]</b><br>" +
			"當日營收：" + brand.DailyRevenueUSD + "<br>" +
			"當日訂單數量：" + brand.DailyOrderCount + "<br>" +
			"當日活躍人數：" + fmt.Sprintf("%d", brand.ActiveUserCount) + "<br>" +
			"當月營收：" + brand.CumulativeRevenueUSD + "<br>"
	}
	log.LogInfo(message)

	token, err := getToken(config.LetsTalk)
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

	groupKeys := []string{"PG Daily Report"}
	roomTokens := []string{}
	for _, room := range rooms {
		for _, key := range groupKeys {
			if room.Title == key {
				roomTokens = append(roomTokens, room.Token)
			}
		}
	}

	err = sendMessage(token, roomTokens, message)
	if err != nil {
		log.LogError(err.Error())
	}
}

func sendMessage(token string, rooms []string, message string) error {
	apiURL := "https://message.biatalk.cc/bot/v3/message/multi-chatroom"
	// Define the request body
	requestBody := map[string]interface{}{
		"receivers": rooms,
		"msg": map[string]string{
			"type": "html",
			"text": message,
		},
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	// Create a new HTTP POST request
	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Set the content type to application/json
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)
	// Set any additional headers here, such as Authorization if required

	// Execute the HTTP request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return err
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	var apiResponse SendMessageResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}

	// Print the struct to verify
	fmt.Printf("%+v\n", apiResponse)
	return nil
}

func getToken(config util.LetsTalkConfig) (string, error) {
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

// Define the struct to match the response structure
type SendMessageResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	GlobalIds []struct {
		Receiver string `json:"reciever"` // Note: There's a typo in the JSON key; it should be "receiver" based on standard spelling, but must match the JSON.
		GlobalId int64  `json:"globalId"`
		Guid     string `json:"guid"`
	} `json:"globalIds"`
	FailList []interface{} `json:"failList"` // Use []interface{} if the list's content can vary; otherwise, specify the type.
}

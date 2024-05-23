package letstalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendMessage(token string, rooms []string, message string) error {
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

	fmt.Printf("===========apiRequest=========\n")
	fmt.Printf("token: " + token + "\n")

	for _, room := range rooms {
		fmt.Printf("room: " + room + "\n")
	}

	// Print the struct to verify
	fmt.Printf("===========apiResponse=========\n")
	fmt.Printf("%+v\n", apiResponse)
	return nil
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

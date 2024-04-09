package letstalk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRooms(token string) ([]Room, error) {
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

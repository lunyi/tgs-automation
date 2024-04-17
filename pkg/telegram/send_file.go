package telegram

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func SendFile(botToken, chatID, filePath string) {
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Add the chat ID to the form-data
	_ = multipartWriter.WriteField("chat_id", chatID)

	// Open the file to send
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a form file part
	part, err := multipartWriter.CreateFormFile("document", filepath.Base(file.Name()))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}

	// Important to close the writer or the request will be missing the terminating boundary.
	multipartWriter.Close()

	// Create and send the request
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+botToken+"/sendDocument", &requestBody)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("failed to send document")
	}
}

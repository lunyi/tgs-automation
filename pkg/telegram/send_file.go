package telegram

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func SendFile(botToken, chatID, filePath string) error {
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Add the chat ID to the form-data
	if err := multipartWriter.WriteField("chat_id", chatID); err != nil {
		return err
	}

	// Open the file to send
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a form file part
	part, err := multipartWriter.CreateFormFile("document", filepath.Base(file.Name()))
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	// Important to close the writer or the request will be missing the terminating boundary.
	if err = multipartWriter.Close(); err != nil {
		return err
	}

	// Create and send the request
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+botToken+"/sendDocument", &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Reading the body might help with debugging.
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return readErr
		}
		return fmt.Errorf("failed to send document: %s", string(body))
	}

	return nil
}

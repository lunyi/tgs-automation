package util

import (
	"io/ioutil"
	"net/http"
	"strings"
	"tgs-automation/internal/model"

	log "github.com/sirupsen/logrus" // 假設 log 包使用的是 logrus
)

func Call(requestMsg model.HttpRequestMsg) string {
	client := &http.Client{}
	req, err := http.NewRequest(requestMsg.Method, requestMsg.Url, strings.NewReader(requestMsg.Body))
	if err != nil {
		log.Error("Error creating new request: ", err.Error())
		return ""
	}

	for key, value := range requestMsg.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error performing request: ", err.Error())
		return ""
	}
	defer func() {
		if resp.Body != nil {
			err := resp.Body.Close()
			if err != nil {
				log.Error("Error closing response body: ", err.Error())
			}
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body: ", err.Error())
		return ""
	}

	log.Info("Response body: ", string(body))
	return string(body)
}

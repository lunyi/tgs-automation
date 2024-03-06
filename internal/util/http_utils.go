package util

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/model"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Call(requestMsg model.HttpRequestMsg) string {
	client := &http.Client{}
	req, err := http.NewRequest(requestMsg.Method, requestMsg.Url, strings.NewReader(requestMsg.Body))
	if err != nil {
		log.LogError(err.Error())
	}
	for key := range requestMsg.Headers {
		req.Header.Set(key, requestMsg.Headers[key])
	}
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.LogError(err.Error())
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.LogError(err.Error())
	}
	log.LogInfo(string(body))
	return string(body)
}

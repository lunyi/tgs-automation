package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"tgs-automation/internal/constant"
	"tgs-automation/internal/log"
	"tgs-automation/internal/model"
	"tgs-automation/internal/util"
	"time"
)

type BasicConfig struct {
	Uri    string
	Method string
}

func TransferRequestMsg(parameter BasicConfig, config util.CdnConfig) model.HttpRequestMsg {

	var requestMsg = model.HttpRequestMsg{Params: map[string]string{}, Headers: map[string]string{}}
	requestMsg.Uri = parameter.Uri
	requestMsg.Method = parameter.Method
	requestMsg.Url = constant.HttpRequestPrefix + parameter.Uri
	if len(config.CdnEndPoint) == 0 || "{endPoint}" == config.CdnEndPoint {
		requestMsg.Host = constant.HttpRequestDomain
		requestMsg.Url = "https://" + constant.HttpRequestDomain + parameter.Uri
	} else {
		requestMsg.Host = config.CdnEndPoint
		requestMsg.Url = "https://" + config.CdnEndPoint + parameter.Uri
	}
	return requestMsg
}

func Invoke(url BasicConfig, jsonStr string) string {
	config := util.GetConfig()

	var requestMsg = TransferRequestMsg(url, config.CdnNetwork)

	if url.Method == "POST" || url.Method == "PUT" || url.Method == "PATCH" || url.Method == "DELETE" {
		requestMsg.Body = jsonStr
	}

	dateStr := getRfc1123Date()
	authorization := "Basic " + authorization(config.CdnNetwork.CdnUserName, hmacsha1(dateStr, config.CdnNetwork.CdnApiKey))
	log.LogInfo(fmt.Sprintf("authorization %s", authorization))
	requestMsg.Headers["Date"] = dateStr
	requestMsg.Headers["Host"] = requestMsg.Host
	requestMsg.Headers["Content-Type"] = constant.ApplicationJson
	requestMsg.Headers[constant.Authorization] = authorization
	requestMsg.Headers[constant.XCncAuthMethod] = constant.BASIC
	return util.Call(requestMsg)
}

func authorization(userName string, passwd string) string {
	result := base64.StdEncoding.EncodeToString([]byte(userName + ":" + passwd))
	return result
}

func hmacsha1(date string, apikey string) string {
	key := []byte(apikey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(date))
	value := mac.Sum(nil)
	result := base64.StdEncoding.EncodeToString(value)
	return result
}

func getRfc1123Date() string {
	date := time.Now().UTC() // UTC time
	return date.Format(http.TimeFormat)
}

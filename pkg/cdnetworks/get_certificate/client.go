package get_certificate

import (
	"cdnetwork/internal/auth"
	"cdnetwork/pkg/cdnetworks/get_certificate/models"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

func GetCertificate(domainName string) (string, error) {
	//call api
	queryCertificateListRequest := models.QueryCertificateListRequest{}

	config := auth.BasicConfig{
		Uri:    "/api/ssl/certificate",
		Method: "GET",
	}
	response := auth.Invoke(config, queryCertificateListRequest.String())
	//使用response結果解析xml數據
	xmlData := response
	var sslCerts models.SslCertificates
	err := xml.Unmarshal([]byte(xmlData), &sslCerts)
	if err != nil {
		fmt.Printf("Error unmarshalling XML: %v\n", err)
		return "", err
	}

	// 處理並轉換數據
	var certsJSON []models.CertificateJSON
	for _, cert := range sslCerts.SslCertificate { // 注意這裡的字段名與結構體中的字段名相匹配
		splitName := strings.Split(cert.Name, "_")
		domainName := splitName[len(splitName)-1] // 取"_"後的部分

		certsJSON = append(certsJSON, models.CertificateJSON{
			CertificateID: cert.CertificateID,
			DomainName:    domainName,
		})
	}

	// 將處理後的數據轉換为JSON
	jsonData, err := json.MarshalIndent(certsJSON, "", "    ")
	if err != nil {
		fmt.Printf("Error marshalling to JSON: %v\n", err)
		return "", err
	}
	// 輸出JSON
	fmt.Println(string(jsonData))

	for _, cert := range certsJSON {
		if cert.DomainName == domainName {
			return cert.CertificateID, nil
			break
		}
	}
	return "", fmt.Errorf("No certficate id for %s", domainName)
}

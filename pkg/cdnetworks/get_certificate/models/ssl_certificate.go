package models

import "encoding/xml"

type SslCertificate struct {
	CertificateID string `xml:"certificate-id"`
	Name          string `xml:"name"`
	// 省略其他字段...
}

type SslCertificates struct {
	XMLName        xml.Name         `xml:"ssl-certificates"`
	SslCertificate []SslCertificate `xml:"ssl-certificate"` // 這裡的字段名改為與XML標籤相匹配
}

// 定義JSON結構體
type CertificateJSON struct {
	CertificateID string `json:"certificate_id"`
	DomainName    string `json:"domain_name"`
}

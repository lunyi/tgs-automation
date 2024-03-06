package models

import "github.com/alibabacloud-go/tea/tea"

type QueryCertificateListResponse struct {
	// {"en":"Certificate list information", "zh_CN":"证书列表信息"}
	SslCertificates []*QueryCertificateListResponseSslCertificates `json:"ssl-certificates,omitempty" xml:"ssl-certificates,omitempty" require:"true" type:"Repeated"`
}

func (s QueryCertificateListResponse) String() string {
	return tea.Prettify(s)
}

func (s QueryCertificateListResponse) GoString() string {
	return s.String()
}

func (s *QueryCertificateListResponse) SetSslCertificates(v []*QueryCertificateListResponseSslCertificates) *QueryCertificateListResponse {
	s.SslCertificates = v
	return s
}

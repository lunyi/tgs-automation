package models

import "github.com/alibabacloud-go/tea/tea"

type QueryCertificateListResponseSslCertificatesRelatedDomains struct {
	// {"en":"Accelerated domain name ID", "zh_CN":"加速域名ID"}
	DomainId *string `json:"domain-id,omitempty" xml:"domain-id,omitempty" require:"true"`
	// {"en":"Name of accelerated domain name", "zh_CN":"加速域名的名称"}
	DomainName *string `json:"domain-name,omitempty" xml:"domain-name,omitempty" require:"true"`
}

func (s QueryCertificateListResponseSslCertificatesRelatedDomains) String() string {
	return tea.Prettify(s)
}

func (s QueryCertificateListResponseSslCertificatesRelatedDomains) GoString() string {
	return s.String()
}

func (s *QueryCertificateListResponseSslCertificatesRelatedDomains) SetDomainId(v string) *QueryCertificateListResponseSslCertificatesRelatedDomains {
	s.DomainId = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificatesRelatedDomains) SetDomainName(v string) *QueryCertificateListResponseSslCertificatesRelatedDomains {
	s.DomainName = &v
	return s
}

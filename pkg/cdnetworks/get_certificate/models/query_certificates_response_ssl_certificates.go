package models

import "github.com/alibabacloud-go/tea/tea"

type QueryCertificateListResponseSslCertificates struct {
	// {"en":"Certificate ID", "zh_CN":"证书ID"}
	CertificateId *string `json:"certificate-id,omitempty" xml:"certificate-id,omitempty" require:"true"`
	// {"en":"Certificate name, unique to customer granularity", "zh_CN":"证书名称，客户粒度下是唯一的"}
	Name *string `json:"name,omitempty" xml:"name,omitempty" require:"true"`
	// {"en":"Remarks on cerfiticate file", "zh_CN":"证书文件的备注"}
	Comment *string `json:"comment,omitempty" xml:"comment,omitempty" require:"true"`
	// {"en":"Shared, optional values are true and false, true represents shared certificates, false represents unshared certificates, default is false
	// This certificate allows cross-customer use when share-ssl is true. (The API does not support cross-customer use certificates. Contact customer service for manual configuration if required.)", "zh_CN":"是否共享，true表示共享证书，false表示非共享证书"}
	ShareSsl *string `json:"share-ssl,omitempty" xml:"share-ssl,omitempty" require:"true"`
	// {"en":"Certificate effective start time (CST), such as 2016-08-01 07:00:00", "zh_CN":"证书有效期的起始时间（CST时区），例如：2016-08-01 07:00:00"}
	CertificateValidityFrom *string `json:"certificate-validity-from,omitempty" xml:"certificate-validity-from,omitempty" require:"true"`
	// {"en":"Certificate effective end time (CST), such as 2018-08-01 19:00:00", "zh_CN":"证书有效期的到期时间（CST时区），例如：2018-08-01 19:00:00"}
	CertificateValidityTo *string `json:"certificate-validity-to,omitempty" xml:"certificate-validity-to,omitempty" require:"true"`
	// {"en":"List of domain names using the current certificate", "zh_CN":"使用当前证书的域名列表"}
	RelatedDomains []*QueryCertificateListResponseSslCertificatesRelatedDomains `json:"related-domains,omitempty" xml:"related-domains,omitempty" require:"true" type:"Repeated"`
	// {"en":"dns-names", "zh_CN":"授权域名列表，证书使用者可选名称，父标签"}
	DnsNames []*string `json:"dns-names,omitempty" xml:"dns-names,omitempty" require:"true" type:"Repeated"`
	// {"en":"The CRT certificate serial number", "zh_CN":"crt证书序列号"}
	CertificateSerial *string `json:"certificate-serial,omitempty" xml:"certificate-serial,omitempty" require:"true"`
	// {"en":"The MD5 value of CRT file.", "zh_CN":"CRT文件内容的md5值。"}
	CrtMd5 *string `json:"crt-md5,omitempty" xml:"crt-md5,omitempty" require:"true"`
	// {"en":"The MD5 value of KEY file.", "zh_CN":"KEY文件内容的md5值。"}
	KeyMd5 *string `json:"key-md5,omitempty" xml:"key-md5,omitempty" require:"true"`
	// {"en":"The MD5 value of CA file.", "zh_CN":"CA。"}
	CaMd5 *string `json:"ca-md5,omitempty" xml:"ca-md5,omitempty" require:"true"`
	// {"en":"The CRT certificate issuer", "zh_CN":"crt证书颁布者"}
	CertificateIssuer *string `json:"certificate-issuer,omitempty" xml:"certificate-issuer,omitempty" require:"true"`
}

func (s QueryCertificateListResponseSslCertificates) String() string {
	return tea.Prettify(s)
}

func (s QueryCertificateListResponseSslCertificates) GoString() string {
	return s.String()
}

func (s *QueryCertificateListResponseSslCertificates) SetCertificateId(v string) *QueryCertificateListResponseSslCertificates {
	s.CertificateId = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetName(v string) *QueryCertificateListResponseSslCertificates {
	s.Name = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetComment(v string) *QueryCertificateListResponseSslCertificates {
	s.Comment = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetShareSsl(v string) *QueryCertificateListResponseSslCertificates {
	s.ShareSsl = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetCertificateValidityFrom(v string) *QueryCertificateListResponseSslCertificates {
	s.CertificateValidityFrom = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetCertificateValidityTo(v string) *QueryCertificateListResponseSslCertificates {
	s.CertificateValidityTo = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetRelatedDomains(v []*QueryCertificateListResponseSslCertificatesRelatedDomains) *QueryCertificateListResponseSslCertificates {
	s.RelatedDomains = v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetDnsNames(v []*string) *QueryCertificateListResponseSslCertificates {
	s.DnsNames = v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetCertificateSerial(v string) *QueryCertificateListResponseSslCertificates {
	s.CertificateSerial = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetCrtMd5(v string) *QueryCertificateListResponseSslCertificates {
	s.CrtMd5 = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetKeyMd5(v string) *QueryCertificateListResponseSslCertificates {
	s.KeyMd5 = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetCaMd5(v string) *QueryCertificateListResponseSslCertificates {
	s.CaMd5 = &v
	return s
}

func (s *QueryCertificateListResponseSslCertificates) SetCertificateIssuer(v string) *QueryCertificateListResponseSslCertificates {
	s.CertificateIssuer = &v
	return s
}

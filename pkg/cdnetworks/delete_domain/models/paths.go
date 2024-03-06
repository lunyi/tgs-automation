package models

import "github.com/alibabacloud-go/tea/tea"

type Paths struct {
	// {"en":"", "zh_CN":"域名名称或域名id，在请求的url后面"}
	DomainName *string `json:"domain-name,omitempty" xml:"domain-name,omitempty" require:"true"`
}

func (s Paths) String() string {
	return tea.Prettify(s)
}

func (s Paths) GoString() string {
	return s.String()
}

func (s *Paths) SetDomainName(v string) *Paths {
	s.DomainName = &v
	return s
}

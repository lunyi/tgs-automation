package models

import "github.com/alibabacloud-go/tea/tea"

type QueryApiDomainListServiceResponse struct {
	// {"en":"domain list", "zh_CN":"域名列表"}
	DomainList []*QueryApiDomainListServiceResponseDomainList `json:"domain-list,omitempty" xml:"domain-list,omitempty" require:"true" type:"Repeated"`
}

func (s QueryApiDomainListServiceResponse) String() string {
	return tea.Prettify(s)
}

func (s QueryApiDomainListServiceResponse) GoString() string {
	return s.String()
}

func (s *QueryApiDomainListServiceResponse) SetDomainList(v []*QueryApiDomainListServiceResponseDomainList) *QueryApiDomainListServiceResponse {
	s.DomainList = v
	return s
}

package models

import "github.com/alibabacloud-go/tea/tea"

type QueryApiDomainListServiceRequest struct {
}

func (s QueryApiDomainListServiceRequest) String() string {
	return tea.Prettify(s)
}

func (s QueryApiDomainListServiceRequest) GoString() string {
	return s.String()
}

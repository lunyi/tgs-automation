package models

import "github.com/alibabacloud-go/tea/tea"

type EnableSingleDomainServiceRequest struct {
}

func (s EnableSingleDomainServiceRequest) String() string {
	return tea.Prettify(s)
}

func (s EnableSingleDomainServiceRequest) GoString() string {
	return s.String()
}

package models

import "github.com/alibabacloud-go/tea/tea"

type QueryCertificateListRequest struct {
}

func (s QueryCertificateListRequest) String() string {
	return tea.Prettify(s)
}

func (s QueryCertificateListRequest) GoString() string {
	return s.String()
}

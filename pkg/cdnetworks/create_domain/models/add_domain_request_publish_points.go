package models

import "github.com/alibabacloud-go/tea/tea"

type AddCdnDomainRequestPublishPoints struct {
	// {"en":"Publish point, support multiple, do not pass the system by default to generate a publishing point uri for slash /", "zh_CN":"发布点，支持多个，不传系统默认生成一条发布点uri为“/”"}
	Uri *string `json:"uri,omitempty" xml:"uri,omitempty"`
}

func (s AddCdnDomainRequestPublishPoints) String() string {
	return tea.Prettify(s)
}

func (s AddCdnDomainRequestPublishPoints) GoString() string {
	return s.String()
}

func (s *AddCdnDomainRequestPublishPoints) SetUri(v string) *AddCdnDomainRequestPublishPoints {
	s.Uri = &v
	return s
}

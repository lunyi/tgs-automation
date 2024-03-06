package models

import "github.com/alibabacloud-go/tea/tea"

type Paths struct {
	// {"en":"", "zh_CN":"加速域名在系统中对应的ID
	// 1. 参看请求示例中的url，123344对应的就是domain-id
	// 2. 可以通过【获取域名配置】和【获取域名列表】接口查询到domain-id"}
	DomainId *int `json:"domainId,omitempty" xml:"domainId,omitempty" require:"true"`
}

func (s Paths) String() string {
	return tea.Prettify(s)
}

func (s Paths) GoString() string {
	return s.String()
}

func (s *Paths) SetDomainId(v int) *Paths {
	s.DomainId = &v
	return s
}

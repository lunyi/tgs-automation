package models

import "github.com/alibabacloud-go/tea/tea"

type Parameters struct {
	// {"en":"Public CNAME alias, optional entry, does not indicate all domain names under the query account number
	// The customer has the demand that the domain cname share more than one level, so we introduce the cname-label identifier in the interface, which is a set of domain cname with the same cname-label, and share the first level of cname.", "zh_CN":"共用一级别名标示，可选入参，不选表示查询账号下所有域名
	// 客户存在较多一级域名共用的需求，因此在接口中引入cname-label标识，即拥有相同cname-label的一组域名，共用一级cname。关于cname-label的具体使用方式和注意事项，请参看【创建加速域名】和【修改域名配置】接口"}
	CnameLabel *string `json:"cname-label,omitempty" xml:"cname-label,omitempty"`
}

func (s Parameters) String() string {
	return tea.Prettify(s)
}

func (s Parameters) GoString() string {
	return s.String()
}

func (s *Parameters) SetCnameLabel(v string) *Parameters {
	s.CnameLabel = &v
	return s
}

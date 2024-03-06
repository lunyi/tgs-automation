package models

import "github.com/alibabacloud-go/tea/tea"

type DeleteApiDomainServiceRequest struct {
	// {"en":"Accelerated domain name, choose one from domain-id. Accelerate the ID of the domain name in the system
	// Note:
	// 1. See the url in the request example, 123344 for domain-id
	// 2、After the domain name is successfully submitted, the location access url in the return parameter can be queried to the domain-id of the domain name; You can also query domain-id through the Get domain Configuration and Get domain List interfaces", "zh_CN":"加速域名与domain-id二选一。
	// domain-id：加速域名在系统中对应的ID
	// domain-name：加速的域名
	// 注意：
	// 1、参看请求示例中的url，123344对应的就是domain-id
	// 2、创建域名成功提交后，返回参数中的location访问url中，能够查询到域名对应的domain-id；也可以通过【获取域名配置】和【获取域名列表】接口查询到domain-id"}
	DomainName *string `json:"domainName,omitempty" xml:"domainName,omitempty" require:"true"`
	// {"en":"Accelerated domain name, choose one from domain-id. Accelerate the ID of the domain name in the system
	// Note:
	// 1. See the url in the request example, 123344 for domain-id
	// 2、After the domain name is successfully submitted, the location access url in the return parameter can be queried to the domain-id of the domain name; You can also query domain-id through the Get domain Configuration and Get domain List interfaces", "zh_CN":"加速域名与domain-id二选一。
	// domain-id：加速域名在系统中对应的ID
	// domain-name：加速的域名
	// 注意：
	// 1、参看请求示例中的url，123344对应的就是domain-id
	// 2、创建域名成功提交后，返回参数中的location访问url中，能够查询到域名对应的domain-id；也可以通过【获取域名配置】和【获取域名列表】接口查询到domain-id"}
	DomainId *string `json:"domainId,omitempty" xml:"domainId,omitempty" require:"true"`
}

func (s DeleteApiDomainServiceRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteApiDomainServiceRequest) GoString() string {
	return s.String()
}

func (s *DeleteApiDomainServiceRequest) SetDomainName(v string) *DeleteApiDomainServiceRequest {
	s.DomainName = &v
	return s
}

func (s *DeleteApiDomainServiceRequest) SetDomainId(v string) *DeleteApiDomainServiceRequest {
	s.DomainId = &v
	return s
}

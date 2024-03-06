package models

import "github.com/alibabacloud-go/tea/tea"

type AddCdnDomainRequestOriginConfig struct {
	// {"en":"Origin site address, which can be an IP or a domain name.
	// 1. IP is separated by semicolons and multiple IPs are supported.
	// 2. Only one domain name can be entered. IP and domain names cannot be entered at the same time.
	// 3. Maximum character limit is 500.", "zh_CN":"回源地址，可以是IP或域名。 1. IP以分号分隔，支持多个。 2. 域名只能输入一个。IP与域名不能同时输入。 3.限制最大不能超过500个字符长度。"}
	OriginIps *string `json:"origin-ips,omitempty" xml:"origin-ips,omitempty"`
	// {"en":"The Origin HOST for changing the HOST field in the return source HTTP request header.
	// Note: It should be domain or IP format. For domain name format, each segement separated by a dot, does not exceed 62 characters, the total length should not exceed 128 characters.", "zh_CN":"回源HOST，用于更改回源HTTP请求头中的HOST字段。支持格式为: 域名或ip。
	// 注意：必须符合ip/域名格式规范。如果是域名，则域名每段（点号分隔）长度小于等于62，域名总长度小于等于128。"}
	DefaultOriginHostHeader *string `json:"default-origin-host-header,omitempty" xml:"default-origin-host-header,omitempty"`
	// {"en":"Origin port", "zh_CN":"回源请求端口"}
	OriginPort *string `json:"origin-port,omitempty" xml:"origin-port,omitempty"`
}

func (s AddCdnDomainRequestOriginConfig) String() string {
	return tea.Prettify(s)
}

func (s AddCdnDomainRequestOriginConfig) GoString() string {
	return s.String()
}

func (s *AddCdnDomainRequestOriginConfig) SetOriginIps(v string) *AddCdnDomainRequestOriginConfig {
	s.OriginIps = &v
	return s
}

func (s *AddCdnDomainRequestOriginConfig) SetDefaultOriginHostHeader(v string) *AddCdnDomainRequestOriginConfig {
	s.DefaultOriginHostHeader = &v
	return s
}

func (s *AddCdnDomainRequestOriginConfig) SetOriginPort(v string) *AddCdnDomainRequestOriginConfig {
	s.OriginPort = &v
	return s
}

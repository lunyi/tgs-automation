package models

import "github.com/alibabacloud-go/tea/tea"

type AddCdnDomainResponse struct {
	// {"en":"httpstatus=202; Indicates that the new domain API was successfully invoked, and the current deployment of the new domain can be viewed using x-cnc-request-id in the header", "zh_CN":"httpstatus=202;   表示成功调用新增域名接口，可使用header中的x-cnc-request-id查看当前新增域名的部署情况"}
	HttpStatus *int `json:"http status code,omitempty" xml:"http status code,omitempty" require:"true"`
	// {"en":"Uniquely identified id for querying tasks per request (for all API)", "zh_CN":"唯一标示的id，用于查询每次请求的任务 （适用全部接口）"}
	XCncRequestId *string `json:"x-cnc-request-id,omitempty" xml:"x-cnc-request-id,omitempty" require:"true"`
	// {"en":"The URL used to access the domain information, where domain-id is the unique token generated by our cloud platform for the domain name and whose value is a string.", "zh_CN":"响应信用于访问该域名信息的URL，其中domain-id为我司云平台为该域名生成的唯一标示，其值为字符串。"}
	Location *string `json:"location,omitempty" xml:"location,omitempty" require:"true"`
	// {"en":"The name of the service domain automatically generated by the My company, for example: xxxx.cdn30.com", "zh_CN":"由我司自动生成的服务域名名称，例如：xxxx.cdn30.com"}
	Cname *string `json:"cname,omitempty" xml:"cname,omitempty" require:"true"`
	// {"en":"Request result code", "zh_CN":"请求结果状态码"}
	Code *string `json:"code,omitempty" xml:"code,omitempty" require:"true"`
	// {"en":"Response information, when success is successful", "zh_CN":"响应信息，成功时为success"}
	Message *string `json:"message,omitempty" xml:"message,omitempty" require:"true"`
}

func (s AddCdnDomainResponse) String() string {
	return tea.Prettify(s)
}

func (s AddCdnDomainResponse) GoString() string {
	return s.String()
}

func (s *AddCdnDomainResponse) SetHttpStatus(v int) *AddCdnDomainResponse {
	s.HttpStatus = &v
	return s
}

func (s *AddCdnDomainResponse) SetXCncRequestId(v string) *AddCdnDomainResponse {
	s.XCncRequestId = &v
	return s
}

func (s *AddCdnDomainResponse) SetLocation(v string) *AddCdnDomainResponse {
	s.Location = &v
	return s
}

func (s *AddCdnDomainResponse) SetCname(v string) *AddCdnDomainResponse {
	s.Cname = &v
	return s
}

func (s *AddCdnDomainResponse) SetCode(v string) *AddCdnDomainResponse {
	s.Code = &v
	return s
}

func (s *AddCdnDomainResponse) SetMessage(v string) *AddCdnDomainResponse {
	s.Message = &v
	return s
}

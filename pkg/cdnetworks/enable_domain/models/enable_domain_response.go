package models

import "github.com/alibabacloud-go/tea/tea"

type EnableSingleDomainServiceResponse struct {
	// {"en":"Error code, which appears when HTTPStatus is not 202, represents the error type of the current request call", "zh_CN":"错误代码，当HTTPStatus不为202时出现，表示当前请求调用的错误类型"}
	Code *string `json:"code,omitempty" xml:"code,omitempty" require:"true"`
	// {"en":"Response information, success when successful", "zh_CN":"响应信息，成功时为success"}
	Message *string `json:"message,omitempty" xml:"message,omitempty" require:"true"`
	// {"en":"httpstatus=202; Indicates that the new domain API was successfully invoked, and the current deployment of the new domain can be viewed using x-cnc-request-id in the header", "zh_CN":"httpstatus=202;   表示成功调用新增域名接口，可使用header中的x-cnc-request-id查看当前新增域名的部署情况"}
	HttpStatus *int `json:"http status code,omitempty" xml:"http status code,omitempty" require:"true"`
	// {"en":"Uniquely identified id for querying tasks per request (for all API)", "zh_CN":"唯一标示的id，用于查询每次请求的任务 （适用全部接口）"}
	XCncRequestId *string `json:"x-cnc-request-id,omitempty" xml:"x-cnc-request-id,omitempty" require:"true"`
}

func (s EnableSingleDomainServiceResponse) String() string {
	return tea.Prettify(s)
}

func (s EnableSingleDomainServiceResponse) GoString() string {
	return s.String()
}

func (s *EnableSingleDomainServiceResponse) SetCode(v string) *EnableSingleDomainServiceResponse {
	s.Code = &v
	return s
}

func (s *EnableSingleDomainServiceResponse) SetMessage(v string) *EnableSingleDomainServiceResponse {
	s.Message = &v
	return s
}

func (s *EnableSingleDomainServiceResponse) SetHttpStatus(v int) *EnableSingleDomainServiceResponse {
	s.HttpStatus = &v
	return s
}

func (s *EnableSingleDomainServiceResponse) SetXCncRequestId(v string) *EnableSingleDomainServiceResponse {
	s.XCncRequestId = &v
	return s
}

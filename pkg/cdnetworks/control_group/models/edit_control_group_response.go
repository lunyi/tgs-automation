package models

type EditControlGroupResponse struct {
	// {"en":"Status Code", "zh_CN":"错误具体状态码"}
	Code *int `json:"code,omitempty" xml:"code,omitempty" require:"true"`
	// {"en":"Message", "zh_CN":"消息提示"}
	Msg *string `json:"msg,omitempty" xml:"msg,omitempty" require:"true"`
	// {"en":"Success Mark", "zh_CN":"成功标记"}
	Success *bool `json:"success,omitempty" xml:"success,omitempty" require:"true"`
	// {"en":"Control Group Name", "zh_CN":"Control Group名称"}
	ControlGroupName *string `json:"controlGroupName,omitempty" xml:"controlGroupName,omitempty" require:"true"`
}

func (s *EditControlGroupResponse) SetCode(v int) *EditControlGroupResponse {
	s.Code = &v
	return s
}

func (s *EditControlGroupResponse) SetMsg(v string) *EditControlGroupResponse {
	s.Msg = &v
	return s
}

func (s *EditControlGroupResponse) SetSuccess(v bool) *EditControlGroupResponse {
	s.Success = &v
	return s
}

func (s *EditControlGroupResponse) SetControlGroupName(v string) *EditControlGroupResponse {
	s.ControlGroupName = &v
	return s
}

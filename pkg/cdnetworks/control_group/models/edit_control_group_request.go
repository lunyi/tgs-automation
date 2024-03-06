package models

import "github.com/alibabacloud-go/tea/tea"

type EditControlGroupRequest struct {
	// {"en":"Control Group name, which only the User Customized type Control Group can be modified, customer type Control Group and product type Control Group can not be modified. User Customized type Control Group keeps the original Control Group name if no value is passed", "zh_CN":"Control Group名称，只有自定义类型的Control Group可做修改，若是客户类型与合同类型Control Group则不做修改。自定义类型Control Group若不传值则保持原来的Control Group名称"}
	ControlGroupName *string `json:"controlGroupName,omitempty" xml:"controlGroupName,omitempty"`
	// {"en":"Account object array,Used to specify accounts with permission.  all types of Control Group can be modified, if no value is passed, the original accountList will be emptied", "zh_CN":"账号对象数组, 用来指定有权限访问的账号。客户类型，合同类型与自定义类型的Control Group都可以做修改，若不传值则将原accountList清空"}
	AccountList []*LoginName `json:"accountList,omitempty" xml:"accountList,omitempty" type:"Repeated"`
	// {"en":"Domain array, which only the User Customized type Control Group can be modified, customer type Control Group and product type Control Group can not be modified.User Customized type Control Group empties the original domainList if no value is passed", "zh_CN":"域名字符串数组，只有自定义类型的Control Group可做修改，若是客户类型与合同类型Control Group则不做修改。自定义类型Control Group若不传值则将原domainList清空"}
	DomainList []*string `json:"domainList,omitempty" xml:"domainList,omitempty" type:"Repeated"`
	// {"en":"Whether to add:
	// 1. Do not pass or pass false: rewrite method;
	// 2. Pass true: append method.", "zh_CN":"是否追加:
	// 1.不传或false：覆盖方式;
	// 2.传true：追加方式."}
	IsAdd *bool `json:"isAdd,omitempty" xml:"isAdd,omitempty"`
}

func (s EditControlGroupRequest) String() string {
	return tea.Prettify(s)
}

func (s *EditControlGroupRequest) SetControlGroupName(v string) *EditControlGroupRequest {
	s.ControlGroupName = &v
	return s
}

func (s *EditControlGroupRequest) SetAccountList(v []*LoginName) *EditControlGroupRequest {
	s.AccountList = v
	return s
}

func (s *EditControlGroupRequest) SetDomainList(v []*string) *EditControlGroupRequest {
	s.DomainList = v
	return s
}

func (s *EditControlGroupRequest) SetIsAdd(v bool) *EditControlGroupRequest {
	s.IsAdd = &v
	return s
}

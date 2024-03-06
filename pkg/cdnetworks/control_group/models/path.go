package models

type Paths struct {
	// {"en":"Control Group Code", "zh_CN":"Control Group 编号"}
	ControlGroupCode *string `json:"ControlGroupCode,omitempty" xml:"ControlGroupCode,omitempty" require:"true"`
}

func (s *Paths) SetControlGroupCode(v string) *Paths {
	s.ControlGroupCode = &v
	return s
}

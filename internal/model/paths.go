package model

import (
	"github.com/alibabacloud-go/tea/tea"
)

type Paths struct {
}

func (s Paths) String() string {
	return tea.Prettify(s)
}

func (s Paths) GoString() string {
	return s.String()
}

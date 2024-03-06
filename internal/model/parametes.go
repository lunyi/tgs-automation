package model

import (
	"github.com/alibabacloud-go/tea/tea"
)

type Parameters struct {
}

func (s Parameters) String() string {
	return tea.Prettify(s)
}

func (s Parameters) GoString() string {
	return s.String()
}

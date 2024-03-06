package model

import (
	"github.com/alibabacloud-go/tea/tea"
)

type ResponseHeader struct {
}

func (s ResponseHeader) String() string {
	return tea.Prettify(s)
}

func (s ResponseHeader) GoString() string {
	return s.String()
}

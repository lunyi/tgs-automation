package model

import (
	"github.com/alibabacloud-go/tea/tea"
)

type RequestHeader struct {
}

func (s RequestHeader) String() string {
	return tea.Prettify(s)
}

func (s RequestHeader) GoString() string {
	return s.String()
}

package foo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// REF: https://blog.kenwsc.com/posts/2020/discover-go-unit-test-stub-and-mock/

type dbMock struct{}

func newDbMock() *dbMock {
	return &dbMock{}
}

func (d *dbMock) First(out interface{}) {
	out.(*User).Age = 1
}

func TestFooDatabaseCustomMock(t *testing.T) {
	want := 1
	m := newDbMock()
	actual := fooDatabaseCaseIndirectCall(m)
	assert.Equal(t, want, actual)
}

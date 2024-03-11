package foo

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// REF: https://blog.kenwsc.com/posts/2020/discover-go-unit-test-stub-and-mock/

func TestFooDatabaseGomock(t *testing.T) {
	expect := 1
	var user User
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockDatabase(ctrl)
	m.EXPECT().First(gomock.Eq(&user)).SetArg(0, User{Age: expect})
	actual := fooDatabaseCaseIndirectCall(m)
	assert.Equal(t, expect, actual)
}

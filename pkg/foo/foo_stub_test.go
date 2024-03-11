package foo

import (
	"testing"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
)

func TestFooDatabaseByValueFunc(t *testing.T) {
	expect := 1
	stub := gostub.Stub(&getUserAgeValueFunc, func() int {
		return 1
	})
	defer stub.Reset()
	actual := getUserAgeValueFunc()
	assert.Equal(t, expect, actual)
}

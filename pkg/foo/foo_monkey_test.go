package foo

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFooDatabaseByMonkeyPatch(t *testing.T) {
	want := 1
	user := User{Age: 1}
	db := &gorm.DB{}
	patch := monkey.Patch(gorm.Open, func(string, ...interface{}) (*gorm.DB, error) {
		return db, nil
	})
	patchFirst := monkey.PatchInstanceMethod(reflect.TypeOf(db), "First", func(_ *gorm.DB, out interface{}, _ ...interface{}) *gorm.DB {
		val := reflect.ValueOf(out).Elem()
		substitute := reflect.ValueOf(user)
		val.Set(substitute)
		return db
	})
	patchClose := monkey.PatchInstanceMethod(reflect.TypeOf(db), "Close", func(*gorm.DB) error {
		return nil
	})
	defer func() {
		patch.Restore()
		patchFirst.Restore()
		patchClose.Restore()
	}()
	actual := fooDatabaseCaseDirectCall()
	assert.Equal(t, want, actual)
}

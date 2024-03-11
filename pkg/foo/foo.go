package foo

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Age      int
	Birthday time.Time
}

type Database interface {
	First(interface{})
}

func fooDatabaseCaseIndirectCall(db Database) int {
	var user User
	db.First(&user)
	return user.Age
}

func fooDatabaseCaseByValueFunc() int {
	return getUserAgeValueFunc()
}

var getUserAgeValueFunc = func() int {
	var user User
	db, err := gorm.Open("postgres", "host=myhost user=gorm dbname=gorm sslmode=disable password=mypassword")
	if err != nil {
		panic("connect fail")
	}
	res := db.First(&user)
	if res.Error != nil {
		panic("error")
	}
	db.Close()
	return user.Age
}

func fooDatabaseCaseDirectCall() int {
	var user User
	db, err := gorm.Open("postgres", "host=myhost user=gorm dbname=gorm sslmode=disable password=mypassword")
	if err != nil {
		panic("connect fail")
	}
	res := db.First(&user)
	if res.Error != nil {
		panic("error")
	}
	db.Close()
	return user.Age
}

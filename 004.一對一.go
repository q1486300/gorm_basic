package main

import (
	"fmt"
	"gorm.io/gorm"
)

type User1 struct {
	ID       uint
	Name     string
	Age      int
	Gender   bool
	UserInfo UserInfo // 通過 UserInfo 可以拿到用戶詳細資料
}

type UserInfo struct {
	User1ID uint   // 外鍵
	User1   *User1 // 改成指針，不然會嵌套引用
	ID      uint
	Addr    string
	Like    string
}

func main() {
	DB.AutoMigrate(&User1{}, &UserInfo{})

	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	// Insert
	DB.Create(&User1{
		Name:   "名稱1",
		Age:    26,
		Gender: true,
		UserInfo: UserInfo{
			Addr: "XXX",
			Like: "coding",
		},
	})

	// Select
	var user User1
	DB.Preload("UserInfo").Take(&user)
	fmt.Println(user)

	var userInfo UserInfo
	DB.Preload("User1").Take(&userInfo)
	fmt.Println(userInfo)

	// Delete
	//var user User1
	DB.Take(&user, 1)
	DB.Select("UserInfo").Delete(&user)
}

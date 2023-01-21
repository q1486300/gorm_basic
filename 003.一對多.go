package main

import (
	"fmt"
	"gorm.io/gorm"
)

// User 用戶表，一個用戶擁有多篇文章
type User struct {
	ID   uint   `gorm:"size:4"`
	Name string `gorm:"size:8;index"` // 用此欄位關聯時，加上索引(index)通過外鍵約束
	// 用戶擁有的文章列表
	Articles []Article // `gorm:"foreignKey:UserName;references:Name"` // references 外鍵關聯的欄位，關聯的欄位需要有索引，不然外鍵約束會報錯
}

// Article 文章表，一篇文章屬於一個用戶
type Article struct {
	ID     uint   `gorm:"size:4"`
	Title  string `gorm:"size:16"`
	UserID uint   `gorm:"size:4"` // 屬於
	// 屬於
	User User // `gorm:"foreignKey:UserName;references:Name"` // references 外鍵關聯的欄位
}

func main() {
	DB.AutoMigrate(&User{}, &Article{})

	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	//Insert()

	//Select()

	//Delete()
}

func Insert() {
	// 創建用戶，帶上文章
	DB.Create(&User{
		Name: "名稱1",
		Articles: []Article{
			{Title: "Golang"},
			{Title: "Java"},
		},
	})

	// 創建文章，關聯已有用戶
	DB.Create(&Article{
		Title:  "歡迎來看Golang",
		UserID: 1,
	})

	DB.Create(&Article{
		Title: "歡迎來看Gorm",
		User: User{
			Name: "名稱2",
		},
	})

	var user User
	DB.Take(&user, 2)
	DB.Create(&Article{
		Title: "名稱2寫的書",
		User:  user,
	})

	// 給已有用戶綁定文章
	var article Article
	DB.Take(&user, 1)
	DB.Take(&article, 6)
	//user.Articles = []Article{article}
	//DB.Save(&user)
	DB.Model(&user).Association("Articles").Append(&article)

	// 給文章綁定用戶
	DB.Take(&user, 1)
	DB.Take(&article, 6)
	//article.UserID = user.ID
	//DB.Save(&article)
	DB.Model(&article).Association("User").Append(&user)
}

func Select() {
	var article Article
	DB.Preload("User").Take(&article)
	fmt.Println(article)

	var user User
	DB.Preload("Articles.User").Take(&user)
	fmt.Println(user)

	//var user User
	DB.Preload("Articles", "id > ?", 2).Take(&user)
	fmt.Println(user)

	//var user User
	DB.Preload("Articles", func(db *gorm.DB) *gorm.DB {
		return db.Where("id > ?", 2)
	}).Take(&user)
	fmt.Println(user)
}

func Delete() {
	// 清除文章的外鍵關係
	var user User
	DB.Preload("Articles").Take(&user, 2)
	DB.Model(&user).Association("Articles").Delete(&user.Articles)
	// 刪除用戶
	DB.Delete(&user)

	// 級聯刪除
	//var user User
	DB.Take(&user, 1)
	DB.Select("Articles").Delete(&user)
}

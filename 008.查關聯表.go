package main

import (
	"fmt"
	"time"
)

type UserModel struct {
	ID       uint
	Name     string
	Collects []ArticleModel2 `gorm:"many2many:user_collect_models;joinForeignKey:UserID;joinReferences:ArticleID"`
}

type ArticleModel2 struct {
	ID    uint
	Title string
	// 這裡也可以反向引用，根據文章查那些用戶收藏了
}

// UserCollectModel 用戶收藏文章表
type UserCollectModel struct {
	UserID        uint          `gorm:"primaryKey"`
	UserModel     UserModel     `gorm:"foreignKey:UserID"`
	ArticleID     uint          `gorm:"primaryKey"`
	ArticleModel2 ArticleModel2 `gorm:"foreignKey:ArticleID"`
	CreatedAt     time.Time
}

func main() {
	DB.SetupJoinTable(&UserModel{}, "Collects", &UserCollectModel{})
	DB.AutoMigrate(&UserModel{}, &ArticleModel2{}, &UserCollectModel{})

	//var user UserModel
	//DB.Preload("Collects").Take(&user, "name = ?", "名稱1")
	//fmt.Println(user)

	var userCollects []UserCollectModel
	DB.Preload("ArticleModel2").Preload("UserModel").Find(&userCollects, 1)
	fmt.Println(userCollects)
}

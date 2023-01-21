package main

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleModel struct {
	ID    uint
	Title string
	Tags  []TagModel `gorm:"many2many:article_tags;joinForeignKey:ArticleID;joinReferences:TagID"`
}

type TagModel struct {
	ID       uint
	Name     string
	Articles []ArticleModel `gorm:"many2many:article_tags;joinForeignKey:TagID;joinReferences:ArticleID"`
}

type ArticleTagModel struct {
	ArticleID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

func main() {
	DB.SetupJoinTable(&ArticleModel{}, "Tags", &ArticleTagModel{})
	DB.SetupJoinTable(&TagModel{}, "Articles", &ArticleTagModel{})
	DB.AutoMigrate(&ArticleModel{}, &TagModel{}, &ArticleTagModel{})

	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	// Insert
	DB.Create(&ArticleModel{
		Title: "Gin基礎",
		Tags: []TagModel{
			{Name: "gin"},
			{Name: "web"},
		},
	})

	// Select
	var articles []ArticleModel
	DB.Preload("Tags").Find(&articles)
	fmt.Println(articles)

}

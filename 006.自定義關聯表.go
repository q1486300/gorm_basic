package main

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Article2 struct {
	ID    uint
	Title string
	Tags  []Tag2 `gorm:"many2many:article_tag2"`
}

type Tag2 struct {
	ID       uint
	Name     string
	Articles []Article2 `gorm:"many2many:article_tag2"`
}

type ArticleTag2 struct {
	Article2ID uint      `gorm:"primaryKey"`
	Tag2ID     uint      `gorm:"primaryKey"`
	CreatedAt  time.Time // CreatedAt 欄位，Gorm 會自動填入目前時間當作預設值，所以不用自己添加 BeforeCreate 的 Hook 函數
}

func main() {
	// 設置 Article 的 Tags 關聯表為 ArticleTag
	DB.SetupJoinTable(&Article2{}, "Tags", &ArticleTag2{})
	// 如果 Tag 要反向引用 Article，那麼也得加上
	DB.SetupJoinTable(&Tag2{}, "Articles", &ArticleTag2{})
	DB.AutoMigrate(&Article2{}, &Tag2{}, &ArticleTag2{})

	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	//添加文章、添加標籤，並自動關聯
	DB.Create(&Article2{
		Title: "Golang基礎",
		Tags: []Tag2{
			{Name: "golang"},
			{Name: "後端"},
			{Name: "web"},
		},
	})

	// 添加文章，關聯已有標籤
	var tags []Tag2
	DB.Find(&tags, "name in (?)", []string{"golang", "web"})
	DB.Create(&Article2{
		Title: "Gorm",
		Tags:  tags,
	})

	// 給已有文章關聯標籤
	var article Article2
	DB.Take(&article, "title = ?", "Gin")
	//var tags []Tag2
	DB.Find(&tags, "name in (?)", []string{"golang", "gin"})
	DB.Model(&article).Association("Tags").Append(tags)

	// 替換已有文章的標籤
	//var article Article2
	DB.Take(&article, "title = ?", "golang基礎")
	//var tags []Tag2
	DB.Find(&tags, "name in (?)", []string{"golang", "後端"})
	DB.Model(&article).Association("Tags").Replace(tags)

	// 查詢文章列表，顯示標籤
	var articles []Article2
	DB.Preload("Tags").Find(&articles)
	fmt.Println(articles)
}

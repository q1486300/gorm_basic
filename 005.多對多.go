package main

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Article1 struct {
	ID    uint
	Title string
	Tags  []Tag `gorm:"many2many:article_tags"`
}

type Tag struct {
	ID       uint
	Name     string
	Articles []Article1 `gorm:"many2many:article_tags"` // 用於反向引用
}

type ArticleTag struct {
	Article1ID uint      `gorm:"primaryKey"`
	TagID      uint      `gorm:"primaryKey"`
	CreateAt   time.Time `json:"create_at"`
}

func (a *ArticleTag) BeforeCreate(db *gorm.DB) error {
	a.CreateAt = time.Now()
	return nil
}

func main() {
	DB.SetupJoinTable(&Article1{}, "Tags", &ArticleTag{}) // 自定義關聯表，使用預設關聯表不用做此設置
	DB.SetupJoinTable(&Tag{}, "Articles", &ArticleTag{})  // 自定義關聯表，使用預設關聯表不用做此設置
	DB.AutoMigrate(&Article1{}, &Tag{}, &ArticleTag{})

	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	// --------- 自定義關聯表的部分 start (多了自定義的 ArticleTag 結構體，作為關聯表) -----------------
	// 添加文章，關聯已有 tag
	var tags []Tag
	DB.Find(&tags, []int{1, 2})
	DB.Create(&Article1{
		Title: "gin基礎",
		Tags:  tags,
	})

	// 給已有文章替換 tag
	var article Article1
	DB.Preload("Tags").Take(&article, 2)
	//var tags []Tag
	DB.Find(&tags, []int{3})
	DB.Model(&article).Association("Tags").Replace(&tags)
	// --------- 自定義關聯表的部分 end -----------------

	// ----------- 以下使用 Gorm 預設的關聯表 -----------
	//Many2ManyInsert()

	//Many2ManySelect()

	//Many2ManyUpdate()
}

func Many2ManyInsert() {
	// Insert
	DB.Create(&Article1{
		Title: "golang基礎",
		Tags: []Tag{
			{Name: "golang"},
			{Name: "後端"},
		},
	})

	// 添加文章關聯已有標籤
	var tags []Tag
	DB.Find(&tags, "name in (?)", []string{"後端"})
	DB.Create(&Article1{
		Title: "Gorm",
		Tags:  tags,
	})
}

func Many2ManySelect() {
	// Select
	var article Article1
	DB.Preload("Tags").Take(&article)
	fmt.Println(article)
}

func Many2ManyUpdate() {
	// 多對多更新
	// 先刪除原有的
	var article Article1
	DB.Preload("Tags").Take(&article, 1)
	DB.Model(&article).Association("Tags").Delete(article.Tags)
	// 再添加新的關聯
	var tag Tag
	DB.Take(&tag, 1)
	DB.Model(&article).Association("Tags").Append(&tag)

	// 使用 gorm 的 Replace()
	//var article Article1
	DB.Preload("Tags").Take(&article, 1)
	//var tag Tag
	DB.Take(&tag, 2)
	DB.Model(&article).Association("Tags").Replace(&tag)
}

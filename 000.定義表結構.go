package main

/*
屬性標籤:
	type			定義欄位類型
	size			定義欄位大小
	column			自定義欄位名
	primaryKey		將欄位定義為主鍵
	unique			將欄位定義為唯一鍵
	default			定義欄位的預設值
	not null		不可為空
	embedded		嵌套欄位
	embeddedPrefix	嵌套欄位前綴
	comment			註解
*/

type _Student struct {
	ID    uint    `gorm:"size:10"`
	Name  string  `gorm:"size:16"`
	Age   int     `gorm:"size:3"`
	Email *string `gorm:"size:128"` // 為了儲存空值，才使用指針
	Type  string  `gorm:"column:_type;size:4"`
	Date  string  `gorm:"default:2023-01-01;comment:日期"`
}

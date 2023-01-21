package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var mysqlLogger logger.Interface

// go get -u gorm.io/driver/mysql
// go get -u gorm.io/gorm
func init() {
	userName := "root"
	password := "your_password"
	host := "127.0.0.1"
	port := 3306
	DBname := "gorm"
	timeout := "10s"

	mysqlLogger = logger.Default.LogMode(logger.Info)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		userName, password, host, port, DBname, timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//SkipDefaultTransaction: true,
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix:   "f_", // 表名前綴
		//	SingularTable: true, // 表名單數
		//	NoLowerCase:   true, // 不要轉換為小寫
		//},
		//Logger: mysqlLogger,
	})
	if err != nil {
		panic("連接資料庫失敗, error= " + err.Error())
	}
	// 連接成功
	DB = db
}

func main() {
	//部分 session 顯示日誌
	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	//Debug() 只想顯示某些語法的日誌
	DB.Debug().AutoMigrate(&_Student{})
}

package main

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
)

type Student struct {
	ID     uint    `gorm:"size:3" json:"id"`
	Name   string  `gorm:"size:8" json:"name"`
	Age    int     `gorm:"size:3" json:"age"`
	Gender bool    `json:"gender"`
	Email  *string `gorm:"size:32" json:"email"`
}

func main() {
	DB.AutoMigrate(&Student{})

	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	//initData()

	var studentList []Student
	DB.Where("name = ?", "名稱5").Find(&studentList)
	fmt.Println(studentList)
	DB.Find(&studentList, "name = ?", "名稱5")
	fmt.Println(studentList)

	fmt.Println(DB.Not("name = ?", "名稱5").Find(&studentList).RowsAffected)
	fmt.Println(DB.Where("name != ?", "名稱5").Find(&studentList).RowsAffected)

	DB.Where("name IN (?)", []string{"名稱2", "名稱8"}).Find(&studentList)
	fmt.Println(studentList)

	DB.Where("name LIKE ?", "名%").Find(&studentList)
	fmt.Println(studentList)

	DB.Where("age > 23 AND email LIKE ?", "%@gmail.com").Find(&studentList)
	DB.Where("age > 23").Where("email LIKE ?", "%@gmail.com").Find(&studentList)
	fmt.Println(studentList)

	DB.Where("gender = ? OR email LIKE ?", false, "%@gmail.com").Find(&studentList)
	DB.Where("gender = ?", false).Or("email LIKE ?", "%@gmail.com").Find(&studentList)
	fmt.Println(studentList)

	// 使用結構體會過濾零值，並且條件用 AND 串接
	DB.Where(&Student{Name: "名稱1", Age: 0}).Find(&studentList)
	// 使用 map 不會過濾零值，並且條件用 AND 串接
	DB.Where(map[string]any{"name": "名稱1", "age": 0}).Find(&studentList)
	fmt.Println(studentList)

	// 只查詢指定欄位
	DB.Select("name", "age").Find(&studentList)
	DB.Select([]string{"name", "age"}).Find(&studentList)
	fmt.Println(studentList)

	// 把查詢結果掃描到結構體中
	type User struct {
		// Scan() 是根據 column 欄位名進行掃描的
		Name123 string `gorm:"column:name"`
		Age     int
	}
	var userList []User
	DB.Select("name", "age").Find(&studentList).Scan(&userList)
	DB.Model(&Student{}).Select("name", "age").Scan(&userList)
	DB.Table("students").Select("name", "age").Scan(&userList)
	fmt.Println(userList)

	// 排序
	DB.Order("age DESC").Find(&studentList)
	fmt.Println(studentList)

	// 分頁
	DB.Limit(2).Offset(0).Find(&studentList)
	DB.Limit(2).Offset(2).Find(&studentList)
	DB.Limit(2).Offset(4).Find(&studentList)
	DB.Limit(2).Offset(6).Find(&studentList)
	DB.Limit(2).Offset(8).Find(&studentList)
	fmt.Println(studentList)

	limit := 2
	page := 1
	DB.Limit(limit).Offset((page - 1) * limit).Find(&studentList)
	fmt.Println(studentList)

	// 去重
	var ageList []int
	DB.Model(&Student{}).Select("age").Distinct("age").Scan(&ageList)
	DB.Model(&Student{}).Select("DISTINCT age").Scan(&ageList)
	fmt.Println(ageList)

	// Group() 分組
	type Group struct {
		Count    int
		Gender   bool
		NameList string
	}
	var groupList []Group
	DB.Model(&Student{}).Select("COUNT(id) AS count", "gender").Group("gender").Scan(&groupList)
	DB.Model(&Student{}).
		Select(
			"GROUP_CONCAT(name) AS name_list",
			"COUNT(id) AS count",
			"gender",
		).
		Group("gender").
		Scan(&groupList)
	fmt.Println(groupList)

	// 原生 SQL 語法
	DB.Raw("SELECT GROUP_CONCAT(name) AS name_list, COUNT(id) AS count, gender FROM students GROUP BY gender;").
		Scan(&groupList)
	fmt.Println(groupList)

	// 子查詢
	DB.Raw("SELECT * FROM students WHERE age > (SELECT AVG(age) FROM students);").Scan(&studentList)
	DB.Model(&Student{}).Where("age > (?)", DB.Model(&Student{}).Select("AVG(age)")).Find(&studentList)
	fmt.Println(studentList)

	// 命名參數
	DB.Where("name = ? AND age = ?", "名稱7", 26).Find(&studentList)
	DB.Where("name = @name AND age = @age",
		sql.Named("name", "名稱7"),
		sql.Named("age", 26)).Find(&studentList)
	DB.Where("name = @name AND age = @age",
		map[string]any{"name": "名稱7", "age": 26}).Find(&studentList)
	fmt.Println(studentList)

	// 查詢引用 Scopes()
	var res []map[string]any
	DB.Model(&Student{}).Scopes(AgeGreaterThan23).Find(&res)
	fmt.Println(res)
}

func initData() {
	var studentList []Student
	DB.Find(&studentList).Delete(&studentList)

	studentList = []Student{
		{ID: 1, Name: "名稱1", Age: rand.Intn(30) + 1, Gender: true, Email: PtrString("test01@gmail.com")},
		{ID: 2, Name: "名稱2", Age: rand.Intn(30) + 1, Gender: true, Email: PtrString("test02@gmail.com")},
		{ID: 3, Name: "名稱3", Age: rand.Intn(30) + 1, Gender: true, Email: PtrString("test03@gmail.com")},
		{ID: 4, Name: "名稱4", Age: rand.Intn(30) + 1, Gender: false, Email: PtrString("test04@gmail.com")},
		{ID: 5, Name: "名稱5", Age: rand.Intn(30) + 1, Gender: true, Email: PtrString("test05@gmail.com")},
		{ID: 6, Name: "名稱6", Age: rand.Intn(30) + 1, Gender: false, Email: PtrString("test06@gmail.com")},
		{ID: 7, Name: "名稱7", Age: rand.Intn(30) + 1, Gender: true, Email: PtrString("test07@gmail.com")},
		{ID: 8, Name: "名稱8", Age: rand.Intn(30) + 1, Gender: false, Email: PtrString("test08@gmail.com")},
		{ID: 9, Name: "名稱9", Age: rand.Intn(30) + 1, Gender: true, Email: PtrString("test09@gmail.com")},
	}
	DB.Create(&studentList)
}

func PtrString(s string) *string {
	return &s
}

func AgeGreaterThan23(db *gorm.DB) *gorm.DB {
	return db.Where("age > ?", 23)
}

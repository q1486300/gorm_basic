package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type __Student struct {
	ID     uint    `gorm:"size:3" json:"id"`
	Name   string  `gorm:"size:8" json:"name"`
	Age    int     `gorm:"size:3" json:"age"`
	Gender bool    `json:"gender"`
	Email  *string `gorm:"size:32" json:"email"`
}

// Hook Function
func (s *__Student) BeforeCreate(db *gorm.DB) error {
	email := "test@gmail.com"
	s.Email = &email
	return nil
}

func main() {
	DB.AutoMigrate(&__Student{})

	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	//insertFunc()

	//selectFunc()

	//updateFunc()

	//deleteFunc()

	// test Hook
	DB.Create(&__Student{Name: "測試 Hook", Age: 22})
}

func insertFunc() {
	// 添加紀錄
	email := "abc123@gmail.com"
	s1 := __Student{
		Name:   "名稱1",
		Age:    26,
		Gender: true,
		Email:  &email,
	}
	err := DB.Create(&s1).Error
	fmt.Printf("%#v\n", s1)
	fmt.Println(err)

	// 批量插入
	var studentList []__Student

	for i := 0; i < 10; i++ {
		studentList = append(studentList, __Student{
			Name:   fmt.Sprintf("名稱%d", i+1),
			Age:    26 + i + 1,
			Gender: true,
			Email:  &email,
		})
	}

	err = DB.Create(&studentList).Error
	fmt.Println(err)
}

func selectFunc() {
	// 單條紀錄的查詢
	var student __Student
	// SELECT * FROM `students` LIMIT 1
	DB.Take(&student)
	fmt.Println(student)
	student = __Student{}
	// SELECT * FROM `students` ORDER BY `students`.`id` LIMIT 1
	DB.First(&student)
	fmt.Println(student)
	student = __Student{}
	// SELECT * FROM `students` ORDER BY `students`.`id` DESC LIMIT 1
	DB.Last(&student)
	fmt.Println(student)

	err := DB.Take(&student, 30).Error
	fmt.Println(err == gorm.ErrRecordNotFound)
	fmt.Println(student)

	// 使用 ? 防止 SQL注入
	DB.Take(&student, "name = ?", "名稱9")
	fmt.Println(student)

	// 根據主鍵值查詢
	student.ID = 2
	err = DB.Take(&student).Error
	switch err {
	case gorm.ErrRecordNotFound:
		fmt.Println("沒有找到")
	default:
		fmt.Println("sql 錯誤")
	}
	fmt.Println(student)

	// 查詢多條紀錄
	var studentList []__Student

	count := DB.Find(&studentList).RowsAffected
	fmt.Println(count)
	for _, student := range studentList {
		fmt.Println(student)
	}
	// JSON 序列化之後，Email 會被轉化為指針內的值
	data, _ := json.Marshal(studentList)
	fmt.Println(string(data))

	DB.Find(&studentList, []int{3, 6, 9})
	fmt.Println(studentList)

	DB.Find(&studentList, "name in (?)", []string{"名稱1", "名稱2"})
	fmt.Println(studentList)
}

func updateFunc() {
	// save 更新所有欄位的值
	var student __Student
	DB.Take(&student, 11)
	student.Age = 20
	student.Name = "改名1"
	DB.Save(&student)
	// Select() 只更新指定欄位
	DB.Select("name").Save(&student)

	// 批量更新
	var studentList []__Student
	DB.Find(&studentList, []int{11, 12, 13}).Update("gender", false)
	DB.Model(&__Student{}).Where("id in (?)", []int{11, 12, 13}).Update("gender", true)

	// 更新多欄
	// 使用結構體不會更新零值
	DB.Model(&__Student{}).Where("id in (?)", []int{11, 12, 13}).Updates(__Student{
		Age:    20,
		Gender: false, // 不會被更新
	})
	// 使用 Select() 可以更新結構體零值
	DB.Model(&__Student{}).Where("id in (?)", []int{11, 12, 13}).
		Select("age", "gender").Updates(__Student{
		Age:    20,
		Gender: false,
	})

	// 使用 map 可以更新零值
	DB.Model(&__Student{}).Where("id in (?)", []int{11, 12, 13}).Updates(map[string]any{
		"age":    30,
		"gender": false,
	})
}

func deleteFunc() {
	// 刪除
	var student __Student
	DB.Delete(&student, 13)
	DB.Delete(&student, []int{11, 12})

	DB.Take(&student)
	DB.Delete(&student)
}

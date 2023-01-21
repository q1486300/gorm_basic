package main

import (
	"fmt"
	"gorm.io/gorm"
)

type _User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Money int    `json:"money"`
}

func main() {
	DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&_User{})

	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	DB.Create([]_User{
		{Name: "名稱1", Money: 1000},
		{Name: "名稱2", Money: 1000},
	})

	generalTransaction()
	manualTransaction()
}

func generalTransaction() {
	var user1, user2 _User
	DB.Take(&user1, "name = ?", "名稱1")
	DB.Take(&user2, "name = ?", "名稱2")

	// user1 給 user2 轉帳 100 元
	DB.Transaction(func(tx *gorm.DB) error {
		// 先給 user1 -100
		user1.Money -= 100
		err := tx.Model(&user1).Update("money", user1.Money).Error
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 再給 user2 + 100
		user2.Money += 100
		err = tx.Model(&user2).Update("money", user2.Money).Error
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 提交 Transaction
		return nil
	})
}

func manualTransaction() {
	var user1, user2 _User
	DB.Take(&user1, "name = ?", "名稱1")
	DB.Take(&user2, "name = ?", "名稱2")

	// user2 給 user1 轉帳 100 元
	tx := DB.Begin()

	// 先給 user2 -100
	user2.Money -= 100
	err := tx.Model(&user2).Update("money", user2.Money).Error
	if err != nil {
		tx.Rollback()
	}

	// 再給 user1 +100
	user1.Money += 100
	err = tx.Model(&user1).Update("money", user1.Money).Error
	if err != nil {
		tx.Rollback()
	}

	// 提交 Transaction
	tx.Commit()
}

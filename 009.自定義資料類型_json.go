package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Info struct {
	Status string `json:"status"`
	Addr   string `json:"addr"`
	Age    int    `json:"age"`
}

// Scan 從資料庫中讀取出來
func (i *Info) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal Json value:", value))
	}
	err := json.Unmarshal(bytes, i)
	return err
}

// Value 存入資料庫
func (i Info) Value() (driver.Value, error) {
	fmt.Printf("入庫前，%#v, %T\n", i, i)
	return json.Marshal(i)
}

type AuthModel struct {
	ID   uint
	Name string
	Info Info `gorm:"type:string"`
}

func main() {
	DB.AutoMigrate(&AuthModel{})

	DB.Create(&AuthModel{
		Name: "名稱1",
		Info: Info{
			Status: "success",
			Addr:   "台灣",
			Age:    26,
		},
	})

	var auth AuthModel
	DB.Take(&auth, "name = ?", "名稱1")
	fmt.Println(auth)
}

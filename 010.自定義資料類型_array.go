package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

type Array []string

func (a *Array) Scan(value any) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("解析失敗: %v %T", value, value))
	}
	*a = strings.Split(string(data), "|")
	return nil
}

func (a Array) Value() (driver.Value, error) {
	return strings.Join(a, "|"), nil
}

type HostModel struct {
	ID    uint
	IP    string
	Ports Array `gorm:"type:string"`
}

func main() {
	DB.AutoMigrate(&HostModel{})

	DB.Create(&HostModel{
		IP:    "192.168.100.1",
		Ports: []string{"8080", "3306"},
	})

	var host []HostModel
	DB.Find(&host)
	fmt.Println(host)
}

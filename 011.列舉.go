package main

import (
	"encoding/json"
	"fmt"
)

const (
	Running Status = iota + 1
	OffLine
	Except
)

type Status int

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s Status) String() string {
	var str string
	switch s {
	case Running:
		str = "Running"
	case OffLine:
		str = "OffLine"
	case Except:
		str = "Except"
	}
	return str
}

type Host struct {
	ID     uint   `json:"id"`
	IP     string `json:"ip"`
	Status Status `gorm:"size:8" json:"status"`
}

func main() {
	DB.AutoMigrate(&Host{})

	DB.Create(&Host{
		IP:     "192.168.100.1",
		Status: Running,
	})

	var host Host
	DB.Take(&host)
	
	fmt.Println(host)
	fmt.Printf("Status: %#v %T\n", host.Status, host.Status)

	data, _ := json.Marshal(host)
	fmt.Println(string(data))
}

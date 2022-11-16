package models

import (
	"time"
)

type Dependences struct {
	Id        int       `xorm:"pk autoincr INTEGER"`
	Name      string    `xorm:"VARCHAR(255)"`
	Type      string    `xorm:"NUMBER"`
	Timestamp string    `xorm:"VARCHAR(255)"`
	Status    string    `xorm:"NUMBER"`
	Log       string    `xorm:"JSON"`
	Remark    string    `xorm:"VARCHAR(255)"`
	Createdat time.Time `xorm:"not null DATETIME created"`
	Updatedat time.Time `xorm:"not null DATETIME updated"`
}

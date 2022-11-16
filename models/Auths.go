package models

import (
	"time"
)

type Auths struct {
	Id        int       `xorm:"pk autoincr INTEGER"`
	Ip        string    `xorm:"VARCHAR(255)"`
	Type      string    `xorm:"VARCHAR(255)"`
	Info      string    `xorm:"JSON"`
	Createdat time.Time `xorm:"not null DATETIME created"`
	Updatedat time.Time `xorm:"not null DATETIME updated"`
}

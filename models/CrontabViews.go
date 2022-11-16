package models

import (
	"time"
)

type Crontabviews struct {
	Id         int       `xorm:"pk autoincr INTEGER"`
	Name       string    `xorm:"VARCHAR(255) unique"`
	Position   string    `xorm:"TINYINT(1)"`
	Isdisabled string    `xorm:"TINYINT(1)"`
	Filters    []any     `xorm:"JSON"`
	Sorts      []any     `xorm:"JSON"`
	Createdat  time.Time `xorm:"not null DATETIME created"`
	Updatedat  time.Time `xorm:"not null DATETIME updated"`
}

package models

import (
	"time"
)

type Apps struct {
	Id           int       `xorm:"pk autoincr INTEGER"`
	Name         string    `xorm:"VARCHAR(255) unique"`
	Scopes       string    `xorm:"JSON"`
	ClientId     string    `xorm:"VARCHAR(255)"`
	ClientSecret string    `xorm:"VARCHAR(255)"`
	Tokens       string    `xorm:"JSON"`
	Createdat    time.Time `xorm:"not null DATETIME created"`
	Updatedat    time.Time `xorm:"not null DATETIME updated"`
}

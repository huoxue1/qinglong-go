package models

import (
	"time"
)

type Subscriptions struct {
	Id               int       `xorm:"pk autoincr INTEGER"`
	Name             string    `xorm:"VARCHAR(255)"`
	Url              string    `xorm:"VARCHAR(255)"`
	Schedule         string    `xorm:"VARCHAR(255)"`
	IntervalSchedule string    `xorm:"JSON"`
	Type             string    `xorm:"VARCHAR(255)"`
	Whitelist        string    `xorm:"VARCHAR(255)"`
	Blacklist        string    `xorm:"VARCHAR(255)"`
	Status           string    `xorm:"NUMBER"`
	Dependences      string    `xorm:"VARCHAR(255)"`
	Extensions       string    `xorm:"VARCHAR(255)"`
	SubBefore        string    `xorm:"VARCHAR(255)"`
	SubAfter         string    `xorm:"VARCHAR(255)"`
	Branch           string    `xorm:"VARCHAR(255)"`
	PullType         string    `xorm:"VARCHAR(255)"`
	PullOption       string    `xorm:"JSON"`
	Pid              string    `xorm:"NUMBER"`
	IsDisabled       string    `xorm:"NUMBER"`
	LogPath          string    `xorm:"VARCHAR(255)"`
	ScheduleType     string    `xorm:"VARCHAR(255)"`
	Alias            string    `xorm:"VARCHAR(255) unique"`
	Createdat        time.Time `xorm:"not null DATETIME created"`
	Updatedat        time.Time `xorm:"not null DATETIME updated"`
}

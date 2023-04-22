package models

import (
	"fmt"
	"xorm.io/builder"
)

type Subscriptions struct {
	Id               int            `xorm:"pk autoincr INTEGER" json:"id"`
	Name             string         `xorm:"TEXT" json:"name"`
	Url              string         `xorm:"TEXT" json:"url"`
	Schedule         string         `xorm:"TEXT" json:"schedule"`
	IntervalSchedule map[string]any `xorm:"JSON" json:"interval_schedule"`
	Type             string         `xorm:"TEXT" json:"type"`
	Whitelist        string         `xorm:"TEXT" json:"whitelist"`
	Blacklist        string         `xorm:"TEXT" json:"blacklist"`
	Status           int            `xorm:"INTEGER default(1)" json:"status"`
	Dependences      string         `xorm:"TEXT" json:"dependences"`
	Extensions       string         `xorm:"TEXT" json:"extensions"`
	SubBefore        string         `xorm:"TEXT" json:"sub_before"`
	SubAfter         string         `xorm:"TEXT" json:"sub_after"`
	Branch           string         `xorm:"TEXT" json:"branch"`
	PullType         string         `xorm:"TEXT" json:"pull_type"`
	PullOption       string         `xorm:"JSON" json:"pull_option"`
	Pid              int            `xorm:"INTEGER" json:"pid"`
	IsDisabled       int            `xorm:"INTEGER" json:"is_disabled"`
	LogPath          string         `xorm:"TEXT" json:"log_path"`
	ScheduleType     string         `xorm:"TEXT" json:"schedule_type"`
	Alias            string         `xorm:"TEXT unique" json:"alias"`
	Createdat        string         `xorm:"not null DATETIME created" json:"createdat"`
	Updatedat        string         `xorm:"not null DATETIME updated" json:"updatedat"`
}

func QuerySubscription(searchValue string) ([]*Subscriptions, error) {
	subscription := make([]*Subscriptions, 0)
	session := engine.Table(new(Subscriptions)).
		Where(
			builder.Like{"name", "%" + searchValue + "%"}.
				Or(builder.Like{"url", "%" + searchValue + "%"}))
	err := session.Find(&subscription)
	if err != nil {
		return nil, err
	}
	return subscription, err

}

func AddSubscription(subscription *Subscriptions) (int, error) {
	_, err := engine.Table(subscription).Insert(subscription)
	if err != nil {
		return 0, err
	}
	_, _ = engine.Where("name=?", subscription.Name).Get(subscription)
	return subscription.Id, err
}

func GetSubscription(id int) (*Subscriptions, error) {
	env := new(Subscriptions)
	_, err := engine.ID(id).Get(env)
	return env, err
}

func UpdateSubscription(subscription *Subscriptions) error {
	_, err := engine.Table(subscription).ID(subscription.Id).AllCols().Update(subscription)
	return err
}

func DeleteSubscription(id int) error {
	_, err := engine.Table(new(Subscriptions)).Delete(&Subscriptions{Id: id})
	return err
}

func (s *Subscriptions) GetCron() string {
	if s.ScheduleType == "interval" {
		t := s.IntervalSchedule["type"].(string)
		return fmt.Sprintf("@every %v%s", s.IntervalSchedule["value"], string(t[0]))
	} else {
		return s.Schedule
	}
}

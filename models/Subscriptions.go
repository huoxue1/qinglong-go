package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"xorm.io/builder"
)

type Subscriptions struct {
	Id               int            `xorm:"pk autoincr INTEGER" json:"id"`
	Name             string         `json:"name"`
	Url              string         `json:"url"`
	Schedule         string         `json:"schedule"`
	IntervalSchedule map[string]any `xorm:"JSON" json:"interval_schedule"`
	Type             string         `json:"type"`
	Whitelist        string         `json:"whitelist"`
	Blacklist        string         `json:"blacklist"`
	Status           int            `xorm:"INTEGER default(1)" json:"status"`
	Dependences      string         `json:"dependences"`
	Extensions       string         `json:"extensions"`
	SubBefore        string         `json:"sub_before"`
	SubAfter         string         `json:"sub_after"`
	Branch           string         `json:"branch"`
	PullType         string         `json:"pull_type"`
	PullOption       string         `xorm:"JSON" json:"pull_option"`
	Pid              int            `xorm:"INTEGER" json:"pid"`
	IsDisabled       int            `xorm:"INTEGER" json:"is_disabled"`
	LogPath          string         `json:"log_path"`
	ScheduleType     string         `json:"schedule_type"`
	Alias            string         `xorm:"varchar(255) unique" json:"alias"`
	Createdat        string         `xorm:"not null DATETIME created" json:"createdat"`
	Updatedat        string         `xorm:"not null DATETIME updated" json:"updatedat"`

	File io.WriteCloser `json:"-" xorm:"-"`
}

func (s *Subscriptions) Close() error {
	return s.File.Close()
}

func (s *Subscriptions) Write(p []byte) (n int, err error) {
	if s.File == nil {
		s.LogPath = "data/log/" + time.Now().Format("2006-01-02") + "/" + s.Alias + "_" + uuid.New().String() + ".log"
		s.Status = 1
		_ = UpdateSubscription(s)
		_ = os.MkdirAll(filepath.Dir(s.LogPath), 0666)
		s.File, _ = os.OpenFile(s.LogPath, os.O_CREATE|os.O_RDWR, 0666)
	}
	return s.File.Write(p)
}
func (s *Subscriptions) WriteString(data string) (n int, err error) {
	p := []byte(data)
	if s.File == nil {
		s.LogPath = "data/log/" + time.Now().Format("2006-01-02") + "/" + s.Alias + "_" + uuid.New().String() + ".log"
		s.Status = 1
		_ = UpdateSubscription(s)
		_ = os.MkdirAll(filepath.Dir(s.LogPath), 0666)
		s.File, _ = os.OpenFile(s.LogPath, os.O_CREATE|os.O_RDWR, 0666)
	}
	return s.File.Write(p)
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

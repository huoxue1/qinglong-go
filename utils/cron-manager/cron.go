package cron_manager

import (
	"errors"
	"github.com/huoxue1/qinglong-go/utils/log"
	"github.com/robfig/cron/v3"
	"strings"
	"sync"
)

var (
	manager     sync.Map
	defaultCron *cron.Cron
	SixCron     *cron.Cron
)

type mapValue struct {
	en cron.EntryID
	c  *cron.Cron
}

func init() {
	defaultCron = cron.New(cron.WithChain(cron.Recover(&log.CronLog{})))
	SixCron = cron.New(cron.WithChain(cron.Recover(&log.CronLog{})), cron.WithParser(
		cron.NewParser(cron.Second|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)))
	defaultCron.Start()
	SixCron.Start()
}

func AddCron(id string, value string, task func()) error {
	if value == "7 7 7 7 7" {
		value = "7 7 7 7 6"
	}
	crons := strings.Split(value, " ")
	cronCmd := defaultCron
	if len(crons) == 6 {
		cronCmd = SixCron
	}
	en, err := cronCmd.AddFunc(value, task)
	if err != nil {
		return err
	}
	manager.Store(id, &mapValue{en, cronCmd})
	return nil
}

func DeleteCron(id string) error {
	value, loaded := manager.LoadAndDelete(id)
	if !loaded {
		return errors.New("the cron " + id + " not found!")
	}
	v := value.(*mapValue)
	v.c.Remove(v.en)
	return nil
}

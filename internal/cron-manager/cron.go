package cron_manager

import (
	"errors"
	"github.com/huoxue1/qinglong-go/utils/log"
	"github.com/robfig/cron/v3"
	"sync"
)

var (
	manager sync.Map
	SixCron *cron.Cron
)

type mapValue struct {
	en cron.EntryID
	c  *cron.Cron
}

func init() {
	SixCron = cron.New(cron.WithChain(cron.Recover(&log.CronLog{})), cron.WithParser(
		cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)))
	SixCron.Start()
}

func AddCron(id string, value string, task func()) error {
	if value == "7 7 7 7 7" {
		value = "7 7 7 7 6"
	}

	en, err := SixCron.AddFunc(value, task)
	if err != nil {
		return err
	}
	manager.Store(id, &mapValue{en, SixCron})
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

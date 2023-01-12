package cron

import (
	"github.com/huoxue1/qinglong-go/models"
	"github.com/robfig/cron/v3"
	"time"
)

func GetCrons(page, size int, searchValue string, sorter map[string]string, filter string) ([]*models.Crontabs, error) {
	crontabs, err := models.QueryCron(page, size, searchValue, sorter["field"], sorter["type"])
	return crontabs, err
}

func AddCron(cron *models.Crontabs) (int, error) {
	AddTask(cron)
	return models.AddCron(cron)
}

func DeleteCron(ids []int) error {
	for _, id := range ids {

		c, _ := manager.Load(id)
		if c != nil {
			c.(*cron.Cron).Stop()
		}

		err := models.DeleteCron(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateCron(c1 *models.Crontabs) error {
	crontabs, _ := models.GetCron(c1.Id)
	crontabs.Name = c1.Name
	crontabs.Command = c1.Command
	crontabs.Labels = c1.Labels
	crontabs.Schedule = c1.Schedule
	crontabs.Updatedat = time.Now().Format(time.RFC3339)

	c, _ := manager.Load(c1.Id)
	if c != nil {
		c.(*cron.Cron).Stop()
	}
	AddTask(c1)

	return models.UpdateCron(crontabs)
}

func DisableCron(ids []int) error {
	for _, id := range ids {

		c, _ := manager.Load(id)
		c.(*cron.Cron).Stop()

		cron, err := models.GetCron(id)
		if err != nil {
			continue
		}
		cron.Isdisabled = 1
		err = models.UpdateCron(cron)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnableCron(ids []int) error {
	for _, id := range ids {
		cron, err := models.GetCron(id)
		if err != nil {
			continue
		}
		AddTask(cron)
		cron.Isdisabled = 0
		err = models.UpdateCron(cron)
		if err != nil {
			return err
		}
	}
	return nil
}

func PinCron(ids []int) error {
	for _, id := range ids {
		cron, err := models.GetCron(id)
		if err != nil {
			continue
		}
		cron.Ispinned = 1
		err = models.UpdateCron(cron)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnPinCron(ids []int) error {
	for _, id := range ids {
		cron, err := models.GetCron(id)
		if err != nil {
			continue
		}
		cron.Ispinned = 0
		err = models.UpdateCron(cron)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunCron(ids []int) error {
	for _, id := range ids {
		crontab, err := models.GetCron(id)
		if err != nil {
			continue
		}
		runCron(crontab)
	}
	return nil
}

func StopCron(ids []int) error {
	for _, id := range ids {
		crontab, err := models.GetCron(id)
		if err != nil {
			continue
		}
		stopCron(crontab)
	}
	return nil
}

func GetCron(id int) *models.Crontabs {
	cron, err := models.GetCron(id)
	if err != nil {
		return nil
	}
	return cron
}

package models

import (
	"errors"
	"xorm.io/builder"
)

type Crontabs struct {
	Id                int      `xorm:"pk autoincr INTEGER" json:"id"`
	Name              string   `xorm:"VARCHAR(255)" json:"name"`
	Command           string   `xorm:"VARCHAR(255)" json:"command"`
	Schedule          string   `xorm:"VARCHAR(255)" json:"schedule"`
	Timestamp         string   `xorm:"VARCHAR(255)" json:"timestamp"`
	Saved             bool     `xorm:"TINYINT(1)" json:"saved"`
	Status            int      `xorm:"TINYINT(1)" json:"status"`
	Issystem          int      `xorm:"TINYINT(1)" json:"isSystem"`
	Pid               int      `xorm:"TINYINT(1)" json:"pid"`
	Isdisabled        int      `xorm:"TINYINT(1)" json:"isDisabled"`
	Ispinned          int      `xorm:"TINYINT(1)" json:"isPinned"`
	LogPath           string   `xorm:"VARCHAR(255)" json:"log_path"`
	Labels            []string `xorm:"JSON" json:"labels"`
	LastRunningTime   int64    `xorm:"NUMBER" json:"last_running_time"`
	LastExecutionTime int64    `xorm:"NUMBER" json:"last_execution_time"`
	Createdat         string   `xorm:"not null TEXT created" json:"createdAt"`
	Updatedat         string   `xorm:"not null TEXT updated" json:"updatedAt"`
}

func QueryCron(page int, size int, searchValue string, orderField string, orderType string) ([]*Crontabs, error) {
	crontabs := make([]*Crontabs, 0)
	session := engine.Table(new(Crontabs)).Limit(size, (page-1)*size).Where(builder.Like{"name", "%" + searchValue + "%"}.Or(builder.Like{"command", "%" + searchValue + "%"}))
	if orderType == "DESC" {
		session.Desc(orderField)
	} else if orderType == "ASC" {
		session.Asc(orderField)
	}
	err := session.Find(&crontabs)
	return crontabs, err
}

func SetAllCornStop() {
	_, _ = engine.Table(new(Crontabs)).Where("status=?", 0).Update(map[string]any{"status": 1})
}

func QueryRunningCron() ([]*Crontabs, error) {
	crontabs := make([]*Crontabs, 0)
	session := engine.Table(new(Crontabs)).Where("status=?", 0)
	err := session.Find(&crontabs)
	return crontabs, err
}

func QueryCronByDir(dir string) ([]*Crontabs, error) {
	crontabs := make([]*Crontabs, 0)
	session := engine.Table(new(Crontabs)).Where(builder.Like{"command", "task " + dir + "%"})
	err := session.Find(&crontabs)
	return crontabs, err
}

func FindAllEnableCron() []*Crontabs {
	crontabs := make([]*Crontabs, 0)
	err := engine.Table(new(Crontabs)).Where("isdisabled=?", 0).Find(&crontabs)
	if err != nil {
		return nil
	}
	return crontabs
}

func GetCronByCommand(command string) (*Crontabs, error) {
	cron := new(Crontabs)
	count, _ := engine.Where("command=?", command).Count(cron)
	if count < 1 {
		return nil, errors.New("not found")
	}
	_, err := engine.Where("command=?", command).Get(cron)
	return cron, err
}

func GetCron(id int) (*Crontabs, error) {
	cron := new(Crontabs)
	_, err := engine.ID(id).Get(cron)
	return cron, err
}

func AddCron(cron *Crontabs) (int, error) {
	_, err := engine.InsertOne(cron)
	if err != nil {
		return 0, err
	}
	_, _ = engine.Where("name=?", cron.Name).Get(cron)
	return cron.Id, err
}

func UpdateCron(cron *Crontabs) error {
	_, err := engine.Table(cron).ID(cron.Id).AllCols().Update(cron)
	return err
}

func DeleteCron(id int) error {
	_, err := engine.Table(new(Crontabs)).Delete(&Crontabs{Id: id})
	return err
}

func Count(searchValue string) int64 {
	count, _ := engine.Table(new(Crontabs)).
		Where(
			builder.Like{"name", "%" + searchValue + "%"}.
				Or(builder.Like{"command", "%" + searchValue + "%"})).
		Count()
	return count
}

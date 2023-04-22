package models

import (
	"time"
	"xorm.io/builder"
)

type Envs struct {
	Id          int       `xorm:"pk autoincr INTEGER" json:"id"`
	Value       string    `xorm:"TEXT" json:"value"`
	Timestamp   string    `xorm:"TEXT" json:"timestamp"`
	Status      int       `xorm:"TINYINT(1)" json:"status"`
	Position    string    `xorm:"TINYINT(1)" json:"position"`
	Name        string    `xorm:"TEXT" json:"name"`
	Remarks     string    `xorm:"TEXT" json:"remarks"`
	Createdat   time.Time `xorm:"not null DATETIME created" json:"createdAt"`
	Updatedat   time.Time `xorm:"not null DATETIME updated" json:"updatedAt"`
	SerialIndex int64     `xorm:"serial_index INTEGER" json:"serialIndex"`
}

func QueryEnv(searchValue string) ([]*Envs, error) {
	envs := make([]*Envs, 0)
	session := engine.Table(new(Envs)).
		Where(
			builder.Like{"name", "%" + searchValue + "%"}.
				Or(builder.Like{"value", "%" + searchValue + "%"}).
				Or(builder.Like{"remarks", "%" + searchValue + "%"})).Asc("serial_index")
	err := session.Find(&envs)
	if err != nil {
		return nil, err
	}
	return envs, err

}

func QueryEnvByIndex(from, to int64) ([]*Envs, error) {
	envs := make([]*Envs, 0)
	err := engine.Table(new(Envs)).Where(builder.Between{
		Col:     "serial_index",
		LessVal: from,
		MoreVal: to,
	}).Find(&envs)
	return envs, err
}

func AddEnv(envs *Envs) (int, error) {
	count, _ := engine.Table(envs).Count()
	envs.SerialIndex = count
	_, err := engine.Table(envs).Insert(envs)
	if err != nil {
		return 0, err
	}
	_, _ = engine.Where("name=? and value=?", envs.Name, envs.Value).Get(envs)
	return envs.Id, err
}

func GetEnv(id int) (*Envs, error) {
	env := new(Envs)
	_, err := engine.ID(id).Get(env)
	return env, err
}

func UpdateEnv(env *Envs) error {
	_, err := engine.Table(env).ID(env.Id).AllCols().Update(env)
	return err
}

func DeleteEnv(id int) error {
	_, err := engine.Table(new(Envs)).Delete(&Envs{Id: id})
	return err
}

package models

import (
	"time"
	"xorm.io/builder"
)

type Envs struct {
	Id        int       `xorm:"pk autoincr INTEGER" json:"id,omitempty"`
	Value     string    `xorm:"TEXT" json:"value,omitempty"`
	Timestamp string    `xorm:"TEXT" json:"timestamp,omitempty"`
	Status    *int      `xorm:"TINYINT(1)" json:"status,omitempty"`
	Position  string    `xorm:"TINYINT(1)" json:"position,omitempty"`
	Name      string    `xorm:"TEXT" json:"name,omitempty"`
	Remarks   string    `xorm:"TEXT" json:"remarks,omitempty"`
	Createdat time.Time `xorm:"not null DATETIME created" json:"createdAt"`
	Updatedat time.Time `xorm:"not null DATETIME updated" json:"updatedAt"`
}

func QueryEnv(searchValue string) ([]*Envs, error) {
	envs := make([]*Envs, 0)
	session := engine.Table(new(Envs)).
		Where(
			builder.Like{"name", "%" + searchValue + "%"}.
				Or(builder.Like{"value", "%" + searchValue + "%"}).
				Or(builder.Like{"remarks", "%" + searchValue + "%"}))
	err := session.Find(&envs)
	if err != nil {
		return nil, err
	}
	return envs, err

}

func AddEnv(envs *Envs) (int, error) {
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

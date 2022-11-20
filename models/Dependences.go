package models

import (
	"time"
	"xorm.io/builder"
)

const (
	NODE   = 0
	PYTHON = 1
	LINUX  = 2
)

type Dependences struct {
	Id        int       `xorm:"pk autoincr INTEGER" json:"id,omitempty"`
	Name      string    `xorm:"TEXT" json:"name,omitempty"`
	Type      int       `xorm:"INTEGER" json:"type,omitempty"`
	Timestamp string    `xorm:"TEXT" json:"timestamp,omitempty"`
	Status    int       `xorm:"INTEGER" json:"status,omitempty"`
	Log       []string  `xorm:"JSON" json:"log,omitempty"`
	Remark    string    `xorm:"TEXT" json:"remark,omitempty"`
	Createdat time.Time `xorm:"not null DATETIME created" json:"createdAt"`
	Updatedat time.Time `xorm:"not null DATETIME updated" json:"updatedAt"`
}

func QueryDependences(searchValue string) ([]*Dependences, error) {
	dep := make([]*Dependences, 0)
	session := engine.Table(new(Dependences)).
		Where(
			builder.Like{"name", "%" + searchValue + "%"})
	err := session.Find(&dep)
	if err != nil {
		return nil, err
	}
	return dep, err

}

func AddDependences(dep *Dependences) (int, error) {
	_, err := engine.Table(dep).Insert(dep)
	if err != nil {
		return 0, err
	}
	_, _ = engine.Where("name=?", dep.Name).Get(dep)
	return dep.Id, err
}

func GetDependences(id int) (*Dependences, error) {
	env := new(Dependences)
	_, err := engine.ID(id).Get(env)
	return env, err
}

func UpdateDependences(dep *Dependences) error {
	_, err := engine.Table(dep).ID(dep.Id).AllCols().Update(dep)
	return err
}

func DeleteDependences(id int) error {
	_, err := engine.Table(new(Dependences)).Delete(&Dependences{Id: id})
	return err
}

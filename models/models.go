package models

import (
	log2 "github.com/huoxue1/qinglong-go/utils/log"
	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var (
	engine *xorm.Engine
)

func init() {
	en, err := xorm.NewEngine("sqlite", "./data/db/database.sqlite")
	if err != nil {
		log.Errorln("[sql] " + err.Error())
		return
	}
	_ = en.Sync2(new(Apps), new(Auths), new(Crontabs), new(Crontabviews), new(Dependences), new(Envs), new(Subscriptions))
	en.ShowSQL(true)
	en.SetLogger(new(log2.MyLog))
	engine = en
}

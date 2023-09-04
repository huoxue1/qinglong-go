package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/huoxue1/go-utils/base/log"
	xLog "github.com/huoxue1/go-utils/base/log/xorm"
	_ "modernc.org/sqlite"
	"os"
	"xorm.io/xorm"
)

var (
	engine *xorm.Engine
)

func InitModels() {
	_ = os.MkdirAll("data/db", 0666)
	en, err := xorm.NewEngine("sqlite", "data/db/database.sqlite")
	//en, err := xorm.NewEngine("mysql", "root:123@tcp(127.0.0.1:3306)/ql?charset=utf8")

	if err != nil {
		log.Errorln("[sql] " + err.Error())
		return
	}
	err = en.Sync2(new(Apps), new(Auths), new(Crontabs), new(Crontabviews), new(Dependences), new(Envs), new(Subscriptions))
	if err != nil {
		log.Errorln("[sql] " + err.Error())
		return
	}
	en.SetLogger(xLog.GetXormLogger(log.StandardLogger(), "info", false))
	engine = en
}
